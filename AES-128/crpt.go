package crpt

import (
	"crypto/aes"
	"crypto/cipher"
	rn "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	rand "math/rand"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

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

var dhInfo DH

func Encrypt(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	fmt.Println(len(iv))

	if _, err = io.ReadFull(rn.Reader, iv); err != nil {

		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}

func Decrypt(key []byte, securemess string) (decodedmess string, err error) {
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

func DeriveKey(passphrase string, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		rand.Read(salt)
	}
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 16, sha256.New), salt
}

func GenRandomNum(size int) []uint64 {

	if size <= 0 {
		println("Size is bellow zero or zero")
		log.Fatal()
	}
	a := make([]uint64, size)
	for i := 0; i < size; i++ {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		a[i] = r1.Uint64()
		time.Sleep(100)
	}
	return a
}
