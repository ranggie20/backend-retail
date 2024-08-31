package transactionhistory

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
)

func (h *Handler) CreateTransactionHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	var req TransactionHistoryRequest
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

	// Convert bool to string for IsPaid
	// isPaidStr := strconv.FormatBool(req.IsPaid)

	// Convert and process the input fields
	// subscriptionsID := util.SqlInt32(req.SubscriptionsID)
	// totalAmount := util.ConvertStringIDToInt32(strconv.FormatFloat(req.TotalAmount, 'f', -1, 64))
	// subscriptionsStartDate := req.SubscriptionsStartDate
	// proof := util.SqlString(req.Proof)

	// Save TransactionHistory to the database
	err = h.db.CreateTransactionHistory(ctx, repo.CreateTransactionHistoryParams{
		SubscriptionID:        util.SqlInt32(req.SubscriptionsID),
		Quantity:              req.Quantity,
		TotalAmount:           util.ConvertStringIDToInt32(strconv.FormatFloat(req.TotalAmount, 'f', -1, 64)),
		IsPaid:                req.IsPaid,
		SubcriptionsStartDate: util.SqlTime(util.ConvertStringToDate(req.SubscriptionsStartDate)),
		Proof:                 util.SqlString(req.Proof),
		CreatedAt:             util.SqlTime(time.Now()),
	})

	if err != nil {
		log.Println("error storing transaction history to db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Transaction history created successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) GetAllTransactionHistory(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetAllTransactionHistory(r.Context())
	if err != nil {
		log.Println("error fetching all transaction history:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []TransactionHistory
	for _, t := range data {
		res = append(res, TransactionHistory{
			TransactionHistoryID:   t.TransactionHistoryID,
			SubscriptionsID:        t.SubscriptionID.Int32,
			Quantity:               t.Quantity,
			TotalAmount:            float64(t.TotalAmount),
			IsPaid:                 t.IsPaid,
			SubscriptionsStartDate: t.SubcriptionsStartDate.Time.Format("2006-01-02"),
			Proof:                  t.Proof.String,
			CreatedAt:              t.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:              t.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:              t.DeletedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}
