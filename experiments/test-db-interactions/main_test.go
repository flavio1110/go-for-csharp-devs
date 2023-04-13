package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSearchPeople(t *testing.T) {
	ctx := context.Background()
	connString, terminate, err := startTestDB(ctx)
	if err != nil {
		t.Error(err)
	}
	defer terminate(t)

	db, err := openDB(connString)
	if err != nil {
		t.Fatal("fail to open DB", err)
	}

	if err := migrateDB(ctx, db); err != nil {
		t.Fatal("fail to migrate DB", err)
	}

	insert := `insert into people values
			('Flavio', 'Silva', 'Olbia'),
			('Joost', 'Van Huis', 'Amsterdam'),
			('Aldben', 'Arimeritin', 'Istambul'),
			('Nando', 'Pelect', 'Perth');`

	if _, err := db.ExecContext(ctx, insert); err != nil {
		t.Fatal("fail to insert initial data", err)
	}

	t.Run("without filters", func(t *testing.T) {
		people, err := searchPeople(ctx, db, searchParams{})
		assert.NoError(t, err)
		assert.Equal(t, 4, len(people))
	})

	t.Run("filtering by first name", func(t *testing.T) {
		people, err := searchPeople(ctx, db, searchParams{firstName: strPtr("Joost")})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(people))
		assert.Equal(t, "Joost", people[0].firstName)
	})
	t.Run("filtering by last name", func(t *testing.T) {
		people, err := searchPeople(ctx, db, searchParams{lastName: strPtr("Arimeritin")})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(people))
		assert.Equal(t, "Arimeritin", people[0].lastName)
	})
	t.Run("filtering by city", func(t *testing.T) {
		people, err := searchPeople(ctx, db, searchParams{city: strPtr("Olbia")})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(people))
		assert.Equal(t, "Olbia", people[0].city)
	})
}

func startTestDB(ctx context.Context) (string, func(t *testing.T), error) {
	var envVars = map[string]string{
		"POSTGRES_USER":     "user",
		"POSTGRES_PASSWORD": "super-secret",
		"POSTGRES_DB":       "people",
		"PORT":              "5432/tcp",
	}

	getConnString := func(host string, port nat.Port) string {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			envVars["POSTGRES_USER"],
			envVars["POSTGRES_PASSWORD"],
			host,
			port.Port(),
			envVars["POSTGRES_DB"])
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:14",
		ExposedPorts: []string{envVars["PORT"]},
		Env:          envVars,
		WaitingFor:   wait.ForSQL(nat.Port(envVars["PORT"]), "pgx", getConnString).WithStartupTimeout(time.Second * 15),
	}
	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to start db container :%w", err)
	}
	port, err := pgC.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return "", nil, fmt.Errorf("failed to get mapped port :%w", err)
	}
	host, err := pgC.Host(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get host :%w", err)
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		envVars["POSTGRES_USER"],
		envVars["POSTGRES_PASSWORD"],
		host,
		port.Int(),
		envVars["POSTGRES_DB"])

	terminate := func(t *testing.T) {
		if err := pgC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}
	return connString, terminate, nil
}

func migrateDB(ctx context.Context, db *sql.DB) error {
	// IRL we would use something like Go Migrate (https://github.com/golang-migrate/migrate)
	// to maintain and apply all of your migrations in your real and test DBs.
	_, err := db.ExecContext(ctx, "create table if not exists people (first_name text,last_name text,city text)")
	if err != nil {
		return fmt.Errorf("failed to migrate DB: %w", err)
	}
	return nil
}

func strPtr(v string) *string {
	return &v
}
