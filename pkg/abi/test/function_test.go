package test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/rioam2/witigo/pkg/abi"
)

func TestCall(t *testing.T) {
	tests := []struct {
		name           string
		opts           abi.AbiOptions
		fnName         string
		params         []uint64
		expectedRet    uint32
		expectedErr    error
		postReturnErr  error
		callResults    []uint64
		callError      error
		postCallError  error
		postCallCalled bool
	}{
		{
			name: "successful call",
			opts: abi.AbiOptions{
				Context: context.Background(),
				Call: func(ctx context.Context, name string, params ...uint64) ([]uint64, error) {
					if name == "test_function" {
						return []uint64{42}, nil
					}
					return nil, nil
				},
			},
			fnName:         "test_function",
			params:         []uint64{1, 2, 3},
			expectedRet:    42,
			expectedErr:    nil,
			postCallCalled: false,
		},
		{
			name: "call function not defined",
			opts: abi.AbiOptions{
				Context: context.Background(),
			},
			fnName:         "test_function",
			params:         []uint64{1, 2, 3},
			expectedRet:    0,
			expectedErr:    fmt.Errorf("call function is not defined in AbiOptions"),
			postCallCalled: false,
		},
		{
			name: "call function returns error",
			opts: abi.AbiOptions{
				Context: context.Background(),
				Call: func(ctx context.Context, name string, params ...uint64) ([]uint64, error) {
					return nil, errors.New("call error")
				},
			},
			fnName:         "test_function",
			params:         []uint64{1, 2, 3},
			expectedRet:    0,
			expectedErr:    fmt.Errorf("function call test_function failed: call error"),
			postCallCalled: false,
		},
		{
			name: "empty result",
			opts: abi.AbiOptions{
				Context: context.Background(),
				Call: func(ctx context.Context, name string, params ...uint64) ([]uint64, error) {
					return []uint64{}, nil
				},
			},
			fnName:         "test_function",
			params:         []uint64{1, 2, 3},
			expectedRet:    0,
			expectedErr:    nil,
			postCallCalled: false,
		},
		{
			name: "post call success",
			opts: abi.AbiOptions{
				Context: context.Background(),
				Call: func(ctx context.Context, name string, params ...uint64) ([]uint64, error) {
					if name == "test_function" {
						return []uint64{42}, nil
					} else if name == "cabi_post_test_function" {
						return []uint64{}, nil
					}
					return nil, errors.New("unexpected function")
				},
			},
			fnName:         "test_function",
			params:         []uint64{1, 2, 3},
			expectedRet:    42,
			expectedErr:    nil,
			postReturnErr:  nil,
			postCallCalled: true,
		},
		{
			name: "post call error",
			opts: abi.AbiOptions{
				Context: context.Background(),
				Call: func(ctx context.Context, name string, params ...uint64) ([]uint64, error) {
					if name == "test_function" {
						return []uint64{42}, nil
					} else if name == "cabi_post_test_function" {
						return nil, errors.New("post call error")
					}
					return nil, errors.New("unexpected function")
				},
			},
			fnName:         "test_function",
			params:         []uint64{1, 2, 3},
			expectedRet:    42,
			expectedErr:    nil,
			postReturnErr:  errors.New("post call error"),
			postCallCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret, postReturn, err := abi.Call(tt.opts, tt.fnName, tt.params...)

			if tt.expectedErr == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedErr)
				} else if err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			}

			if ret != tt.expectedRet {
				t.Errorf("expected return %d, got %d", tt.expectedRet, ret)
			}

			if tt.postCallCalled {
				postErr := postReturn()
				if (tt.postReturnErr == nil && postErr != nil) ||
					(tt.postReturnErr != nil && postErr == nil) ||
					(tt.postReturnErr != nil && postErr != nil && tt.postReturnErr.Error() != postErr.Error()) {
					t.Errorf("expected post return error %v, got %v", tt.postReturnErr, postErr)
				}
			}
		})
	}
}
