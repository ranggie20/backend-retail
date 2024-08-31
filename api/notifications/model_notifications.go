package notifications

type (
	// Model Notification yang sesuai dengan tabel notification
	Notification struct {
		NotificationID int32  `json:"notification_id"` // ID notifikasi
		UserID         int32  `json:"user_id"`         // ID pengguna yang menerima notifikasi
		CourseID       int32  `json:"course_id"`       // ID kursus terkait
		Message        string `json:"message"`         // Pesan notifikasi
		IsRead         string `json:"is_read"`         // Status baca notifikasi (misalnya 'true' atau 'false')
	}

	// Model NotificationRequest untuk request input
	NotificationRequest struct {
		UserID   int32  `json:"user_id" validate:"required"` // ID pengguna (wajib diisi)
		CourseID int32  `json:"course_id"`                   // ID kursus (opsional)
		Message  string `json:"message" validate:"required"` // Pesan notifikasi (wajib diisi)
		IsRead   string `json:"is_read" validate:"required"` // Status baca notifikasi (wajib diisi)
	}
)
