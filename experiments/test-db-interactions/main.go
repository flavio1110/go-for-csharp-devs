package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	db, err := openDB("postgres://user:super-secret@localhost:5432/people?sslmode=disable")
	if err != nil {
		fmt.Printf("fail to open DB: %s\n", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	people, err := searchPeople(context.Background(), db, searchParams{})
	if err != nil {
		fmt.Printf("fail to search people: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Found %d people\n", len(people))
	for _, p := range people {
		fmt.Printf("Name: %s %s\tCity:%s\n", p.firstName, p.lastName, p.city)
	}
}

func openDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("unable to open DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to connect to DB: %w", err)
	}

	return db, nil
}

type person struct {
	firstName string
	lastName  string
	city      string
}

type searchParams struct {
	firstName *string
	lastName  *string
	city      *string
}

func searchPeople(ctx context.Context, db *sql.DB, params searchParams) ([]person, error) {
	query := `select first_name, last_name, city
		 from people
		 where ($1::text is null or first_name = $1)
		   and ($2::text is null or last_name = $2)
		   and ($3::text is null or city = $3)
		 order by first_name asc`

	rows, err := db.QueryContext(ctx, query, params.firstName, params.lastName, params.city)
	if err != nil {
		return nil, fmt.Errorf("fail to query people:%w", err)
	}

	var people []person
	var person person

	for rows.Next() {
		err := rows.Scan(&person.firstName, &person.lastName, &person.city)
		if err != nil {
			return nil, fmt.Errorf("fail to scan people results: %w", err)
		}
		people = append(people, person)
	}

	return people, nil
}
