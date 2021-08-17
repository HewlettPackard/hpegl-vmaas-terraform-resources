// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

//go:generate go run github.com/golang/mock/mockgen -source ./retry.go -package utils -destination ./retry_mock.go

package utils

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func Test_retry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	count := 0

	type args struct {
		meta    interface{}
		fn      RetryFunc
		cRetry  CustomRetry
		tClient token
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *Mocktoken)
		want    interface{}
		wantErr bool
	}{
		{
			name: "Normal test case 1 - no retry",
			args: args{
				meta: "mock meta",
				fn: func(ctx context.Context) (interface{}, error) {
					return "success", nil
				},
				cRetry: CustomRetry{
					RetryCount: 1,
					Timeout:    defaultTimeout,
					Cond:       defaultCond,
				},
			},
			given: func(m *Mocktoken) {
				m.EXPECT().setScmClientToken(gomock.Any(), "mock meta")
			},
			want: "success",
		},
		{
			name: "Normal test case 2 - 2 retry",
			args: args{
				meta: "mock meta",
				fn: func(ctx context.Context) (interface{}, error) {
					if count == 2 {
						return "success", nil
					}
					count++

					return "", errors.New("error")
				},
				cRetry: CustomRetry{
					RetryCount: 3,
					Timeout:    defaultTimeout,
					Cond:       defaultCond,
				},
			},
			given: func(m *Mocktoken) {
				m.EXPECT().setScmClientToken(gomock.Any(), "mock meta").Times(3)
			},
			want: "success",
		},
		{
			name: "Normal test case 3 - with timeout",
			args: args{
				meta: "mock meta",
				fn: func(ctx context.Context) (interface{}, error) {
					return "success", nil
				},
				cRetry: CustomRetry{
					RetryCount: noRetryCount,
					Timeout:    time.Millisecond * 10,
					Cond:       defaultCond,
				},
			},
			given: func(m *Mocktoken) {
				m.EXPECT().setScmClientToken(gomock.Any(), "mock meta")
			},
			want: "success",
		},
		{
			name: "Failed test case 1 - retry count exceeded",
			args: args{
				meta: "mock meta",
				fn: func(ctx context.Context) (interface{}, error) {
					return nil, errors.New("error")
				},
				cRetry: CustomRetry{
					RetryCount: 3,
					Timeout:    defaultTimeout,
					Cond:       defaultCond,
				},
			},
			given: func(m *Mocktoken) {
				m.EXPECT().setScmClientToken(gomock.Any(), "mock meta").Times(3)
			},
			wantErr: true,
		},
		{
			name: "Failed test case 2 - retry timeout exceeds",
			args: args{
				meta: "mock meta",
				fn: func(ctx context.Context) (interface{}, error) {
					return nil, errors.New("error")
				},
				cRetry: CustomRetry{
					RetryCount: 3,
					Timeout:    time.Millisecond * 5,
					Cond:       defaultCond,
					RetryDelay: time.Second,
				},
			},
			given: func(m *Mocktoken) {
				m.EXPECT().setScmClientToken(gomock.Any(), "mock meta").Times(3)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMocktoken(ctrl)
			tt.args.tClient = m
			count = 0

			tt.given(m)
			got, err := retry(ctx, tt.args.meta, tt.args.fn, tt.args.cRetry, tt.args.tClient)
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
