package mongohelper

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindOne wraps the mongo.Database.Collection.FindOne() method
// It returns a SingleResult for one document in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// returned. If the filter does not match any documents, a SingleResult with an error set to
// ErrNoDocuments will be returned. If the filter matches multiple documents, one will be selected from the matched set.
func (l *Link) FindOne(database, collection string, filter interface{}, dest interface{}) error {
	if l.client == nil {
		return fmt.Errorf("mongohelper is not connected")
	}

	if dest == nil {
		return fmt.Errorf(`given "dest" is null`)
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.execTimeout())

	defer cancel()

	if filter == nil {
		filter = bson.M{}
	}

	rs := l.client.Database(database).Collection(collection).FindOne(ctx, filter, options.FindOne())

	if err := rs.Err(); err != nil {
		// If not connected, try once again, reconnecting. otherwise, just return/leave
		if !errors.Is(err, mongo.ErrClientDisconnected) {
			return err
		}

		if err := l.connect(); err != nil {
			return err
		}

		ctx2, cancel2 := context.WithTimeout(context.Background(), l.execTimeout())

		defer cancel2()

		rs = l.client.Database(database).Collection(collection).FindOne(ctx2, filter, options.FindOne())

		if err := rs.Err(); err != nil {
			return err
		}
	}

	if err := rs.Decode(dest); err != nil {
		return err
	}

	return nil
}
