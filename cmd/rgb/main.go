package main

import (
	"rgb/internal/cli"
	"rgb/internal/conf"
	"rgb/internal/server"
)

func main() {
	env := cli.Parse()
	server.Start(conf.NewConfig(env))
}
