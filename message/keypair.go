package message

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/maxwellhealth/go-gpg"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	_ "golang.org/x/crypto/ripemd160"
	"io/ioutil"
	"os"
)

//gpg key pair
type KeyPair struct {
	//gpg pub key
	PublicKey []byte

	//gpg private key
	PrivateKey []byte

	//gpg private key passphrase
	Passphrase []byte
}

//encode data using kp's public key
func (kp *KeyPair) Encode(data []byte) ([]byte, error) {
	//check publickey existence
	if kp.PublicKey == nil {
		return nil, fmt.Errorf("no public key provided")
	}

	//create buffers
	inputBuffer := bytes.NewBuffer(data)
	var outputBuffer bytes.Buffer

	//encode
	//fmt.Println(string(kp.PublicKey))
	//fmt.Println(string(data))
	err := gpg.Encode(kp.PublicKey, inputBuffer, &outputBuffer)
	if err != nil {
		return nil, fmt.Errorf("error encoding data: %s", err)
	}

	return ioutil.ReadAll(&outputBuffer)
}

//decode data using kp's private key
func (kp *KeyPair) Decode(data []byte) ([]byte, error) {
	//check privatekey existence
	if kp.PrivateKey == nil {
		return nil, fmt.Errorf("no private key provided")
	}

	//create buffers
	inputBuffer := bytes.NewBuffer(data)
	var outputBuffer bytes.Buffer

	//decode
	err := gpg.Decode(kp.PrivateKey, kp.Passphrase, inputBuffer, &outputBuffer)
	if err != nil {
		return nil, fmt.Errorf("error decoding data: %s", err)
	}

	return ioutil.ReadAll(&outputBuffer)
}

//get address to send to
func (kp *KeyPair) GetBase58Address() string {
	sum := md5.Sum(kp.PublicKey)
	return base58.Encode(sum[:])
}

//string representation
func (kp *KeyPair) String() string {
	return "pub:" + string(kp.PublicKey) + "\npriv:" + string(kp.PrivateKey) + "\npassphrase:" + string(kp.Passphrase)
}

func KeyPairFromFile(publicKeyFile, privateKeyFile, passphrase string) (*KeyPair, error) {
	pub, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("keypair from file cannot open %s: %s", publicKeyFile, err)
	}
	priv, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("keypair from file cannot open %s: %s", privateKeyFile, err)
	}
	_, err = openpgp.ReadArmoredKeyRing(bytes.NewBuffer(pub))
	if err != nil {
		return nil, fmt.Errorf("public key read error: %s", err)
	}
	_, err = openpgp.ReadArmoredKeyRing(bytes.NewBuffer(priv))
	if err != nil {
		return nil, fmt.Errorf("private key read error: %s", err)
	}
	//TODO fix this crutch
	_, err = os.Stat(passphrase)
	if err == nil {
		pass, err := ioutil.ReadFile(passphrase)
		if err != nil {
			return nil, fmt.Errorf("passphrase read error: %s", err)
		}
		passphrase = string(pass)
	}
	return &KeyPair{PrivateKey: priv, PublicKey: pub, Passphrase: []byte(passphrase)}, nil
}

func (kp *KeyPair) SaveToFile(name string) error {
	pubFileName := name + ".pub"
	privFileName := name + ".priv"
	passFileName := name + ".pass"

	err := ioutil.WriteFile(pubFileName, kp.PublicKey, 0660)
	if err != nil {
		return fmt.Errorf("public key write error: %s", err)
	}

	err = ioutil.WriteFile(privFileName, kp.PrivateKey, 0660)
	if err != nil {
		return fmt.Errorf("private key write error: %s", err)
	}

	err = ioutil.WriteFile(passFileName, kp.Passphrase, 0660)
	if err != nil {
		return fmt.Errorf("passphrase write error: %s", err)
	}
	return nil
}

func GenerateKeyPair(name, comment, email, passphrase string) (*KeyPair, error) {
	//unused, because openpgp doesn't create encrypted keys
	_ = passphrase
	//create new entity
	entity, err := openpgp.NewEntity(name, comment, email, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot generate keypair with %s %s %s: %s", name, comment, email, err)
	}

	//sign ourselves
	for _, id := range entity.Identities {
		err := id.SelfSignature.SignUserId(id.UserId.Id, entity.PrimaryKey, entity.PrivateKey, nil)
		if err != nil {
			return nil, fmt.Errorf("cannot sign key pair %s", err)
		}
	}

	kp := &KeyPair{}
	kp.Passphrase = []byte("")

	//serialize private key
	buf := &bytes.Buffer{}
	encoder, err := armor.Encode(buf, openpgp.PrivateKeyType, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create private armor encoder: %s", err)
	}
	err = entity.SerializePrivate(encoder, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot serialize private key: %s", err)
	}
	encoder.Close()
	kp.PrivateKey, err = ioutil.ReadAll(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot write private key: %s", err)
	}

	//serialize public key
	buf = &bytes.Buffer{}
	encoder, err = armor.Encode(buf, openpgp.PublicKeyType, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot creare public armor encoder: %s", err)
	}
	err = entity.Serialize(encoder)
	if err != nil {
		return nil, fmt.Errorf("cannot serialize public key: %s", err)
	}
	encoder.Close()
	kp.PublicKey, err = ioutil.ReadAll(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot write public key: %s", err)
	}

	//fmt.Println(kp)
	return kp, nil
}
