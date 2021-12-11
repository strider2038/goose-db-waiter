package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4"
)

const (
	defaultMigrationsTable = "goose_db_version"
)

func main() {
	expectedVersion := getExpectedVersion()

	connection, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}
	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("failed to ping postgres:", err)
	}

	actualVersion, err := getActualVersion(connection)
	if err != nil {
		log.Fatal(err)
	}

	if actualVersion < int64(expectedVersion) {
		log.Fatalf("verification failed: actual migrations version %d is earlier than expected %d", actualVersion, expectedVersion)
	}

	log.Printf("verification succeeded: actual migrations version %d is later or equal than expected %d\n", actualVersion, expectedVersion)
}

func getExpectedVersion() int {
	versionString := os.Getenv("MIGRATION_VERSION")
	if versionString == "" {
		log.Fatal("env variable MIGRATION_VERSION is required")
	}
	version, err := strconv.Atoi(versionString)
	if err != nil {
		log.Fatal("MIGRATION_VERSION should be a valid integer:", err)
	}

	return version
}

func getActualVersion(connection *pgx.Conn) (int64, error) {
	tableName := os.Getenv("MIGRATIONS_TABLE")
	if tableName == "" {
		tableName = defaultMigrationsTable
	}

	row := connection.QueryRow(context.Background(), "SELECT version_id FROM "+tableName+" ORDER BY version_id DESC LIMIT 1")
	var actualVersion int64
	err := row.Scan(&actualVersion)
	if err != nil {
		return 0, fmt.Errorf("failed to select migrations version from %s: %s", tableName, err)
	}

	return actualVersion, nil
}
