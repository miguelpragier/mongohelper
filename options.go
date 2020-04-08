package mongodbhelper

import "time"

// Link Options
type Options struct {
	// printLogMessages if true allow the engine to print log messages
	printLogMessages bool
	// connectionTimeoutInSeconds is quite obvious
	connectionTimeoutInSeconds uint
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
// insistOnFail If can't connect on first attempt, if should retry
// reconnectTimeInSeconds Seconds btween each connection attempt
// reconnecAttemptsLimit maximum number of (re)connection attempts or 0 for infinite
// reconnectAttemptsLimitMinutes maximum time ( in minutes ) trying to (re)connect or 0 for infinite
// logMessages if true allow the engine to print out log messages to stdout
func OptionsNew(connectTimeoutInSeconds uint, insistOnFail bool, reconnectTimeInSeconds, reconnecAttemptsLimit, reconnectAttemptsLimitMinutes uint, logMessages bool) *Options {
	if reconnectTimeInSeconds < SecondsBetweenAttemptsMinDefault {
		reconnectTimeInSeconds = SecondsBetweenAttemptsMinDefault
	}

	if connectTimeoutInSeconds == 0 {
		connectTimeoutInSeconds = ConnectionTimeoutInSecondsDefault
	}

	return &Options{
		connectionTimeoutInSeconds:         connectTimeoutInSeconds,
		reconnectionInsistOnFail:           insistOnFail,
		reconnectionSecondsBetweenAttempts: reconnectTimeInSeconds,
		reconnectionAttemptsLimit:          reconnecAttemptsLimit,
		reconnectionAttemptsLimitMinutes:   reconnectAttemptsLimitMinutes,
		printLogMessages:                   logMessages,
	}
}

