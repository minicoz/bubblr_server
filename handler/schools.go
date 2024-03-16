package handler

import (
	"bubblr/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetSchoolsWithTier will return all schools with tiers
func (h *Handler) GetSchoolsWithTier(w http.ResponseWriter, r *http.Request) {
	schools, err := h.d.GetSchools()
	if err != nil {
		msg := fmt.Sprintf("unable to get user %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var res []models.School
	for _, s := range schools {
		res = append(res, models.School{
			ID:     s.ID,
			School: s.School,
			Tier:   s.Tier,
		})
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error with get schools", http.StatusBadRequest)
	}
}

// GetSchools will return all schools without tiers
func (m *Handler) GetSchools(w http.ResponseWriter, r *http.Request) {
	schools, err := m.d.GetSchools()
	if err != nil {
		msg := fmt.Sprintf("unable to get user %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var res []models.School
	for _, s := range schools {
		res = append(res, models.School{
			ID:     s.ID,
			School: s.School,
		})
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error with get schools", http.StatusBadRequest)
	}
}
