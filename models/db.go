package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

var schema = `
CREATE TABLE digests (
    id INTEGER PRIMARY KEY,
    feed_name text NOT NULL,
    title text NOT NULL,
	content text,
	created_at datetime
);
`

// NewDB creates a new DB
func NewDB(filename string) (DB, error) {
	database := DB{}
	var err error
	log.Printf("Connecting to database %q", filename)
	if _, err := os.Stat(filename); err == nil {
		database.db, err = sqlx.Connect("sqlite3", filename)
		return database, err
	}

	log.Println("Database does not exist, initializing schema")
	database.db, err = sqlx.Connect("sqlite3", filename)
	if err != nil {
		return database, err
	}

	return database, database.InitializeSchema()
}

// DB holds the database. Don't drop it
type DB struct {
	db *sqlx.DB
}

func (d *DB) InitializeSchema() error {
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	_, err := d.db.Exec(schema)
	if err != nil {
		fmt.Println("Error with schema:", err)
	}

	return err
}

func (d *DB) InsertDigest(digest Digest) error {
	_, err := d.db.NamedExec(`INSERT INTO digests (feed_name, title, content, created_at) VALUES (:feed_name, :title, :content, :created_at)`,
		digest)

	return err
}

func (d *DB) GetDigests() ([]Digest, error) {
	digests := []Digest{}
	err := d.db.Select(&digests, "SELECT * FROM digests ORDER BY id DESC")
	return digests, err
}

func (d *DB) GetDigestsByFeed(name string) ([]Digest, error) {
	digests := []Digest{}
	err := d.db.Select(&digests, "SELECT * FROM digests WHERE feed_name=$1 ORDER BY datetime(created_at) DESC", name)
	return digests, err
}

func (d *DB) GetDigestByID(id string) (Digest, error) {
	digest := Digest{}
	err := d.db.Get(&digest, "SELECT * FROM digests WHERE id=$1", id)
	return digest, err
}

func (d *DB) CountDigestsByFeed(name string) (int, error) {
	var count int
	err := d.db.Get(&count, "SELECT count(*) FROM digests WHERE feed_name=$1", name)
	return count, err
}
