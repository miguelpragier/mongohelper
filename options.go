package mongohelper

import (
	"log"
	"time"
)

// Link Options
type Options struct {
	// printLogMessages if true allow the engine to print log messages
	printLogMessages bool
	// connectionTimeoutInSeconds is quite obvious
	connectionTimeoutInSeconds uint
	// executionTimeoutInSeconds equals to how much time the engine waits before return an error
	executionTimeoutInSeconds uint
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
// connectTimeoutInSeconds is quite obvious
// reconnectTimeInSeconds Seconds btween each connection attempt
// reconnecAttemptsLimit maximum number of (re)connection attempts or 0 for infinite
// reconnectAttemptsLimitMinutes maximum time ( in minutes ) trying to (re)connect or 0 for infinite
// insistOnFail If can't connect on first attempt, if should retry
// logMessages if true allow the engine to print out log messages to stdout
func OptionsNew(connectTimeoutInSeconds, execTimeoutInSeconds uint, reconnectTimeInSeconds, reconnecAttemptsLimit, reconnectAttemptsLimitMinutes uint, insistOnFail, logMessages bool) *Options {
	if reconnectTimeInSeconds < SecondsBetweenAttemptsMinDefault {
		if logMessages {
			log.Printf("value too low for reconnectTimeInSeconds: %d, when minimum allowed is %d; using default mongohelper.SecondsBetweenAttemptsMinDefault: %d instead\n", reconnectTimeInSeconds, SecondsBetweenAttemptsMinDefault, SecondsBetweenAttemptsMinDefault)
		}

		reconnectTimeInSeconds = SecondsBetweenAttemptsMinDefault
	}

	if connectTimeoutInSeconds < ConnectionTimeoutInSecondsMinDefault {
		if logMessages {
			log.Printf("value too low for connectTimeoutInSeconds: %d, when minimum allowed is %d; using default mongohelper.ConnectionTimeoutInSecondsDefault: %d instead\n", connectTimeoutInSeconds, SecondsBetweenAttemptsMinDefault, ConnectionTimeoutInSecondsDefault)
		}

		connectTimeoutInSeconds = ConnectionTimeoutInSecondsDefault
	}

	if execTimeoutInSeconds < ExecutionTimeoutInSecondsMinDefault {
		if logMessages {
			log.Printf("value too low for execTimeoutInSeconds: %d, when minimum allowed is %d; using default mongohelper.ExecutionTimeoutInSecondsDefault: %d instead\n", execTimeoutInSeconds, SecondsBetweenAttemptsMinDefault, ExecutionTimeoutInSecondsDefault)
		}

		execTimeoutInSeconds = ExecutionTimeoutInSecondsDefault
	}

	return &Options{
		connectionTimeoutInSeconds:         connectTimeoutInSeconds,
		executionTimeoutInSeconds:          execTimeoutInSeconds,
		reconnectionInsistOnFail:           insistOnFail,
		reconnectionSecondsBetweenAttempts: reconnectTimeInSeconds,
		reconnectionAttemptsLimit:          reconnecAttemptsLimit,
		reconnectionAttemptsLimitMinutes:   reconnectAttemptsLimitMinutes,
		printLogMessages:                   logMessages,
	}
}
