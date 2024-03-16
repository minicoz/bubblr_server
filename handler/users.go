package handler

import (
	"bubblr/datastore"
	"bubblr/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

func filterList(list1 []models.UserAll, list2 []datastore.Matches) []models.UserAll {
	result := make([]models.UserAll, 0)

	// Create a map to store elements from list2 for faster lookup
	list2Map := make(map[string]bool)
	for _, value := range list2 {
		list2Map[value.Matches] = true
	}

	// Iterate through elements of list1
	for _, value := range list1 {
		// If the element is not present in list2, add it to the result
		if !list2Map[value.UserID] {
			result = append(result, value)
		}
	}

	return result
}

// GetUser will return a single user by its user_id
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	// Access individual query parameters
	userID := queryParams.Get("user_id")
	if len(userID) == 0 {
		http.Error(w, "No user ID sent!", http.StatusBadRequest)
		return
	}
	user, err := h.doGetLogic(userID)
	if err != nil {
		msg := fmt.Sprintf("unable to get user %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "error with user", http.StatusBadRequest)
	}
}

func (h *Handler) doGetLogic(userID string) (*models.UserAll, error) {
	dbUser, err := h.d.GetUser(userID)
	if err != nil {
		return nil, err
	}
	dbPics, err := h.d.GetPicsByUserID(dbUser.UserID.String)
	if err != nil {
		return nil, err
	}

	var pics []string
	for _, p := range dbPics {
		pics = append(pics, p.URL.String)
	}

	dbMatches, err := h.d.GetMatches(dbUser.UserID.String)
	if err != nil {
		return nil, err
	}

	dbSchool, err := h.d.GetSchoolById(dbUser.SchoolID.Int16)
	if err != nil {
		return nil, err
	}

	var matches []string
	for _, p := range dbMatches {
		matches = append(matches, p.Matches)
	}
	var matchesWNames []models.MatchesWName

	var dbMatchesWName []datastore.User
	if len(matches) > 0 {
		dbMatchesWName, err = h.d.GetMatchesUserName(matches)
		if err != nil {
			return nil, err
		}
		for _, p := range dbMatchesWName {
			matchesWNames = append(matchesWNames, models.MatchesWName{
				UserID:    p.UserID.String,
				FirstName: p.FirstName.String,
				LastName:  p.LastName.String,
			})
		}
	}

	//convert from db user to json user
	user := models.UserAll{
		User: models.User{
			UserID:    dbUser.UserID.String,
			FirstName: dbUser.FirstName.String,
			LastName:  dbUser.LastName.String,
			Email:     dbUser.Email,
			SchoolID:  int(dbUser.SchoolID.Int16),
			DobDay:    int(dbUser.DobDay.Int16),
			DobMonth:  int(dbUser.DobMonth.Int16),
			DobYear:   int(dbUser.DobYear.Int16),
			IsMale:    dbUser.IsMale.Bool,
			GradYear:  int(dbUser.GradYear.Int16),
			Verified:  dbUser.Verified.Bool,
			About:     dbUser.About.String,
			CreatedAt: dbUser.CreatedAt,
		},
		Pics:    pics,
		Matches: matchesWNames,
		School:  dbSchool.School,
	}
	return &user, nil
}

// GetUser will return a single user by its user_id
func (h *Handler) GetGenderedUser(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	isMale := queryParams.Get("is_male") == "true"
	userID := queryParams.Get("user_id")
	dbGenderInterestUsers, err := h.d.GetGenderedUser(isMale)
	if err != nil {
		msg := fmt.Sprintf("unable to get user %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var res []models.UserAll
	for _, u := range dbGenderInterestUsers {
		dbPics, err := h.d.GetPicsByUserID(u.UserID.String)
		if err != nil {
			msg := fmt.Sprintf("unable to get user pics %v", err)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		var pics []string
		for _, p := range dbPics {
			pics = append(pics, p.URL.String)
		}

		dbSchool, err := h.d.GetSchoolById(u.SchoolID.Int16)
		if err != nil {
			msg := fmt.Sprintf("unable to get user school %v", err)
			http.Error(w, msg, http.StatusBadRequest)
		}

		res = append(res, models.UserAll{
			User: models.User{
				UserID:    u.UserID.String,
				FirstName: u.FirstName.String,
				LastName:  u.LastName.String,
				Email:     u.Email,
				SchoolID:  int(u.SchoolID.Int16),
				DobDay:    int(u.DobDay.Int16),
				DobMonth:  int(u.DobMonth.Int16),
				DobYear:   int(u.DobYear.Int16),
				IsMale:    u.IsMale.Bool,
				GradYear:  int(u.GradYear.Int16),
				Verified:  u.Verified.Bool,
				About:     u.About.String,
				CreatedAt: u.CreatedAt,
			},
			Pics:   pics,
			School: dbSchool.School,
		})
	}

	for i := range res {
		j := rand.Intn(i + 1)
		res[i], res[j] = res[j], res[i]
	}

	userMatches, err := h.d.GetAllMatches(userID)
	if err != nil {
		msg := fmt.Sprintf("unable to get matches %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// filter list for previous seen matches
	result := filterList(res, userMatches)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "error with user", http.StatusBadRequest)
	}
}

// GetUser will return a single user by its id
func (h *Handler) CompleteUser(w http.ResponseWriter, r *http.Request) {
	var u *models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload, %v", err), http.StatusBadRequest)
		return
	}

	u, err := h.d.CompleteUser(u)
	if err != nil {
		msg := fmt.Sprintf("unable to get user %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	dbUser, err := h.d.GetUser(u.UserID)
	if err != nil {
		msg := fmt.Sprintf("unable to get user %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	user := models.User{
		UserID:    dbUser.UserID.String,
		FirstName: dbUser.FirstName.String,
		LastName:  dbUser.LastName.String,
		Email:     dbUser.Email,
		SchoolID:  int(dbUser.SchoolID.Int16),
		DobDay:    int(dbUser.DobDay.Int16),
		DobMonth:  int(dbUser.DobMonth.Int16),
		DobYear:   int(dbUser.DobYear.Int16),
		IsMale:    dbUser.IsMale.Bool,
		GradYear:  int(dbUser.GradYear.Int16),
		Verified:  dbUser.Verified.Bool,
		About:     dbUser.About.String,
		CreatedAt: dbUser.CreatedAt,
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "error with user", http.StatusBadRequest)
	}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	var userIdsList []string
	userIdsParam := queryParams.Get("userIds")
	err := json.Unmarshal([]byte(userIdsParam), &userIdsList)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var res []*models.UserAll
	for _, userID := range userIdsList {
		u, err := h.doGetLogic(userID)
		if err != nil {
			msg := fmt.Sprintf("unable to get user %v", err)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		res = append(res, u)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error with users", http.StatusBadRequest)
	}

}

func (h *Handler) GetUnverifiedUsers(w http.ResponseWriter, r *http.Request) {

	userIdsList, err := h.d.GetUnverifiedUsers()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var res []*models.UserAll
	for _, u := range userIdsList {
		u, err := h.doGetLogic(u.UserID.String)
		if err != nil {
			msg := fmt.Sprintf("unable to get user %v", err)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		res = append(res, u)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "error with users", http.StatusBadRequest)
	}

}
