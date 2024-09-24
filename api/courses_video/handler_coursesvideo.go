package coursesvideo

import (
	"database/sql"
	"encoding/json"
	"io"
	"os"
	"path"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
)

func (h *Handler) CreateCourseVideo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		log.Printf("error parsing form data: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing form data", struct{}{}).WriteResponse(w, r)
		return
	}

	// Retrieve form values
	var req CourseVideoRequest
	req.CourseVideoName = r.FormValue("course_video_name")

	// Convert course_id to int32
	courseIDStr := r.FormValue("course_id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		log.Printf("error converting course_id to int32: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid course_id", struct{}{}).WriteResponse(w, r)
		return
	}
	req.CourseID = int32(courseID)

	// Retrieve the video file from form
	file, handler, err := r.FormFile("path_video")
	if err != nil {
		log.Printf("error retrieving the file: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error retrieving the file", struct{}{}).WriteResponse(w, r)
		return
	}
	defer file.Close()

	// Save the video file
	basePath, err := os.Getwd()
	if err != nil {
		log.Printf("error getting current working directory: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error getting working directory", struct{}{}).WriteResponse(w, r)
		return
	}
	publicPath := path.Join(basePath, "public", "videos")
	videoPath := path.Join(publicPath, handler.Filename)

	dst, err := os.Create(videoPath)
	if err != nil {
		log.Printf("error creating video file: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating video file", struct{}{}).WriteResponse(w, r)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("error copying the file: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error copying the video file", struct{}{}).WriteResponse(w, r)
		return
	}

	// Update the request with the video path
	req.PathVideo = handler.Filename

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		log.Println("error validating request:", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, err.Error(), struct{}{}).WriteResponse(w, r)
		return
	}

	// Save course video data to the database
	err = h.db.CreateCourseVideo(ctx, repo.CreateCourseVideoParams{
		CourseID:        sql.NullInt32{Int32: req.CourseID, Valid: true},
		CourseVideoName: req.CourseVideoName,
		PathVideo:       req.PathVideo,
	})

	if err != nil {
		log.Println("error storing course video to db:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error storing course video", struct{}{}).WriteResponse(w, r)
		return
	}

	// Create a response with the course video data
	responseData := map[string]interface{}{
		"course_video_name": req.CourseVideoName,
		"course_id":         req.CourseID,
		"path_video":        req.PathVideo,
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Course video created successfully", responseData).WriteResponse(w, r)
}

func (h *Handler) GetCourseVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Execute the query
	data, err := h.db.GetCourseVideo(r.Context())
	if err != nil {
		log.Println("error fetching course videos:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Process the data
	var res []GetCourseVideoRow
	for _, d := range data {
		res = append(res, GetCourseVideoRow{
			CourseID:          d.CourseID,
			CourseName:        d.CourseName,
			CourseDescription: d.CourseDescription,
			CategoryName:      d.CategoryName,
			CourseVideoID:     d.CourseVideoID,
			CourseVideoName:   d.CourseVideoName,
			PathVideo:         d.PathVideo,
		})
	}

	// Write the response
	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) GetAllCourseVideos(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetAllCourseVideos(r.Context())
	if err != nil {
		log.Println("error fetching all data item:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []CourseVideo
	for _, d := range data {
		res = append(res, CourseVideo{
			CoursesVideoID:  d.CourseVideoID,
			CourseID:        d.CourseID.Int32,
			CourseVideoName: d.CourseVideoName,
			PathVideo:       d.PathVideo,
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) GetCourseVideoByID(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter ID dari URL
	vars := chi.URLParam(r, "id")

	// Parsing ID dari string ke int32
	id, err := strconv.ParseInt(vars, 10, 32)
	if err != nil {
		log.Println("error parsing ID:", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid ID format", struct{}{}).WriteResponse(w, r)
		return
	}

	// Mendapatkan data CourseVideo dari database berdasarkan ID
	data, err := h.db.GetCourseVideoByID(r.Context(), int32(id))
	if err == sql.ErrNoRows {
		util.NewResponse(http.StatusUnauthorized, http.StatusUnauthorized, "Invalid ID", struct{}{}).WriteResponse(w, r)
		return
	} else if err != nil {
		log.Println("error fetching course video by ID:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Internal server error", struct{}{}).WriteResponse(w, r)
		return
	}

	// Mapping data CourseVideo ke struktur respons
	res := CourseVideo{
		CoursesVideoID:  data.CourseVideoID,
		CourseID:        data.CourseID.Int32,
		CourseVideoName: data.CourseVideoName,
		PathVideo:       data.PathVideo,
	}

	// Mengirimkan respons
	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) UpdateCourseVideo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	// Get the course video ID from the URL parameters
	videoID := chi.URLParam(r, "id")
	if videoID == "" {
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Course video ID is required", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Convert videoID to int32
	id, err := strconv.Atoi(videoID)
	if err != nil {
		log.Println("invalid course video ID:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid course video ID", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	var req CourseVideoRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error parsing request:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing request", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		log.Println("error validation request:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, err.Error(), struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Update the course video in the database
	err = h.db.UpdateCourseVideo(ctx, repo.UpdateCourseVideoParams{
		CourseVideoID:   int32(id),
		CourseID:        util.SqlInt32(req.CourseID), // Ensure CourseID is included
		CourseVideoName: req.CourseVideoName,
		PathVideo:       req.PathVideo,
	})

	if err != nil {
		log.Println("error updating course video in db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Course video updated successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) DeleteCourseVideo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	// Mendapatkan ID CourseVideo dari parameter URL
	videoID := chi.URLParam(r, "id")
	if videoID == "" {
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Course video ID is required", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Mengonversi videoID ke int32
	id, err := strconv.Atoi(videoID)
	if err != nil {
		log.Println("invalid course video ID:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid course video ID", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Menghapus CourseVideo dari database
	err = h.db.DeleteCourseVideo(ctx, int32(id))
	if err != nil {
		log.Println("error deleting course video from db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Course video deleted successfully"
	resp.WriteResponse(w, r)
}
