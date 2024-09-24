package courses

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
)

func (h *Handler) CreateCourses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse form data
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		log.Printf("error parsing form data: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing form data", struct{}{}).WriteResponse(w, r)
		return
	}

	// Ambil data dari form
	var req CourseRequest
	req.CourseName = r.FormValue("course_name")
	req.CourseDescription = r.FormValue("course_description")

	// Convert category_id to int32
	categoryIDStr := r.FormValue("category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Printf("error converting category_id to int32: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid category_id", struct{}{}).WriteResponse(w, r)
		return
	}
	req.CategoryID = int32(categoryID) // Convert int to int32

	// Convert price to float64
	priceStr := r.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Printf("error converting price to float64: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid price", struct{}{}).WriteResponse(w, r)
		return
	}
	req.Price = int32(price)

	// Ambil file thumbnail dari form
	file, handler, err := r.FormFile("thumbnail")
	if err != nil {
		log.Printf("error retrieving the file: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error retrieving the file", struct{}{}).WriteResponse(w, r)
		return
	}
	defer file.Close()

	// Simpan file thumbnail
	basePath, err := os.Getwd()
	if err != nil {
		log.Printf("error getting current working directory: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error getting working directory", struct{}{}).WriteResponse(w, r)
		return
	}
	publicPath := path.Join(basePath, "public", "course")
	thumbnailPath := path.Join(publicPath, handler.Filename)
	dst, err := os.Create(thumbnailPath)
	if err != nil {
		log.Printf("error creating file: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating file", struct{}{}).WriteResponse(w, r)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("error copying the file: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error copying the file", struct{}{}).WriteResponse(w, r)
		return
	}

	// Ambil file video dari form
	fileVideo, videoHandler, err := r.FormFile("video")
	if err != nil {
		log.Printf("error retrieving the file: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error retrieving the video", struct{}{}).WriteResponse(w, r)
		return
	}
	defer fileVideo.Close()

	// Simpan file video
	publicPath = path.Join(basePath, "public", "video", "course")
	videoPath := path.Join(publicPath, videoHandler.Filename)
	videoDst, err := os.Create(videoPath)
	if err != nil {
		log.Printf("error creating file: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating file", struct{}{}).WriteResponse(w, r)
		return
	}
	defer videoDst.Close()

	if _, err := io.Copy(videoDst, fileVideo); err != nil {
		log.Printf("error copying the file: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error copying the file", struct{}{}).WriteResponse(w, r)
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		log.Println("error validating request:", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, err.Error(), struct{}{}).WriteResponse(w, r)
		return
	}

	now := time.Now()
	// Save course data to the database
	err = h.db.CreateCourse(ctx, repo.CreateCourseParams{
		CourseName:        req.CourseName,
		CourseDescription: req.CourseDescription,
		CategoryID:        util.SqlInt32(req.CategoryID),
		Price:             req.Price,
		Thumbnail:         util.SqlString(path.Join("static", "course", handler.Filename)),
		DeletedAt:         sql.NullTime{},
		CreatedAt:         sql.NullTime{Time: now, Valid: true},
		UpdatedAt:         sql.NullTime{Time: now, Valid: true},
	})

	if err != nil {
		log.Println("error storing course to db:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating course", struct{}{}).WriteResponse(w, r)
		return
	}

	// Get final course id
	courseID, err := h.db.GetLastCourseID(ctx)
	if err != nil {
		log.Println("error getting last course id:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating course", struct{}{}).WriteResponse(w, r)
	}

	// Save course video to the database
	err = h.db.CreateCourseVideo(ctx, repo.CreateCourseVideoParams{
		CourseID:        util.SqlInt32(courseID),
		CourseVideoName: req.CourseName,
		PathVideo:       path.Join("static", "video", "course", videoHandler.Filename),
		CreatedAt:       util.SqlTime(time.Now()),
		UpdatedAt:       util.SqlTime(time.Now()),
	})
	if err != nil {
		log.Println("error storing course video to db:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating course video", struct{}{}).WriteResponse(w, r)
		return
	}

	// Create a response with the course data
	responseData := map[string]interface{}{
		"course_name":        req.CourseName,
		"course_description": req.CourseDescription,
		"category_id":        req.CategoryID,
		"price":              req.Price,
		"thumbnail":          path.Join("static", "course", handler.Filename),
		"video":              path.Join("static", "video", "course", videoHandler.Filename),
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Course created successfully", responseData).WriteResponse(w, r)
}

func (h *Handler) GetMyCourse(w http.ResponseWriter, r *http.Request) {
	// Extract course_id from the URL path
	courseIDParam := chi.URLParam(r, "course_id")
	courseID, err := strconv.ParseInt(courseIDParam, 10, 32)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	// Extract user ID from the token
	userIDContext := r.Context().Value("user_id")
	if userIDContext == nil {
		http.Error(w, "User ID missing", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDContext.(int32)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Prepare the parameters for the query
	params := repo.GetMyCourseParams{
		CourseID: sql.NullInt32{Int32: int32(courseID), Valid: true},
		UserID:   sql.NullInt32{Int32: userID, Valid: true},
	}

	// Call the method that executes the GetMyCourse query
	data, err := h.db.GetMyCourse(r.Context(), params)
	if err == sql.ErrNoRows {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert the database result to the response format
	res := MyCoursePage{
		SubscriptionID:    data.SubscriptionID,
		UserID:            data.UserID.Int32,
		CourseID:          data.CourseID.Int32,
		CourseName:        data.CourseName.String,
		CourseDescription: data.CourseDescription.String,
		Thumbnail:         data.Thumbnail.String,
		Video:             data.PathVideo.String,
	}

	// Send the response
	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
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
		Video:             c.PathVideo.String,
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
	data, err := h.db.GetAllMyCourse(r.Context(), util.SqlInt32(userID))
	if err != nil {
		log.Println("error fetching my courses:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert the database result to the response format
	var res []MyCoursePage
	for _, c := range data {
		res = append(res, MyCoursePage{
			SubscriptionID:    c.SubscriptionID,
			UserID:            c.UserID.Int32,
			CourseID:          c.CourseID.Int32,
			CourseName:        c.CourseName.String,
			CourseDescription: c.CourseDescription.String,
			Thumbnail:         c.Thumbnail.String,
		})
	}

	// Send the response
	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	courseID := chi.URLParam(r, "id")
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", struct{}{})
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

	// Parse form data
	err = r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		log.Printf("error parsing form data: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing form data", struct{}{}).WriteResponse(w, r)
		return
	}

	// Ambil data dari form
	var req CourseRequest
	req.CourseName = r.FormValue("course_name")
	req.CourseDescription = r.FormValue("course_description")

	// Convert category_id to int32
	categoryIDStr := r.FormValue("category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Printf("error converting category_id to int32: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid category_id", struct{}{}).WriteResponse(w, r)
		return
	}
	req.CategoryID = int32(categoryID) // Convert int to int32

	// Convert price to float64
	priceStr := r.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Printf("error converting price to float64: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid price", struct{}{}).WriteResponse(w, r)
		return
	}
	req.Price = int32(price)

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		log.Println("error validating request:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, err.Error(), struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Ambil file thumbnail dari form
	file, handler, err := r.FormFile("thumbnail")
	fileExist := false
	filePath := ""
	if err == nil {
		defer file.Close()
		fileExist = true

		// Simpan file thumbnail
		basePath, err := os.Getwd()
		if err != nil {
			log.Printf("error getting current working directory: %v", err)
			util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error getting working directory", struct{}{}).WriteResponse(w, r)
			return
		}
		publicPath := path.Join(basePath, "public", "course")
		thumbnailPath := path.Join(publicPath, handler.Filename)
		dst, err := os.Create(thumbnailPath)
		if err != nil {
			log.Printf("error creating file: %v", err)
			util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating file", struct{}{}).WriteResponse(w, r)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			log.Printf("error copying the file: %v", err)
			util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error copying the file", struct{}{}).WriteResponse(w, r)
			return
		}

		filePath = path.Join("static", "course", handler.Filename)
	}

	// Ambil file video dari form
	fileVideo, videoHandler, err := r.FormFile("video")
	videoFileExist := false
	videoFilePath := ""
	if err == nil {
		defer fileVideo.Close()

		// Simpan file video
		videoFileExist = true

		// Simpan file thumbnail
		basePath, err := os.Getwd()
		if err != nil {
			log.Printf("error getting current working directory: %v", err)
			util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error getting working directory", struct{}{}).WriteResponse(w, r)
			return
		}
		publicPath := path.Join(basePath, "public", "video", "course")
		videoPath := path.Join(publicPath, videoHandler.Filename)
		videoDst, err := os.Create(videoPath)
		if err != nil {
			log.Printf("error creating file: %v", err)
			util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating file", struct{}{}).WriteResponse(w, r)
			return
		}
		defer videoDst.Close()

		if _, err := io.Copy(videoDst, fileVideo); err != nil {
			log.Printf("error copying the file: %v", err)
			util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error copying the file", struct{}{}).WriteResponse(w, r)
			return
		}

		videoFilePath = path.Join("static", "video", "course", videoHandler.Filename)
	}

	// Get current data
	course, err := h.db.GetCourseByID(r.Context(), int32(id))
	if err != nil {
		log.Println("error getting course in db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	if !fileExist {
		filePath = course.Thumbnail.String
	}

	fmt.Println(filePath)
	fmt.Println(id)

	// Update the course in the database
	err = h.db.UpdateCourse(r.Context(), repo.UpdateCourseParams{
		CourseID:          int32(id),
		CourseName:        req.CourseName,
		CourseDescription: req.CourseDescription,
		CategoryID:        util.SqlInt32(int32(id)),
		Price:             req.Price,
		Thumbnail:         util.SqlString(filePath),
		UpdatedAt:         util.SqlTime(time.Now()),
	})
	if err != nil {
		log.Println("error updating course in db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	videoCourse, err := h.db.GetCourseVideoByCourseID(r.Context(), util.SqlInt32(course.CourseID))
	if err != nil {
		log.Println("error getting course video in db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	if !videoFileExist {
		videoFilePath = videoCourse.PathVideo
	}

	err = h.db.UpdateCourseVideo(r.Context(), repo.UpdateCourseVideoParams{
		CourseID:        util.SqlInt32(course.CourseID),
		CourseVideoName: videoCourse.CourseVideoName,
		PathVideo:       videoFilePath,
		UpdatedAt:       util.SqlTime(time.Now()),
		CourseVideoID:   videoCourse.CourseVideoID,
	})
	if err != nil {
		log.Println("error updating course video in db:", err)
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
