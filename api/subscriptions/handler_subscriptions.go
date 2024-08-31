package subscriptions

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
)

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	var req SubscriptionRequest
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
	// Save Subscription item to database
	err = h.db.CreateSubscription(ctx, repo.CreateSubscriptionParams{
		UserID:    util.SqlInt32(req.UserID),
		CourseID:  util.SqlInt32(req.CourseID),
		CartID:    util.SqlInt32(req.CartID),
		IsCorrect: req.IsCorrect,
		CreatedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	})

	if err != nil {
		log.Println("error storing subscription to db:", err)
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
			CartID:         s.CartID.Int32,
			IsCorrect:      s.IsCorrect,
			CreatedAt:      s.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:      s.UpdatedAt.Time.Format(time.RFC3339),
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}
