package mongohelper

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertMany wraps the mongo.Database.Collection.InsertMany() method
// It returns an array with generated ObjectIDs and an error
func (l *Link) InsertMany(database, collection string, document []interface{}) (string, error) {
	if l.client == nil {
		return ``, fmt.Errorf("mongohelper is not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.execTimeout())

	defer cancel()

	rs, err := l.client.Database(database).Collection(collection).InsertMany(ctx, document, options.InsertMany())

	if err != nil {
		// If not connected, try once again
		if errors.Is(err, mongo.ErrClientDisconnected) {
			if err = l.connect(); err != nil {
				return ``, err
			}

			ctx2, cancel2 := context.WithTimeout(context.Background(), l.execTimeout())

			defer cancel2()

			if rs, err = l.client.Database(database).Collection(collection).InsertMany(ctx2, document, options.InsertMany()); err != nil {
				return ``, err
			}
		} else {
			return ``, err
		}
	}

	return fmt.Sprintf("%v", rs.InsertedIDs), nil
}