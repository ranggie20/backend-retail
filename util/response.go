package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/online-bnsp/backend/constant"

	"github.com/google/uuid"
)

type (
	Response struct {
		Code      int         `json:"code"`
		RequestID uuid.UUID   `json:"request_id"`
		Username  string      `json:"username,omitempty"`
		Status    int         `json:"status"`
		Message   string      `json:"message"`
		Data      interface{} `json:"data,omitempty"`
		Meta      Meta        `json:"meta"`
	}

	Meta struct {
		Latency    string `json:"latency"`
		NextCursor int    `json:"next_cursor,omitempty"`
	}
)

func NewResponse(code, status int, msg string, data interface{}, username ...string) *Response {
	resp := &Response{
		Code:    code,
		Status:  status,
		Message: msg,
		Data:    data,
		Meta:    Meta{},
	}

	if len(username) > 0 {
		resp.Username = username[0]
	}

	return resp
}

// WriteResponse - write response to the client
func (resp *Response) WriteResponse(w http.ResponseWriter, r *http.Request) {
	birthTime := r.Context().Value(constant.ContextBirthTime).(time.Time)
	latency := time.Since(birthTime).Seconds() * 1000
	resp.Meta.Latency = fmt.Sprintf("%.2f ms", latency)

	requestID := r.Context().Value(middleware.RequestIDKey).(uuid.UUID)
	resp.RequestID = requestID

	// return ctx.JSON(r.Status, r)
	responseJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)
	w.Write(responseJSON)
}
