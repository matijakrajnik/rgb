package store

import (
	"rgb/internal/conf"
	"rgb/internal/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDBConnection(t *testing.T) {
	dbOptions := database.NewDBOptions(conf.NewTestConfig())
	assert.NotPanics(t, func() { SetDBConnection(dbOptions) })
	assert.NotNil(t, db)
	assert.Equal(t, dbOptions.Addr, db.Options().Addr)
	assert.Equal(t, dbOptions.User, db.Options().User)
	assert.Equal(t, dbOptions.Password, db.Options().Password)
	assert.Equal(t, dbOptions.Database, db.Options().Database)
}

func TestSetDBConnectionNilOptions(t *testing.T) {
	assert.Panics(t, func() { SetDBConnection(nil) })
}

func TestGetDBConnection(t *testing.T) {
	dbOptions := database.NewDBOptions(conf.NewTestConfig())
	SetDBConnection(dbOptions)
	fetched := GetDBConnection()
	assert.NotNil(t, fetched)
	assert.Equal(t, dbOptions.Addr, fetched.Options().Addr)
	assert.Equal(t, dbOptions.User, fetched.Options().User)
	assert.Equal(t, dbOptions.Password, fetched.Options().Password)
	assert.Equal(t, dbOptions.Database, fetched.Options().Database)
}
