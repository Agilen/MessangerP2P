package dh

import (
	"log"
	"math/rand"
	"time"

	bigintegers "example.com/DiffieHellmanGO/BigIntegers"
)

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

func testDH() ([]uint64, []uint64) {

	a := GenRandomNum(2)
	b := GenRandomNum(2)
	g := GenRandomNum(2)
	p := GenRandomNum(2)
	A := PowMod(g, a, p)
	B := PowMod(g, b, p)
	S1 := PowMod(B, a, p)
	S2 := PowMod(A, b, p)

	return S1, S2

}
