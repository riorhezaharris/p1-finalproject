package main

import (
	"p1finalproject/cli"
	"p1finalproject/config"
	"p1finalproject/handler"
)

func main() {
	db, err := config.InitDb()
	if err != nil {
		panic(err)
	}
	handlerObject := handler.NewHandler(db)
	cli := cli.NewCli(*handlerObject)
	cli.Init()
}
