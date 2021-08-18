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
	DB "github.com/Agilen/MessangerP2P/DB"
	dh "github.com/Agilen/MessangerP2P/DH"
)

type Message struct {
	Message string
	Salt    []byte
}

type Peer struct {
	Name string
	Port int
}

type Messanger struct {
	Name       string
	UrPort     int
	PortToCon  int
	peer       Peer
	message    Message
	connection net.Conn
}

func NewMessanger() *Messanger {
	mes := &Messanger{
		Name:   "Anon",
		UrPort: 8000,
	}
	return mes
}

// Start messanger
func (mes *Messanger) Start(wg *sync.WaitGroup, dhInfo *dh.DHContext, feed *DB.Feed) {
	fmt.Println("Welcom")
	wg.Add(1)
	go mes.Listen(wg, dhInfo, feed)
	if mes.PortToCon != 0 {
		go mes.SendRequest(wg, dhInfo, feed)
	}
	wg.Wait()
	go mes.Write(dhInfo, feed)
	mes.Read(dhInfo, feed)

}

//Listen for new connection
func (mes *Messanger) Listen(wg *sync.WaitGroup, dhInfo *dh.DHContext, feed *DB.Feed) {
	listener, _ := net.Listen("tcp", ":"+strconv.Itoa(mes.UrPort))

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go mes.OnConnection(conn, wg, dhInfo, feed)

	}
}

//Send request to make connection and send info to init DH(Diffie Hellman)
func (mes *Messanger) SendRequest(wg *sync.WaitGroup, dhInfo *dh.DHContext, feed *DB.Feed) {

	connection, _ := net.Dial("tcp", ":"+strconv.Itoa(mes.PortToCon))
	for connection == nil {
		connection, _ = net.Dial("tcp", ":"+strconv.Itoa(mes.PortToCon))
	}
	feed.AddNewAdr("127.0.0.1", mes.peer.Port, mes.peer.Name)
	mes.connection = connection
	fmt.Printf("New connection from: %v", connection.RemoteAddr().String())
	data := mes.DataPreparation(dhInfo)
	connection.Write([]byte(data + "\n"))

	message, _ := bufio.NewReader(connection).ReadString('\n')
	M := strings.Fields(message)

	response := []byte(M[0] + " " + M[1] + " " + M[2])

	peerData := []byte(M[3])
	err := json.Unmarshal(peerData, &mes.peer)
	if err != nil {

		log.Fatal(err)
	}
	err = json.Unmarshal(response, &dhInfo.DHParams)
	if err != nil {

		log.Fatal(err)
	}

	dhInfo.CalculateSharedSecret()

	wg.Done()
}

//Processes the connection request
func (mes *Messanger) OnConnection(conn net.Conn, wg *sync.WaitGroup, dhInfo *dh.DHContext, feed *DB.Feed) { // обработка запроса
	fmt.Printf("New connection from: %v", conn.RemoteAddr().String())
	mes.connection = conn

	message, _ := bufio.NewReader(conn).ReadString('\n') //жду запрос
	M := strings.Fields(message)
	data := []byte(M[0] + " " + M[1] + " " + M[2])

	peerData := []byte(M[3])
	err := json.Unmarshal(data, &dhInfo.DHParams)
	if err != nil {

		log.Fatal(err)
	}
	err = json.Unmarshal(peerData, &mes.peer)
	if err != nil {

		log.Fatal(err)
	}

	feed.AddNewAdr("127.0.0.1", mes.peer.Port, mes.peer.Name)
	dhInfo.GenerateDHPrivateKey()
	dhInfo.CalculateSharedSecret()
	dhInfo.CalculateDHPublicKey()

	r := dh.Params{dhInfo.DHParams.G, dhInfo.DHParams.P, dhInfo.DHParams.PublicKey}
	json_data, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	pr := Peer{mes.Name, mes.UrPort}
	json_peerData, err := json.Marshal(pr)
	if err != nil {
		log.Fatal(err)
	}
	response := string(json_data) + " " + string(json_peerData)
	conn.Write([]byte(response + "\n"))

	wg.Done()
}

//Prepares data for request
func (mes *Messanger) DataPreparation(dhInfo *dh.DHContext) string { // Подготовка данных для запроса на подключение
	*dhInfo = *dh.NewDHContext()
	dhInfo.GenerateDHPrivateKey()
	dhInfo.CalculateDHPublicKey()
	dh := dh.Params{dhInfo.DHParams.G, dhInfo.DHParams.P, dhInfo.DHParams.PublicKey}
	json_data, err := json.Marshal(dh)
	if err != nil {
		log.Fatal(err)
	}
	pr := Peer{mes.Name, mes.UrPort}
	json_peerData, err := json.Marshal(pr)
	if err != nil {
		log.Fatal(err)
	}
	data := string(json_data) + " " + string(json_peerData)
	return data
}

//Read the line, encrypt the message and send it
func (mes *Messanger) Write(dhInfo *dh.DHContext, feed *DB.Feed) {

	fmt.Print("\n" + feed.GetHistory())
	for {
		if mes.connection != nil {

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			mes.Commands(text)
			key, salt := crpt.DeriveKey([]uint64(dhInfo.SharedSecret), nil)
			textTosend, _ := crpt.Encrypt(key, text)

			feed.EditChatHistory(mes.peer.Port, mes.Name+":"+text)
			send := Message{textTosend, (salt)}
			json_data, err := json.Marshal(send)
			if err != nil {
				log.Fatal(err)
			}
			mes.connection.Write([]byte(string(json_data) + "\n"))
		}
	}
}

//Waiting for message, later it decrypt it and write it in console
func (mes *Messanger) Read(dhInfo *dh.DHContext, feed *DB.Feed) {
	for {
		if mes.connection != nil {

			message, _ := bufio.NewReader(mes.connection).ReadString('\n')
			err := json.Unmarshal([]byte(message), &mes.message)
			if err != nil {

				log.Fatal(err)
			}
			key, _ := crpt.DeriveKey([]uint64(dhInfo.SharedSecret), mes.message.Salt)
			messageToRead, _ := crpt.Decrypt(key, mes.message.Message)
			feed.EditChatHistory(mes.peer.Port, mes.peer.Name+":"+messageToRead)
			fmt.Print(mes.peer.Name + ":" + messageToRead)
		}
	}
}
func (mes *Messanger) Commands(str string) bool {
	str = strings.TrimSpace(str)
	str = strings.Trim(str, "\n")
	if str[0:1] != "/" {
		return false
	}
	str = str[1:2]
	switch str {
	case "g":
		{
			return true
		}
	case "e":
		{
			os.Exit(0)
		}
	case "p":
		{
			str = strings.TrimSpace(str[2:])
			if len(str) > 5 || len(str) == 0 {
				fmt.Println("Incorect input")
			} else {
				port, err := strconv.Atoi(str)
				if err != nil {
					fmt.Println("Some problem with value")
				}
				mes.UrPort = port
			}

		}
	case "n":
		{
			str = strings.TrimSpace(str[2:])
			if len(str) == 0 {
				fmt.Println("Space is not nickname")
			} else {
				mes.Name = str
			}
		}
	case "j":
		{
			str = strings.TrimSpace(str[2:])
			if len(str) > 5 || len(str) == 0 {
				fmt.Println("Incorect input")
			} else {
				port, err := strconv.Atoi(str)
				if port < 2000 || port > 65534 {
					fmt.Println("Incorect input")
				} else {
					if err != nil {
						fmt.Println("Some problem with value")
					}
					mes.PortToCon = port
				}
			}
		}
	default:
		{
			fmt.Println("ia ne znau chto eto")
		}
	}
	return false
}
