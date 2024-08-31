package subscriptions

import (
	"log"
	"net/http"
	"time"

	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
)

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	// Get user ID
	userID := ctx.Value("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userIDInt, ok := userID.(int32)
	if !ok {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Get last payment ID
	paymentID, err := h.db.GetLastPayment(ctx, util.SqlInt32(userIDInt))
	if err != nil {
		http.Error(w, "cannot get last payment", http.StatusInternalServerError)
		return
	}

	// Get courses from cart
	courses, err := h.db.GetCartByUserID(ctx, util.SqlInt32(userIDInt))
	if err != nil {
		http.Error(w, "cannot get courses", http.StatusInternalServerError)
		return
	}

	now := time.Now()

	for _, c := range courses {
		// Insert to subscription
		err = h.db.CreateSubscription(ctx, repo.CreateSubscriptionParams{
			UserID:    util.SqlInt32(userIDInt),
			CourseID:  util.SqlInt32(c.CourseID.Int32),
			IsCorrect: "yes",
			PaymentID: util.SqlInt32(paymentID),
			CreatedAt: util.SqlTime(now),
			UpdatedAt: util.SqlTime(now),
		})

		if err != nil {
			log.Println("error storing subscription to db:", err)
			resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
			resp.WriteResponse(w, r)
			return
		}

		// Get the subscription id
		subscriptionID, err := h.db.GetLastSubscription(ctx, util.SqlInt32(userIDInt))
		if err != nil {
			log.Println("error getting last subscription id:", err)
			resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
			resp.WriteResponse(w, r)
			return
		}

		// Insert to transaction history
		err = h.db.CreateTransactionHistory(ctx, repo.CreateTransactionHistoryParams{
			SubscriptionID:        util.SqlInt32(subscriptionID),
			Quantity:              c.Quantity.Int32,
			TotalAmount:           c.TotalAmount.Int32,
			IsPaid:                "yes",
			SubcriptionsStartDate: util.SqlTime(now),
			Proof:                 util.SqlString("-"),
			CreatedAt:             util.SqlTime(now),
			UpdatedAt:             util.SqlTime(now),
		})
		if err != nil {
			log.Println("error saving to transaction history:", err)
			resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
			resp.WriteResponse(w, r)
			return
		}
	}

	// Delete cart because we have processeed payment
	err = h.db.DeleteCart(ctx, util.SqlInt32(userIDInt))
	if err != nil {
		log.Println("error deleting cart: ", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Subscription has been created successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetAllSubscriptions(r.Context())
	if err != nil {
		log.Println("error fetching all subscriptions:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []Subscription
	for _, s := range data {
		res = append(res, Subscription{
			SubscriptionID: s.SubscriptionID,
			UserID:         s.UserID.Int32,
			CourseID:       s.CourseID.Int32,
			IsCorrect:      s.IsCorrect,
			CreatedAt:      s.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      s.UpdatedAt.Time.Format(time.RFC3339),
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}
