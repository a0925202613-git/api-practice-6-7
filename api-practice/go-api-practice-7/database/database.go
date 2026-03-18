package database

import (
	"database/sql"

	"go-api-practice-7/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const createTablesSQL = `
CREATE TABLE IF NOT EXISTS books (
	id SERIAL PRIMARY KEY,
	title VARCHAR(500) NOT NULL,
	isbn VARCHAR(20),
	available BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS borrowals (
	id SERIAL PRIMARY KEY,
	book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE RESTRICT,
	user_name VARCHAR(255) NOT NULL,
	borrowed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	returned_at TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
`

const seedBooksSQL = `
DO $$
BEGIN
  IF (SELECT COUNT(*) FROM books) = 0 THEN
    INSERT INTO books (title, isbn, available) VALUES
      ('Go 程式設計', '978-986-123456-0', true),
      ('PostgreSQL 實戰', '978-986-123456-1', true),
      ('API 設計入門', '978-986-123456-2', false);
  END IF;
END $$;
`

func Connect() error {
	var err error
	DB, err = sql.Open("postgres", config.DatabaseURL())
	if err != nil {
		return err
	}
	if err := DB.Ping(); err != nil {
		return err
	}
	if _, err = DB.Exec(createTablesSQL); err != nil {
		return err
	}
	_, err = DB.Exec(seedBooksSQL)
	return err
}
