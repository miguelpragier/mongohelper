package mongohelper

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Find cs wraps the mongo.Database.Collection.Find() method
// It returns a Cursor over the matching documents in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select which documents are
// included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all documents.
func (l *Link) Find(database, collection string, filter interface{}, dest interface{}) error {
	if l.client == nil {
		return fmt.Errorf("mongohelper is not connected")
	}

	if dest == nil {
		return fmt.Errorf(`given "dest" is null`)
	}

	if filter == nil {
		filter = bson.M{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.execTimeout())

	defer cancel()

	rs, err := l.client.Database(database).Collection(collection).Find(ctx, filter, options.Find())

	if err != nil {
		// If not connected, try once again, reconnecting. otherwise, just return/leave
		if !errors.Is(err, mongo.ErrClientDisconnected) {
			return err
		}

		if err := l.connect(); err != nil {
			return err
		}

		ctx2, cancel2 := context.WithTimeout(context.Background(), l.execTimeout())

		defer cancel2()

		rs, err = l.client.Database(database).Collection(collection).Find(ctx2, filter, options.Find())

		if err != nil {
			return err
		}
	}

	return rs.All(context.TODO(), &dest)
}
