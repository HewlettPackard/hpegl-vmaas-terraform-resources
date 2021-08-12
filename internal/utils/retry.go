// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/auth"
)

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
	count int,
	retryDelay time.Duration,
	fn RetryFunc,
	cond CondFunc,
	timeout time.Duration,
) (interface{}, error) {
	var err error
	var resp interface{}
	errorChannel := make(chan error)
	responseChannel := make(chan interface{})
	go func(ctx context.Context, meta interface{}, errorChannel chan error, responseChannel chan interface{}) {
		for i := 0; ; i++ {
			if i == count {
				break
			}
			select {
			case <-ctx.Done():
				errorChannel <- fmt.Errorf("context timed out")
				responseChannel <- nil
				return
			case <-time.After(timeout):
				errorChannel <- fmt.Errorf("retry timed out")
				responseChannel <- nil
				return
			default:
				auth.SetScmClientToken(&ctx, meta)
				resp, err = fn(ctx)
				c, err := cond(resp, err)
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
				time.Sleep(retryDelay)
			}
		}
	}(ctx, meta, errorChannel, responseChannel)
	err = <-errorChannel
	resp = <-responseChannel
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// Retry with default count and timeout
func Retry(ctx context.Context, meta interface{}, fn RetryFunc) (interface{}, error) {
	return retry(ctx, meta, defaultRetryCount, defaultRetryTimeout, fn, defaultCond, defaultTimeout)
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
	if c.RetryDelay == 0 {
		c.RetryDelay = defaultRetryTimeout
	}
	if c.Cond == nil {
		c.Cond = defaultCond
	}
	if c.Timeout != 0 {
		c.RetryCount = -1
	}
	time.Sleep(c.Delay)

	return retry(ctx, meta, c.RetryCount, c.RetryDelay, fn, c.Cond, c.Timeout)
}
