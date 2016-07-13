package network

import (
	"testing"
	"crypto/sha256"
	"fmt"
)

func TestShaHash(t *testing.T) {
	sum1 := sha256.Sum256([]byte("blockchain-chat"))
	sum2 := sha256.Sum256(sum1[:])
	shahash1, err := ShaHashFromSum256(sum2[:])
	if err != nil {
		t.Errorf("hash from sum error %s", err)
	}
	fmt.Println(shahash1)
	shahash, err := ShaHashFromString(shahash1.String())
	if err != nil {
		t.Errorf("hash from string error %s", err)
	}
	fmt.Println(shahash)
	shahash2, err := ShaHashFromData([]byte("blockchain-chat"))
	if err != nil {
		t.Errorf("hash from data error %s", err)
	}
	fmt.Println(shahash2)
}
