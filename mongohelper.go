package mongohelper

import (
	"fmt"
)

const (
	// SecondsBetweenAttemptsMinDefault When retrying connection, minimum time betwwen attempts
	SecondsBetweenAttemptsMinDefault uint = 5
	// ConnectionTimeoutInSecondsDefault limits the time waiting from a connection request
	ConnectionTimeoutInSecondsDefault uint = 30
	// ConnectionTimeoutInSecondsMinDefault limits the minimum time waiting from a connection request
	ConnectionTimeoutInSecondsMinDefault uint = 3
	// ExecutionTimeoutInSecondsDefault limits the time waiting from an exection request
	ExecutionTimeoutInSecondsDefault uint = 10
	// ExecutionTimeoutInSecondsMinDefault limits the minimum time waiting from an exection request
	ExecutionTimeoutInSecondsMinDefault uint = 1
)

func New(connectionString string, opts *Options) (*Link, error) {
	if connectionString == "" {
		return nil, fmt.Errorf("empty connection string")
	}

	if opts == nil {
		opts = &Options{}
	}

	link := Link{
		options:          *opts,
		connectionString: connectionString,
	}

	if err := link.connect(); err != nil {
		return nil, err
	}

	return &link, nil
}
