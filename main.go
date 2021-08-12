package main

import (
	"database/sql"
	"log"
	"os"
	"sync"

	client "github.com/Agilen/MessangerP2P/Client"
	DB "github.com/Agilen/MessangerP2P/DB"
	dh "github.com/Agilen/MessangerP2P/DH"
	"github.com/urfave/cli"
)

func main() {
	db, _ := sql.Open("sqlite3", "./ChatHistory.db")
	feed := *DB.NewFeed(db)
	info := *client.NewInfo()
	var dhInfo dh.DHContext
	var wg sync.WaitGroup
	app := cli.NewApp() // &cli.App{}
	app.Name = "&"
	app.Usage = "hmmmm"
	app.Description = "help urself"

	app.Flags = []cli.Flag{
		&cli.StringFlag{Destination: &info.Name, Name: "name", Value: "Anon", Usage: "It is your nickname"},
		&cli.IntFlag{Destination: &info.UrPort, Name: "port", Usage: "ur port"},
		&cli.IntFlag{Destination: &info.PortToCon, Name: "conn", Usage: "port to con"},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	info.Start(&wg, &dhInfo, &feed)
}
