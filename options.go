package mongohelper

import (
	"fmt"
	"log"
	"time"
)

// Options contains all necessary parameters to connect and mantain database client
type Options struct {
	// appName helps to identify the source application on database logs
	appName string
	// URI with directives for mongodb connection
	connString string
	// printLogMessages if true allow the engine to print log messages
	printLogMessages bool
	// connTimeoutSeconds is quite obvious
	connTimeoutSeconds uint
	// execTimeoutSeconds equals to how much time the engine waits before return an error
	execTimeoutSeconds uint
	// InsistOnFail Is can't connect on first attempt, if should retry
	reconnectionInsistOnFail bool
	// SecondsBetweenAttempts Seconds btween each connection attempt
	reconnectionSecondsBetweenAttempts uint
	// ReconnectionAttemptsLimit maximum number of (re)connection attempts or 0 for infinite
	reconnectionAttemptsLimit uint
	// ReconnectionAttemptsLimitMinutes maximum time ( in minutes ) trying to (re)connect or 0 for infinite
	reconnectionAttemptsLimitMinutes uint
	// firstAttempt stores the time.Time when last connection succeeded. it restarts when a connection succeeds
	lastConnection time.Time
	// number of (re)connection attempts. it restarts when a connection succeeds
	attempts uint
}

// OptionsNew returns a pointer to mongohelper.Options instance.
// Why a pointer? In this case is because you can send nil instead, and the engine provides default values.
// connString, a well-formed URI for mongodb. Attention: is mandatory
// connectTimeoutInSeconds is quite obvious
// reconnectTimeInSeconds Seconds btween each connection attempt
// reconnecAttemptsLimit maximum number of (re)connection attempts or 0 for infinite
// reconnectAttemptsLimitMinutes maximum time ( in minutes ) trying to (re)connect or 0 for infinite
// insistOnFail If can't connect on first attempt, if should retry
// logMessages if true allow the engine to print out log messages to stdout
func OptionsNew(appName, connectionString string, connectTimeoutInSeconds, execTimeoutInSeconds, reconnectTimeInSeconds, reconnecAttemptsLimit, reconnectAttemptsLimitMinutes uint, insistOnFail, logMessages bool) *Options {
	logIfAllowed := func(msg string) {
		if logMessages {
			log.Println(msg)
		}
	}

	if appName == "" {
		logIfAllowed("empty application name")
	}

	if connectionString == "" {
		logIfAllowed("empty connection string")

		return nil
	}

	if connectTimeoutInSeconds < ConnectionTimeoutSecondsMin {
		logIfAllowed(fmt.Sprintf("value too low for connectTimeoutInSeconds: %d, when minimum allowed is %d; using default mongohelper.ConnectionTimeoutSecondsDefault: %d instead\n", connectTimeoutInSeconds, SecondsBetweenAttemptsMin, ConnectionTimeoutSecondsDefault))

		connectTimeoutInSeconds = ConnectionTimeoutSecondsDefault
	}

	if reconnectTimeInSeconds < SecondsBetweenAttemptsMin {
		logIfAllowed(fmt.Sprintf("value too low for reconnectTimeInSeconds: %d, when minimum allowed is %d; using default mongohelper.SecondsBetweenAttemptsMin: %d instead\n", reconnectTimeInSeconds, SecondsBetweenAttemptsMin, SecondsBetweenAttemptsMin))

		reconnectTimeInSeconds = SecondsBetweenAttemptsMin
	}

	if execTimeoutInSeconds < ExecutionTimeoutSecondsMin {
		logIfAllowed(fmt.Sprintf("value too low for execTimeoutInSeconds: %d, when minimum allowed is %d; using default mongohelper.ExecutionTimeoutSecondsDefault: %d instead\n", execTimeoutInSeconds, SecondsBetweenAttemptsMin, ExecutionTimeoutSecondsDefault))

		execTimeoutInSeconds = ExecutionTimeoutSecondsDefault
	}

	return &Options{
		connString:                         connectionString,
		connTimeoutSeconds:                 connectTimeoutInSeconds,
		execTimeoutSeconds:                 execTimeoutInSeconds,
		reconnectionInsistOnFail:           insistOnFail,
		reconnectionSecondsBetweenAttempts: reconnectTimeInSeconds,
		reconnectionAttemptsLimit:          reconnecAttemptsLimit,
		reconnectionAttemptsLimitMinutes:   reconnectAttemptsLimitMinutes,
		printLogMessages:                   logMessages,
	}
}
