module github.com/utrack/clay/integration

replace github.com/utrack/clay/transport/middlewares v0.0.0 => ../transport/middlewares

replace github.com/utrack/clay/transport/server v0.0.0 => ../transport/server

replace github.com/utrack/clay/server v0.0.0 => ../server

replace github.com/utrack/clay/transport/v2 v2.0.0 => ../transport

replace github.com/utrack/clay/log v0.0.0 => ../log

require (
	github.com/utrack/clay/log v0.0.0
	github.com/utrack/clay/transport/middlewares v0.0.0
	github.com/utrack/clay/transport/server v0.0.0
)
