package mongohelper

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DeleteMany wraps the mongo.Database.Collection.DeleteMany() method
// It returns the number of affected records and an error
func (l *Link) DeleteMany(database, collection string, filter interface{}) (int64, error) {
	if l.client == nil {
		return 0, fmt.Errorf("mongohelper is not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.execTimeout())

	defer cancel()

	rs, err := l.client.Database(database).Collection(collection).DeleteMany(ctx, filter, options.Delete())

	if err != nil {
		// If not connected, try once again
		if errors.Is(err, mongo.ErrClientDisconnected) {
			if err = l.connect(); err != nil {
				return 0, err
			}

			ctx2, cancel2 := context.WithTimeout(context.Background(), l.execTimeout())

			defer cancel2()

			if rs, err = l.client.Database(database).Collection(collection).DeleteMany(ctx2, filter, options.Delete()); err != nil {
				return 0, err
			}
		}

		return 0, err
	}

	return rs.DeletedCount, nil
}
