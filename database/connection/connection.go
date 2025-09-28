package connection

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var DB *sqlx.DB

func ConnectDatabase() {
	var dsn string
	if viper.IsSet("DATABASE_URL") {
		dsn = viper.GetString("DATABASE_URL")
	} else {
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			viper.GetString("PGHOST"),
			viper.GetInt("PGPORT"),
			viper.GetString("PGUSER"),
			viper.GetString("PGPASSWORD"),
			viper.GetString("PGDATABASE"),
			viper.GetString("PGSSLMODE"),
		)
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("gagal menghubungkan ke database: ", err)
	}

	DB = db
	log.Println("database terhubung!")
}
