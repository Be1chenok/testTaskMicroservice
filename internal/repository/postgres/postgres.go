package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Be1chenok/testTaskMicroservice/internal/config"
	_ "github.com/lib/pq"
)

func New(conf *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.DBName,
		conf.Database.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
