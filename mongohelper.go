package mongohelper

import (
	"fmt"
)

const (
	// SecondsBetweenAttemptsMin When retrying connection, minimum time betwwen attempts
	SecondsBetweenAttemptsMin uint = 5
	// ConnectionTimeoutSecondsDefault limits the time waiting from a connection request
	ConnectionTimeoutSecondsDefault uint = 30
	// ConnectionTimeoutSecondsMin limits the minimum time waiting from a connection request
	ConnectionTimeoutSecondsMin uint = 3
	// ExecutionTimeoutSecondsDefault limits the time waiting from an execution request
	ExecutionTimeoutSecondsDefault uint = 10
	// ExecutionTimeoutSecondsMin limits the minimum time waiting from an execution request
	ExecutionTimeoutSecondsMin uint = 1
)

// New returns an instance of mongohelper, ugins given options
// It the connection conditions are ok, it comes alerady connected and tested with .Ping()
// You may prefer to create the options with .OptionsNew() function
func New(opts *Options) (*Link, error) {
	if opts == nil {
		return nil, fmt.Errorf("uninitialized options")
	}

	link := Link{
		options: *opts,
	}

	if err := link.connect(); err != nil {
		return nil, err
	}

	return &link, nil
}
