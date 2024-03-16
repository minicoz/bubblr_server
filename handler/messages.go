package handler

import (
	"bubblr/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	userID := queryParams.Get("userId")
	correspondingUserId := queryParams.Get("correspondingUserId")

	m, err := h.d.GetMessages(userID, correspondingUserId)
	if err != nil {
		msg := fmt.Sprintf("unable to get user %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	var res []*models.Messages
	for _, msg := range m {
		res = append(res, &models.Messages{
			ToUserID:   msg.ToUserID,
			FromUserID: msg.FromUserID,
			TxtMessage: msg.TxtMessage,
			CreatedAt:  msg.CreatedAt})
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error with users", http.StatusBadRequest)
	}

}

func (h *Handler) GetLatestMessages(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	userID := queryParams.Get("userId")
	correspondingUserId := queryParams.Get("correspondingUserIds")
	var uuids []string

	// Unmarshal the JSON string into the slice
	err := json.Unmarshal([]byte(correspondingUserId), &uuids)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	var res []*models.Messages
	for _, u := range uuids {
		msg, err := h.d.GetLatestMessages(userID, string(u))
		if err != nil {
			msg := fmt.Sprintf("unable to get user %v", err)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		res = append(res, &models.Messages{
			ToUserID:   msg.ToUserID,
			FromUserID: msg.FromUserID,
			TxtMessage: msg.TxtMessage,
			CreatedAt:  msg.CreatedAt})
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error with users", http.StatusBadRequest)
	}

}

func (h *Handler) AddMessage(w http.ResponseWriter, r *http.Request) {

	var msg models.Messages
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.d.AddMessage(msg.FromUserID, msg.ToUserID, msg.TxtMessage)
	if err != nil {
		msg := fmt.Sprintf("unable to write message %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
