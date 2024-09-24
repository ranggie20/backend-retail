package categories

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
)

// CreateCategory handles the creation of a new category
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CategoryRequest

	// Parse form data
	err := r.ParseMultipartForm(10 << 20) // batasan ukuran file (10MB)
	if err != nil {
		log.Printf("error parsing form data: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing form data", struct{}{}).WriteResponse(w, r)
		return
	}

	// Ambil data dari form
	req.CategoryName = r.FormValue("category_name")

	// Ambil file icon dari form
	file, handler, err := r.FormFile("icon")
	if err != nil {
		log.Printf("error retrieving the file: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error retrieving the file", struct{}{}).WriteResponse(w, r)
		return
	}
	defer file.Close()

	// Simpan file icon
	basePath, _ := os.Getwd()
	publicPath := path.Join(basePath, "public", "category")
	iconPath := path.Join(publicPath, handler.Filename)
	dst, err := os.Create(iconPath)
	if err != nil {
		log.Printf("error saving the file: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error saving the file", struct{}{}).WriteResponse(w, r)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("error copying the file: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error copying the file", struct{}{}).WriteResponse(w, r)
		return
	}

	// Update the request with the icon path
	req.Icon = handler.Filename

	// Validate request
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		var errMsg strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			errMsg.WriteString(err.Field() + " is invalid; ")
		}
		log.Printf("validation error: %v", errMsg.String())
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid input: "+errMsg.String(), struct{}{}).WriteResponse(w, r)
		return
	}

	now := time.Now()
	// Save category data to the database
	err = h.db.CreateCategory(r.Context(), repo.CreateCategoryParams{
		CategoryName: req.CategoryName,
		Icon:         util.SqlString(path.Join("static", "category", handler.Filename)).String, // Save the path to the icon
		CreatedAt:    sql.NullTime{Time: now, Valid: true},
		UpdatedAt:    sql.NullTime{Time: now, Valid: true},
	})

	if err != nil {
		log.Printf("error storing category to db: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating category", struct{}{}).WriteResponse(w, r)
		return
	}

	// Create a response with the category data
	responseData := map[string]interface{}{
		"category_name": req.CategoryName,
		"icon":          req.Icon,
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Category created successfully", responseData).WriteResponse(w, r)
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

func (h *Handler) GetCoursesByCategoryID(w http.ResponseWriter, r *http.Request) {
	// Get the category ID from the URL parameters
	categoryIDParam := chi.URLParam(r, "category_id")

	// Convert categoryID to int32
	categoryID, err := strconv.Atoi(categoryIDParam)
	if err != nil {
		log.Println("invalid category ID:", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid category ID", struct{}{}).WriteResponse(w, r)
		return
	}

	nullCategoryID := sql.NullInt32{
		Int32: int32(categoryID),
		Valid: true,
	}

	course, err := h.db.GetCoursesByCategoryID(r.Context(), nullCategoryID)
	if err == sql.ErrNoRows {
		util.NewResponse(http.StatusNotFound, http.StatusNotFound, "No courses found for this category", struct{}{}).WriteResponse(w, r)
		return
	} else if err != nil {
		log.Println("error fetching courses by category ID:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Internal server error", struct{}{}).WriteResponse(w, r)
		return
	}

	// Prepare the response
	var res []Course
	for _, c := range course {
		res = append(res, Course{
			CourseID:          c.CourseID,
			CourseName:        c.CourseName,
			CourseDescription: c.CourseDescription,
			CategoryID:        c.CategoryID.Int32,
			Price:             float64(c.Price),
			Thumbnail:         sql.NullString{String: c.Thumbnail.String, Valid: true},
			CreatedAt:         c.CreatedAt.Time,
			UpdatedAt:         c.UpdatedAt.Time,
			DeletedAt:         sql.NullTime{Time: c.DeletedAt.Time, Valid: true},
		})
	}

	// Send the response
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
