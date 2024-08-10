package helper

import (
	"crypto/rand"
	"io"
	mrand "math/rand"
	"time"
)

func GenerateId() uint {
	mrand.Seed(time.Now().UnixNano())
	id := mrand.Intn(100) + 1
	return uint(id)
}

func GenerateAccountNumber() string {
	var m = 6
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, m)
	number, err := io.ReadAtLeast(rand.Reader, b, m)

	if number != m {
		panic(err)
	}

	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func GenerateAccountBranch() string {
	var m = 4
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, m)
	branch, err := io.ReadAtLeast(rand.Reader, b, m)

	if branch != m {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
