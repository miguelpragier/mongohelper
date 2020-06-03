package mongohelper

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CountDocs wraps the mongo.Database.Collection.CountDocuments() method
// It returns the number of documents in the collection and an error
//
// The filter parameter must be a document and can be used to select which documents contribute to the count. It
// cannot be nil. An empty document (e.g. bson.D{}) should be used to count all documents in the collection. This will
// result in a full collection scan.
func (l *Link) CountDocs(database, collection string, filter interface{}) (int64, error) {
	if err := l.linkCheck("link.CountDocs"); err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.execTimeout())

	defer cancel()

	n, err := l.client.Database(database).Collection(collection).CountDocuments(ctx, filter, options.Count())

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
		} else {
			return 0, err
		}
	}

	return n, nil
}
