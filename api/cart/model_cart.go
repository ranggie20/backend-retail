package cart

import "database/sql"

type (
	// Model Cart yang sesuai dengan tabel cart
	Cart struct {
		CartID      int32  `json:"cart_id"`      // ID cart
		UserID      int32  `json:"user_id"`      // ID pengguna yang memiliki cart
		CourseID    int32  `json:"course_id"`    // ID kursus yang ada di cart
		CourseName  string `json:"course_name"`  // Nama kursus
		Price       int32  `json:"price"`        // Harga kursus
		Quantity    int32  `json:"quantity"`     // Jumlah kursus yang dibeli
		TotalAmount int32  `json:"total_amount"` // Total biaya
	}

	// Model CartRequest untuk request input
	CartRequest struct {
		UserID      int32  `json:"user_id" validate:"required"`      // ID pengguna (wajib diisi)
		CourseID    int32  `json:"course_id" validate:"required"`    // ID kursus (wajib diisi)
		CourseName  string `json:"course_name" validate:"required"`  // Nama kursus (wajib diisi)
		Price       int32  `json:"price" validate:"required"`        // Harga kursus (wajib diisi)
		Quantity    int32  `json:"quantity" validate:"required"`     // Jumlah kursus (wajib diisi)
		TotalAmount int32  `json:"total_amount" validate:"required"` // Total biaya (wajib diisi)
	}
	GetCartRow struct {
		CourseID    sql.NullInt32  `json:"course_id"`
		Thumbnail   sql.NullString `json:"thumbnail"`
		CourseName  sql.NullString `json:"course_name"`
		Price       sql.NullInt32  `json:"price"`
		Quantity    sql.NullInt32  `json:"quantity"`
		TotalAmount sql.NullInt32  `json:"total_amount"`
	}
)
