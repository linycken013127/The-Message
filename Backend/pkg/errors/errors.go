package errors

import (
	"errors"
	"fmt"
)

// 錯誤類型定義
var (
	// ErrNotFound 找不到資源
	ErrNotFound = errors.New("resource not found")
	// ErrInvalidInput 無效輸入
	ErrInvalidInput = errors.New("invalid input")
	// ErrUnauthorized 未授權
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden 禁止存取
	ErrForbidden = errors.New("forbidden")
	// ErrInternalServer 內部伺服器錯誤
	ErrInternalServer = errors.New("internal server error")
)

// AppError 應用程式錯誤
type AppError struct {
	Code    string
	Message string
	Err     error
}

// Error 實作 error 介面
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 解包內部錯誤
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewNotFoundError 建立找不到資源錯誤
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: message,
		Err:     ErrNotFound,
	}
}

// NewValidationError 建立驗證錯誤
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Err:     ErrInvalidInput,
	}
}

// NewForbiddenError 建立禁止存取錯誤
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    "FORBIDDEN",
		Message: message,
		Err:     ErrForbidden,
	}
}

// NewInternalError 建立內部伺服器錯誤
func NewInternalError(message string) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: message,
		Err:     ErrInternalServer,
	}
}

// Wrap 包裝錯誤
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// IsNotFoundError 檢查是否為找不到資源錯誤
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsValidationError 檢查是否為驗證錯誤
func IsValidationError(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsForbiddenError 檢查是否為禁止存取錯誤
func IsForbiddenError(err error) bool {
	return errors.Is(err, ErrForbidden)
}
