package main

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"sync"

	client "github.com/Agilen/MessangerP2P/Client"
	DB "github.com/Agilen/MessangerP2P/DB"
	dh "github.com/Agilen/MessangerP2P/DH"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
)

func main() {
	db, err := sql.Open("sqlite3", "./ChatHistory.db")
	if err != nil {
		log.Fatal(err)
	}
	feed := *DB.NewFeed(db)
	mes := *client.NewMessanger()
	var dhInfo dh.DHContext
	var wg sync.WaitGroup
	app := cli.NewApp() // &cli.App{}
	app.Name = "&"
	app.Usage = "hmmmm"
	app.Description = "help urself"

	app.Flags = []cli.Flag{
		&cli.StringFlag{Destination: &mes.Name, Name: "name", Value: "Anon", Usage: "It is your nickname"},
		&cli.IntFlag{Destination: &mes.UrPort, Name: "port", Usage: "ur port"},
		&cli.IntFlag{Destination: &mes.PortToCon, Name: "conn", Usage: "port to con"},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		flag := mes.Commands(text)
		if flag {
			break
		}
	}

	mes.Start(&wg, &dhInfo, &feed)
}
