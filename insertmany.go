package mongohelper

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertMany wraps the mongo.Database.Collection.InsertMany() method
// It returns an array with generated ObjectIDs and an error
func (l *Link) InsertMany(database, collection string, document []interface{}) ([]string, error) {
	if err := l.linkCheck("link.InsertMany"); err != nil {
		return []string{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.execTimeout())

	defer cancel()

	rs, err := l.client.Database(database).Collection(collection).InsertMany(ctx, document, options.InsertMany())

	if err != nil {
		// If not connected, try once again
		if errors.Is(err, mongo.ErrClientDisconnected) {
			if err = l.connect(); err != nil {
				return []string{}, err
			}

			ctx2, cancel2 := context.WithTimeout(context.Background(), l.execTimeout())

			defer cancel2()

			if rs, err = l.client.Database(database).Collection(collection).InsertMany(ctx2, document, options.InsertMany()); err != nil {
				return []string{}, err
			}
		} else {
			return []string{}, err
		}
	}

	var oidHex []string

	for _, o := range rs.InsertedIDs {
		if oid, ok := o.(primitive.ObjectID); ok {
			oidHex = append(oidHex, oid.Hex())
		}
	}

	return oidHex, nil
}
