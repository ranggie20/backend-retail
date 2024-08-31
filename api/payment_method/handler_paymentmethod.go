package paymentmethod

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

func (h *Handler) CreatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	var req PaymentMethodRequest

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

	// Store payment method in the database
	err = h.db.CreatePaymentMethod(r.Context(), repo.CreatePaymentMethodParams{
		PaymentMethodName: req.PaymentMethodName,
		CreatedAt:         util.SqlTime(time.Now()),
	})

	if err != nil {
		log.Printf("error creating payment method: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating payment method", struct{}{}).WriteResponse(w, r)
		return
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Payment method created successfully", struct{}{}).WriteResponse(w, r)
}

func (h *Handler) GetAllPaymentMethod(w http.ResponseWriter, r *http.Request) {
	// Fetch all payment methods from the database
	paymentMethods, err := h.db.GetAllPaymentMethod(r.Context())
	if err != nil {
		log.Printf("error fetching payment methods: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error fetching payment methods", struct{}{}).WriteResponse(w, r)
		return
	}

	var res []PaymentMethod
	for _, d := range paymentMethods {
		res = append(res, PaymentMethod{
			PaymentMethodID:   d.PaymentMethodID,
			PaymentMethodName: d.PaymentMethodName,
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Payment methods retrieved successfully", res).WriteResponse(w, r)
}

func (h *Handler) GetPaymentMethodById(w http.ResponseWriter, r *http.Request) {
	// Get the payment method ID from the URL parameters
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("invalid payment method ID: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid payment method ID", struct{}{}).WriteResponse(w, r)
		return
	}

	// Fetch the payment method from the database
	paymentMethod, err := h.db.GetPaymentMethodByID(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			util.NewResponse(http.StatusNotFound, http.StatusNotFound, "Payment method not found", struct{}{}).WriteResponse(w, r)
		} else {
			log.Printf("error fetching payment method: %v", err)
			util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error fetching payment method", struct{}{}).WriteResponse(w, r)
		}
		return
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Payment method retrieved successfully", paymentMethod).WriteResponse(w, r)
}
