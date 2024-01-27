package database

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Dot struct {
	X        int
	Y        int
	Filename string
	Content  []byte
}

func (d *Dot) GetContent() string {
	var sb strings.Builder
	for _, b := range d.Content {
		sb.WriteByte(b)
	}
	return sb.String()
}

type DBMethods interface {
	noEmptyError() (err error)
	createTable(query string) (err error)
	createTables() (err error)

	Init() (err error)
	StoreDotsToDB(d *Dot)
	GetDotsFromDB() Dot
}

type Database struct {
	con    *pgxpool.Pool
	IsInit bool
}

func (db *Database) noEmptyError(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "no rows in result set") {
		return nil
	}
	return err
}

func (db *Database) createTable(query string) (err error) {
	err = db.con.QueryRow(context.Background(), query).Scan()
	if db.noEmptyError(err) != nil {
		return err
	}
	return nil
}

func (db *Database) Init() (err error) {
	DB_URL := `postgres://postgres:postgres@localhost:5432/postgres`
	db.con, err = pgxpool.Connect(context.Background(), DB_URL)
	if err != nil {
		return err
	}
	err = db.createTables()
	if err != nil {
		return err
	}
	db.IsInit = true
	return
}

func (db *Database) createTables() (err error) {
	if err = db.createTable(`CREATE TABLE IF NOT EXISTS dots (id SERIAL PRIMARY KEY, filename VARCHAR(40), x INT NOT NULL, y INT NOT NULL, dot bytea NOT NULL);`); err != nil {
		return err
	}
	return nil
}

func (db *Database) StoreDotsToDB(d *Dot) (err error) {
	var id int
	fmt.Println(d.Content)
	hexContent := "\\x" + hex.EncodeToString(d.Content)
	err = db.con.QueryRow(context.Background(), "INSERT INTO dots (x, y, filename, dot) VALUES ($1, $2, $3, $4::bytea) RETURNING id", d.X, d.Y, d.Filename, hexContent).Scan(&id)
	return err
}
