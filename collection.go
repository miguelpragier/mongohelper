package mongohelper

import "go.mongodb.org/mongo-driver/mongo"

// Collection returns a collection from the target database
func (l Link) Collection(database, collection string) (*mongo.Collection, error) {
	if err := l.linkCheck("link.Collection"); err != nil {
		return nil, err
	}

	return l.client.Database(database).Collection(collection), nil
}
