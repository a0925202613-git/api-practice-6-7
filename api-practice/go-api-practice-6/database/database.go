package database

import (
	"database/sql"

	"go-api-practice-6/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const createTablesSQL = `
CREATE TABLE IF NOT EXISTS menus (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	category VARCHAR(100),
	price INTEGER NOT NULL CHECK (price >= 0),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
	menu_id INTEGER NOT NULL REFERENCES menus(id) ON DELETE RESTRICT,
	quantity INTEGER NOT NULL CHECK (quantity >= 1),
	status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'cancelled')),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
`

const seedMenusSQL = `
DO $$
BEGIN
  IF (SELECT COUNT(*) FROM menus) = 0 THEN
    INSERT INTO menus (name, category, price) VALUES
      ('招牌滷肉飯', '主食', 45),
      ('牛肉麵', '主食', 120),
      ('燙青菜', '小菜', 40),
      ('滷蛋', '小菜', 15);
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
	_, err = DB.Exec(seedMenusSQL)
	return err
}
