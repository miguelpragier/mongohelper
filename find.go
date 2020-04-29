package mongohelper

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Find cs wraps the mongo.Database.Collection.Find() method
// It returns a Cursor over the matching documents in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select which documents are
// included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all documents.
func (l *Link) Find(database, collection string, filter interface{}) (int64, error) {
	if l.client == nil {
		return 0, fmt.Errorf("mongohelper is not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.execTimeout())

	defer cancel()

	n, err := l.client.Database(database).Collection(collection).Find(ctx, filter, options.Find())

	if err != nil {
		// If not connected, try once again
		if errors.Is(err, mongo.ErrClientDisconnected) {
			if err = l.connect(); err != nil {
				return 0, err
			}

			ctx2, cancel2 := context.WithTimeout(context.Background(), l.execTimeout())

			defer cancel2()

			if n, err = l.client.Database(database).Collection(collection).CountDocuments(ctx2, filter, options.Count()); err != nil {
				return 0, err
			}
		}

		return 0, err
	}

	return n, nil
}
