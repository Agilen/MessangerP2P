package db

import (
	"database/sql"
	"fmt"
	"log"
)

type Feed struct {
	DB *sql.DB
	id int
}
type Text struct {
	chatHistory string
}

func NewFeed(db *sql.DB) *Feed {
	stmt, err := db.Prepare(`

		CREATE TABLE  IF NOT EXISTS "Client"(
		"id"	INTEGER NOT NULL UNIQUE,
		ip TEXT,
		port INTEGER UNIQUE,
		"nickname" TEXT,
		PRIMARY KEY("id" AUTOINCREMENT)
		);`)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()
	stmt, err = db.Prepare(`

		CREATE TABLE IF NOT EXISTS "ChatHistory"(
		"id" INTEGER,
		chatHistory TEXT,
		clientId INTEGER NOT NULL UNIQUE,
		FOREIGN KEY (clientId) REFERENCES Client(id)
		PRIMARY KEY("id" AUTOINCREMENT)
	 
	);`)
	if err != nil {
		log.Fatal(err)
	}

	stmt.Exec()
	return &Feed{
		DB: db,
	}
}
func (feed *Feed) AddNew(port int) {
	fmt.Println("hi1")
	stmt, err := feed.DB.Prepare(`INSERT INTO Client (ip,port,nickname) values(?,?,?)`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("hi2")
	stmt.Exec(port)

}
func (feed *Feed) AddNewAdr(ip string, port int, name string) {
	feed.DB.Exec("INSERT INTO Client (ip, port, nickname) values($1,$2,$3)", ip, port, name)

	ID := feed.DB.QueryRow("select id from Client where port=$1 ", port)
	ID.Scan(&feed.id)
	feed.DB.Exec("INSERT INTO ChatHistory (clientId) values($1)", feed.id)
}
func (feed *Feed) EditChatHistory(port int, chat string) {
	text := feed.GetHistory()
	text = text + chat
	feed.DB.Exec("update ChatHistory set chatHistory=$1 where clientId=$2 ", text, feed.id)
}
func (feed *Feed) GetHistory() string {
	var T Text

	Tex := feed.DB.QueryRow("select chatHistory from ChatHistory where clientId=$1 ", feed.id)
	Tex.Scan(&T.chatHistory)
	return T.chatHistory
}
