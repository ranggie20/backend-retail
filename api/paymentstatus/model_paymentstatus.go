package paymentstatus

import (
	"time"
)

type PaymentStatus struct {
	PaymentStatusID   int32     `json:"payment_status_id"`
	PaymentStatusName string    `json:"payment_status_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type PaymentStatusRequest struct {
	PaymentStatusName string `json:"payment_status_name" validate:"required"`
	CreatedAt         string `json:"created_at"`
}
