package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	_ "modernc.org/sqlite"
)

var schema = `
CREATE TABLE IF NOT EXISTS job (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  method TEXT NOT NULL,
  headers TEXT NOT NULL,
  body TEXT,
  query TEXT,
  created_at INTEGER NOT NULL,
  scheduled_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  FOREIGN KEY (job_id) REFERENCES job(id)
);

CREATE TABLE IF NOT EXISTS job_run (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  job_id INTEGER NOT NULL,
  status text NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  FOREIGN KEY (job_id) REFERENCES job(id)
);
`

func GetDatabase(name string) *sqlx.DB {

	dbName := fmt.Sprintf("file:./data/%s.db", name)
	db, err := sqlx.Open("sqlite", dbName)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to open database")
		return nil
	}

	if _, err = db.Exec(schema); err != nil {
		return nil
	}
	return db

}
