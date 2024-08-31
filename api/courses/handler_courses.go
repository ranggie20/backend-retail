package courses

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
)

func (h *Handler) CreateCourses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	var req CourseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error parsing request:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing request", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		log.Println("error validating request:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, err.Error(), struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	now := time.Now()
	// Save Data item to database
	err = h.db.CreateCourse(ctx, repo.CreateCourseParams{
		CourseName:        req.CourseName,
		CourseDescription: req.CourseDescription,
		CategoryID:        util.SqlInt32(req.CategoryID),
		Price:             req.Price,
		Thumbnail:         sql.NullString{String: req.Thumbnail, Valid: req.Thumbnail != ""},
		DeletedAt:         sql.NullTime{},
		CreatedAt:         sql.NullTime{Time: now, Valid: true},
		UpdatedAt:         sql.NullTime{Time: now, Valid: true},
	})

	if err != nil {
		log.Println("error storing category to db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Category has been created successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) GetAllCourses(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetAllCourse(r.Context())
	if err != nil {
		log.Println("error fetching all courses:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []Course
	for _, c := range data {
		// Convert sql.NullString to string
		thumbnail := ""
		if c.Thumbnail.Valid {
			thumbnail = c.Thumbnail.String
		}

		// Convert sql.NullInt32 to int32
		categoryID := int32(0)
		if c.CategoryID.Valid {
			categoryID = c.CategoryID.Int32
		}

		res = append(res, Course{
			CourseID:          c.CourseID,
			CourseName:        c.CourseName,
			CourseDescription: c.CourseDescription,
			CategoryID:        categoryID, // Use converted value
			Price:             c.Price,
			Thumbnail:         thumbnail, // Use converted value
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) GetCourseByID(w http.ResponseWriter, r *http.Request) {
	// Extract the course ID from the URL path
	courseIDParam := chi.URLParam(r, "course_id")
	courseID, err := strconv.Atoi(courseIDParam)
	if err != nil {
		log.Println("invalid course ID:", err)
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	// Fetch the course by ID from the database
	c, err := h.db.GetCourseByID(r.Context(), int32(courseID))
	if err != nil {
		log.Println("error fetching course by ID:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert sql.NullString to string
	thumbnail := ""
	if c.Thumbnail.Valid {
		thumbnail = c.Thumbnail.String
	}

	// Convert sql.NullInt32 to int32
	categoryID := int32(0)
	if c.CategoryID.Valid {
		categoryID = c.CategoryID.Int32
	}

	// Create the response object
	course := Course{
		CourseID:          c.CourseID,
		CourseName:        c.CourseName,
		CourseDescription: c.CourseDescription,
		CategoryID:        categoryID,
		Price:             c.Price,
		Thumbnail:         thumbnail,
	}

	// Send the response
	util.NewResponse(http.StatusOK, http.StatusOK, "", course).WriteResponse(w, r)
}

func (h *Handler) GetCoursePrice(w http.ResponseWriter, r *http.Request) {
	// Call the method that executes the GetCoursePrice query
	data, err := h.db.GetCoursePrice(r.Context())
	if err != nil {
		log.Println("error fetching courses by price:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert the database result to the response format
	var res []CoursePriceRow
	for _, c := range data {
		res = append(res, CoursePriceRow{
			CourseID:          c.CourseID,
			CourseName:        c.CourseName,
			CourseDescription: c.CourseDescription,
			Price:             float64(c.Price),
		})
	}

	// Send the response
	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) GetPopularCourses(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetPopularCourse(r.Context())
	if err != nil {
		log.Println("error fetching popular courses:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []GetPopularCourseRow
	for _, course := range data {
		res = append(res, GetPopularCourseRow{
			CourseID:         course.CourseID,
			CourseName:       course.CourseName,
			TotalEnrollments: course.TotalEnrollments,
			Thumbnail:        course.Thumbnail.String,
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) GetMyCoursePage(w http.ResponseWriter, r *http.Request) {
	// Extract the course ID from the URL path
	userIDContext := r.Context().Value("user_id")

	if userIDContext == nil {
		http.Error(w, "User ID incorrect", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDContext.(int32)
	if !ok {
		log.Println("invalid user ID")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Call the method that executes the GetMyCoursePage query
	data, err := h.db.GetMyCoursePage(r.Context(), util.SqlInt32(userID))
	if err != nil {
		log.Println("error fetching my courses:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert the database result to the response format
	var res []MyCoursePageRow
	for _, c := range data {
		res = append(res, MyCoursePageRow{
			CourseID:          c.CourseID.Int32,
			CourseName:        c.CourseName.String,
			CourseDescription: c.CourseDescription.String,
		})
	}

	// Send the response
	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	// Get the course ID from the URL parameters
	courseID := chi.URLParam(r, "id")
	if courseID == "" {
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Course ID is required", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Convert courseID to int32
	id, err := strconv.Atoi(courseID)
	if err != nil {
		log.Println("invalid course ID:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid course ID", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	var req CourseRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error parsing request:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing request", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		log.Println("error validating request:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, err.Error(), struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Convert Thumbnail to sql.NullString
	thumbnail := sql.NullString{
		String: req.Thumbnail,
		Valid:  req.Thumbnail != "",
	}

	// Update the course in the database
	err = h.db.UpdateCourse(ctx, repo.UpdateCourseParams{
		CourseID:          int32(id),
		CourseName:        req.CourseName,
		CourseDescription: req.CourseDescription,
		CategoryID:        util.SqlInt32(int32(id)),
		Price:             req.Price,
		Thumbnail:         thumbnail,
	})

	if err != nil {
		log.Println("error updating course in db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Course updated successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	// Get the course ID from the URL parameters
	courseID := chi.URLParam(r, "id")
	if courseID == "" {
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Course ID is required", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Convert courseID to int32
	id, err := strconv.Atoi(courseID)
	if err != nil {
		log.Println("invalid course ID:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid course ID", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Delete the course from the database
	err = h.db.DeleteCourse(ctx, int32(id))
	if err != nil {
		log.Println("error deleting course from db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Course deleted successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) GetCourseByNew(w http.ResponseWriter, r *http.Request) {
	// Use the existing connection pool instead of opening a new connection
	db, err := sql.Open("postgres", "user=postgres password=25112004 dbname=postgres sslmode=disable")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	// Query to fetch courses ordered by creation date (newest first)
	query := `
		SELECT course_id, category_id, course_name, course_description, price, thumbnail, created_at, deleted_at, updated_at 
		FROM courses 
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Query execution error", http.StatusInternalServerError)
		log.Println("Query execution error:", err)
		return
	}
	defer rows.Close()

	var courses []Course

	// Iterate over the rows and add each course to the slice
	for rows.Next() {
		var course Course
		if err := rows.Scan(
			&course.CourseID,
			&course.CategoryID,
			&course.CourseName,
			&course.CourseDescription,
			&course.Price,
			&course.Thumbnail,
			&course.CreatedAt,
			&course.DeletedAt,
			&course.UpdatedAt,
		); err != nil {
			http.Error(w, "Error scanning course data", http.StatusInternalServerError)
			log.Println("Error scanning course data:", err)
			return
		}
		courses = append(courses, course)
	}

	// Check for any errors that occurred during iteration
	if err := rows.Err(); err != nil {
		http.Error(w, "Rows iteration error", http.StatusInternalServerError)
		log.Println("Rows iteration error:", err)
		return
	}

	// Convert courses slice to JSON and send it as a response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(courses); err != nil {
		http.Error(w, "Error encoding response to JSON", http.StatusInternalServerError)
		log.Println("Error encoding response to JSON:", err)
	}
}
