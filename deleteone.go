package mongohelper

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DeleteOne wraps the mongo.Database.Collection.DeleteOne() method
// It returns the number of affected records and an error
func (l *Link) DeleteOne(database, collection string, filter interface{}) (int64, error) {
	if err := l.linkCheck("link.DeleteOne"); err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.execTimeout())

	defer cancel()

	rs, err := l.client.Database(database).Collection(collection).DeleteOne(ctx, filter, options.Delete())

	if err != nil {
		// If not connected, try once again
		if errors.Is(err, mongo.ErrClientDisconnected) {
			if err = l.connect(); err != nil {
				return 0, err
			}

			ctx2, cancel2 := context.WithTimeout(context.Background(), l.execTimeout())

			defer cancel2()

			if rs, err = l.client.Database(database).Collection(collection).DeleteOne(ctx2, filter, options.Delete()); err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	return rs.DeletedCount, nil
}
