# Mongo Helper
##### A simple wrapper and concentrator that helps to deal with connection and to recover from unexpected disconnection
Made over [MongoDB Official Driver](https://github.com/mongodb/mongo-go-driver)
<br>
![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg) 
![GitHub](https://img.shields.io/badge/goDoc-Yes!-blue.svg) 
[![Go Report Card](https://goreportcard.com/badge/github.com/miguelpragier/mongohelper?update)](https://goreportcard.com/report/github.com/miguelpragier/mongohelper) 
<br>
[Check the docs here](https://pkg.go.dev/github.com/miguelpragier/mongohelper?tab=doc)
<br>
### Examples
```golang
package main

import (
	"fmt"
	"github.com/go-acme/lego/log"
	"github.com/miguelpragier/mongohelper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	testDB         = `mongohelper`
	testCollection = `testsuite`
)

var testConnectionString = fmt.Sprintf("mongodb://127.0.0.1:27017/%s", testDB)

type testDocStruct struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name,omitempty"`
	TS   time.Time          `bson:"ts,omitempty"`
	N    int                `bson:"n,omitempty"`
}

func main() {
	const (
        appName = "my app"
		connTimeoutSeconds         = 10 // Time to wait for the first connection
		execTimeoutSeconds         = 10 // Time to wait for execution
		reconnTimeoutSeconds       = 10 // Time between reconnection attempts
		reconnAttemptsLimit        = 0  // Top limit for (re)connection attempts
		reconnAttemptsLimitMinutes = 5  // limit time trying to reconnect
		insistOnFailure            = false
		logDebugMessages           = true
	)

	var (
		mdb             mongohelper.Link
		lastInsertedOID primitive.ObjectID
	)

	log.Println("Connecting db...")

	if _m, err := mongohelper.New(mongohelper.OptionsNew(appName,testConnectionString,connTimeoutSeconds, execTimeoutSeconds, reconnTimeoutSeconds, reconnAttemptsLimit, reconnAttemptsLimitMinutes, insistOnFailure, logDebugMessages)); err != nil {
		log.Fatal(err)
	} else {
		mdb = *_m
	}

	x := testDocStruct{
		ID:   primitive.NewObjectID(),
		Name: "testing insertone",
		TS:   time.Now(),
		N:    125,
	}

	log.Println("Inserting one doc")

	if s, err := mdb.InsertOne(testDB, testCollection, &x); err != nil {
		log.Println(err)
	} else {
		log.Printf("inserted oID: %s\n", s)
	}

	var a []interface{}

	for i := 0; i < 100; i++ {
		x := testDocStruct{
			ID:   primitive.NewObjectID(),
			Name: fmt.Sprintf("test #%d", i),
			TS:   time.Now().Add(time.Duration(i) * time.Hour),
			N:    i * 2,
		}

		a = append(a, x)

		lastInsertedOID = x.ID
	}

	log.Println("Inserting many doc")

	if ioids, err := mdb.InsertMany(testDB, testCollection, a); err != nil {
		log.Println(err)
	} else {
		log.Printf("inserted oIDs: %d\n", len(ioids))
	}

	log.Println("testing .CountDocs() with empty filter")

	if n, err := mdb.CountDocs(testDB, testCollection, bson.M{}); err != nil {
		log.Println(err)
	} else {
		log.Printf("found %d docs\n", n)
	}

	log.Println("testing .CountDocs() with ObjectId filter")

	if n, err := mdb.CountDocs(testDB, testCollection, bson.M{"_id": lastInsertedOID}); err != nil {
		log.Println(err)
	} else {
		log.Printf("searched oID %s and found %d docs\n", lastInsertedOID.Hex(), n)
	}

	log.Println("updating one doc")

	if n, err := mdb.UpdateOne(testDB, testCollection, bson.M{"_id": lastInsertedOID}, bson.M{"$set": bson.M{"xyz": "abc"}}); err != nil {
		log.Println(err)
	} else {
		log.Printf("%d doc updated\n", n)
	}

	var docsToDelete int64

	log.Println("updating many docs")

	if n, err := mdb.UpdateMany(testDB, testCollection, bson.M{"n": 8}, bson.M{"$set": bson.M{"name": "delete me!"}}); err != nil {
		log.Println(err)
	} else {
		docsToDelete = n

		log.Printf("updated docs flagged to be deleted: %d\n", docsToDelete)
	}

	log.Println("deleting one doc")

	if n, err := mdb.DeleteOne(testDB, testCollection, bson.M{"xyz": "abc"}); err != nil {
		log.Println(err)
	} else {
		log.Printf("delete %d document %s\n", n, lastInsertedOID.Hex())
	}

	log.Println("deleting many docs")

	if n, err := mdb.DeleteMany(testDB, testCollection, bson.M{"name": "delete me!"}); err != nil {
		log.Println(err)
	} else {
		log.Printf("deleted docs %d\n", n)

		if n != docsToDelete {
			log.Printf("expected deleteable docs: %d\n", docsToDelete)
		}
	}
}
```