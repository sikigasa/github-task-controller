package model

import "errors"

// ErrNotFound はリソースが見つからない場合のエラー
var ErrNotFound = errors.New("resource not found")

// ErrUnauthorized は認証されていない場合のエラー
var ErrUnauthorized = errors.New("unauthorized")

// ErrForbidden は権限がない場合のエラー
var ErrForbidden = errors.New("forbidden")

// ErrInvalidInput は入力が不正な場合のエラー
var ErrInvalidInput = errors.New("invalid input")

// ErrConflict はリソースが競合している場合のエラー
var ErrConflict = errors.New("resource conflict")

// ErrInternalServer は内部サーバーエラー
var ErrInternalServer = errors.New("internal server error")
