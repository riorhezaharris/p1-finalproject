package main

import (
	"log"

	"p1finalproject/cli"
	"p1finalproject/config"
	"p1finalproject/handler"
)

func main() {
	db, err := config.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	h := handler.NewHandler(db)
	c := cli.NewCLI(h)
	c.Init()
}
