package mongohelper

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func (l Link) ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), l.connTimeout())

	defer cancel()

	return l.client.Ping(ctx, readpref.Primary())
}

// connect tries to conect database using the given options
func (l *Link) connect() error {
	var ctx context.Context

	// Context with timeout can't be used in loops, because they expire before the loop complete its job
	if l.insistOnFail() {
		ctx = context.Background()
	} else {
		timeout := time.Duration(l.options.connTimeoutSeconds) * time.Second

		_ctx, cancel := context.WithTimeout(context.Background(), timeout)

		ctx = _ctx

		defer cancel()
	}

	for {
		var err error

		opts := options.Client().ApplyURI(l.connectionString())
		opts.SetConnectTimeout(l.connTimeout())
		opts.SetMaxConnIdleTime(8 * time.Hour)
		opts.SetSocketTimeout(l.execTimeout())
		opts.SetMinPoolSize(10)
		opts.SetAppName(l.appName())

		l.client, err = mongo.Connect(ctx, opts)

		if err != nil {
			l.log("mongo.Connect", err.Error())
		}

		err = l.client.Ping(context.Background(), readpref.Primary())

		if err != nil {
			l.log("mongo.Ping", err.Error())
		} else {
			l.notifyConnection()

			return nil
		}

		if l.insistOnFail() {
			if l.canInsist() {
				l.wait()
				l.increment()
				continue
			}
		}

		return err
	}
}
