package db

import (
	"github.com/poslegm/blockchain-chat/network"
	"github.com/boltdb/bolt"
	"fmt"
	"encoding/json"
	"encoding/binary"
	"github.com/poslegm/blockchain-chat/message"
	"os"
	"errors"
)

var db *bolt.DB;

var (
	knownAddresses = []byte("knownAddresses") //addresses of known network peers
	messages = []byte("messages") //stored messages
	blocks = []byte("blocks") //stored blockchain blocks
	keys = []byte("keys") //user's encryption keys
	contacts = []byte("contacts") //contact list
	baseName = "data.db" //database file name
	testBaseName = "test.db" //database file name for testing
	neededBuckets = [][]byte{knownAddresses, messages, blocks, keys, contacts} //needed buckets
)

//database initialization, creating top-level buckets if they are not present
//func InitDB(path string, mode os.FileMode, options *bolt.Options) (err error) {
func initDB(fileName string) (err error) {
	db, err = bolt.Open(fileName, 0660, nil)
	if (err != nil) {
		return fmt.Errorf("db init: %s", err)
	}
	return db.Update(func(tx *bolt.Tx) error {
		//create bucket for each needed
		for _, bucket := range neededBuckets {
			_, terr := tx.CreateBucketIfNotExists(bucket)
			if (terr != nil) {
				return fmt.Errorf("create bucket %s: %s", bucket, terr)
			}
		}
		return nil
	})
}

func InitDB() error {
	return initDB(baseName)
}

//database closing
func CloseDB() error {
	err := db.Close()
	if err != nil {
		return fmt.Errorf("db close: %s", err)
	}
	return nil
}

//for testing
func tInitDB() error {
	return initDB(testBaseName)
}
func tCloseDB() error {
	err := CloseDB()
	if err != nil {
		return err
	}
	os.Remove(testBaseName)
	return nil
}

//get known peers to connect to them
func GetKnownAddresses() (data []network.NetAddress, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(knownAddresses)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			addr := network.NetAddress{}
			terr := json.Unmarshal(v, &addr)
			if (terr != nil) {
				return fmt.Errorf("get known adrresses unmarshal: %s", terr)
			}
			data = append(data, addr)
		}
		return nil
	})
	return
}

//convert int to byte array
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

//add new known peer
func AddKnownAddresses(data []network.NetAddress) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(knownAddresses)
		for _, v := range data {
			buf, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("add known addresses marshal: %s", err)
			}
			bid, _ := b.NextSequence()
			id := int(bid)
			err = b.Put(itob(id), buf)
			if err != nil {
				return fmt.Errorf("add known addresses db put: %s", err)
			}
		}
		return nil
	})
}

//add new messages
func AddMessages(data []network.NetworkMessage) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(messages)
		for _, v := range data {
			buf, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("add messages marshal: %s", err)
			}
			bid, _ := b.NextSequence()
			id := int(bid)
			err = b.Put(itob(id), buf)
			if err != nil {
				return fmt.Errorf("add messages db put: %s", err)
			}
		}
		return nil
	})
}

//get stored messages
func GetAllMessages() (data []network.NetworkMessage, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(messages)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			msg := network.NetworkMessage{}
			terr := json.Unmarshal(v, &msg)
			if terr != nil {
				return fmt.Errorf("get all messages unmarshal: %s", terr)
			}
			data = append(data, msg)
		}
		return nil
	})
	return
}

func GetPublicKey() (string, error) {
	keyPairs, err := GetAllKeys()
	if err != nil {
		fmt.Println("db.GetPublicKey: can't get keys from db ", err.Error())
		return "", err
	} else if len(keyPairs) == 0 {
		fmt.Println("db.GetPublicKey: there is no key pairs in db")
		return "", errors.New("There is no key pairs in db")
	}

	return keyPairs[0].GetBase58Address(), nil
}

func AddKeys(data []*message.KeyPair) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(keys)
		for _, v := range data {
			buf, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("add keys marshal: %s", err)
			}
			key := []byte(v.GetBase58Address())
			err = b.Put(key, buf)
			if err != nil {
				return fmt.Errorf("add keys db put: %s", err)
			}
		}
		return nil
	})
}

func GetAllKeys() (data []*message.KeyPair, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(keys)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			kp := &message.KeyPair{}
			terr := json.Unmarshal(v, kp)
			if terr != nil {
				return fmt.Errorf("get all keys unmarshal: %s", terr)
			}
			data = append(data, kp)
		}
		return nil
	})
	return
}

func GetKeyByAddress(address string) (kp *message.KeyPair, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(keys)
		buf := b.Get([]byte(address))
		if buf == nil {
			return nil
		}
		kp = &message.KeyPair{}
		err := json.Unmarshal(buf, kp)
		if err != nil {
			return fmt.Errorf("get key by addrress unmarshal: %s", err)
		}
		return nil
	})
	return
}