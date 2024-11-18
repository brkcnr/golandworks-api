package apierror_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/brkcnr/golandworks-api/internal/apierror"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		code        int
		message     string
		wantCode    int
		wantMessage string
	}{
		{
			name:        "create new error",
			code:        http.StatusBadRequest,
			message:     "test error",
			wantCode:    http.StatusBadRequest,
			wantMessage: "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apierror.New(tt.code, tt.message)
			if err.Code != tt.wantCode {
				t.Errorf("New() code = %v, want %v", err.Code, tt.wantCode)
			}
			if err.Message != tt.wantMessage {
				t.Errorf("New() message = %v, want %v", err.Message, tt.wantMessage)
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *apierror.APIError
		want    string
		inner   error
		message string
	}{
		{
			name:    "error without inner error",
			err:     apierror.New(http.StatusBadRequest, "test error"),
			want:    "test error",
			inner:   nil,
			message: "test error",
		},
		{
			name:    "error with inner error",
			err:     apierror.Wrap(errors.New("inner error"), http.StatusBadRequest, "test error"),
			want:    "test error: inner error",
			inner:   errors.New("inner error"),
			message: "test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	tests := []struct {
		name    string
		err     *apierror.APIError
		wantErr error
	}{
		{
			name:    "unwrap nil inner error",
			err:     apierror.New(http.StatusBadRequest, "test error"),
			wantErr: nil,
		},
		{
			name:    "unwrap inner error",
			err:     apierror.Wrap(innerErr, http.StatusBadRequest, "test error"),
			wantErr: innerErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.err.Unwrap(); err != tt.wantErr {
				t.Errorf("APIError.Unwrap() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	innerErr := errors.New("inner error")
	tests := []struct {
		name        string
		err         error
		code        int
		message     string
		wantCode    int
		wantMessage string
		wantInner   error
	}{
		{
			name:        "wrap error",
			err:         innerErr,
			code:        http.StatusBadRequest,
			message:     "test error",
			wantCode:    http.StatusBadRequest,
			wantMessage: "test error",
			wantInner:   innerErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apierror.Wrap(tt.err, tt.code, tt.message)
			if got.Code != tt.wantCode {
				t.Errorf("Wrap() code = %v, want %v", got.Code, tt.wantCode)
			}
			if got.Message != tt.wantMessage {
				t.Errorf("Wrap() message = %v, want %v", got.Message, tt.wantMessage)
			}
			if got.Inner != tt.wantInner {
				t.Errorf("Wrap() inner = %v, want %v", got.Inner, tt.wantInner)
			}
		})
	}
} 