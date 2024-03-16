package datastore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres golang driver
)

type Datastore struct {
	db *sqlx.DB
}

func dsn() string {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	hostname := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	if len(user) == 0 || len(password) == 0 || len(hostname) == 0 || len(dbName) == 0 {
		log.Println("Missing env variables!")
		a := map[string]string{
			"user":     user,
			"password": password,
			"hostname": hostname,
			"dbName":   dbName,
			"port":     port,
		}
		bs, _ := json.Marshal(a)
		fmt.Println(string(bs))
		panic("stopping!")
	}
	_ = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, hostname, port, dbName)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostname, port, user, password, dbName,
	)
	return psqlInfo
}

// create connection with postgres db
func NewDB() (*Datastore, error) {
	conn := dsn()
	fmt.Printf("-----> \t%s \n", conn)

	// Retry for a certain duration until the database is ready
	retryInterval := 5 * time.Second
	maxRetries := 12 // Adjust based on your needs

	for i := 0; i < maxRetries; i++ {
		log.Println("Trying to connect to db")
		db, err := sql.Open("postgres", conn)
		fmt.Println(err)
		if err == nil {
			err = db.Ping()
		}
		_, _ = db.Exec(`set search_path=public`)
		if err == nil {
			log.Println("Connected to the database")
			return &Datastore{
				db: sqlx.NewDb(db, "postgres"),
			}, nil
		}
		log.Printf("Could not connect to db trying again. Retry %d", i)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to connect to the database after %d retries", maxRetries)
}
