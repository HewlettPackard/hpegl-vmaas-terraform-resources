// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/auth"
)

type token interface {
	setScmClientToken(ctx *context.Context, meta interface{})
}

type tokenStruct struct{}

func (t *tokenStruct) setScmClientToken(ctx *context.Context, meta interface{}) {
	auth.SetScmClientToken(ctx, meta)
}

// CondFunc function accepts response and error of the RetryFunc. If any error returns
// retry will terminated and returns the error
type CondFunc func(response interface{}, ResponseErr error) (bool, error)

// RetryFunc accepts ctx as parameters and return response and error
type RetryFunc func(ctx context.Context) (interface{}, error)

func defaultCond(resp interface{}, err error) (bool, error) {
	return err == nil, nil
}

func retry(
	ctx context.Context,
	meta interface{},
	fn RetryFunc,
	params CustomRetryParameters,
	tClient token,
) (interface{}, error) {
	errorChannel := make(chan error)
	responseChannel := make(chan interface{})
	go func(ctx context.Context, meta interface{}, errorChannel chan error, responseChannel chan interface{}) {
		for i := 0; ; i++ {
			if i == params.retryCount {
				break
			}
			select {
			case <-ctx.Done():
				errorChannel <- fmt.Errorf("context timed out")
				responseChannel <- nil

				return
			case <-time.After(params.timeout):
				errorChannel <- fmt.Errorf("retry timed out")
				responseChannel <- nil

				return
			default:
				tClient.setScmClientToken(&ctx, meta)
				resp, err := fn(ctx)
				c, err := params.cond(resp, err)
				if err != nil {
					errorChannel <- err
					responseChannel <- nil

					return
				}
				if c {
					errorChannel <- nil
					responseChannel <- resp

					return
				}
				log.Printf("[WARN] on API call got error: %+v, response: %+v. Retrying", err, resp)
				time.Sleep(params.retryDelay)
			}
		}
	}(ctx, meta, errorChannel, responseChannel)
	err := <-errorChannel
	resp := <-responseChannel
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Retry with default count and timeout
func Retry(ctx context.Context, meta interface{}, fn RetryFunc) (interface{}, error) {
	params := CustomRetryParameters{
		retryCount: defaultRetryCount,
		retryDelay: defaultRetryTimeout,
		cond:       defaultCond,
		timeout:    defaultTimeout,
	}
	tStruct := tokenStruct{}

	return retry(ctx, meta, fn, params, &tStruct)
}

type CustomRetryParameters struct {
	retryCount int
	retryDelay time.Duration
	cond       CondFunc
	timeout    time.Duration
}

// CustomRetry allows developers to configure the timeout, retry count and delay
type CustomRetry struct {
	RetryCount int
	RetryDelay time.Duration
	Delay      time.Duration
	Cond       CondFunc
	Timeout    time.Duration
}

// Retry with custom count, timeout and delay
func (c *CustomRetry) Retry(
	ctx context.Context,
	meta interface{},
	fn RetryFunc,
) (interface{}, error) {
	if c.RetryCount == 0 {
		c.RetryCount = defaultRetryCount
	}
	if c.RetryDelay <= 0 {
		c.RetryDelay = defaultRetryTimeout
	}
	if c.Cond == nil {
		c.Cond = defaultCond
	}
	if c.Timeout != 0 {
		c.RetryCount = -1
	}
	time.Sleep(c.Delay)
	params := CustomRetryParameters{
		retryCount: c.RetryCount,
		retryDelay: c.RetryDelay,
		cond:       c.Cond,
		timeout:    c.Timeout,
	}
	tStruct := tokenStruct{}

	return retry(ctx, meta, fn, params, &tStruct)
}
