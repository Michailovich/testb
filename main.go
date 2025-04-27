package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"testb/gql"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	schemaFile, err := os.ReadFile("./gql/schema.graphql")
	if err != nil {
		log.Fatal(err)
	}

	schema, err := graphql.ParseSchema(string(schemaFile), &gql.Resolver{DB: db})
	if err != nil {
		log.Fatalf("Failed to parse schema: %v", err)
	}

	http.Handle("/graphql", &relay.Handler{Schema: schema})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
