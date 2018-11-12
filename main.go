package main

import (
	"fmt"
	"github.com/f26401004/LifeGamer_Database/database"
)

func main() {
	fmt.Println("Hello world.")
	db, err := database.InitDB("map")
	if err != nil {
		fmt.Println(err)
	}
	newChuckErr := db.NewChunk(1, 1)
	if newChuckErr != nil {
		fmt.Println(err)
	}
	deleteChunkErr := db.DeleteChunk(1, 1)
	if deleteChunkErr != nil {
		fmt.Println(deleteChunkErr)
	}
	/*
	db, err := leveldb.OpenFile("./dbtest", nil)
	if err != nil {
		fmt.Println("create LevelDB database failed!!", err)
	}
	fmt.Println("create LevelDB database success!!", err)
	
	err = db.Put([]byte("key1"), []byte("value1"), nil)
	
	data, err := db.Get([]byte("key1"), nil)
	fmt.Println("Key1's value: ", string(data))

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Println("Key: ", string(key))
		fmt.Println("Value: ", string(value))
	}
	iter.Release()
	err = iter.Error()
	*/

}
