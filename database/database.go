package database

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"container/list"
	"encoding/gob"
	"math/rand"
	"time"
	"bytes"
	"errors"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))


type Pos struct {
	x int
	y int
}

type Area struct {
	width int
	height int
}

type Frame struct {
	Category int
	Area
	Pos
	Image string
	Color string
	Owner string
}

type Chunk struct {
	Area
	Pos
	Owner string
	Frames list.List
}


type Database struct {
	instance *leveldb.DB
}

func InitDB(dbname string) (*Database) {
	var db Database
	instance, err := leveldb.OpenFile("./"+dbname, nil)
	if err != nil {
		fmt.Println("Create database failed!!")
	}
	fmt.Println("Create database success!!")
	db.instance = instance
	gob.Register(Chunk{})
	return &db
}

func (d Database) FindChunk(x int, y int) (Chunk, error) {
	buffer := bytes.NewBuffer(nil)
	decoder := gob.NewDecoder(buffer)
	iter := d.instance.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Println(iter.Key())
		value := iter.Value()
		buffer.Reset()
		buffer.Write(value)
		// decode the value
		var targetChunk Chunk
		decoder.Decode(&targetChunk)
		fmt.Println(targetChunk)

		if targetChunk.Pos.x == x && targetChunk.Pos.y == y {
			return targetChunk, nil
		}
	}

	return Chunk{}, nil
}

func (d Database) NewChunk(x int, y int) (error) {
	flag := false
	// check there is not same chunk on x, y
	iter := d.instance.NewIterator(nil, nil)
	for iter.Next() {
		// decode the value.
		var targetChunk Chunk
		err := decode(iter.Value(), &targetChunk)
		if err != nil {
			return errors.New("Decode the value failed.")
		}
		fmt.Println(targetChunk)
		if targetChunk.Pos.x == x && targetChunk.Pos.y == y {
			flag = true
			break
		}
	}
	if flag {
		return errors.New("The location exist a chunk, you can not create a chunk on the same location.")
	}
	// new the chunk and store in database.
	var newChunk Chunk
	newChunk.Area = Area{16, 16}
	newChunk.Pos = Pos{x, y}
	result, err_encode := encode(newChunk)
	if err_encode != nil {
		return errors.New("Encode value failed")
	}
	err := d.instance.Put([]byte(randomString(10, charset)), result, nil)
	if err != nil {
		return errors.New("Put data in database failed.")
	}
	return nil
}
/*

func UpdateChuck(x int, y int, data struct) {
}

func DeleteChunk(x int, y int) {

}
*/

func encode(data interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func decode(data []byte, to interface{}) (error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(to)
}

func randomString(length int, charset string) (string) {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}


