package teachers

import (
	"time"
)

type (
	// Model Teacher yang sesuai dengan tabel teachers
	Teacher struct {
		TeacherID   int32     `json:"teacher_id"` // Menggunakan int32 untuk mencocokkan tipe SERIAL
		UserID      int32     `json:"user_id"`
		TeacherName string    `json:"teachername`
		CreatedAt   time.Time `json:"created_at"`
		DeletedAt   time.Time `json:"deleted_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	// Model TeacherRequest untuk request input
	TeacherRequest struct {
		TeacherName string `json:"Nama" validate:"required"`
	}
)
