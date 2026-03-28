package response

import (
	"encoding/json"
	"net/http"

	"intelligent_guidance_system_v2/common/pkg/errors"
)

const (
	CodeSuccess = 0
	CodeError   = 1
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int32       `json:"page"`
	PageSize int32       `json:"page_size"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	}
}

func SuccessWithMessage(message string, data interface{}) *Response {
	return &Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	}
}

func Error(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}

func ErrorWithData(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func FromAppError(err *errors.AppError) *Response {
	return &Response{
		Code:    err.Code,
		Message: err.Message,
		Data:    err.Metadata,
	}
}

func FromError(err error) *Response {
	appErr := errors.GetAppError(err)
	return FromAppError(appErr)
}

func PagedList(list interface{}, total int64, page, pageSize int32) *Response {
	return Success(&PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func JSON(w http.ResponseWriter, statusCode int, resp *Response) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(resp)
}

func JSONSuccess(w http.ResponseWriter, data interface{}) error {
	return JSON(w, http.StatusOK, Success(data))
}

func JSONError(w http.ResponseWriter, err error) error {
	appErr := errors.GetAppError(err)
	return JSON(w, appErr.HTTPStatus, FromAppError(appErr))
}

func JSONBadRequest(w http.ResponseWriter, message string) error {
	return JSON(w, http.StatusBadRequest, Error(errors.CodeInvalidArgument, message))
}

func JSONUnauthorized(w http.ResponseWriter, message string) error {
	return JSON(w, http.StatusUnauthorized, Error(errors.CodeUnauthenticated, message))
}

func JSONForbidden(w http.ResponseWriter, message string) error {
	return JSON(w, http.StatusForbidden, Error(errors.CodePermissionDenied, message))
}

func JSONNotFound(w http.ResponseWriter, resource string) error {
	return JSON(w, http.StatusNotFound, Error(errors.CodeNotFound, resource+" not found"))
}

func JSONInternalError(w http.ResponseWriter, message string) error {
	return JSON(w, http.StatusInternalServerError, Error(errors.CodeInternal, message))
}