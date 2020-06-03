package mongohelper

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Link is a concentrator wrapper for mongodb client
type Link struct {
	client  *mongo.Client
	options Options
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

// wait N seconds before next (9)re)connection attempt
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

	if l.options.printLogMessages {
		l.log("link.notifyConnection", "mongodb connected")
	}
}

// log print log message if allowed by programmer in options
func (l Link) log(routine, message string) {
	if l.options.printLogMessages {
		log.Printf("%s - mongohelper %s - %s\n", time.Now().Format(time.RFC3339), routine, message)
	}
}

func (l Link) appName() string {
	return l.options.appName
}

func (l Link) connectionString() string {
	return l.options.connString
}

func (l Link) connTimeout() time.Duration {
	return time.Duration(l.options.connTimeoutSeconds) * time.Second
}

func (l Link) execTimeout() time.Duration {
	return time.Duration(l.options.execTimeoutSeconds) * time.Second
}
