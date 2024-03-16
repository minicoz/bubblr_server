package datastore

import (
	"database/sql"
	"fmt"
	"log"
)

// get all schools
func (d *Datastore) GetSchools() ([]School, error) {
	var schools []School
	sqlStatement := `SELECT * FROM schools`
	if err := d.db.Select(&schools, sqlStatement); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return nil, err
		default:
			log.Panicf("Unable to scan the row. %v", err)
			return nil, err
		}
	}
	return schools, nil
}

func (d *Datastore) GetSchoolById(id int16) (*School, error) {
	var school School
	sqlStatement := `SELECT * FROM schools where id = $1`
	if err := d.db.Get(&school, sqlStatement, id); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return nil, err
		default:
			log.Panicf("Unable to scan the row. %v", err)
			return nil, err
		}
	}
	return &school, nil
}
