// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"time"

	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
)

type CondFunc func(interface{}, error) bool

func defaultCond(resp interface{}, err error) bool {
	return err == nil
}

func retry(count int, timeout time.Duration, fn func() (interface{}, error), cond CondFunc) (interface{}, error) {
	var err error
	var resp interface{}
	for i := 0; i < count; i++ {
		resp, err = fn()
		if cond(resp, err) {
			break
		}
		logger.Error("warning: ", err, ". Response: ", resp, ". retrying")
		time.Sleep(timeout)
	}

	return resp, err
}

// Retry with default count and timeout
func Retry(fn func() (interface{}, error)) (interface{}, error) {
	return retry(defaultRetryCount, defaultTimeout, fn, defaultCond)
}

// CustomRetry allows developers to configure the timeout, retry count and delay
type CustomRetry struct {
	RetryCount   int
	RetryTimeout time.Duration
	Delay        time.Duration
	Cond         CondFunc
}

// Retry with custom count, timeout and delay
func (c *CustomRetry) Retry(fn func() (interface{}, error)) (interface{}, error) {
	if c.RetryCount == 0 {
		c.RetryCount = defaultRetryCount
	}
	if c.RetryTimeout == 0 {
		c.RetryTimeout = defaultTimeout
	}
	if c.Cond == nil {
		c.Cond = defaultCond
	}
	time.Sleep(c.Delay)

	return retry(c.RetryCount, c.RetryTimeout, fn, c.Cond)
}
