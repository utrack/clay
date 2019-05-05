package server

import (
	"net"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/soheilhy/cmux"
)

const (
	listenRetryWait     = 500 * time.Millisecond
	listenRetryDuration = 10 * time.Second
)

type listenerSet struct {
	mainListener cmux.CMux // nil or CMux. If nil - don't listen
	HTTP         net.Listener
	GRPC         net.Listener
}

func newListenerSet(opts *serverOpts) (*listenerSet, error) {
	liSet := &listenerSet{}
	var err error

	liSet.GRPC, err = newListener(opts.Host, opts.RPCPort)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create main listener")
	}

	if opts.RPCPort == opts.HTTPPort {
		mux := cmux.New(liSet.GRPC)
		liSet.GRPC = mux.Match(cmux.HTTP2())
		liSet.HTTP = mux.Match(cmux.Any())
		liSet.mainListener = mux
	} else {
		liSet.HTTP, err = newListener(opts.Host, opts.HTTPPort)
	}
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create HTTP listener")
	}

	return liSet, nil
}

// newListener start net.Listener on a port.
// It keeps retrying if port is already in use.
func newListener(host string, port int) (net.Listener, error) {
	var listener net.Listener
	var err error
	start := time.Now()
	for time.Since(start) < listenRetryDuration {
		listener, err = net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
		if err == nil {
			return listener, nil
		}
		time.Sleep(listenRetryWait)
	}
	return nil, err
}
