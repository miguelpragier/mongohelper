package mongohelper

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func (l Link) defineLink(opts *options.ClientOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), l.connTimeout())

	defer cancel()

	c, err := mongo.Connect(ctx, opts)

	if err != nil {
		l.log("link.defineLink.mongo.Connect", err.Error())

		return err
	}

	l.client = c

	return nil
}

func (l Link) ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), l.connTimeout())

	defer cancel()

	err := l.client.Ping(ctx, readpref.Primary())

	if err != nil {
		l.log("link.ping", err.Error())
	}

	return err
}

// connect tries to conect database using the given options
func (l *Link) connect() error {
	opts := options.Client().ApplyURI(l.connectionString())
	opts.SetConnectTimeout(l.connTimeout())
	opts.SetMaxConnIdleTime(8 * time.Hour)
	opts.SetSocketTimeout(l.execTimeout())
	opts.SetMinPoolSize(10)
	opts.SetAppName(l.appName())

	// It's not possible to restore from errors in options validation
	if err := l.defineLink(opts); err != nil {
		return err
	}

	for {
		err := l.ping()

		if err != nil {
			l.log("link.mongo.Ping", err.Error())

			if l.insistOnFail() && l.canInsist() {
				l.wait()
				l.increment()
				continue
			} else {
				return err
			}
		}

		l.notifyConnection()

		return nil
	}
}
