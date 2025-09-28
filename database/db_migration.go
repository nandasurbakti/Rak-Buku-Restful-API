package database

import (
	"embed"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migration/*.sql
var dbMigrations embed.FS

var DB *sqlx.DB

func DBMigrate(dbParam *sqlx.DB) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root: "sql_migration",
	}

	n, err := migrate.Exec(dbParam.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal("migrate failed: ", err)
	}
	
	fmt.Println("Migration sucess, applied", n, "migrations!")
}