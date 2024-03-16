package datastore

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// add a match to db
func (d *Datastore) AddMatch(userID string, matched_user_id string, matched bool) (bool, error) {

	insertQuery := `
		INSERT INTO matches (user_id, matched_user_id, liked_at, matched)
		VALUES ($1, $2, $3, $4)
	`
	_, err := d.db.Exec(insertQuery, userID, matched_user_id, time.Now(), matched)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return false, err
	case err != nil:
		return false, err
	}

	return matched, nil
}

// rm a match to db
func (d *Datastore) RemoveMatch(userID string, matched_user_id string) error {

	insertQuery := `
		DELETE FROM matches 
		WHERE user_id = $1 and matched_user_id = $2
	`
	if _, err := d.db.Exec(insertQuery, userID, matched_user_id); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return err
		default:
			log.Printf("Unable to scan the row. %v", err)
			return err
		}
	}
	return nil
}

func (d *Datastore) CheckForPrevMatched(userAID string, userBID string) (bool, error) {
	query := `
		SELECT matched FROM matches
		WHERE user_id = $2 AND matched_user_id = $1
	`

	var matched bool
	err := d.db.Get(&matched, query, userAID, userBID)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return false, nil
	case err != nil:
		log.Printf("Unable to scan the row. %v", err)
		return false, err
	}

	if matched {
		// yes a prev row! match found
		update_query := `
		UPDATE matches SET matched=true
		WHERE user_id = $2 AND matched_user_id = $1
		`
		if _, err := d.db.Exec(update_query, userAID, userBID); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (d *Datastore) GetMatches(userID string) ([]Matches, error) {
	insertQuery := `
		SELECT * FROM matches WHERE user_id = $1 AND matched is true
	`
	var m []Matches
	if err := d.db.Select(&m, insertQuery, userID); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return nil, err
		default:
			log.Printf("Unable to scan the row. %v", err)
			return nil, err
		}
	}
	return m, nil
}

func (d *Datastore) GetAllMatches(userID string) ([]Matches, error) {
	insertQuery := `
		SELECT * FROM matches WHERE user_id = $1 
	`
	var m []Matches
	if err := d.db.Select(&m, insertQuery, userID); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return nil, err
		default:
			log.Printf("Unable to scan the row. %v", err)
			return nil, err
		}
	}
	return m, nil
}
