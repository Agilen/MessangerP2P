package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	crpt "github.com/Agilen/MessangerP2P/AES-128"
	bigintegers "github.com/Agilen/MessangerP2P/BigIntegers"
	dh "github.com/Agilen/MessangerP2P/DH"
	"github.com/urfave/cli"
)

type Message struct {
	Message string
	Salt    []byte
}
type DH struct {
	UrSecret     string
	SharedSecret string
	Params       DHParams
}
type DHParams struct {
	PublicSecret string
	Module       string
	G            string
}
type Peer struct {
	Name string
}

type Info struct {
	Name       string
	UrPort     int
	PortToCon  int
	connection net.Conn
}

var mes Message
var info Info
var peer Peer
var dhInfo DH

func init() {
	info = Info{
		Name:   "Alex",
		UrPort: 8000,
	}
}

func main() {

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
	fmt.Println(info.UrPort, info.PortToCon, info, info.Name)

	var wg sync.WaitGroup
	fmt.Println("Welcom")
	time.Sleep(1000)
	wg.Add(1)
	go Listen(&wg)
	if info.PortToCon != 0 {
		go SendRequest(&wg)
	}
	time.Sleep(1000)
	wg.Wait()
	go Write()
	Read()
}

func Listen(wg *sync.WaitGroup) {
	listener, _ := net.Listen("tcp", ":"+strconv.Itoa(info.UrPort))

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go onConnection(conn, wg)

	}
}
func SendRequest(wg *sync.WaitGroup) {
	connection, _ := net.Dial("tcp", ":"+strconv.Itoa(info.PortToCon))
	for connection == nil {
		connection, _ = net.Dial("tcp", ":"+strconv.Itoa(info.PortToCon))
	}
	info.connection = connection
	data := DataPreparation()

	connection.Write([]byte(data + "\n"))

	message, _ := bufio.NewReader(connection).ReadString('\n')
	M := strings.Fields(message)
	fmt.Println(M)
	response := []byte(M[0])
	peerData := []byte(M[1])
	err := json.Unmarshal(peerData, &peer)
	if err != nil {

		log.Fatal(err)
	}
	err = json.Unmarshal(response, &dhInfo.Params)
	if err != nil {

		log.Fatal(err)
	}

	dhInfo.SharedSecret = bigintegers.ToHex(bigintegers.LongModPowerBarrett(bigintegers.ReadHex(dhInfo.Params.PublicSecret), bigintegers.ReadHex(dhInfo.UrSecret), bigintegers.ReadHex(dhInfo.Params.Module)))

	fmt.Println(
		"\nSharedSecret:", dhInfo.SharedSecret,
	)
	wg.Done()
}
func onConnection(conn net.Conn, wg *sync.WaitGroup) { // обработка запроса
	fmt.Printf("New connection from: %v", conn.RemoteAddr().String())
	info.connection = conn
	message, _ := bufio.NewReader(conn).ReadString('\n') //жду запрос
	M := strings.Fields(message)
	data := []byte(M[0])
	peerData := []byte(M[1])
	err := json.Unmarshal(data, &dhInfo.Params)
	if err != nil {

		log.Fatal(err)
	}
	err = json.Unmarshal(peerData, &peer)
	if err != nil {

		log.Fatal(err)
	}
	dhInfo.UrSecret = bigintegers.ToHex(dh.GenRandomNum(2))
	dhInfo.SharedSecret = bigintegers.ToHex(bigintegers.LongModPowerBarrett(bigintegers.ReadHex(dhInfo.Params.PublicSecret), bigintegers.ReadHex(dhInfo.UrSecret), bigintegers.ReadHex(dhInfo.Params.Module)))
	dhInfo.Params.PublicSecret = bigintegers.ToHex(bigintegers.LongModPowerBarrett(bigintegers.ReadHex(dhInfo.Params.G), bigintegers.ReadHex(dhInfo.UrSecret), bigintegers.ReadHex(dhInfo.Params.Module)))

	r := DHParams{dhInfo.Params.PublicSecret, dhInfo.Params.Module, dhInfo.Params.G}
	json_data, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	pr := Peer{info.Name}
	json_peerData, err := json.Marshal(pr)
	if err != nil {
		log.Fatal(err)
	}
	response := string(json_data) + " " + string(json_peerData)

	conn.Write([]byte(response + "\n"))

	fmt.Println(
		"\nSharedSecret:", dhInfo.SharedSecret,
	)
	wg.Done()
}
func DataPreparation() string { // Подготовка данных для запроса на подключение
	dhInfo.UrSecret = bigintegers.ToHex(dh.GenRandomNum(2))
	dhInfo.Params.Module = bigintegers.ToHex(dh.GenRandomNum(2))
	dhInfo.Params.G = bigintegers.ToHex(dh.GenRandomNum(2))
	dhInfo.Params.PublicSecret = bigintegers.ToHex(bigintegers.LongModPowerBarrett(bigintegers.ReadHex(dhInfo.Params.G), bigintegers.ReadHex(dhInfo.UrSecret), bigintegers.ReadHex(dhInfo.Params.Module)))

	// data := "Hello " + info.Name + " " + strconv.Itoa(info.UrPort) + " " + dhInfo.Params.PublicSecret + " " + dhInfo.Params.G + " " + dhInfo.Params.Module
	dh := DHParams{dhInfo.Params.PublicSecret, dhInfo.Params.Module, dhInfo.Params.G}
	json_data, err := json.Marshal(dh)
	if err != nil {
		log.Fatal(err)
	}
	pr := Peer{info.Name}
	json_peerData, err := json.Marshal(pr)
	if err != nil {
		log.Fatal(err)
	}
	data := string(json_data) + " " + string(json_peerData)
	return data
}
func Write() {
	for {
		if info.connection != nil {

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			key, salt := crpt.DeriveKey(dhInfo.SharedSecret, nil)
			textTosend, _ := crpt.Encrypt(key, text)
			fmt.Println("enc textToSend:", textTosend)
			send := Message{textTosend, (salt)}
			json_data, err := json.Marshal(send)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(key)
			fmt.Println(salt)
			info.connection.Write([]byte(string(json_data) + "\n"))
			fmt.Println("i send")
		}
	}
}
func Read() {
	for {
		if info.connection != nil {

			message, _ := bufio.NewReader(info.connection).ReadString('\n')
			err := json.Unmarshal([]byte(message), &mes)
			if err != nil {

				log.Fatal(err)
			}
			fmt.Println("i got it")

			key, _ := crpt.DeriveKey(dhInfo.SharedSecret, mes.Salt)
			fmt.Println("message to ddec", message)
			messageToRead, _ := crpt.Decrypt(key, mes.Message)
			fmt.Println("decr mes", messageToRead)
			fmt.Println(key)
			fmt.Println(mes.Salt)
			fmt.Print(peer.Name + ":" + messageToRead)
		}
	}
}
