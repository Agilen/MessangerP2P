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
		go onConnection(conn)
		wg.Done()
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
func onConnection(conn net.Conn) { // обработка запроса
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
			key, salt := deriveKey(dhInfo.SharedSecret, nil)
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
func Read() {
	for {
		if info.connection != nil {

			message, _ := bufio.NewReader(info.connection).ReadString('\n')
			err := json.Unmarshal([]byte(message), &mes)
			if err != nil {

				log.Fatal(err)
			}
			fmt.Println("i got it")

			key, _ := deriveKey(dhInfo.SharedSecret, mes.Salt)
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

func deriveKey(passphrase string, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		// http://www.ietf.org/rfc/rfc2898.txt
		// Salt.
		rand.Read(salt)
	}
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 16, sha256.New), salt
}
