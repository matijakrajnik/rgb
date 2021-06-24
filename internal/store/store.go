package store

import (
	"errors"
	"fmt"
	"regexp"
	"rgb/internal/conf"
	"rgb/internal/database"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
)

// Database connector
var db *pg.DB

func SetDBConnection(dbOpts *pg.Options) {
	if dbOpts == nil {
		log.Panic().Msg("DB options can't be nil")
	} else {
		db = pg.Connect(dbOpts)
	}
}

func GetDBConnection() *pg.DB { return db }

func ResetTestDatabase() {
	// Connect to test database
	SetDBConnection(database.NewDBOptions(conf.NewTestConfig()))

	// Empty all tables and restart sequence counters
	tables := []string{"users", "posts"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s;", table))
		if err != nil {
			log.Panic().Err(err).Str("table", table).Msg("Error clearing test database")
		}

		_, err = db.Exec(fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART;", table))
	}
}

func dbError(_err interface{}) error {
	if _err == nil {
		return nil
	}
	switch _err.(type) {
	case pg.Error:
		err := _err.(pg.Error)
		switch err.Field(82) {
		case "_bt_check_unique":
			return errors.New(extractColumnName(err.Field(110)) + " already exists.")
		}
	case error:
		err := _err.(error)
		switch err.Error() {
		case "pg: no rows in result set":
			return errors.New("Not found.")
		}
		return err
	}
	return errors.New(fmt.Sprint(_err))
}

func extractColumnName(text string) string {
	reg := regexp.MustCompile(`.+_(.+)_.+`)
	if reg.MatchString(text) {
		return strings.Title(reg.FindStringSubmatch(text)[1])
	}
	return "Unknown"
}
