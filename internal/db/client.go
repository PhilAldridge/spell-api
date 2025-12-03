package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"entgo.io/ent/dialect"
	"github.com/PhilAldridge/spell-api/ent"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

// Load .env file (ignore error if already loaded elsewhere)
func NewClient() *ent.Client {
	_ = godotenv.Load()

    user := os.Getenv("MYSQL_USER")
    pass := os.Getenv("MYSQL_PASS")
    host := os.Getenv("MYSQL_HOST")
    port := os.Getenv("MYSQL_PORT")
    dbname := os.Getenv("MYSQL_DB")

    if user == "" || pass == "" || host == "" || dbname == "" {
        log.Fatal("missing database environment variables")
    }

    // Build MySQL DSN
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
        user,
        pass,
        host,
        port,
        dbname,
    )

    client, err := ent.Open(dialect.MySQL, dsn)
    if err != nil {
        log.Fatalf("failed opening mysql: %v", err)
    }

    // Auto-migrate schema
    if err := client.Schema.Create(context.Background()); err != nil {
        log.Fatalf("failed creating schema: %v", err)
    }

    return client
}