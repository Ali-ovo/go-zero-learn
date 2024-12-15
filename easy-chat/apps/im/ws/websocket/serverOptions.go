package websocket

import "time"

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authentication
	patten string

	maxConnectionIdle time.Duration
}

func newServerOptions(opts ...ServerOptions) serverOption {

	o := serverOption{
		Authentication:    new(authentication),
		maxConnectionIdle: defaultMaxConnectionIdle,
		patten:            "/ws",
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func WithServerAuthentication(auth Authentication) ServerOptions {
	return func(opt *serverOption) {
		opt.Authentication = auth
	}
}

func WithServerPatten(patten string) ServerOptions {
	return func(opt *serverOption) {
		opt.patten = patten
	}
}

func WithServerMaxConnectionIdle(d time.Duration) ServerOptions {
	return func(opt *serverOption) {
		if d > 0 {
			opt.maxConnectionIdle = d
		}
	}
}