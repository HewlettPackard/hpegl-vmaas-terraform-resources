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

// scmTokenInterface will wrap up setScmClientToken function to help on writing
// unit test
type scmTokenInterface interface {
	setScmClientToken(ctx *context.Context, meta interface{})
}

// tokenStruct implements scmTokenInterface
type tokenStruct struct{}

type continueStruct struct {
	resp    interface{}
	respErr error
}

// retryChan used as an arguments for groutine function retryRoutineFunc
type retryChan struct {
	errChan      chan error
	respChan     chan interface{}
	continueChan chan continueStruct
}

// setScmClientToken calls auth.SetScmClientToken
func (t *tokenStruct) setScmClientToken(ctx *context.Context, meta interface{}) {
	auth.SetScmClientToken(ctx, meta)
}

// CondFunc function accepts response and error of the RetryFunc. If any error returns
// retry will terminated and returns the error
type CondFunc func(response interface{}, ResponseErr error) (bool, error)

// RetryFunc accepts ctx as parameters and return response and error
type RetryFunc func(ctx context.Context) (interface{}, error)

// defaultCond default condition check for 'Retry' function
func defaultCond(resp interface{}, err error) (bool, error) {
	return err == nil, nil
}

// retry supports both retry with count and timeout as well and returns result as
// interface{}. This result can converted to proper struct/model afterwards
func retry(
	ctx context.Context,
	meta interface{},
	fn RetryFunc,
	cRetry CustomRetry,
	tClient scmTokenInterface,
) (interface{}, error) {
	rChan := retryChan{
		errChan:      make(chan error),
		respChan:     make(chan interface{}),
		continueChan: make(chan continueStruct),
	}

	timeoutTimer := time.NewTimer(cRetry.Timeout)
	// trigger retry initially and then wait for channels.
	go retryRoutineFunc(ctx, meta, rChan, tClient, cRetry, fn)
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context timed out")
		case <-timeoutTimer.C:
			return nil, fmt.Errorf("retry timed out")
		case continueChan := <-rChan.continueChan:
			// check exit condition before invoking next retry
			if i == cRetry.RetryCount-1 {
				return nil, fmt.Errorf(
					"maximum retry limit reached, with Error: %#v, Response: %#v",
					continueChan.respErr,
					continueChan.resp,
				)
			}
			go retryRoutineFunc(ctx, meta, rChan, tClient, cRetry, fn)
		case err := <-rChan.errChan:
			return nil, err
		case resp := <-rChan.respChan:
			return resp, nil
		}
		if i == cRetry.RetryCount-1 {
			return nil, fmt.Errorf("maximum retry limit reached")
		}
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

// Retry supports extra arguments. nitialDelay will put a delay before invoking the function.
// RetryCount supports customized retry count. If Timeout specified then RetryCount will be
// skipped. RetryDelay will put a delay in between each retrys. If any of these values are
// not specified then default value will be assigned.
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

// retryRoutineFunc implements logic for retry, this will check custom condition and call given function
func retryRoutineFunc(
	ctx context.Context,
	meta interface{},
	sChan retryChan,
	tClient scmTokenInterface,
	cRetry CustomRetry,
	fn RetryFunc,
) {
	tClient.setScmClientToken(&ctx, meta)
	resp, respErr := fn(ctx)
	c, err := cRetry.Cond(resp, respErr)
	if err != nil {
		sChan.errChan <- err

		return
	}
	if c {
		sChan.respChan <- resp

		return
	}
	log.Printf("[WARN] on API execution. Error:%#v, Response: %#v", err, resp)
	time.Sleep(cRetry.RetryDelay)
	// continue retry
	sChan.continueChan <- continueStruct{
		resp:    resp,
		respErr: respErr,
	}
}
