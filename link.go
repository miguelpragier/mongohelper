package mongohelper

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

// Link is a concentrator wrapper for mongodb client
type Link struct {
	client           *mongo.Client
	connectionString string
	options          Options
}

// insistOnFail returns l.options.reconnectionInsistOnFail value
func (l Link) insistOnFail() bool {
	return l.options.reconnectionInsistOnFail
}

// canInsist checks if this engine can retry to connect database, considering the options rules
func (l Link) canInsist() bool {
	if l.options.reconnectionAttemptsLimit > 0 && l.options.attempts < l.options.reconnectionAttemptsLimit {
		return true
	}

	if l.options.reconnectionAttemptsLimitMinutes > 0 {
		expiration := l.options.lastConnection.Add(time.Duration(l.options.reconnectionAttemptsLimitMinutes) * time.Minute)

		if time.Now().After(expiration) {
			return true
		}
	}

	return false
}

// wait N seconds before next 9re)connection attempt
func (l Link) wait() {
	timeout := time.Duration(l.options.reconnectionSecondsBetweenAttempts) * time.Second

	time.Sleep(timeout)
}

// increment increments in one the connection attempt counter
func (l *Link) increment() {
	if l.options.reconnectionAttemptsLimit > 0 {
		l.options.attempts++
	}
}

// notifyConnection set attempts to zero and lastConnection to NOW
func (l *Link) notifyConnection() {
	if l.options.reconnectionAttemptsLimit > 0 {
		l.options.attempts = 0
	}

	if l.options.reconnectionAttemptsLimitMinutes > 0 {
		l.options.lastConnection = time.Now()
	}
}

// log print log message if allowed by programmer in options
func (l Link) log(routine, message string) {
	if l.options.printLogMessages {
		log.Printf("%s - mongohelper %s - %s\n", time.Now().Format(time.RFC3339), routine, message)
	}
}

// connect tries to conect database using the given options
func (l *Link) connect() error {
	var ctx context.Context

	// Context with timeout can't be used in loops, because they expire before the loop complete its job
	if l.insistOnFail() {
		ctx = context.Background()
	} else {
		timeout := time.Duration(l.options.connectionTimeoutInSeconds) * time.Second

		_ctx, cancel := context.WithTimeout(context.Background(), timeout)

		ctx = _ctx

		defer cancel()
	}

	for {
		var err error

		l.client, err = mongo.Connect(ctx, options.Client().ApplyURI(l.connectionString))

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

func (l Link) connTimeout() time.Duration {
	return time.Duration(l.options.connectionTimeoutInSeconds) * time.Second
}

func (l Link) execTimeout() time.Duration {
	return time.Duration(l.options.executionTimeoutInSeconds) * time.Second
}

// quickPing tries to reach the database in 10 seconds
func (l Link) quickPing() error {
	ctx, cancel := context.WithTimeout(context.Background(), l.connTimeout())

	defer cancel()

	return l.client.Ping(ctx, readpref.Primary())
}

// Collection returns a collection from the target database
func (l Link) Collection(database, collection string) (*mongo.Collection, error) {
	if l.client == nil {
		return nil, fmt.Errorf("use of uninitialized connection")
	}

	if err := l.quickPing(); err != nil {
		return nil, err
	}

	return l.client.Database(database).Collection(collection), nil
}
