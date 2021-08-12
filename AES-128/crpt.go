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

	bigintegers "github.com/Agilen/MessangerP2P/BigIntegers"
	"golang.org/x/crypto/pbkdf2"
)

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

func DeriveKey(passphrase []uint64, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		// http://www.ietf.org/rfc/rfc2898.txt
		// Salt.
		rand.Read(salt)
	}
	return pbkdf2.Key([]byte(bigintegers.ToHex(passphrase)), salt, 1000, 16, sha256.New), salt
}
