package sqliteImplementation

import (
	"context"
	"database/sql"
	"fmt"
	storage "tgBot/storage"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(path string) (storage.Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Can't find database: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't ping database: %s", err.Error())
	}

	return &Storage{
		db: db,
	}, nil

}

func (d *Storage) Init(ctx context.Context) error {
	query := `create table if not exists missions(id integer primary key, text TEXT, deadline text)`

	_, err := d.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Can't create database: %s", err.Error())
	}

	return nil
}

func (d *Storage) CreateMission(ctx context.Context, mission *storage.Mission) error {
	query := `insert into missions(text, deadline) values(?, ?)`

	_, err := d.db.ExecContext(ctx, query, mission.Text, mission.Deadline)
	if err != nil {
		return fmt.Errorf("Can't add row to a database: %s", err.Error())
	}

	return nil
}

func (d *Storage) ReadLatestMissions(ctx context.Context) ([]storage.Mission, error) {
	var buffer = make([]storage.Mission, 0)

	query :=
		`SELECT id, text, deadline
		FROM missions
		ORDER BY deadline asc
		LIMIT 5;
	`
	rows, err := d.db.QueryContext(ctx, query, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return []storage.Mission{}, fmt.Errorf("Can't get rows from a database: %s", err.Error())
	}

	if err == sql.ErrNoRows {
		return []storage.Mission{}, nil
	}

	for rows.Next() {
		var id int
		var text string
		var deadline string

		err = rows.Scan(&id, &text, &deadline)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		rightDate, err := time.Parse("2006-01-02 15:04:05", deadline)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		buffer = append(buffer, storage.Mission{Id: id, Text: text, Deadline: rightDate})
	}

	return buffer, nil
}

func (d *Storage) RemoveMission(ctx context.Context, id int) error {
	query := `delete from missions where id = ?`

	_, err := d.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("can't remove from database: %s", err.Error())
	}

	return nil
}
