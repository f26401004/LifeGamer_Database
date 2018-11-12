package database

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"reflect"
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
	X int
	Y int
}

type Area struct {
	Width int
	Height int
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
	Frames []Frame
}


type Database struct {
	Instance *leveldb.DB
}

func InitDB(dbname string) (*Database, error) {
	// gob register concrete type
	gob.Register(Chunk{})
	var db Database
	instance, err := leveldb.OpenFile("./"+dbname, nil)
	if err != nil {
		fmt.Println("Create database failed!!")
		return nil, err
	}
	fmt.Println("Create database success!!")
	db.Instance = instance
	return &db, nil
}

func (d *Database) FindChunk(x int, y int) (Chunk, error) {
	buffer := bytes.NewBuffer(nil)
	decoder := gob.NewDecoder(buffer)
	iter := d.Instance.NewIterator(nil, nil)
	// check current iter value
	value := iter.Value()
	buffer.Reset()
	buffer.Write(value)
	// decode the value
	var targetChunk Chunk
	decoder.Decode(&targetChunk)
	if targetChunk.Pos.X == x && targetChunk.Pos.Y == y {
		return targetChunk, nil
	}
	for iter.Next() {
		value := iter.Value()
		buffer.Reset()
		buffer.Write(value)
		// decode the value
		var targetChunk Chunk
		decoder.Decode(&targetChunk)
		if targetChunk.Pos.X == x && targetChunk.Pos.Y == y {
			return targetChunk, nil
		}
	}
	iter.Release()

	return Chunk{}, nil
}

func (d *Database) NewChunk(x int, y int) (error) {
	flag := false
	// check there is no same chunk on x, y
	iter := d.Instance.NewIterator(nil, nil)
	for iter.Next() {
		// decode the value.
		var targetChunk Chunk
		err := decode(iter.Value(), &targetChunk)
		if err != nil {
			return errors.New("Decode the value failed.")
		}
		if targetChunk.Pos.X == x && targetChunk.Pos.Y == y {
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
		return err_encode
	}
	err := d.Instance.Put([]byte(randomString(10, charset)), result, nil)
	if err != nil {
		return errors.New("Put data in database failed.")
	}
	return nil
}

func (d *Database) UpdateChuck(x int, y int, data struct) (error) {
	flag := false
	// check there is chunk on x, y
	iter := d.Instance.NewIterator(nil, nil)
	for iter.Next() {
		// decode the value.
		var targetChunk Chunk
		err := decode(iter.Value(), &targetChunk)
		if err != nil {
			return errors.New("Decode the value failed.")
		}
		if targetChunk.Pos.X == x && targetChunk.Pos.Y == y {
			// TODO: update the new value of the chunk.
			flag = true
			break
		}
	}
	if !flag {
		return errors.New("Theere is no exist chunk.")
	}
	return nil
}

func (d *Database) DeleteChunk(x int, y int) (error) {
	flag := false
	// check there is chunk on x, y
	iter := d.Instance.NewIterator(nil, nil)
	for iter.Next() {
		// decode the value.
		var targetChunk Chunk
		err := decode(iter.Value(), &targetChunk)
		if err != nil {
			return errors.New("Decode the value failed.")
		}
		if targetChunk.Pos.X == x && targetChunk.Pos.Y == y {
			errDelete := d.Instance.Delete(iter.Key(), nil)
			if errDelete != nil {
				return errDelete
			}
			flag = true
			break
		}
	}
	if !flag {
		return errors.New("Theere is no exist chunk.")
	}
	return nil
}


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


