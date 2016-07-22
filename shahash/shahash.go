package shahash

import (
	"encoding/hex"
	"fmt"
	"crypto/sha256"
)

//sha256 hash size in bytes
const HashSize = 32

//hash string length
const HashStringSize = HashSize * 2

//number of first zero bytes
const Difficulty = 2

//hash type
type ShaHash [HashSize]byte

//convert hash to string
func (hash ShaHash) String() string {
	return hex.EncodeToString(hash[:])
}

//create hash from string
func ShaHashFromString(hash string) (ShaHash, error) {
	if len(hash) > HashStringSize {
		return ShaHash{}, fmt.Errorf("hash string too long")
	}

	if len(hash) % 2 != 0 {
		hash = "0" + hash
	}

	buf, err := hex.DecodeString(hash)
	if err != nil {
		return ShaHash{}, err
	}
	var ret ShaHash
	copy(ret[:], buf)
	return ret, nil
}

//create hash from precounted sum
func ShaHashFromSum256(sum []byte) (ShaHash, error) {
	if len(sum) != HashSize {
		return ShaHash{}, fmt.Errorf("hash is %v bytes long, provided %v", HashSize, len(sum))
	}
	var ret ShaHash
	copy(ret[:], sum)
	return ret, nil
}

//create hash from byte array
func ShaHashFromData(data []byte) (ShaHash, error) {
	sum1 := sha256.Sum256(data)
	sum2 := sha256.Sum256(sum1[:])
	return ShaHashFromSum256(sum2[:])
}

func (h ShaHash) Equal(h2 ShaHash) bool {
	for i := 0; i < HashSize; i++ {
		if h[i] != h2[i] {
			return false
		}
	}
	return true
}

func (h ShaHash) Check() bool {
	for i := 0; i < Difficulty; i++ {
		if h[i] != 0 {
			return false
		}
	}

	return true
}