package notifications

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
)

func (h *Handler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var req NotificationRequest

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

	// Store notification in the database
	err = h.db.CreateNotification(r.Context(), repo.CreateNotificationParams{
		UserID: sql.NullInt32{
			Int32: req.UserID,
			Valid: true,
		},
		CourseID: sql.NullInt32{
			Int32: req.CourseID,
			Valid: true,
		},
		Message: sql.NullString{
			String: req.Message,
			Valid:  true,
		},
		IsRead: req.IsRead,
	})

	if err != nil {
		log.Printf("error creating notification: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating notification", struct{}{}).WriteResponse(w, r)
		return
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Notification created successfully", struct{}{}).WriteResponse(w, r)
}

func (h *Handler) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetAllNotifications(r.Context())
	if err != nil {
		log.Println("error fetching all notifications:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []Notification
	for _, d := range data {
		res = append(res, Notification{
			NotificationID: d.NotificationID,
			UserID:         d.UserID.Int32,   // Mengakses nilai Int32
			CourseID:       d.CourseID.Int32, // Mengakses nilai Int32
			Message:        d.Message.String, // Mengakses nilai String
			IsRead:         d.IsRead,         // Mengakses nilai String
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) UpdateNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	// Get the notification ID from the URL parameters
	notificationID := chi.URLParam(r, "id")
	if notificationID == "" {
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Notification ID is required", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Convert notificationID to int32
	id, err := strconv.Atoi(notificationID)
	if err != nil {
		log.Println("invalid notification ID:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid notification ID", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	var req NotificationRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error parsing request:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing request", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Validate request
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		var errMsg strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			if field == "UserID" {
				errMsg.WriteString("User ID is required; ")
			} else if field == "Message" {
				errMsg.WriteString("Message is required; ")
			} else {
				errMsg.WriteString(field + " is invalid; ")
			}
		}
		log.Println("validation error:", errMsg.String())
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid input: "+errMsg.String(), struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Update the notification in the database
	err = h.db.UpdateNotification(ctx, repo.UpdateNotificationParams{
		NotificationID: int32(id),
		UserID:         sql.NullInt32{Int32: req.UserID, Valid: true},
		CourseID:       sql.NullInt32{Int32: req.CourseID, Valid: true},
		Message:        sql.NullString{String: req.Message, Valid: true},
		IsRead:         req.IsRead,
	})

	if err != nil {
		log.Println("error updating notification in db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "Notification updated successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	// Get the notification ID from the URL parameters
	notificationID := chi.URLParam(r, "id")

	// Convert notificationID to int32
	id, err := strconv.ParseInt(notificationID, 10, 32)
	if err != nil {
		log.Println("error parsing ID:", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid ID format", struct{}{}).WriteResponse(w, r)
		return
	}

	// Delete the notification from the database
	err = h.db.DeleteNotification(r.Context(), int32(id))
	if err != nil {
		log.Println("error deleting notification:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error deleting notification", struct{}{}).WriteResponse(w, r)
		return
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "Notification deleted successfully", struct{}{}).WriteResponse(w, r)
}
