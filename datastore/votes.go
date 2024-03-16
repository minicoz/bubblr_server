package datastore

import (
	"database/sql"
	"fmt"
)

func (d *Datastore) GetVoteCount(user_id string) (int, error) {
	var vc int
	sqlStatement := `
		select 
		SUM(CASE WHEN vote_yes THEN 1 ELSE 0 END) AS yes_count
		from votes
		where prospective_user_id=$1
	`
	err := d.db.Get(&vc, sqlStatement, user_id)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return 0, nil
	case err != nil:
		return 0, nil
	}
	return vc, err
}

func (d *Datastore) GetRejectionVoteCount(user_id string) (int, error) {
	var vc int
	sqlStatement := `
		select 
		SUM(CASE WHEN vote_yes THEN 0 ELSE 1 END) AS yes_count
		from votes
		where prospective_user_id=$1
	`
	err := d.db.Get(&vc, sqlStatement, user_id)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return 0, nil
	case err != nil:
		return 0, nil
	}
	return vc, err
}

func (d *Datastore) AddVote(user_id string, matched_id string) error {

	return nil
}

func (d *Datastore) RmVote(user_id string, matched_id string) error {
	return nil
}
