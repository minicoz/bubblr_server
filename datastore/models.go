package datastore

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int64          `db:"id"`
	UserID         sql.NullString `db:"user_id"`
	FirstName      sql.NullString `db:"first_name"`
	LastName       sql.NullString `db:"last_name"`
	Email          string         `db:"email"`
	HashedPassword string         `db:"hashed_password"`
	SchoolID       sql.NullInt16  `db:"school_id"`
	DobDay         sql.NullInt16  `db:"dob_day"`
	DobMonth       sql.NullInt16  `db:"dob_month"`
	DobYear        sql.NullInt16  `db:"dob_year"`
	IsMale         sql.NullBool   `db:"is_male"`
	GradYear       sql.NullInt16  `db:"grad_year"`
	Verified       sql.NullBool   `db:"verified"`
	About          sql.NullString `db:"about"`
	CreatedAt      *time.Time     `db:"created_at"`
}

type School struct {
	ID     int    `db:"id"`
	School string `db:"school"`
	Tier   int    `db:"tier"`
}

type Pictures struct {
	ID     int            `db:"id"`
	UserID string         `db:"user_id"`
	URL    sql.NullString `db:"url"`
}

type Matches struct {
	ID      int64     `db:"id"`
	UserID  string    `db:"user_id"`
	Matches string    `db:"matched_user_id"`
	LikedAt time.Time `db:"liked_at"`
	Matched bool      `db:"matched"`
}

type Messages struct {
	ID         int    `db:"id"`
	ToUserID   string `db:"to_user_id"`
	FromUserID string `db:"from_user_id"`
	TxtMessage string `db:"txt_message"`
	CreatedAt  string `db:"created_at"`
}
