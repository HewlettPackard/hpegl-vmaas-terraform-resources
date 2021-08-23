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

const testRetrySuccess = "success"

func Test_retry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	count := 0

	type args struct {
		meta    interface{}
		fn      RetryFunc
		cRetry  CustomRetry
		tClient scmTokenInterface
	}
	tests := []struct {
		name    string
		args    args
		given   func(m *MockscmTokenInterface)
		want    interface{}
		wantErr bool
	}{
		{
			name: "Normal test case 1 - no retry",
			args: args{
				meta: "mock meta",
				fn: func(ctx context.Context) (interface{}, error) {
					return testRetrySuccess, nil
				},
				cRetry: CustomRetry{
					RetryCount: 1,
					Timeout:    defaultTimeout,
					Cond:       defaultCond,
				},
			},
			given: func(m *MockscmTokenInterface) {
				m.EXPECT().setScmClientToken(gomock.Any(), "mock meta")
			},
			want: testRetrySuccess,
		},
		{
			name: "Normal test case 2 - 2 retry",
			args: args{
				meta: "mock meta",
				fn: func(ctx context.Context) (interface{}, error) {
					if count == 2 {
						return testRetrySuccess, nil
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
			given: func(m *MockscmTokenInterface) {
				m.EXPECT().setScmClientToken(gomock.Any(), "mock meta").Times(3)
			},
			want: testRetrySuccess,
		},
		{
			name: "Normal test case 3 - with timeout",
			args: args{
				meta: "mock meta",
				fn: func(ctx context.Context) (interface{}, error) {
					return testRetrySuccess, nil
				},
				cRetry: CustomRetry{
					RetryCount: noRetryCount,
					Timeout:    time.Millisecond * 10,
					Cond:       defaultCond,
				},
			},
			given: func(m *MockscmTokenInterface) {
				m.EXPECT().setScmClientToken(gomock.Any(), "mock meta")
			},
			want: testRetrySuccess,
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
			given: func(m *MockscmTokenInterface) {
				m.EXPECT().setScmClientToken(gomock.Any(), "mock meta").Times(3)
			},
			wantErr: true,
		},
		{
			name: "Failed test case 2 - retry count exceeds",
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
			given: func(m *MockscmTokenInterface) {
				m.EXPECT().setScmClientToken(gomock.Any(), "mock meta").AnyTimes()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMockscmTokenInterface(ctrl)
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

func TestRetryRoutineStruct_WaitForRetryRoutine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	meta := "mock meta"
	count := 0

	type args struct {
		fn RetryFunc
	}

	tests := []struct {
		name    string
		args    args
		given   func(m *MockscmTokenInterface)
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Normal test case 1 - no retry",
			args: args{
				fn: func(ctx context.Context) (interface{}, error) {
					return testRetrySuccess, nil
				},
			},
			given: func(m *MockscmTokenInterface) {
				m.EXPECT().setScmClientToken(gomock.Any(), meta)
			},
			want: testRetrySuccess,
		},
		{
			name: "Normal test case 2 - Default retry values",
			args: args{
				fn: func(ctx context.Context) (interface{}, error) {
					if count == defaultRetryCount-1 {
						return testRetrySuccess, nil
					}
					count++

					return "", errors.New("error")
				},
			},
			given: func(m *MockscmTokenInterface) {
				m.EXPECT().setScmClientToken(gomock.Any(), meta).Times(defaultRetryCount)
			},
			want: testRetrySuccess,
		},
		{
			name: "Failed test case 1 - Default retry count exceeds",
			args: args{
				fn: func(ctx context.Context) (interface{}, error) {
					return nil, errors.New("error")
				},
			},
			given: func(m *MockscmTokenInterface) {
				m.EXPECT().setScmClientToken(gomock.Any(), meta).Times(defaultRetryCount)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMockscmTokenInterface(ctrl)
			a := RetryRoutineStruct{
				retryTokenStruct: m,
			}

			tt.given(m)
			count = 0
			a.RetryRoutine(ctx, meta, tt.args.fn)
			got, err := a.WaitForRetryRoutine()
			if (err != nil) != tt.wantErr {
				t.Errorf("RetryRoutineStruct.WaitForRetryRoutine() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetryRoutineStruct.WaitForRetryRoutine() = %v, want %v", got, tt.want)
			}
		})
	}
}
