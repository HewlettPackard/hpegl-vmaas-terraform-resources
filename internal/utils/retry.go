// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"context"
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
	// If timeout is set then this code will be skipped
	for i := 0; i < count; i++ {
		auth.SetScmClientToken(&ctx, meta)
		resp, err = fn(ctx)
		c, err := cond(resp, err)
		if err != nil {
			return nil, err
		}
		if c {
			break
		}

		log.Printf("[WARN] on API call got error: %+v, response: %+v. Retrying", err, resp)
		time.Sleep(retryDelay)
	}
	// This sleeps for the timeout period and then checks the response and returns
	if timeout != 0 {
		time.Sleep(timeout)
		auth.SetScmClientToken(&ctx, meta)
		resp, err = fn(ctx)
		_, err := cond(resp, err)
		if err != nil {
			return nil, err
		}
	}

	return resp, err
}

// Retry with default count and timeout
func Retry(ctx context.Context, meta interface{}, fn RetryFunc) (interface{}, error) {
	return retry(ctx, meta, defaultRetryCount, defaultTimeout, fn, defaultCond)
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
		c.RetryDelay = defaultTimeout
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
