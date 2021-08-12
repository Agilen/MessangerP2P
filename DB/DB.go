package db

import (
	"database/sql"
	"log"
)

type Feed struct {
	DB *sql.DB
}
type T struct {
	chatHistory string
}

// db, _ := sql.Open("sqlite3", "./ChatHistory.db")
func NewFeed(db *sql.DB) *Feed {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS "ChatHistory" (
		"ID"	INTEGER NOT NULL UNIQUE,
		"port"	INTEGER NOT NULL UNIQUE,
		"chatHistory"	TEXT,
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

func (feed *Feed) AddNewPort(port int) {
	stmt, _ := feed.DB.Prepare(`INSERT INTO ChatHistory (port) values(?)`)
	stmt.Exec(port)
}
func (feed *Feed) EditChatHistory(port int, chat string) {
	text := feed.GetHistory(port)
	text = text + chat
	feed.DB.Exec("update ChatHistory set chatHistory=$1 where port=$2 ", text, port)
}
func (feed *Feed) GetHistory(port int) string {
	var T T
	Text := feed.DB.QueryRow("select chatHistory from ChatHistory where port=$2 ", port)
	Text.Scan(&T.chatHistory)
	return T.chatHistory
}
