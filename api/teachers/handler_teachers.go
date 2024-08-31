package teachers

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

func (h *Handler) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	var req TeacherRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error parsing request:", err)
		resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing request", struct{}{})
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

	// Save Teacher to database
	err = h.db.CreateTeacher(ctx, repo.CreateTeacherParams{
		TeacherName: req.TeacherName,
	})

	if err != nil {
		log.Println("error storing teacher to db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Teacher created successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) GetAllTeachers(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetAllTeacher(r.Context())
	if err != nil {
		log.Println("error fetching all teachers:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []Teacher
	for _, t := range data {
		res = append(res, Teacher{
			TeacherID:   t.TeacherID,
			UserID:      t.UserID.Int32,
			TeacherName: t.TeacherName,
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) GetTeacherByID(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(vars, 10, 32)
	if err != nil {
		log.Println("error parsing ID:", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid ID format", struct{}{}).WriteResponse(w, r)
		return
	}

	data, err := h.db.GetTeacherByID(r.Context(), int32(id))
	if err == sql.ErrNoRows {
		util.NewResponse(http.StatusNotFound, http.StatusNotFound, "Teacher not found", struct{}{}).WriteResponse(w, r)
		return
	} else if err != nil {
		log.Println("error fetching teacher by ID:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Internal server error", struct{}{}).WriteResponse(w, r)
		return
	}

	res := Teacher{
		TeacherID:   data.TeacherID,
		UserID:      data.UserID.Int32,
		TeacherName: data.TeacherName,
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	// Get the teacher ID from the URL parameters
	teacherID := chi.URLParam(r, "id")
	if teacherID == "" {
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Teacher ID is required", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Convert teacherID to int32
	id, err := strconv.Atoi(teacherID)
	if err != nil {
		log.Println("invalid teacher ID:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid teacher ID", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	var req TeacherRequest
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

	// Update the teacher in the database
	err = h.db.UpdateTeacher(ctx, repo.UpdateTeacherParams{
		TeacherID:   int32(id),
		TeacherName: req.TeacherName,
	})

	if err != nil {
		log.Println("error updating teacher in db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Teacher updated successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	// Get the teacher ID from the URL parameters
	teacherID := chi.URLParam(r, "id")
	if teacherID == "" {
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Teacher ID is required", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Convert teacherID to int32
	id, err := strconv.Atoi(teacherID)
	if err != nil {
		log.Println("invalid teacher ID:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid teacher ID", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Delete the teacher from the database
	err = h.db.DeleteTeacher(ctx, int32(id))
	if err != nil {
		log.Println("error deleting teacher from db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Teacher deleted successfully"
	resp.WriteResponse(w, r)
}
