package subscriptions

type (
	Subscription struct {
		SubscriptionID int32  `json:"subscription_id"` // Primary key
		UserID         int32  `json:"user_id"`         // Foreign key to the user
		CourseID       int32  `json:"course_id"`       // Foreign key to the course
		CartID         int32  `json:"cart_id"`         // Foreign key to the cart
		IsCorrect      string `json:"is_correct"`      // Whether the subscription is correct
		CreatedAt      string `json:"created_at"`      // Timestamp of creation
		UpdatedAt      string `json:"updated_at"`      // Timestamp of the last update
		// Timestamp of deletion (optional, nullable)
	}

	SubscriptionRequest struct {
		UserID    int32  `json:"user_id" validate:"required"`    // Required user ID
		CourseID  int32  `json:"course_id" validate:"required"`  // Required course ID
		CartID    int32  `json:"cart_id" validate:"required"`    // Required cart ID
		IsCorrect string `json:"is_correct" validate:"required"` // Required correctness status
	}
)
