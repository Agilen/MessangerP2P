package client

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

	crpt "github.com/Agilen/MessangerP2P/AES-128"
	bigintegers "github.com/Agilen/MessangerP2P/BigIntegers"
	DB "github.com/Agilen/MessangerP2P/DB"
	dh "github.com/Agilen/MessangerP2P/DH"
)

type Message struct {
	Message string
	Salt    []byte
}

type Peer struct {
	Name string
	port int
}

type Info struct {
	Name       string
	UrPort     int
	PortToCon  int
	peer       Peer
	message    Message
	connection net.Conn
}

func NewInfo() *Info {
	info := &Info{
		Name:   "Anon",
		UrPort: 8000,
	}
	return info
}

func (info *Info) Start(wg *sync.WaitGroup, dhInfo *dh.DHContext, feed *DB.Feed) {
	fmt.Println("Welcom")
	wg.Add(1)
	go info.Listen(wg, dhInfo, feed)
	if info.PortToCon != 0 {
		go info.SendRequest(wg, dhInfo, feed)
	}
	wg.Wait()
	go info.Write(dhInfo, feed)
	info.Read(dhInfo, feed)
}
func (info *Info) Listen(wg *sync.WaitGroup, dhInfo *dh.DHContext, feed *DB.Feed) {
	listener, _ := net.Listen("tcp", ":"+strconv.Itoa(info.UrPort))

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go info.OnConnection(conn, wg, dhInfo, feed, &info.peer)

	}
}
func (info *Info) SendRequest(wg *sync.WaitGroup, dhInfo *dh.DHContext, feed *DB.Feed) {

	connection, _ := net.Dial("tcp", ":"+strconv.Itoa(info.PortToCon))
	for connection == nil {
		connection, _ = net.Dial("tcp", ":"+strconv.Itoa(info.PortToCon))
	}
	feed.AddNewPort(info.PortToCon)
	info.connection = connection
	data := info.DataPreparation(dhInfo)
	fmt.Println("PB", dhInfo.DHParams.PublicKey)
	connection.Write([]byte(data + "\n"))

	message, _ := bufio.NewReader(connection).ReadString('\n')
	M := strings.Fields(message)
	fmt.Println("M", M)
	response := []byte(M[0] + " " + M[1] + " " + M[2])
	fmt.Println("asdf", M[0], M[1], M[2])
	peerData := []byte(M[3])
	err := json.Unmarshal(peerData, &info.peer)
	if err != nil {

		log.Fatal(err)
	}
	err = json.Unmarshal(response, &dhInfo.DHParams)
	if err != nil {

		log.Fatal(err)
	}
	fmt.Println("PB", bigintegers.ToHex(dhInfo.DHParams.PublicKey))
	dhInfo.CalculateSharedSecret()
	fmt.Println(
		"\nSharedSecret:", bigintegers.ToHex([]uint64(dhInfo.SharedSecret)),
	)
	wg.Done()
}
func (info *Info) OnConnection(conn net.Conn, wg *sync.WaitGroup, dhInfo *dh.DHContext, feed *DB.Feed, peer *Peer) { // обработка запроса
	fmt.Printf("New connection from: %v", conn.RemoteAddr().String())
	info.connection = conn

	message, _ := bufio.NewReader(conn).ReadString('\n') //жду запрос
	M := strings.Fields(message)
	data := []byte(M[0] + " " + M[1] + " " + M[2])
	peerData := []byte(M[3])
	err := json.Unmarshal(data, &dhInfo.DHParams)
	if err != nil {

		log.Fatal(err)
	}
	err = json.Unmarshal(peerData, &info.peer)
	if err != nil {

		log.Fatal(err)
	}
	feed.AddNewPort(peer.port)
	dhInfo.GenerateDHPrivateKey()
	dhInfo.CalculateSharedSecret()
	dhInfo.CalculateDHPublicKey()

	r := dh.Params{dhInfo.DHParams.G, dhInfo.DHParams.P, dhInfo.DHParams.PublicKey}
	json_data, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	pr := Peer{info.Name, info.UrPort}
	json_peerData, err := json.Marshal(pr)
	if err != nil {
		log.Fatal(err)
	}
	response := string(json_data) + " " + string(json_peerData)
	conn.Write([]byte(response + "\n"))

	fmt.Println(
		"\nSharedSecret:", bigintegers.ToHex([]uint64(dhInfo.SharedSecret)),
	)
	wg.Done()
}
func (info *Info) DataPreparation(dhInfo *dh.DHContext) string { // Подготовка данных для запроса на подключение
	*dhInfo = *dh.NewDHContext()
	dhInfo.GenerateDHPrivateKey()
	dhInfo.CalculateDHPublicKey()
	dh := dh.Params{dhInfo.DHParams.G, dhInfo.DHParams.P, dhInfo.DHParams.PublicKey}
	json_data, err := json.Marshal(dh)
	if err != nil {
		log.Fatal(err)
	}
	pr := Peer{info.Name, info.UrPort}
	json_peerData, err := json.Marshal(pr)
	if err != nil {
		log.Fatal(err)
	}
	data := string(json_data) + " " + string(json_peerData)
	return data
}

func (info *Info) Write(dhInfo *dh.DHContext, feed *DB.Feed) {
	for {
		if info.connection != nil {

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			key, salt := crpt.DeriveKey([]uint64(dhInfo.SharedSecret), nil)
			textTosend, _ := crpt.Encrypt(key, text)
			fmt.Println("enc textToSend:", textTosend)
			feed.EditChatHistory(info.PortToCon, info.Name+":"+textTosend+"\n")
			send := Message{textTosend, (salt)}
			json_data, err := json.Marshal(send)
			if err != nil {
				log.Fatal(err)
			}
			info.connection.Write([]byte(string(json_data) + "\n"))
		}
	}
}
func (info *Info) Read(dhInfo *dh.DHContext, feed *DB.Feed) {
	for {
		if info.connection != nil {

			message, _ := bufio.NewReader(info.connection).ReadString('\n')
			err := json.Unmarshal([]byte(message), &info.message)
			if err != nil {

				log.Fatal(err)
			}
			key, _ := crpt.DeriveKey([]uint64(dhInfo.SharedSecret), info.message.Salt)
			messageToRead, _ := crpt.Decrypt(key, info.message.Message)
			feed.EditChatHistory(info.PortToCon, info.peer.Name+":"+messageToRead)
			fmt.Print(info.peer.Name + ":" + messageToRead)
		}
	}
}
