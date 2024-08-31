package wishlist

type (
	Wishlist struct {
		WishlistID int32  `json:"wishlist_id"`
		UserID     int32  `json:"user_id"`
		CourseID   int32  `json:"course_id"`
		CreatedAt  string `json:"created_at"`
		DeletedAt  string `json:"deleted_at"`
		UpdatedAt  string `json:"updated_at"`
	}

	WishlistRequest struct {
		UserID   int32 `json:"user_id" validate:"required"`
		CourseID int32 `json:"course_id" validate:"required"`
	}
)
