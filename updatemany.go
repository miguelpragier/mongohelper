package mongohelper

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateMany wraps the mongo.Database.Collection.UpdateMany() method
// It returns the number of matched records and an error
// The filter parameter must be a document containing query operators and can be used to select the documents to be
// updated. It cannot be nil. If the filter does not match any documents, the operation will succeed and an UpdateResult
// with a MatchedCount of 0 will be returned.
//
// The update parameter must be a document containing update operators
// (https://docs.mongodb.com/manual/reference/operator/update/) and can be used to specify the modifications to be made
// to the selected documents. It cannot be nil or empty.
func (l *Link) UpdateMany(database, collection string, filter, update interface{}) (int64, error) {
	if err := l.linkCheck("link.UpdateMany"); err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.execTimeout())

	defer cancel()

	rs, err := l.client.Database(database).Collection(collection).UpdateMany(ctx, filter, update, options.Update())

	if err != nil {
		// If not connected, try once again
		if errors.Is(err, mongo.ErrClientDisconnected) {
			if err = l.connect(); err != nil {
				return 0, err
			}

			ctx2, cancel2 := context.WithTimeout(context.Background(), l.execTimeout())

			defer cancel2()

			if rs, err = l.client.Database(database).Collection(collection).UpdateMany(ctx2, filter, update, options.Update()); err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	return rs.MatchedCount, nil
}
