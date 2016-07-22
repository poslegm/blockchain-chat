package message

import (
	"testing"
	"github.com/poslegm/blockchain-chat/shahash"
	"fmt"
)

func TestMessage(t *testing.T) {
	txt := TextMessage{"recv", "send", "text", 10, shahash.ShaHash{}, shahash.ShaHash{}, 0, 0}
	err := txt.Mine()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(txt.MessageHash)
	fmt.Println(txt.Nonce)
	ver, err := txt.Verify()
	if err != nil || !ver {
		t.Fatal("verification error")
	}

}
