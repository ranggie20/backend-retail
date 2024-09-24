package courses

import (
	"database/sql"
)

type (
	// GetPopularCourseRow represents the structure of a popular course with total enrollments.
	GetPopularCourseRow struct {
		CourseID         int32  `json:"course_id"`         // Unique ID of the course
		CourseName       string `json:"course_name"`       // Name of the course
		TotalEnrollments int64  `json:"total_enrollments"` // Total number of enrollments for the course
		Thumbnail        string `json:"thumbnail"`
	}

	// Course represents the structure of a course.
	Course struct {
		CourseID          int32        `json:"course_id"`          // Unique ID of the course
		CourseName        string       `json:"course_name"`        // Name of the course
		CourseDescription string       `json:"course_description"` // Description of the course
		CategoryID        int32        `json:"category_id"`        // Category ID of the course
		Price             int32        `json:"price"`              // Price of the course
		Thumbnail         string       `json:"thumbnail"`          // Thumbnail URL for the course
		Video             string       `json:"video"`
		CreatedAt         sql.NullTime `json:"created_at"` // Timestamp of course creation
		DeletedAt         sql.NullTime `json:"deleted_at"` // Timestamp of course deletion
		UpdatedAt         sql.NullTime `json:"updated_at"` // Timestamp of the last course update
	}

	// CourseRequest represents the structure for creating or updating a course.
	CourseRequest struct {
		CourseName        string `json:"course_name" validate:"required"`        // Name of the course
		CourseDescription string `json:"course_description" validate:"required"` // Description of the course
		CategoryID        int32  `json:"category_id" validate:"required"`        // Category ID of the course
		Price             int32  `json:"price" validate:"required"`              // Price of the course
		Thumbnail         string `json:"thumbnail"`                              // Thumbnail URL (optional)
	}
	GetCourseRow struct {
		CourseID   int32          `json:"course_id"`
		Thumbnail  sql.NullString `json:"thumbnail"`
		CourseName string         `json:"course_name"`
		Price      int32          `json:"price"`
	}
	CoursePriceRow struct {
		CourseID          int32   `json:"course_id"`
		CourseName        string  `json:"course_name"`
		CourseDescription string  `json:"course_description"`
		Price             float64 `json:"price"` // Adjust the type according to your actual price data type
	}
	MyCoursePageRow struct {
		CourseID          int32  `json:"course_id"`
		CourseName        string `json:"course_name"`
		CourseDescription string `json:"course_description"`
	}

	MyCoursePage struct {
		SubscriptionID    int32  `json:"subscription_id"`
		UserID            int32  `json:"user_id"`
		CourseID          int32  `json:"course_id"`
		CourseName        string `json:"course_name"`
		CourseDescription string `json:"course_description"`
		Thumbnail         string `json:"thumbnail"`
		Video             string `json:"video,omitempty"`
	}
)
