package datastore

import (
	"database/sql"
	"fmt"
	"log"
)

// get messages
func (d *Datastore) GetMessages(user_id string, correspondingUserID string) ([]Messages, error) {
	var msgs []Messages
	sqlStatement := `SELECT * FROM messages WHERE to_user_id=$1 AND from_user_id=$2`
	if err := d.db.Select(&msgs, sqlStatement, user_id, correspondingUserID); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return nil, err
		default:
			log.Panicf("Unable to scan the row. %v", err)
			return nil, err
		}
	}
	return msgs, nil
}

// get messages
func (d *Datastore) GetLatestMessages(user_id string, correspondingUserID string) (*Messages, error) {
	var msgs Messages
	sqlStatement := `SELECT * FROM messages WHERE to_user_id=$1 AND from_user_id=$2 order by created_at desc LIMIT 1`
	if err := d.db.Get(&msgs, sqlStatement, correspondingUserID, user_id); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return nil, err
		default:
			log.Panicf("Unable to scan the row. %v", err)
			return nil, err
		}
	}
	return &msgs, nil
}

func (d *Datastore) AddMessage(fromUserID string, toUserID string, message string) error {
	sqlStatement := `INSERT INTO messages (from_user_id, to_user_id, txt_message) VALUES ($1, $2, $3)`

	_, err := d.db.Exec(sqlStatement, fromUserID, toUserID, message)
	if err != nil {
		log.Panicf("Unable to execute the query. %v", err)
		return err
	}
	return nil
}
