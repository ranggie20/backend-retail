package transactionhistory

type (
	TransactionHistory struct {
		TransactionHistoryID   int32   `json:"transaction_history_id"` // Assuming there's a primary key
		SubscriptionsID        int32   `json:"subscriptions_id"`
		Quantity               int32   `json:"quantity"`
		TotalAmount            float64 `json:"total_amount"`
		IsPaid                 string  `json:"is_paid"`
		SubscriptionsStartDate string  `json:"subscriptions_start_date"`
		Proof                  string  `json:"proof"`
		CreatedAt              string  `json:"created_at"`
		UpdatedAt              string  `json:"updated_at"`
		DeletedAt              string  `json:"deleted_at"` // Nullable field
	}

	TransactionHistoryRequest struct {
		SubscriptionsID        int32   `json:"subscriptions_id" validate:"required"`
		Quantity               int32   `json:"quantity" validate:"required"`
		TotalAmount            float64 `json:"total_amount" validate:"required"`
		IsPaid                 string  `json:"is_paid" validate:"required"`
		SubscriptionsStartDate string  `json:"subscriptions_start_date" validate:"required"`
		Proof                  string  `json:"proof"`
		CreatedAt              string  `json:"created_at"`
	}
)
