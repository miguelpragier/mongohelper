package mongohelper

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne wraps the mongo.Database.Collection.InsertOne() method
// It returns the generated ObjectId and an error
func (l *Link) InsertOne(database, collection string, document interface{}) (string, error) {
	if l.client == nil {
		return ``, fmt.Errorf("mongohelper is not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.execTimeout())

	defer cancel()

	rs, err := l.client.Database(database).Collection(collection).InsertOne(ctx, document, options.InsertOne())

	if err != nil {
		// If not connected, try once again
		if errors.Is(err, mongo.ErrClientDisconnected) {
			if err = l.connect(); err != nil {
				return ``, err
			}

			ctx2, cancel2 := context.WithTimeout(context.Background(), l.execTimeout())

			defer cancel2()

			if rs, err = l.client.Database(database).Collection(collection).InsertOne(ctx2, document, options.InsertOne()); err != nil {
				return ``, err
			}
		} else {
			return ``, err
		}
	}

	return fmt.Sprintf("%v", rs.InsertedID), nil
}