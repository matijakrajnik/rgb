package database

import (
	"rgb/internal/conf"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDBOptions(t *testing.T) {
	cfg := conf.NewConfig("dev")
	dbOptions := NewDBOptions(cfg)
	assert.Equal(t, cfg.DbHost+":"+cfg.DbPort, dbOptions.Addr)
	assert.Equal(t, cfg.DbName, dbOptions.Database)
	assert.Equal(t, cfg.DbUser, dbOptions.User)
	assert.Equal(t, cfg.DbPassword, dbOptions.Password)
}

func TestSetTestDBOptions(t *testing.T) {
	testCfg := conf.NewTestConfig()
	dbOptions := NewDBOptions(testCfg)
	assert.Equal(t, testCfg.DbHost+":"+testCfg.DbPort, dbOptions.Addr)
	assert.Equal(t, testCfg.DbName, dbOptions.Database)
	assert.Contains(t, dbOptions.Database, "_test")
	assert.Equal(t, testCfg.DbUser, dbOptions.User)
	assert.Equal(t, testCfg.DbPassword, dbOptions.Password)
}
