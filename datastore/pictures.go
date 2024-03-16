package datastore

import (
	"database/sql"
	"fmt"
	"log"
)

// get Pics by user id
func (d *Datastore) GetPicsByUserID(userID string) ([]Pictures, error) {
	var pics []Pictures
	sqlStatement := `SELECT * FROM pictures where user_id=$1`
	if err := d.db.Select(&pics, sqlStatement, userID); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return nil, err
		default:
			log.Printf("Unable to scan the row. %v", err)
			return nil, err
		}
	}
	return pics, nil
}
