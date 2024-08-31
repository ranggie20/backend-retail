package paymentmethod

type PaymentMethod struct {
	PaymentMethodID   int32  `json:"payment_method_id"`
	PaymentMethodName string `json:"payment_method_name"`
}

type PaymentMethodRequest struct {
	PaymentMethodName string `json:"payment_method_name" validate:"required"`
	CreatedAt         string `json:"created_at"`
}
