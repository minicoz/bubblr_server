package main

import (
	"bubblr/datastore"
	"bubblr/handler"
	"bubblr/router"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	db, err := datastore.NewDB()
	if err != nil {
		panic(err)
	}
	h := handler.NewHandler(db)
	r := router.Router(h)

	port := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	fmt.Printf("Starting server on the port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))

}
