package payment

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
)

func (h *Handler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req PaymentRequest

	// Decode JSON request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("error parsing request: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing request", struct{}{}).WriteResponse(w, r)
		return
	}

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

	// Store payment in the database
	err = h.db.CreatePayment(r.Context(), repo.CreatePaymentParams{
		UserID:          sql.NullInt32{Int32: req.UserID, Valid: true},
		CourseID:        sql.NullInt32{Int32: req.CourseID, Valid: true},
		SubscriptionID:  sql.NullInt32{Int32: req.SubscriptionID, Valid: true},
		PaymentMethodID: sql.NullInt32{Int32: req.PaymentMethodID, Valid: true},
		PaymentStatusID: sql.NullInt32{Int32: req.PaymentStatusID, Valid: true},
		TotalAmount:     util.SqlInt32(req.TotalAmount),
	})

	if err != nil {
		log.Printf("error creating payment: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating payment", struct{}{}).WriteResponse(w, r)
		return
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Payment created successfully", struct{}{}).WriteResponse(w, r)
}

func (h *Handler) GetAllPayment(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetAllPayment(r.Context())
	if err != nil {
		log.Println("error fetching all payments:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []Payment
	for _, d := range data {
		res = append(res, Payment{
			PaymentID:       d.PaymentID,
			UserID:          d.UserID.Int32,          // Access Int32 value
			CourseID:        d.CourseID.Int32,        // Access Int32 value
			SubscriptionID:  d.SubscriptionID.Int32,  // Access Int32 value
			PaymentMethodID: d.PaymentMethodID.Int32, // Access Int32 value
			PaymentStatusID: d.PaymentStatusID.Int32, // Access Int32 value
			TotalAmount:     d.TotalAmount.Int32,
			PaymentDate:     d.PaymentDate.Time,
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) GetPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Execute the query to get payment details
	data, err := h.db.GetPayment(ctx)
	if err != nil {
		log.Println("error fetching payment data:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert query result to GetPaymentRow model
	var res []GetPaymentRow
	for _, d := range data {
		res = append(res, GetPaymentRow{
			CartID:            d.CartID,
			PaymentID:         d.PaymentID,
			CourseName:        d.CourseName,
			Quantity:          d.Quantity,
			TotalAmount:       d.TotalAmount,
			PaymentMethodName: d.PaymentMethodName,
			PaymentDate:       d.PaymentDate,
		})
	}

	// Create response
	resp := util.NewResponse(http.StatusOK, http.StatusOK, "", res)
	resp.WriteResponse(w, r)
}

func (h *Handler) GetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Execute the query to get payment history details
	data, err := h.db.GetPaymentHistory(ctx)
	if err != nil {
		log.Println("error fetching payment history data:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert query result to GetPaymentHistoryRow model
	var res []GetPaymentHistoryRow
	for _, d := range data {
		res = append(res, GetPaymentHistoryRow{
			PaymentID:             d.PaymentID,
			CourseID:              d.CourseID,
			CourseName:            d.CourseName,
			TotalAmount:           d.TotalAmount,
			PaymentStatusName:     d.PaymentStatusName,
			SubcriptionsStartDate: d.SubcriptionsStartDate,
		})
	}

	// Create response
	resp := util.NewResponse(http.StatusOK, http.StatusOK, "", res)
	resp.WriteResponse(w, r)
}

// func (h *Handler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

// 	// Get the payment ID from the URL parameters
// 	paymentID := chi.URLParam(r, "id")
// 	if paymentID == "" {
// 		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Payment ID is required", struct{}{})
// 		resp.WriteResponse(w, r)
// 		return
// 	}

// 	// Convert paymentID to int32
// 	id, err := strconv.Atoi(paymentID)
// 	if err != nil {
// 		log.Println("invalid payment ID:", err)
// 		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid payment ID", struct{}{})
// 		resp.WriteResponse(w, r)
// 		return
// 	}

// 	var req PaymentRequest
// 	err = json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		log.Println("error parsing request:", err)
// 		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing request", struct{}{})
// 		resp.WriteResponse(w, r)
// 		return
// 	}

// 	// Validate request
// 	validate := validator.New()
// 	err = validate.Struct(req)
// 	if err != nil {
// 		var errMsg strings.Builder
// 		for _, err := range err.(validator.ValidationErrors) {
// 			field := err.Field()
// 			if field == "Amount" {
// 				errMsg.WriteString("Amount is required; ")
// 			} else {
// 				errMsg.WriteString(field + " is invalid; ")
// 			}
// 		}
// 		log.Println("validation error:", errMsg.String())
// 		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid input: "+errMsg.String(), struct{}{})
// 		resp.WriteResponse(w, r)
// 		return
// 	}

// 	// Update the payment in the database
// 	err = h.db.UpdatePayment(ctx, repo.UpdatePaymentParams{
// 		PaymentID:       int32(id),
// 		UserID:          sql.NullInt32{Int32: req.UserID, Valid: true},
// 		CourseID:        sql.NullInt32{Int32: req.CourseID, Valid: true},
// 		SubscriptionID:  sql.NullInt32{Int32: req.SubscriptionID, Valid: true},
// 		PaymentMethodID: sql.NullInt32{Int32: req.PaymentMethodID, Valid: true},
// 		PaymentStatusID: sql.NullInt32{Int32: req.PaymentStatusID, Valid: true},
// 		Amount:          req.Amount,
// 		PaymentDate:     req.PaymentDate,
// 	})

// 	if err != nil {
// 		log.Println("error updating payment in db:", err)
// 		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
// 		resp.WriteResponse(w, r)
// 		return
// 	}

// 	resp.Status = http.StatusOK
// 	resp.Code = http.StatusOK
// 	resp.Message = "Payment updated successfully"
// 	resp.WriteResponse(w, r)
// }

// func (h *Handler) DeletePayment(w http.ResponseWriter, r *http.Request) {
// 	// Get the payment ID from the URL parameters
// 	paymentID := chi.URLParam(r, "id")

// 	// Convert paymentID to int32
// 	id, err := strconv.ParseInt(paymentID, 10, 32)
// 	if err != nil {
// 		log.Println("error parsing ID:", err)
// 		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid ID format", struct{}{}).WriteResponse(w, r)
// 		return
// 	}

// 	// Delete the payment from the database
// 	err = h.db.DeletePayment(r.Context(), int32(id))
// 	if err != nil {
// 		log.Println("error deleting payment:", err)
// 		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error deleting payment", struct{}{}).WriteResponse(w, r)
// 		return
// 	}

// 	util.NewResponse(http.StatusOK, http.StatusOK, "Payment deleted successfully", struct{}{}).WriteResponse(w, r)
// }
