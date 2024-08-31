package cart

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
)

func (h *Handler) CreateCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	var req CartRequest
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

	// Calculate TotalAmount
	totalAmount := req.Price * req.Quantity

	// Save Cart item to database
	err = h.db.CreateCart(ctx, repo.CreateCartParams{
		UserID: sql.NullInt32{
			Int32: req.UserID,
			Valid: req.UserID != 0,
		},
		CourseID: sql.NullInt32{
			Int32: req.CourseID,
			Valid: req.CourseID != 0,
		},
		Price: sql.NullInt32{
			Int32: req.Price,
			Valid: true,
		},
		Quantity: sql.NullInt32{
			Int32: req.Quantity,
			Valid: true,
		},
		TotalAmount: sql.NullInt32{
			Int32: totalAmount,
			Valid: true,
		},
	})

	if err != nil {
		log.Println("error storing cart item to db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Cart item has been created successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) GetAllCart(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetAllCart(r.Context())
	if err != nil {
		log.Println("error fetching all cart items:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []Cart
	for _, d := range data {
		cartItem := Cart{
			CartID: d.CartID, // Asumsikan CartID selalu valid
		}

		// Konversi UserID
		if d.UserID.Valid {
			cartItem.UserID = d.UserID.Int32
		}

		// Konversi CourseID
		if d.CourseID.Valid {
			cartItem.CourseID = d.CourseID.Int32
		}

		// Konversi Price
		if d.Price.Valid {
			cartItem.Price = d.Price.Int32 // Konversi int32 ke float64
		}

		// Konversi Quantity
		if d.Quantity.Valid {
			cartItem.Quantity = d.Quantity.Int32
		}

		// Konversi TotalAmount
		if d.TotalAmount.Valid {
			cartItem.TotalAmount = d.TotalAmount.Int32 // Konversi int32 ke float64
		}

		// Menambahkan item ke dalam slice
		res = append(res, cartItem)
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) GetCartByUserID(w http.ResponseWriter, r *http.Request) {
	// Ambil user_id dari parameter URL atau query string
	userID := r.Context().Value("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// Panggil metode yang mengeksekusi query GetCartByUserID
	userIDInt, ok := userID.(int32)
	if !ok {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Panggil metode yang mengeksekusi query GetCartByUserID
	data, err := h.db.GetCartByUserID(r.Context(), sql.NullInt32{
		Int32: int32(userIDInt),
		Valid: true,
	})
	if err != nil {
		log.Println("error fetching cart data:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Konversi hasil database ke format response
	var res []GetCartRow
	for _, c := range data {
		res = append(res, GetCartRow{
			CourseID:    c.CourseID.Int32,
			Thumbnail:   c.Thumbnail.String,
			CourseName:  c.CourseName.String,
			Price:       c.Price.Int32,
			Quantity:    c.Quantity.Int32,
			TotalAmount: c.TotalAmount.Int32,
		})
	}

	// Kirim response
	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	// Ambil user_id dari r.Context().Value
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// Cek apakah userID bisa di-cast ke tipe int32
	userIDInt, ok := userID.(int32)
	if !ok {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Menghapus data cart dari database berdasarkan user_id
	err := h.db.DeleteCart(r.Context(), sql.NullInt32{Int32: userIDInt, Valid: true})
	if err == sql.ErrNoRows {
		util.NewResponse(http.StatusNotFound, http.StatusNotFound, "Cart item not found", struct{}{}).WriteResponse(w, r)
		return
	} else if err != nil {
		log.Println("error deleting cart item:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Internal server error", struct{}{}).WriteResponse(w, r)
		return
	}

	// Kirim response jika berhasil
	util.NewResponse(http.StatusOK, http.StatusOK, "Cart item deleted successfully", struct{}{}).WriteResponse(w, r)
}
