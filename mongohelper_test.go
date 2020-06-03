package mongohelper

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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
	if _m, err := New(OptionsNew("mongohelpertest", testConnectionString, 10, 10, 10, 0, 0, false, true)); err != nil {
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

var lastInsertedOID primitive.ObjectID

func TestLink_InsertMany(t *testing.T) {
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

	if ioids, err := mdb.InsertMany(testDB, testCollection, a); err != nil {
		t.Error(err)
	} else {
		t.Logf("inserted oIDs: %d\n", len(ioids))
	}
}

func TestLink_CountDocs(t *testing.T) {
	t.Log("testing with empty filter")

	if n, err := mdb.CountDocs(testDB, testCollection, bson.M{}); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("found %d docs\n", n)
	}

	t.Log("testing with oID filter")

	if n, err := mdb.CountDocs(testDB, testCollection, bson.M{"_id": lastInsertedOID}); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("searched oID %s and found %d docs\n", lastInsertedOID.Hex(), n)
	}
}

func TestLink_UpdateOne(t *testing.T) {
	if n, err := mdb.UpdateOne(testDB, testCollection, bson.M{"_id": lastInsertedOID}, bson.M{"$set": bson.M{"xyz": "abc"}}); err != nil {
		t.Error(err)
	} else {
		docsToDelete = n

		fmt.Printf("updated docs flagged to be deleted: %d\n", docsToDelete)
	}
}

var docsToDelete int64

func TestLink_UpdateMany(t *testing.T) {
	if n, err := mdb.UpdateMany(testDB, testCollection, bson.M{"n": 8}, bson.M{"$set": bson.M{"name": "delete me!"}}); err != nil {
		t.Error(err)
	} else {
		docsToDelete = n

		fmt.Printf("updated docs flagged to be deleted: %d\n", docsToDelete)
	}
}

func TestLink_DeleteOne(t *testing.T) {
	if n, err := mdb.DeleteOne(testDB, testCollection, bson.M{"xyz": "abc"}); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("delete %d document %s\n", n, lastInsertedOID.Hex())
	}
}

func TestLink_DeleteMany(t *testing.T) {
	if n, err := mdb.DeleteMany(testDB, testCollection, bson.M{"name": "delete me!"}); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("deleted docs %d\n", n)

		if n != docsToDelete {
			t.Errorf("expected deleteable docs: %d\n", docsToDelete)
		}
	}
}
