package coursesvideo

import (
	"database/sql"
	"time"
)

type (
	// Model CourseVideo yang sesuai dengan tabel courses_video
	CourseVideo struct {
		CoursesVideoID  int32     `json:"courses_video_id"`  // Menggunakan int32 untuk mencocokkan tipe SERIAL
		CourseID        int32     `json:"course_id"`         // ID kursus yang terhubung dengan video
		CourseVideoName string    `json:"course_video_name"` // Nama video kursus
		PathVideo       string    `json:"path_video"`        // Path atau lokasi video
		CreatedAt       time.Time `json:"created_at"`        // Waktu pembuatan video kursus
		DeletedAt       time.Time `json:"deleted_at"`        // Waktu penghapusan video kursus (soft delete)
		UpdatedAt       time.Time `json:"updated_at"`        // Waktu update terakhir video kursus
	}

	// Model CourseVideoRequest untuk request input
	CourseVideoRequest struct {
		CourseID        int32  `json:"course_id" validate:"required"`         // ID kursus (wajib diisi)
		CourseVideoName string `json:"course_video_name" validate:"required"` // Nama video kursus (wajib diisi)
		PathVideo       string `json:"path_video" validate:"required"`        // Path atau lokasi video (wajib diisi)
	}
	GetCourseVideoRow struct {
		CourseID          int32          `json:"course_id"`
		CourseName        string         `json:"course_name"`
		CourseDescription string         `json:"course_description"`
		CategoryName      sql.NullString `json:"category_name"`
		CourseVideoID     sql.NullInt32  `json:"course_video_id"`
		CourseVideoName   sql.NullString `json:"course_video_name"`
		PathVideo         sql.NullString `json:"path_video"`
	}
)
