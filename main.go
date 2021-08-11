package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	bigintegers "github.com/Agilen/MessangerP2P/BigIntegers"
	dh "github.com/Agilen/MessangerP2P/DH"
	"github.com/urfave/cli"
	"golang.org/x/crypto/pbkdf2"
)

type Message struct {
	Message string
	Salt    []byte
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

	var dhInfo dh.DHContext
	var wg sync.WaitGroup
	fmt.Println("Welcom")
	time.Sleep(1000)
	wg.Add(1)
	go Listen(&wg, &dhInfo)
	if info.PortToCon != 0 {
		go SendRequest(&wg, &dhInfo)
	}
	time.Sleep(1000)
	wg.Wait()
	go Write(&dhInfo)
	Read(&dhInfo)
}

func Listen(wg *sync.WaitGroup, dhInfo *dh.DHContext) {
	listener, _ := net.Listen("tcp", ":"+strconv.Itoa(info.UrPort))

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go onConnection(conn, wg, dhInfo)

	}
}
func SendRequest(wg *sync.WaitGroup, dhInfo *dh.DHContext) {

	connection, _ := net.Dial("tcp", ":"+strconv.Itoa(info.PortToCon))
	for connection == nil {
		connection, _ = net.Dial("tcp", ":"+strconv.Itoa(info.PortToCon))
	}
	info.connection = connection
	data := DataPreparation(dhInfo)
	fmt.Println("PB", dhInfo.DHParams.PublicKey)
	connection.Write([]byte(data + "\n"))

	message, _ := bufio.NewReader(connection).ReadString('\n')
	M := strings.Fields(message)
	fmt.Println("M", M)
	response := []byte(M[0] + " " + M[1] + " " + M[2])
	fmt.Println("asdf", M[0], M[1], M[2])
	peerData := []byte(M[3])
	err := json.Unmarshal(peerData, &peer)
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
func onConnection(conn net.Conn, wg *sync.WaitGroup, dhInfo *dh.DHContext) { // обработка запроса
	fmt.Printf("New connection from: %v", conn.RemoteAddr().String())
	info.connection = conn

	message, _ := bufio.NewReader(conn).ReadString('\n') //жду запрос
	M := strings.Fields(message)
	data := []byte(M[0])
	peerData := []byte(M[1])
	err := json.Unmarshal(data, &dhInfo.DHParams)
	if err != nil {

		log.Fatal(err)
	}
	err = json.Unmarshal(peerData, &peer)
	if err != nil {

		log.Fatal(err)
	}
	dhInfo.GenerateDHPrivateKey()
	dhInfo.CalculateSharedSecret()
	dhInfo.CalculateDHPublicKey()

	r := dh.Params{dhInfo.DHParams.G, dhInfo.DHParams.P, dhInfo.DHParams.PublicKey}
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
func DataPreparation(dhInfo *dh.DHContext) string { // Подготовка данных для запроса на подключение
	*dhInfo = *dh.NewDHContext()
	dhInfo.GenerateDHPrivateKey()
	dhInfo.CalculateDHPublicKey()
	fmt.Println("PB1", dhInfo.PrivateKey)
	dh := dh.Params{dhInfo.DHParams.G, dhInfo.DHParams.P, dhInfo.DHParams.PublicKey}
	json_data, err := json.Marshal(dh)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PB2", dhInfo.PrivateKey)
	pr := Peer{info.Name}
	json_peerData, err := json.Marshal(pr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PB3", dhInfo.PrivateKey)
	data := string(json_data) + " " + string(json_peerData)
	fmt.Println(data)
	fmt.Println("PB4", dhInfo.PrivateKey)
	return data
}

func Write(dhInfo *dh.DHContext) {
	for {
		if info.connection != nil {

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			key, salt := deriveKey([]uint64(dhInfo.SharedSecret), nil)
			textTosend, _ := encrypt(key, text)
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
func Read(dhInfo *dh.DHContext) {
	for {
		if info.connection != nil {

			message, _ := bufio.NewReader(info.connection).ReadString('\n')
			err := json.Unmarshal([]byte(message), &mes)
			if err != nil {

				log.Fatal(err)
			}
			fmt.Println("i got it")

			key, _ := deriveKey([]uint64(dhInfo.SharedSecret), mes.Salt)
			fmt.Println("message to ddec", message)
			messageToRead, _ := decrypt(key, mes.Message)
			fmt.Println("decr mes", messageToRead)
			fmt.Println(key)
			fmt.Println(mes.Salt)
			fmt.Print(peer.Name + ":" + messageToRead)
		}
	}
}

func encrypt(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	fmt.Println(len(iv))

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}

func decrypt(key []byte, securemess string) (decodedmess string, err error) {
	fmt.Println(securemess)

	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	fmt.Println("promlem is here")
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Println("234567")
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	fmt.Println("12345")
	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}
	fmt.Println("sdfg")
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = string(cipherText)
	return
}

func deriveKey(passphrase []uint64, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		// http://www.ietf.org/rfc/rfc2898.txt
		// Salt.
		rand.Read(salt)
	}
	return pbkdf2.Key([]byte(bigintegers.ToHex(passphrase)), salt, 1000, 16, sha256.New), salt
}
