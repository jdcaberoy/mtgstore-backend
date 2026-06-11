package main

import (
	"fmt"
	"os"

	"github.com/gobuffalo/pop/v6"
)

var DB *pop.Connection

func Connect() error {
	cd := &pop.ConnectionDetails{
		Dialect:  "postgresql",
		Database: os.Getenv("database"),
		User:     os.Getenv("user"),
		Host:     os.Getenv("host"),
		Password: os.Getenv("password"),
	}
	cd.Finalize()
	db, err := pop.NewConnection(cd)
	if err != nil {
		return fmt.Errorf("database: could not create connection: %w", err)
	}

	if err = db.Open(); err != nil {
		return fmt.Errorf("database: could not open connection: %w", err)
	}

	DB = db
	return nil
}
