package bigintegers

import (
	"fmt"
	"strconv"
	"strings"
)

func ReadDec(str string) []uint64 {
	str = strings.TrimSpace(str)
	numStr := strings.Split(str, "")
	var num []uint64
	for i := len(numStr) - 1; i >= 0; i-- {
		n, _ := strconv.ParseUint(numStr[i], 10, 64)
		num = append(num, n)
	}
	res := []uint64{0}
	for i := 0; i < len(num); i++ {
		res = LongAdd(LongMul(Pow(10, i), []uint64{num[i]}), res)
	}
	res = DelNull(res)

	return res
}
func ReadHex(str string) []uint64 {
	str = strings.TrimSpace(str)
	Lenght := 0
	if len(str) < 16 {
		Digit := make([]uint64, 1)
		n, _ := strconv.ParseUint(str, 16, 64)
		Digit[0] = uint64(n)
		return Digit
	}
	if len(str)%16 != 0 {
		Lenght = len(str)/16 + 1
	} else {
		Lenght = len(str) / 16
	}
	Digit := make([]uint64, Lenght)
	for i := 0; i < Lenght-1; i++ {
		s := str[len(str)-(i+1)*16 : len(str)-i*16]
		n, _ := strconv.ParseUint(s, 16, 64)
		Digit[i] = uint64(n)
	}
	if len(str)%16 != 0 {
		s := str[0 : len(str)%16]
		n, _ := strconv.ParseUint(s, 16, 64)
		Digit[Lenght-1] = uint64(n)
	} else if len(str)%16 == 0 {
		s := str[0:16]
		n, _ := strconv.ParseUint(s, 16, 64)
		Digit[len(Digit)-1] = uint64(n)
	}

	return Digit
}
func ReadBin(str string) []uint64 {
	str = strings.TrimSpace(str)
	Lenght := 0
	if len(str) < 64 {
		Digit := make([]uint64, 1)
		n, _ := strconv.ParseUint(str, 2, 64)
		Digit[0] = uint64(n)
		return Digit
	}
	if len(str)%64 != 0 {
		Lenght = len(str)/64 + 1
	} else {
		Lenght = len(str) / 64
	}
	Digit := make([]uint64, Lenght)
	for i := 0; i < Lenght-1; i++ {
		s := str[len(str)-(i+1)*64 : len(str)-i*64]
		n, _ := strconv.ParseUint(s, 2, 64)
		Digit[i] = uint64(n)
	}
	if len(str)%64 != 0 {
		s := str[0 : len(str)%64]
		n, _ := strconv.ParseUint(s, 2, 64)
		Digit[Lenght-1] = uint64(n)
	} else if len(str)%64 == 0 {
		s := str[0:64]
		n, _ := strconv.ParseUint(s, 2, 64)
		Digit[len(Digit)-1] = uint64(n)
	}

	return Digit
}
func DelLeadZero(a string) string {
	if string(a[0]) == "0" {
		for string(a[0]) == "0" {
			a = a[1:]
		}
	}
	return a
}
func DelNull(a []uint64) []uint64 {
	k := 0
	for i := len(a) - 1; i >= 0; i-- {
		if a[i] != 0 {
			k = i
			break
		}
	}
	a = a[0 : k+1]
	return a
}
func ToHex(a []uint64) string {
	result := ""
	digit := ""
	buf := ""
	k := len(a) - 1
	for i := 0; i < len(a); i++ {
		digit = fmt.Sprintf("%X", a[k])
		if len(digit) < 16 {
			for i := 0; i < 16-len(digit); i++ {
				buf += "0"
			}
		}
		buf += digit
		result += buf
		buf = ""
		k--
	}

	return result
}
func ToDec(a []uint64) string {
	ten := []uint64{10}
	output := ""
	for LongCmp(a, ten) == 1 {
		mod := LongMod(a, ten)
		output = strconv.Itoa(int(mod[0])) + output
		a = LongDiv(a, ten)
	}
	output = strconv.Itoa(int(a[0])) + output
	return output
}
func ToUInt32(a []uint64) []uint32 {
	str := ToHex(a)
	Lenght := 0
	if len(str) < 8 {
		Digit := make([]uint32, 1)
		n, _ := strconv.ParseUint(str, 16, 64)
		Digit[0] = uint32(n)
		return Digit
	}
	if len(str)%8 != 0 {
		Lenght = len(str)/8 + 1
	} else {
		Lenght = len(str) / 8
	}
	Digit := make([]uint32, Lenght)
	for i := 0; i < Lenght-1; i++ {
		s := str[len(str)-(i+1)*8 : len(str)-i*8]
		n, _ := strconv.ParseUint(s, 16, 64)
		Digit[i] = uint32(n)
	}
	if len(str)%8 != 0 {
		s := str[0 : len(str)%8]
		n, _ := strconv.ParseUint(s, 16, 64)
		Digit[Lenght-1] = uint32(n)
	} else if len(str)%8 == 0 {
		s := str[0:8]
		n, _ := strconv.ParseUint(s, 16, 64)
		Digit[len(Digit)-1] = uint32(n)
	}

	return Digit
}
func ToUInt64(a []uint32) []uint64 {
	result := ""
	digit := ""
	buf := ""
	k := len(a) - 1
	for i := 0; i < len(a); i++ {
		digit = fmt.Sprintf("%X", a[k])
		if len(digit) < 8 {
			for i := 0; i < 8-len(digit); i++ {
				buf += "0"
			}
		}
		buf += digit
		result += buf
		buf = ""
		k--
	}

	return ReadHex(result)
}
func ToBin(a []uint64) string {
	result := ""
	digit := ""
	buf := ""
	k := len(a) - 1
	for i := 0; i < len(a); i++ {
		digit = strconv.FormatUint((a[k]), 2)
		if len(digit) < 64 {
			for i := 0; i < 64-len(digit); i++ {
				buf += "0"
			}
		}
		buf += digit
		result += buf
		buf = ""
		k--
	}

	return result
}
func ToBinDigit(a uint64) string {
	result := ""
	digit := ""
	buf := ""

	digit = strconv.FormatUint((a), 2)
	if len(digit) < 64 {
		for i := 0; i < 64-len(digit); i++ {
			buf += "0"
		}
	}
	buf += digit
	result += buf
	buf = ""

	return result
}

func Pow(a, k int) []uint64 {
	if k == 0 {
		return []uint64{1}
	}
	n := []uint64{uint64(a)}
	res := []uint64{uint64(a)}
	for i := 0; i < k-1; i++ {
		res = LongMul(res, n)
		res = DelNull(res)
	}
	return res
}
func LongCmp(a, b []uint64) int {
	a, b = SameSize(a, b)

	for i := len(a) - 1; i >= 0; i-- {
		if a[i] > b[i] {
			return 1 //>
		} else if a[i] < b[i] {
			return -1 //<
		}
	}
	return 0 //=
}
func LongAdd(a, b []uint64) []uint64 {
	a, b = SameSize(a, b)
	C := make([]uint64, len(a)+1)
	carry := uint64(0)
	for i := 0; i < len(a); i++ {
		temp := a[i] + b[i] + carry
		C[i] = temp & 0xffffffffffffffff
		carry = isCarryExist(a[i], b[i], C[i])

	}
	C[len(a)] = carry
	return DelNull(C)
}
func LongSub(a, b []uint64) []uint64 {
	a, b = SameSize(a, b)
	c := make([]uint64, len(a))
	borrow := uint64(0)
	for i := 0; i < len(a); i++ {
		c[i] = a[i] - b[i] - (borrow)
		if b[i] != 0 && b[i]+(borrow) == 0 {
			borrow = 1
		} else if a[i] >= b[i]+(borrow) {
			borrow = 0
		} else {
			borrow = 1
		}
	}
	return DelNull(c)
}
func LongMulOneDigit(a []uint32, b uint32) []uint32 {
	c := make([]uint32, len(a)+1)
	carry := uint64(0)
	for i := 0; i < len(a); i++ {
		temp := (uint64(a[i])*uint64(b) + carry)
		c[i] = uint32(temp & 0xffffffff)
		carry = temp >> 32
	}
	c[len(a)] = uint32(carry)
	return c
}
func LongMul(a, b []uint64) []uint64 {
	a, b = SameSize(a, b)
	A := ToUInt32(a)
	B := ToUInt32(b)
	c := make([]uint64, len(a)*2)
	for i := 0; i < len(B); i++ {
		temp := LongMulOneDigit(A, B[i])
		temp64 := ToUInt64(temp)

		temps := LongShiftLeft(temp64, i*32)

		c = LongAdd(c, temps)
	}
	return DelNull(c)
}
func LongShiftLeft(a []uint64, shiftVal int) []uint64 {
	if shiftVal <= -1 {
		return []uint64{1}
	}
	if shiftVal == 0 {
		return DelNull(a)
	}
	a = DelNull(a)
	for i := shiftVal; i > 0; i -= 64 {
		a = append(a, 0)
	}
	A := ToUInt32(a)
	shiftAmount := 32
	buflen := len(A)

	for buflen > 1 && A[buflen-1] == 0 {
		buflen = buflen - 1
	}

	for count := shiftVal; count > 0; {

		if count < shiftAmount {
			shiftAmount = count
		}
		carry := uint64(0)
		for i := 0; i < buflen; i++ {
			val := uint64(A[i]) << uint64(shiftAmount)
			val |= carry

			A[i] = uint32(val & 0xffffffff)
			carry = val >> 32
		}

		if carry != 0 {
			if buflen+1 <= len(A) {
				A[buflen] = uint32(carry)
				buflen++
			}
		}
		count -= shiftAmount
	}

	return DelNull(ToUInt64(A))
}
func LongShiftRight(a []uint64, shiftVal int) []uint64 {
	a = DelNull(a)
	A := ToUInt32(a)
	shiftAmount := 32
	invShift := 0
	buflen := len(A)

	for buflen > 1 && A[buflen-1] == 0 {
		buflen = buflen - 1
	}

	for count := shiftVal; count > 0; {

		if count < shiftAmount {
			shiftAmount = count
			invShift = 32 - shiftAmount
		}
		carry := uint64(0)

		for i := buflen - 1; i >= 0; i-- {
			val := uint64(A[i]) >> uint64(shiftAmount)
			val |= carry

			carry = uint64(A[i]) << uint64(invShift) & 0xffffffff
			A[i] = uint32(val)
		}

		count -= shiftAmount
	}

	return DelNull(ToUInt64(A))
}
func LongShiftRightV2(a []uint64, shiftVal int) []uint64 {
	B := ToBin(a)
	B = B[:len(B)-shiftVal]
	return ReadBin(B)
}
func LongDivMod(a, b []uint64) ([]uint64, []uint64) {
	k := BitLength(b)
	r := a
	t := 0
	var c []uint64
	var q []uint64
	i := 0
	for LongCmp(r, b) == 1 || LongCmp(r, b) == 0 {

		t = BitLength(r)
		c = LongShiftLeft(b, t-k)
		for LongCmp(r, c) == -1 {

			if LongCmp(r, c) == -1 {
				t = t - 1
				c = LongShiftLeft(b, t-k)
			}

		}

		r = LongSub(r, c)
		q = LongAdd(q, LongShiftLeft([]uint64{2}, t-k-1))

		i++

	}
	return DelNull(q), DelNull(r)
}
func LongDiv(a, b []uint64) []uint64 {
	k := BitLength(b)
	r := a
	t := 0
	var c []uint64
	var q []uint64

	for LongCmp(r, b) == 1 || LongCmp(r, b) == 0 {

		t = BitLength(r)
		c = LongShiftLeft(b, t-k)
		for LongCmp(r, c) == -1 {

			if LongCmp(r, c) == -1 {
				t = t - 1
				c = LongShiftLeft(b, t-k)
			}
		}
		r = LongSub(r, c)
		q = LongAdd(q, LongShiftLeft([]uint64{2}, t-k-1))
	}
	return DelNull(q)
}
func LongMod(a, b []uint64) []uint64 {
	k := BitLength(b)
	r := a
	t := 0
	var c []uint64
	for LongCmp(r, b) == 1 || LongCmp(r, b) == 0 {

		t = BitLength(r)
		c = LongShiftLeft(b, t-k)
		for LongCmp(r, c) == -1 {
			if LongCmp(r, c) == -1 {
				t = t - 1
				c = LongShiftLeft(b, t-k)
			}
		}
		r = LongSub(r, c)
	}
	return DelNull(r)
}
func LongModPowerBarrett(a, b, n []uint64) []uint64 {
	zero := []uint64{0}

	if LongCmp(b, zero) == 0 {
		return []uint64{1}
	}

	k := BitLength(n)
	mu := LongDiv(LongShiftLeft([]uint64{1}, 2*k), n)

	B := DelLeadZero(ToBin(b))
	m := len(B)
	c := []uint64{1}

	for i := m - 1; i > -1; i-- {
		if string(B[i]) == "1" {
			buf := LongMul(a, c)
			c = BarrettReduction(buf, n, mu)
		}
		a = BarrettReduction(LongMul(a, a), n, mu)
	}
	return DelNull(c)
}
func BarrettReduction(x, n, mu []uint64) []uint64 {

	if LongCmp(x, n) == -1 {
		return x
	}
	k := BitLength(n)

	q := KillDigits(x, k-1)
	q = LongMul(q, mu)
	q = KillDigits(q, k+1)
	r := LongMul(q, n)
	t := LongShiftLeft([]uint64{1}, k+1)
	r1 := x
	r2 := r
	if LongCmp(r1, r2) == 0 || LongCmp(r1, r2) == 1 {
		r = LongSub(r1, r2)
	} else {
		r = LongSub(LongAdd(t, r1), r2)
	}
	for LongCmp(r, n) == 0 || LongCmp(r, n) == 1 {
		r = LongSub(r, n)
	}
	return r
}

func KillDigits(a []uint64, k int) []uint64 {
	zero := []uint64{0}

	if LongCmp(a, zero) == 0 {
		return a
	}

	return LongShiftRight(a, k)
}
func isCarryExist(a, b, c uint64) uint64 {

	if a>>63 == 1 && b>>63 == 1 {
		return 1
	} else if a>>63 == 0 && b>>63 == 0 {
		return 0
	} else if c>>63 == 0 {
		return 1
	} else {
		return 0
	}
}
func SameSize(a, b []uint64) ([]uint64, []uint64) {
	leng := 0
	if len(a) == len(b) {
		return a, b
	}
	if len(a) > len(b) {
		leng = len(a)
		for i := len(b); i < leng; i++ {
			b = append(b, 0)
		}
	} else if len(a) < len(b) {
		leng = len(b)
		for i := len(a); i < leng; i++ {
			a = append(a, 0)
		}
	}

	return a, b
}
func BitLengtha(a []uint64) int {
	a = DelNull(a)
	for j := 63; j > 0; j-- {
		if a[len(a)-1]>>j&1 == 1 {
			return (j + 1) + 64*(len(a)-1)
		}
	}
	return 0
}
func BitLength(a []uint64) int {
	bit := DelLeadZero(ToBin(a))
	return len(bit)
}
func IsEvenNumber(a []uint64) bool {
	b := true
	if (a[0] & uint64(1)) == 0 {
		return b
	} else {
		return !b
	}
}
