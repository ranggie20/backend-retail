package coursesvideo

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
)

func (h *Handler) CreateCourseVideo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	var req CourseVideoRequest
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

	// Save Data item to database
	err = h.db.CreateCourseVideo(ctx, repo.CreateCourseVideoParams{
		CourseID:        sql.NullInt32{Int32: req.CourseID, Valid: true},
		CourseVideoName: req.CourseVideoName,
		PathVideo:       req.PathVideo,
	})

	if err != nil {
		log.Println("error storing course video to db: ", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Course video has been created successfully"
	resp.WriteResponse(w, r)
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
