// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

//go:generate go run github.com/golang/mock/mockgen -source ./retry.go -package utils -destination ./retry_mock.go

package utils

import (
	"context"
	"reflect"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/retrieve"
)

type RetrySampleReturn struct {
	Name string
}

var metaFunc retrieve.TokenRetrieveFuncCtx = func(ctx context.Context) (string, error) {
	return "", nil
}

func TestCustomRetry_Retry(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	meta := map[string]interface{}{
		"tokenRetrieveFunc": metaFunc,
	}
	tests := []struct {
		name    string
		fn      RetryFunc
		given   func(m *Mocktoken)
		delay   time.Duration
		param   CustomRetryParameters
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Normal Test case 1: Retry with count",
			fn: func(ctx context.Context) (interface{}, error) {
				return RetrySampleReturn{Name: "template"}, nil
			},
			given: func(m *Mocktoken) {
				//m.EXPECT().setScmClientToken(gomock.Any(), meta)
			},
			delay: time.Second * 5,
			param: CustomRetryParameters{
				retryCount: 5,
				retryDelay: time.Second * 5,
				cond:       defaultCond,
				timeout:    0,
			},
			want: RetrySampleReturn{
				Name: "template",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockToken := NewMocktoken(ctrl)
			a := CustomRetry{
				RetryCount: tt.param.retryCount,
				RetryDelay: tt.param.retryDelay,
				Cond:       tt.param.cond,
				Timeout:    tt.param.timeout,
				Delay:      tt.delay,
			}
			tt.given(mockToken)
			got, err := a.Retry(ctx, meta, tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomRetry.Retry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CustomRetry.Retry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_retry(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	meta := map[string]interface{}{
		"tokenRetrieveFunc": metaFunc,
	}

	tests := []struct {
		name    string
		fn      RetryFunc
		given   func(m *Mocktoken)
		delay   time.Duration
		params  CustomRetryParameters
		tClient token
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Normal Test case 1: Retry with count",
			fn: func(ctx context.Context) (interface{}, error) {
				return RetrySampleReturn{Name: "template"}, nil
			},
			given: func(m *Mocktoken) {
				m.EXPECT().setScmClientToken(gomock.Any(), meta)
			},
			delay: time.Second * 5,
			params: CustomRetryParameters{
				retryCount: 5,
				retryDelay: time.Second * 5,
				cond:       defaultCond,
				timeout:    0,
			},
			tClient: &tokenStruct{},
			want: RetrySampleReturn{
				Name: "template",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockToken := NewMocktoken(ctrl)
			tt.given(mockToken)
			got, err := retry(ctx, meta, tt.fn, tt.params, tt.tClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("retry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("retry() = %v, want %v", got, tt.want)
			}
		})
	}
}
