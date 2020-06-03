package mongohelper

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne wraps the mongo.Database.Collection.InsertOne() method
// It returns the generated ObjectId and an error
func (l *Link) InsertOne(database, collection string, document interface{}) (string, error) {
	if err := l.linkCheck("link.InsertOne"); err != nil {
		return "", err
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

	oidHex := ""

	if oid, ok := rs.InsertedID.(primitive.ObjectID); ok {
		oidHex = oid.Hex()
	}

	return oidHex, nil
}
