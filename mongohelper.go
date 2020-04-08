package mongohelper

import (
	"fmt"
)

const (
	// SecondsBetweenAttemptsMinDefault When retrying connection, minimum time betwwen attempts
	SecondsBetweenAttemptsMinDefault uint = 5
	// ConnectionTimeoutInSecondsDefault limits the time waiting from a connection request
	ConnectionTimeoutInSecondsDefault uint = 30
)

func New(connectionString string, opts *Options) (Link, error) {
	if connectionString == "" {
		return Link{}, fmt.Errorf("empty connection string")
	}

	if opts == nil {
		opts = &Options{}
	}

	link := Link{
		options:          *opts,
		connectionString: connectionString,
	}

	if err := link.connect(); err != nil {
		return Link{}, err
	}

	return link, nil
}
