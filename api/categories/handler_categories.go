package categories

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

// CreateCategory handles the creation of a new category
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	var req CategoryRequest
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
	err = h.db.CreateCategory(ctx, repo.CreateCategoryParams{
		CategoryName: req.CategoryName,
		Icon:         req.Icon,
		CreatedAt:    sql.NullTime{Time: now, Valid: true},
		UpdatedAt:    sql.NullTime{Time: now, Valid: true},
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

// GetAllCategories handles retrieving all categories
func (h *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetAllCategories(r.Context())
	if err != nil {
		log.Println("error fetching all categories:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []Category
	for _, d := range data {
		res = append(res, Category{
			CategoryID:   d.CategoryID,
			CategoryName: d.CategoryName,
			Icon:         d.Icon,
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetAllCategories(r.Context()) // Adjust the method name according to your actual implementation
	if err != nil {
		log.Println("error fetching all categories:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []GetCategoryRow
	for _, c := range data {
		res = append(res, GetCategoryRow{
			CourseID:     c.CategoryID,
			CourseName:   c.CategoryName,
			CategoryID:   c.CategoryID,
			CategoryName: c.CategoryName,
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

// GetCategoryByID handles retrieving a category by its ID
func (h *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Get the category ID from the URL parameters
	categoryID := chi.URLParam(r, "id")

	// Convert categoryID to int32
	id, err := strconv.Atoi(categoryID)
	if err != nil {
		log.Println("error parsing ID:", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid ID format", struct{}{}).WriteResponse(w, r)
		return
	}

	// Get the category from the database by ID
	data, err := h.db.GetCategoryByID(r.Context(), int32(id))
	if err == sql.ErrNoRows {
		util.NewResponse(http.StatusUnauthorized, http.StatusUnauthorized, "Invalid ID", struct{}{}).WriteResponse(w, r)
		return
	} else if err != nil {
		log.Println("error fetching category by ID:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Internal server error", struct{}{}).WriteResponse(w, r)
		return
	}

	// Mapping data Category ke struktur respons
	res := Category{
		CategoryID:   data.CategoryID,
		CategoryName: data.CategoryName,
		Icon:         data.Icon,
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

// UpdateCategory handles updating an existing category
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	// 	// Get the category ID from the URL parameters
	categoryID := chi.URLParam(r, "id")
	if categoryID == "" {
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Category ID is required", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Convert categoryID to int32
	CategoryID, err := strconv.Atoi(categoryID)
	if err != nil {
		log.Println("invalid category ID:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid category ID", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	var req CategoryRequest
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

	now := time.Now()
	// Update the category in the database
	err = h.db.UpdateCategory(ctx, repo.UpdateCategoryParams{
		CategoryID:   int32(CategoryID),
		CategoryName: req.CategoryName,
		Icon:         req.Icon,
		UpdatedAt:    sql.NullTime{Time: now, Valid: true},
	})

	if err != nil {
		log.Println("error updating category in db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Category updated successfully"
	resp.WriteResponse(w, r)
}

// DeleteCategory handles deleting an existing category
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	// Get the category ID from the URL parameters
	categoryID := chi.URLParam(r, "id")
	if categoryID == "" {
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Category ID is required", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Convert categoryID to int32
	id, err := strconv.Atoi(categoryID)
	if err != nil {
		log.Println("invalid category ID:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid category ID", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Delete the category from the database
	err = h.db.DeleteCategory(ctx, int32(id))
	if err != nil {
		log.Println("error deleting category from db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Category deleted successfully"
	resp.WriteResponse(w, r)
}
