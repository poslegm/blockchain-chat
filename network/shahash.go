package network

import (
	"encoding/hex"
	"fmt"
	"crypto/sha256"
)

//sha256 hash size in bytes
const HashSize = 32

//hash string length
const HashStringSize = HashSize * 2

//hash type
type ShaHash [HashSize]byte

//convert hash to string
func (hash *ShaHash) String() string {
	return hex.EncodeToString(hash[:])
}

//create hash from string
func ShaHashFromString(hash string) (*ShaHash, error) {
	if len(hash) > HashStringSize {
		return nil, fmt.Errorf("hash string too long")
	}

	if len(hash) % 2 != 0 {
		hash = "0" + hash
	}

	buf, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	var ret ShaHash
	copy(ret[:], buf)
	return &ret, nil
}

//create hash from precounted sum
func ShaHashFromSum256(sum []byte) (*ShaHash, error) {
	if len(sum) != HashSize {
		return nil, fmt.Errorf("hash is %v bytes long, provided %v", HashSize, len(sum))
	}
	var ret ShaHash
	copy(ret[:], sum)
	return &ret, nil
}

//create hash from byte array
func ShaHashFromData(data []byte) (*ShaHash, error) {
	sum1 := sha256.Sum256(data)
	sum2 := sha256.Sum256(sum1[:])
	return ShaHashFromSum256(sum2[:])
}
