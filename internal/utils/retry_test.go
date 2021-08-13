// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/auth"
)

func TestCustomRetry_Retry(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	meta := 1

	tests := []struct {
		name       string
		fn         RetryFunc
		RetryCount int
		RetryDelay time.Duration
		Delay      time.Duration
		Cond       CondFunc
		Timeout    time.Duration
		given      func(m *auth.MockToken)
		want       interface{}
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "Normal Test case 1: Retry with count",
			RetryCount: 5,
			RetryDelay: time.Second * 5,
			Delay:      time.Second * 5,
			given: func(m *auth.MockToken) {
				m.EXPECT().SetScmClientToken(&ctx, meta)
			},
			want: 
		},
		{
			name:       "Normal Test case 2: Retry with timeout",
			RetryDelay: time.Second * 5,
			Delay:      time.Second * 5,
			Timeout:    time.Second * 20,
			given: func(m *auth.MockToken) {
				m.EXPECT().SetScmClientToken(&ctx, meta)
				m.EXPECT().GetToken(ctx, meta)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockToken := auth.NewMockToken(ctrl)
			a := CustomRetry{
				RetryCount: tt.RetryCount,
				RetryDelay: tt.RetryDelay,
				Delay:      tt.Delay,
				Cond:       tt.Cond,
				Timeout:    tt.Timeout,
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
