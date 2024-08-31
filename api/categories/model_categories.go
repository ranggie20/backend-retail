package categories

import (
	"time"
)

type (
	// Model Category yang sesuai dengan tabel categories
	Category struct {
		CategoryID   int32     `json:"category_id"`   // Menggunakan int32 untuk mencocokkan tipe SERIAL
		CategoryName string    `json:"category_name"` // Nama kategori
		Icon         string    `json:"icon"`          // Ikon kategori
		CreatedAt    time.Time `json:"created_at"`    // Waktu pembuatan kategori
		UpdatedAt    time.Time `json:"updated_at"`    // Waktu update terakhir kategori
	}

	// Model CategoryRequest untuk request input
	CategoryRequest struct {
		CategoryName string    `json:"category_name" validate:"required"` // Nama kategori (wajib diisi)
		Icon         string    `json:"icon" validate:"required"`
		UpdatedAt    time.Time `json: "updated_at" validate:"required"` // Ikon kategori (wajib diisi)
	}
	GetCategoryRow struct {
		CourseID     int32  `json:"course_id"`
		CourseName   string `json:"course_name"`
		CategoryID   int32  `json:"category_id"`
		CategoryName string `json:"category_name"`
	}
)
