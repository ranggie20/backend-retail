package wishlist

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

func (h *Handler) CreateWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	var req WishlistRequest
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

	// Save Wishlist item to database
	err = h.db.CreateWishlist(ctx, repo.CreateWishlistParams{
		UserID: sql.NullInt32{
			Int32: req.UserID,
			Valid: true, // Menandakan bahwa nilai ini valid
		},
		CourseID: sql.NullInt32{
			Int32: req.CourseID,
			Valid: true, // Menandakan bahwa nilai ini valid
		},
	})

	if err != nil {
		log.Println("error storing wishlist to db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Wishlist item has been created successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) GetAllWishlist(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetAllWishlists(r.Context())
	if err != nil {
		log.Println("error fetching all wishlist items:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []Wishlist
	for _, d := range data {
		res = append(res, Wishlist{
			WishlistID: d.WishlistID,
			UserID:     d.UserID.Int32,
			CourseID:   d.CourseID.Int32,
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) GetWishlistByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		log.Println("error parsing ID:", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid ID format", struct{}{}).WriteResponse(w, r)
		return
	}

	data, err := h.db.GetWishlistByID(r.Context(), int32(id))
	if err == sql.ErrNoRows {
		util.NewResponse(http.StatusUnauthorized, http.StatusUnauthorized, "Invalid ID", struct{}{}).WriteResponse(w, r)
		return
	} else if err != nil {
		log.Println("error fetching wishlist by ID:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Internal server error", struct{}{}).WriteResponse(w, r)
		return
	}

	res := Wishlist{
		WishlistID: data.WishlistID,
		UserID:     data.UserID.Int32,
		CourseID:   data.CourseID.Int32,
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) UpdateWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	wishlistID := chi.URLParam(r, "id")
	if wishlistID == "" {
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Wishlist ID is required", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	id, err := strconv.Atoi(wishlistID)
	if err != nil {
		log.Println("invalid wishlist ID:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid wishlist ID", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	var req WishlistRequest
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

	// Update Wishlist item in database
	err = h.db.UpdateWishlist(ctx, repo.UpdateWishlistParams{
		WishlistID: int32(id),
		UserID: sql.NullInt32{
			Int32: req.UserID,
			Valid: true, // Menandakan bahwa nilai ini valid
		},
		CourseID: sql.NullInt32{
			Int32: req.CourseID,
			Valid: true, // Menandakan bahwa nilai ini valid
		},
	})

	if err != nil {
		log.Println("error updating wishlist in db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Wishlist item updated successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	// Get the user_id and course_id parameters from the URL
	userIDParam := chi.URLParam(r, "user_id")

	// Parse user_id and course_id from string to int32
	userID, err := strconv.ParseInt(userIDParam, 10, 32)
	if err != nil {
		log.Println("error parsing user ID:", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid user ID format", struct{}{}).WriteResponse(w, r)
		return
	}

	// Create a sql.NullInt32 value
	nullUserID := sql.NullInt32{
		Int32: int32(userID),
		Valid: true, // Set Valid to true since we have a valid user ID
	}

	// Delete the wishlist item from the database
	err = h.db.DeleteWishlist(r.Context(), nullUserID)
	if err == sql.ErrNoRows {
		util.NewResponse(http.StatusNotFound, http.StatusNotFound, "Wishlist item not found", struct{}{}).WriteResponse(w, r)
		return
	} else if err != nil {
		log.Println("error deleting wishlist item:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Internal server error", struct{}{}).WriteResponse(w, r)
		return
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Wishlist item deleted successfully", struct{}{}).WriteResponse(w, r)
}
