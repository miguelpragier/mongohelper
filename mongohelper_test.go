package mongohelper

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

const (
	testDB         = `mongohelper`
	testCollection = `testsuite`
)

var (
	testConnectionString = fmt.Sprintf("mongodb://127.0.0.1:27017/%s", testDB)
	mdb                  *Link
)

type testDocStruct struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name,omitempty"`
	TS   time.Time          `bson:"ts,omitempty"`
	N    int                `bson:"n,omitempty"`
}

func TestNew(t *testing.T) {
	if _m, err := New(testConnectionString, OptionsNew(10, 10, 10, 0, 0, false, true)); err != nil {
		t.Fatal(err)
	} else {
		mdb = _m
	}
}

func TestLink_InsertOne(t *testing.T) {
	x := testDocStruct{
		ID:   primitive.NewObjectID(),
		Name: "testing insertone",
		TS:   time.Now(),
		N:    125,
	}

	if s, err := mdb.InsertOne(testDB, testCollection, &x); err != nil {
		t.Error(err)
	} else {
		fmt.Println(s)
	}
}
