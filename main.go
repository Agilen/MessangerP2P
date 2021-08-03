package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/user"
	"strconv"
	"sync"

	bigintegers "github.com/Agilen/MessangerP2P/BigIntegers"
	dh "github.com/Agilen/MessangerP2P/DH"
)

type Info struct {
	Name         string
	UrPort       int
	PortToCon    int
	UrSecret     []uint64
	PublicSecret []uint64
	P            []uint64
	G            []uint64
	S            []uint64
}

var info Info

func main() {

	// var name string
	// var port int
	// var portTo int
	// // fmt.Print("Enter your nick:")
	// // fmt.Fscan(os.Stdin, &name)
	// // fmt.Print("Enter your port:")
	// // fmt.Fscan(os.Stdin, &port)
	// // fmt.Print("Enter port to con:")
	// // fmt.Fscan(os.Stdin, &portTo)
	// // info.Name = name
	// // info.UrPort = port
	// // info.PortToCon = portTo

	l := make(chan net.Conn)
	c := make(chan net.Conn)

	go listen(l)
	go connection(c)
	conn, con := <-l, <-c
	// SendInfo(con)
	var wg sync.WaitGroup
	fmt.Println("Ready to chat")
	SendInfo(con)
	for {
		wg.Add(1)
		go Write(con, &wg)
		go Read(conn, &wg)

		wg.Wait()
	}
}
func listen(c chan net.Conn) {
	// ln, _ := net.Listen("tcp", "192.168.0.106:8000")
	ln, _ := net.Listen("tcp", ":"+strconv.Itoa(info.UrPort))

	conn, _ := ln.Accept()
	for conn == nil {
		conn, _ = ln.Accept()
	}

	fmt.Println("Port has been open")
	c <- conn
}
func connection(c chan net.Conn) {
	// con, _ := net.Dial("tcp", "192.168.0.105:8000")
	// for con == nil {
	// 	con, _ = net.Dial("tcp", "192.168.0.105:8000")
	// }
	con, _ := net.Dial("tcp", ":"+strconv.Itoa(info.PortToCon))
	for con == nil {
		con, _ = net.Dial("tcp", ":"+strconv.Itoa(info.PortToCon))
	}
	fmt.Println("Conected")
	c <- con
}
func Write(con net.Conn, wg *sync.WaitGroup) {

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if text != "" {
		fmt.Fprintf(con, text+"\n")
		fmt.Print("You: ", text)
	}
	wg.Done()
}
func Read(conn net.Conn, wg *sync.WaitGroup) {
	message, _ := bufio.NewReader(conn).ReadString('\n')
	if message != "" {
		fmt.Print("Message from 8001: " + message)
	}

	wg.Done()
}
func SendInfo(con net.Conn) {
	info.UrSecret = dh.GenRandomNum(16)
	info.G = dh.GenRandomNum(16)
	info.P = dh.GenRandomNum(16)
	info.PublicSecret = bigintegers.LongModPowerBarrett(info.G, info.UrSecret, info.P)
	fmt.Fprintf(con, info.Name+" "+bigintegers.ToHex(info.G)+" "+bigintegers.ToHex(info.P)+" "+bigintegers.ToHex(info.PublicSecret))
	fmt.Fprintf(con, " ")
}
func init() {
	curUser, _ := user.Current()
	hostName, _ := os.Hostname()

	info = Info{
		Name:      curUser.Username + "@" + hostName,
		UrPort:    8000,
		PortToCon: 8001,
	}
}
