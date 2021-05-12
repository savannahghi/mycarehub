package mock

import (
	"context"
	"net/http"
)

// FakeServiceLogin ..
type FakeServiceLogin struct {
	GetLoginFuncFn       func(ctx context.Context) http.HandlerFunc
	GetLogoutFuncFn      func(ctx context.Context) http.HandlerFunc
	GetRefreshFuncFn     func() http.HandlerFunc
	GetVerifyTokenFuncFn func(ctx context.Context) http.HandlerFunc
}

// GetLoginFunc ...
func (l *FakeServiceLogin) GetLoginFunc(ctx context.Context) http.HandlerFunc {
	return l.GetLoginFuncFn(ctx)
}

// GetLogoutFunc ...
func (l *FakeServiceLogin) GetLogoutFunc(ctx context.Context) http.HandlerFunc {
	return l.GetLogoutFuncFn(ctx)
}

// GetRefreshFunc ...
func (l *FakeServiceLogin) GetRefreshFunc() http.HandlerFunc {
	return l.GetRefreshFuncFn()
}

// GetVerifyTokenFunc ...
func (l *FakeServiceLogin) GetVerifyTokenFunc(ctx context.Context) http.HandlerFunc {
	return l.GetVerifyTokenFuncFn(ctx)
}
