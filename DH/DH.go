package dh

import (
	"log"
	"math/rand"
	"time"

	bigintegers "github.com/Agilen/MessangerP2P/BigIntegers"
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
