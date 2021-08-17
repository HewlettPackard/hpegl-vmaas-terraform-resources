// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

//go:generate go run github.com/golang/mock/mockgen -source ./retry.go -package utils -destination ./retry_mock.go

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
	cRetry CustomRetry,
	tClient token,
) (interface{}, error) {
	errorChannel := make(chan error)
	responseChannel := make(chan interface{})
	go func(ctx context.Context, meta interface{}, errorChannel chan error, responseChannel chan interface{}) {
		for i := 0; ; i++ {
			if i == cRetry.RetryCount {
				errorChannel <- fmt.Errorf("maximum retry limit reached")

				return
			}
			select {
			case <-ctx.Done():
				errorChannel <- fmt.Errorf("context timed out")

				return
			case <-time.After(cRetry.Timeout):
				errorChannel <- fmt.Errorf("retry timed out")

				return
			default:
				tClient.setScmClientToken(&ctx, meta)
				resp, respErr := fn(ctx)
				c, err := cRetry.Cond(resp, respErr)
				if err != nil {
					errorChannel <- err

					return
				}
				if c {
					responseChannel <- resp

					return
				}
				log.Printf("[WARN] on API call got error: %#v, response: %#v. Retrying", err, resp)
				time.Sleep(cRetry.RetryDelay)
			}
		}
	}(ctx, meta, errorChannel, responseChannel)

	select {
	case err := <-errorChannel:
		return nil, err
	case resp := <-responseChannel:
		return resp, nil
	}
}

// Retry with default count and timeout
func Retry(ctx context.Context, meta interface{}, fn RetryFunc) (interface{}, error) {
	c := CustomRetry{
		RetryCount: defaultRetryCount,
		RetryDelay: defaultRetryDelay,
		Cond:       defaultCond,
		Timeout:    defaultTimeout,
	}

	return retry(ctx, meta, fn, c, &tokenStruct{})
}

// CustomRetry allows developers to configure the timeout, retry count and delay
type CustomRetry struct {
	RetryCount   int
	RetryDelay   time.Duration
	InitialDelay time.Duration
	Cond         CondFunc
	Timeout      time.Duration
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
		c.RetryDelay = defaultRetryDelay
	}
	if c.Cond == nil {
		c.Cond = defaultCond
	}
	if c.Timeout != 0 {
		c.RetryCount = noRetryCount
	}
	time.Sleep(c.InitialDelay)

	return retry(ctx, meta, fn, *c, &tokenStruct{})
}
