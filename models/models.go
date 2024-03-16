package models

import (
	"time"
)

type User struct {
	UserID         string     `json:"user_id"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	Email          string     `json:"email"`
	HashedPassword string     `json:"hashed_password"`
	SchoolID       int        `json:"school_id"`
	DobDay         int        `json:"dob_day"`
	DobMonth       int        `json:"dob_month"`
	DobYear        int        `json:"dob_year"`
	IsMale         bool       `json:"is_male"`
	GradYear       int        `json:"grad_year"`
	Verified       bool       `json:"verified"`
	About          string     `json:"about"`
	CreatedAt      *time.Time `json:"created_at"`
}

type UserCredentials struct {
	UserID         string `json:"user_id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

type School struct {
	ID     int    `json:"id"`
	School string `json:"school"`
	Tier   int    `json:"tier"`
}

type Pictures struct {
	ID     int    `json:"id"`
	UserID string `json:"user_id"`
	URL    string `json:"url"`
}

type Match struct {
	ID            int64  `json:"id"`
	UserID        string `json:"user_id"`
	MatchedUserID string `json:"matchedUserId"`
}

type Vote struct {
	UserID        string `json:"user_id"`
	MatchedUserID string `json:"matchedUserId"`
}

type Matches struct {
	ID      int64       `json:"id"`
	UserID  string      `json:"user_id"`
	Matches []string    `json:"matches"`
	LikedAt []time.Time `json:"liked_ats"`
}

type UserAll struct {
	User
	Pics    []string       `json:"pics"`
	Matches []MatchesWName `json:"matches"`
	School  string         `json:"school"`
}

type MatchesWName struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserID    string `json:"user_id"`
}

type Messages struct {
	ToUserID   string `json:"to_user_id"`
	FromUserID string `json:"from_user_id"`
	TxtMessage string `json:"txt_message"`
	CreatedAt  string `json:"created_at"`
}
