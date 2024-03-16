package datastore

import (
	"bubblr/models"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// insert one user in the DB
func (d *Datastore) InsertUser(user *models.User) (string, error) {
	sqlStatement := `INSERT INTO users (user_id, email, hashed_password) VALUES ($1, $2, $3) RETURNING user_id`

	var id string
	err := d.db.QueryRow(sqlStatement, user.UserID, user.Email, user.HashedPassword).Scan(&id)

	if err != nil {
		log.Panicf("Unable to execute the query. %v", err)
		return "", err
	}
	fmt.Printf("Inserted a single record %v", id)
	return id, nil
}

// get one user from the DB by its userid
func (d *Datastore) GetUser(user_id string) (*User, error) {
	var user User
	sqlStatement := `SELECT id,
			user_id, first_name, last_name, email, 
			hashed_password, school_id, dob_day, dob_month, 
			dob_year, is_male, grad_year, verified, about, created_at 
		FROM users WHERE user_id=$1`
	if err := d.db.Get(&user, sqlStatement, user_id); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return nil, err
		default:
			log.Panicf("Unable to scan the row. %v", err)
			return nil, err
		}
	}
	return &user, nil
}

// get all user from the DB by its gender_interest
func (d *Datastore) GetGenderedUser(isMale bool) ([]User, error) {
	var user []User
	sqlStatement := `SELECT id,
			user_id, first_name, last_name, email, 
			hashed_password, school_id, dob_day, dob_month, 
			dob_year, is_male, grad_year, verified, about, created_at 
		FROM users WHERE is_male != $1 and verified is true`
	if err := d.db.Select(&user, sqlStatement, isMale); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return nil, err
		default:
			log.Printf("Unable to scan the row. %v", err)
			return nil, err
		}
	}
	return user, nil
}

// get one user credentials from the DB by its email
func (d *Datastore) GetUserByEmail(email string) (User, error) {
	var u User
	sqlStatement := `SELECT * FROM users WHERE email=$1`
	err := d.db.Get(&u, sqlStatement, email)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return u, nil
	case nil:
		return u, nil
	default:
		log.Panicf("Unable to scan the row. %v", err)
	}
	return u, err
}

// insert one user in the DB
func (d *Datastore) CompleteUser(u *models.User) (*models.User, error) {
	sqlStatement := `
			UPDATE users 
			SET
				first_name = $1,
				last_name = $2,
				school_id = $3,
				dob_day = $4,
				dob_month = $5,
				dob_year = $6,
				is_male = $7,
				grad_year = $10,
				verified = false,
				about = $11
			WHERE 
				user_id = $12;
		`
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(sqlStatement)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	row, err := stmt.Exec(u.FirstName, u.LastName, u.SchoolID, u.DobDay, u.DobMonth, u.DobYear, u.IsMale, u.GradYear, u.About, u.UserID)
	if err != nil {
		log.Panicf("Unable to execute the query. %v", err)
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		stmt.Close()
		_ = tx.Rollback()
		return nil, err
	}
	stmt.Close()
	if r, err := row.RowsAffected(); r > 0 && err != nil {
		log.Panicf("No rows were entered. %v", err)
		return nil, err
	}
	return u, err
}

func (d *Datastore) GetMatchesUserName(user_ids []string) ([]User, error) {
	var users []User

	// Convert user_ids to []interface{}
	interfaceUserIDs := make([]interface{}, len(user_ids))
	for i, v := range user_ids {
		interfaceUserIDs[i] = v
	}

	// Use the pgx driver placeholder syntax $1, $2, etc., to handle multiple parameters
	placeholders := make([]string, len(user_ids))
	for i := range user_ids {
		placeholders[i] = "$" + strconv.Itoa(i+1)
	}
	placeholderStr := strings.Join(placeholders, ",")

	sqlStatement := fmt.Sprintf(`SELECT user_id, first_name, last_name
		FROM users WHERE user_id in (%s)`, placeholderStr)

	// Use the ... syntax to expand the slice into individual arguments
	if err := d.db.Select(&users, sqlStatement, interfaceUserIDs...); err != nil {
		switch err {
		default:
			log.Panicf("Unable to scan the row. %v", err)
			return nil, err
		}
	}

	return users, nil
}

func (d *Datastore) VerifyUser(user_id string) error {
	sqlStatement := `
		update users 
		set verified=true
		where user_id=$1
		`
	if _, err := d.db.Exec(sqlStatement, user_id); err != nil {
		return err
	}
	return nil
}

// get one user from the DB by its userid
func (d *Datastore) GetUnverifiedUsers() ([]User, error) {
	var user []User
	sqlStatement := `SELECT id,
			user_id, first_name, last_name, email, 
			hashed_password, school_id, dob_day, dob_month, 
			dob_year, is_male, grad_year, verified, about, created_at 
		FROM users WHERE verified is false`
	if err := d.db.Select(&user, sqlStatement); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return nil, err
		default:
			log.Panicf("Unable to scan the row. %v", err)
			return nil, err
		}
	}
	return user, nil
}
