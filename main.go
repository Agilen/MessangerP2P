package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"

	bigintegers "github.com/Agilen/MessangerP2P/BigIntegers"
	dh "github.com/Agilen/MessangerP2P/DH"
)

type Info struct {
	Name         string
	UrPort       int
	PortToCon    int
	G            []uint64
	UrSecret     []uint64
	PublicSecret []uint64
	Module       []uint64
	SharedSecret []uint64
	connection   net.Conn
}

var info Info

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go Listen(&wg)
	go SendRequest(&wg)
	wg.Wait()
	go Write()
	Read()

}

func Listen(wg *sync.WaitGroup) {
	listener, _ := net.Listen("tcp", ":8000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go onConnection(conn)
		wg.Done()
	}
}

func SendRequest(wg *sync.WaitGroup) {
	connection, _ := net.Dial("tcp", ":8001")
	for connection == nil {
		connection, _ = net.Dial("tcp", ":8001")
	}
	info.connection = connection
	data := DataPreparation()

	connection.Write([]byte(data + "\n"))
	fmt.Println("запрос отправлен")

	message, _ := bufio.NewReader(connection).ReadString('\n')
	fmt.Println("получен ответ")
	buffer := strings.Fields(message)
	h := bigintegers.ReadHex(buffer[1])
	info.SharedSecret = bigintegers.LongModPowerBarrett(h, info.UrSecret, info.Module)

	fmt.Println(
		"\nUrSecret: ", bigintegers.ToHex(info.UrSecret),
		"\nG:", bigintegers.ToHex(info.G),
		"\nModule:", bigintegers.ToHex(info.Module),
		"\nPublicSecret:", bigintegers.ToHex(info.PublicSecret),
		"\nSharedSecret:", bigintegers.ToHex(info.SharedSecret),
	)
	wg.Done()
}
func onConnection(conn net.Conn) { // обработка запроса
	fmt.Printf("New connection from: %v", conn.RemoteAddr().String())
	info.connection = conn
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Получен запрос")
	buffer := strings.Fields(message) //len 6
	info.G = bigintegers.ReadHex(buffer[4])
	info.Module = bigintegers.ReadHex(buffer[5])
	H := bigintegers.ReadHex(buffer[3])
	info.UrSecret = dh.GenRandomNum(2)
	info.PublicSecret = bigintegers.LongModPowerBarrett(info.G, info.UrSecret, info.Module)
	info.SharedSecret = bigintegers.LongModPowerBarrett(H, info.UrSecret, info.Module)

	response := info.Name + " " + bigintegers.ToHex(info.PublicSecret)

	conn.Write([]byte(response + "\n"))
	fmt.Println("Отправлен ответ")

	fmt.Println(
		"\nUrSecret: ", bigintegers.ToHex(info.UrSecret),
		"\nG:", bigintegers.ToHex(info.G),
		"\nModule:", bigintegers.ToHex(info.Module),
		"\nPublicSecret:", bigintegers.ToHex(info.PublicSecret),
		"\nSharedSecret:", bigintegers.ToHex(info.SharedSecret),
	)
}
func DataPreparation() string { // Подготовка данных для запроса на подключение
	info.UrSecret = dh.GenRandomNum(2)
	info.Module = dh.GenRandomNum(2)
	info.G = dh.GenRandomNum(2)
	info.PublicSecret = bigintegers.LongModPowerBarrett(info.G, info.UrSecret, info.Module)

	data := "Hello " + info.Name + " " + strconv.Itoa(info.UrPort) + " " + bigintegers.ToHex(info.PublicSecret) + " " + bigintegers.ToHex(info.G) + " " + bigintegers.ToHex(info.Module)

	return data
}
func Write() {
	for {
		if info.connection != nil {
			fmt.Println("Write")
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			info.connection.Write([]byte(text + "\n"))
		}
	}
}
func Read() {
	for {
		if info.connection != nil {
			fmt.Println("Read")
			message, _ := bufio.NewReader(info.connection).ReadString('\n')
			fmt.Print("Message from 8001: " + message)
		}
	}
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
