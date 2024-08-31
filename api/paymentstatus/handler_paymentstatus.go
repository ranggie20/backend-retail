package paymentstatus

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
)

func (h *Handler) CreatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	var req PaymentStatusRequest

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

	// Store payment status in the database
	err = h.db.CreatePaymentStatus(r.Context(), repo.CreatePaymentStatusParams{
		PaymentStatusName: req.PaymentStatusName,
		CreatedAt:         util.SqlTime(time.Now()),
	})

	if err != nil {
		log.Printf("error creating payment status: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating payment status", struct{}{}).WriteResponse(w, r)
		return
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Payment status created successfully", struct{}{}).WriteResponse(w, r)
}

func (h *Handler) GetAllPaymentStatus(w http.ResponseWriter, r *http.Request) {
	// Fetch all payment statuses from the database
	paymentStatuses, err := h.db.GetAllPaymentStatus(r.Context())
	if err != nil {
		log.Printf("error fetching payment statuses: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error fetching payment statuses", struct{}{}).WriteResponse(w, r)
		return
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Payment statuses retrieved successfully", paymentStatuses).WriteResponse(w, r)
}

func (h *Handler) GetPaymentStatusById(w http.ResponseWriter, r *http.Request) {
	// Get the payment status ID from the URL parameters
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("invalid payment status ID: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid payment status ID", struct{}{}).WriteResponse(w, r)
		return
	}

	// Fetch the payment status from the database
	paymentStatus, err := h.db.GetPaymentStatusByID(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			util.NewResponse(http.StatusNotFound, http.StatusNotFound, "Payment status not found", struct{}{}).WriteResponse(w, r)
		} else {
			log.Printf("error fetching payment status: %v", err)
			util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error fetching payment status", struct{}{}).WriteResponse(w, r)
		}
		return
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Payment status retrieved successfully", paymentStatus).WriteResponse(w, r)
}
