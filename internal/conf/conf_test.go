package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDevConfig(t *testing.T) {
	conf := NewConfig("dev")
	assert.NotEqual(t, "", conf.Host)
	assert.NotEqual(t, "", conf.Port)
	assert.NotEqual(t, "", conf.DbHost)
	assert.NotEqual(t, "", conf.DbPort)
	assert.NotEqual(t, "", conf.DbName)
	assert.NotEqual(t, "", conf.DbPassword)
	assert.Equal(t, "dev", conf.Env)
}

func TestNewProdConfig(t *testing.T) {
	conf := NewConfig("prod")
	assert.NotEqual(t, "", conf.Host)
	assert.NotEqual(t, "", conf.Port)
	assert.NotEqual(t, "", conf.DbHost)
	assert.NotEqual(t, "", conf.DbPort)
	assert.NotEqual(t, "", conf.DbName)
	assert.NotEqual(t, "", conf.DbPassword)
	assert.Equal(t, "prod", conf.Env)
}

func TestNewConfigHostNotSet(t *testing.T) {
	host, ok := os.LookupEnv(hostKey)
	err := os.Setenv(hostKey, "")
	defer os.Setenv(hostKey, host)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Panics(t, func() { NewConfig("dev") })
}

func TestNewConfigPortNotSet(t *testing.T) {
	port, ok := os.LookupEnv(portKey)
	err := os.Setenv(portKey, "")
	defer os.Setenv(portKey, port)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Panics(t, func() { NewConfig("dev") })
}

func TestNewConfigDbHostNotSet(t *testing.T) {
	dbHost, ok := os.LookupEnv(dbHostKey)
	err := os.Setenv(dbHostKey, "")
	defer os.Setenv(dbHostKey, dbHost)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Panics(t, func() { NewConfig("dev") })
}

func TestNewConfigDbPortNotSet(t *testing.T) {
	dbPort, ok := os.LookupEnv(dbPortKey)
	err := os.Setenv(dbPortKey, "")
	defer os.Setenv(dbPortKey, dbPort)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Panics(t, func() { NewConfig("dev") })
}

func TestNewConfigDbNameNotSet(t *testing.T) {
	dbName, ok := os.LookupEnv(dbNameKey)
	err := os.Setenv(dbNameKey, "")
	defer os.Setenv(dbNameKey, dbName)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Panics(t, func() { NewConfig("dev") })
}

func TestNewConfigDbPasswordNotSet(t *testing.T) {
	dbPassword, ok := os.LookupEnv(dbPasswordKey)
	err := os.Setenv(dbPasswordKey, "")
	defer os.Setenv(dbPasswordKey, dbPassword)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Panics(t, func() { NewConfig("dev") })
}

func TestNewTestConfig(t *testing.T) {
	conf := NewConfig("dev")
	testConf := NewTestConfig()
	assert.Equal(t, conf.Host, testConf.Host)
	assert.Equal(t, conf.Port, testConf.Port)
	assert.Equal(t, conf.DbHost, testConf.DbHost)
	assert.Equal(t, conf.DbPort, testConf.DbPort)
	assert.Equal(t, conf.DbName+"_test", testConf.DbName)
	assert.Equal(t, conf.DbPassword, testConf.DbPassword)
	assert.Equal(t, "dev", conf.Env)
}
