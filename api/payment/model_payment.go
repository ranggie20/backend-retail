package payment

import (
	"database/sql"
	"time"
)

type (
	// Model Payment represents the payment record in the database
	Payment struct {
		PaymentID       int32     `json:"payment_id"`
		UserID          int32     `json:"user_id"`
		CourseID        int32     `json:"course_id"`
		SubscriptionID  int32     `json:"subscription_id"`
		PaymentMethodID int32     `json:"payment_method_id"`
		PaymentStatusID int32     `json:"payment_status_id"`
		TotalAmount     int32     `json:"amount"`
		PaymentDate     time.Time `json:"payment_date"` // Use time.Time for timestamp fields
	}

	// Model PaymentRequest is used for creating or updating a payment record
	PaymentRequest struct {
		UserID          int32 `json:"user_id" validate:"required"`
		CourseID        int32 `json:"course_id" validate:"required"`
		SubscriptionID  int32 `json:"subscription_id" validate:"required"`
		PaymentMethodID int32 `json:"payment_method_id" validate:"required"`
		PaymentStatusID int32 `json:"payment_status_id" validate:"required"`
		TotalAmount     int32 `json:"amount" validate:"required"`
	}

	// Model GetPaymentRow represents the result of a query for payment details
	GetPaymentRow struct {
		CartID            sql.NullInt32  `json:"cart_id"`
		PaymentID         int32          `json:"payment_id"`
		CourseName        sql.NullString `json:"course_name"`
		Quantity          sql.NullInt32  `json:"quantity"`
		TotalAmount       sql.NullInt32  `json:"total_amount"`
		PaymentMethodName sql.NullString `json:"payment_method_name"`
		PaymentDate       sql.NullTime   `json:"payment_date"`
	}
	GetPaymentHistoryRow struct {
		PaymentID             sql.NullInt32  `json:"payment_id"`
		CourseID              sql.NullInt32  `json:"course_id"`
		CourseName            sql.NullString `json:"course_name"`
		TotalAmount           sql.NullInt32  `json:"total_amount"`
		PaymentStatusName     sql.NullString `json:"payment_status_name"`
		SubcriptionsStartDate sql.NullTime   `json:"subcriptions_start_date"`
	}
)
