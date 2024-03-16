package datastore

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
)

type DatastoreTest struct {
	db     *sql.DB
	dbName string
}

func SetUpTestDB(t *testing.T, dbName string) (*DatastoreTest, func()) {
	databaseName, err := createDatabaseName(dbName)
	if err != nil {
		return nil, nil
	}

	datastore, err := createOrReplaceDatabase(databaseName)
	if err != nil {
		return nil, nil
	}

	cleanupFunc := func() {
		cleanUpErr := datastore.CleanUp()
		assert.NoError(t, cleanUpErr, "Error cleaning up database.")
	}
	migrationPath := os.Getenv("TEST_BASE") + "/server/db/migrations"

	if err := datastore.Migrate(migrationPath); err != nil {
		t.Fatal(err)
	}
	return datastore, cleanupFunc
}

func (d *DatastoreTest) Migrate(path string) error {
	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		panic(err)
	}
	return nil
}

func createDatabaseName(name string) (string, error) {
	data := fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
	hasher := sha1.New()
	if _, err := hasher.Write([]byte(data)); err != nil {
		return "", fmt.Errorf("database test postgres database hasher: %w", err)
	}
	sha := hex.EncodeToString(hasher.Sum(nil))
	return fmt.Sprintf("%s_%s", name, sha[:8]), nil
}

func dsntest(dbName string) string {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	hostname := os.Getenv("POSTGRES_INSTANCE")
	if len(user) == 0 || len(password) == 0 || len(hostname) == 0 || len(dbName) == 0 {
		log.Println("Missing env variables!")
		a := map[string]string{
			"user":     user,
			"password": password,
			"hostname": hostname,
			"dbName":   dbName,
		}
		bs, _ := json.Marshal(a)
		fmt.Println(string(bs))
		panic("stopping!")
	}
	return fmt.Sprintf("postgresql://%s:%s@%s?sslmode=disable", user, password, hostname)
}

func createOrReplaceDatabase(dbName string) (*DatastoreTest, error) {
	conn := dsntest(dbName)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	fmt.Println(conn)

	var greeting string
	err = db.QueryRow("select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)

	// if _, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)); err != nil {
	// 	return nil, err
	// }

	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)); err != nil {
		log.Fatal(err)
	}

	fmt.Println(conn)

	return &DatastoreTest{
		db:     db,
		dbName: dbName,
	}, nil
}

func (d *DatastoreTest) CleanUp() error {
	if _, err := d.db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", d.dbName)); err != nil {
		return fmt.Errorf("database test mysql drop database %s: %w", d.dbName, err)
	}
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("database test close db: %w", err)
	}
	return nil
}
