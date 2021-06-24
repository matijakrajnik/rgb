package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table users...")
		_, err := db.Exec(`CREATE TABLE users(
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			hashed_password BYTEA NOT NULL,
			salt BYTEA NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table users...")
		_, err := db.Exec(`DROP TABLE users`)
		return err
	})
}
