package model

import (
	"errors"
	"fmt"
)

// 共通エラー定義
var (
	// ErrNotFound はリソースが見つからない場合のエラー
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput は入力が不正な場合のエラー
	ErrInvalidInput = errors.New("invalid input")

	// ErrConflict はリソースの競合が発生した場合のエラー
	ErrConflict = errors.New("resource conflict")

	// ErrInternal は内部エラー
	ErrInternal = errors.New("internal server error")
)

// AppError はアプリケーション固有のエラー情報を保持する
type AppError struct {
	Err     error
	Message string
	Code    string
}

// Error はerrorインターフェースの実装
func (e *AppError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Err.Error()
}

// Unwrap はエラーのアンラップをサポート
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError は新しいAppErrorを作成する
func NewAppError(err error, message string, code string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}
