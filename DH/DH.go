package dh

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"time"

	bigintegers "github.com/Agilen/MessangerP2P/BigIntegers"
)

type BigInteger []uint64

type DHContext struct {
	PrivateKey   BigInteger
	SharedSecret BigInteger
	DHParams     Params
}

type Params struct {
	G         BigInteger
	P         BigInteger
	PublicKey BigInteger
}

func (st Params) MarshalJSON() ([]byte, error) {
	A := bigintegers.ToHex([]uint64(st.G))
	B := bigintegers.ToHex([]uint64(st.P))
	C := bigintegers.ToHex([]uint64(st.PublicKey))
	res := A + " " + B + " " + C
	data, err := json.Marshal(res)
	return data, err
}
func (st *Params) UnmarshalJSON(by []byte) error {
	var str string
	err := json.Unmarshal(by, &str)
	M := strings.Fields(str)
	st.G = bigintegers.ReadHex(M[0])
	st.P = bigintegers.ReadHex(M[1])
	st.PublicKey = bigintegers.ReadHex(M[2])

	return err
}

func NewDHContext() *DHContext {
	context := &DHContext{
		DHParams: Params{
			G: GenRandomNum(2),
			P: GenRandomNum(2),
		},
	}

	return context
}
func (dh *DHContext) GenerateDHPrivateKey() {
	dh.PrivateKey = GenRandomNum(2)
}
func (dh *DHContext) CalculateDHPublicKey() {
	dh.DHParams.PublicKey = PowMod([]uint64(dh.DHParams.G), []uint64(dh.PrivateKey), []uint64(dh.DHParams.P))
}
func (dh *DHContext) CalculateSharedSecret() {

	dh.SharedSecret = PowMod(dh.DHParams.PublicKey, dh.PrivateKey, dh.DHParams.P)
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

func PowMod(a []uint64, pow []uint64, p []uint64) []uint64 {
	return bigintegers.LongModPowerBarrett(a, pow, p)
}
