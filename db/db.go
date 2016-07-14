package db

import (
	"github.com/poslegm/blockchain-chat/network"
	"github.com/boltdb/bolt"
	"fmt"
	"encoding/json"
	"encoding/binary"
)

var db *bolt.DB;

var (
	knownAddresses = []byte("knownAddresses") //addresses of known network peers
	messages = []byte("messages") //stored messages
	blocks = []byte("blocks") //stored blockchain blocks
	keys = []byte("keys") //user's encryption keys
	contacts = []byte("contacts") //contact list
	baseName = "data.db" //database file name
	neededBuckets = [][]byte{knownAddresses, messages, blocks, keys, contacts} //needed buckets
)

//database initialization, creating top-level buckets if they are not present
//func InitDB(path string, mode os.FileMode, options *bolt.Options) (err error) {
func InitDB() (err error) {
	db, err = bolt.Open(baseName, 0660, nil)
	if (err != nil) {
		return err
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
			if (err != nil) {
				return fmt.Errorf("add known addresses marshal: %s", err)
			}
			bid, _ := b.NextSequence()
			id := int(bid)
			err = b.Put(itob(id), buf)
			if (err != nil) {
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
			if (err != nil) {
				return fmt.Errorf("add message marshal: %s", err)
			}
			bid, _ := b.NextSequence()
			id := int(bid)
			err = b.Put(itob(id), buf)
			if (err != nil) {
				return fmt.Errorf("add message db put: %s", err)
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
			if (terr != nil) {
				return fmt.Errorf("get messages unmarshal: %s", terr)
			}
			data = append(data, msg)
		}
		return nil
	})
	return
}

func GetPublicKey() string {
	return "123456789"
}