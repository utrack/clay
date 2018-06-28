module github.com/utrack/clay/integration

replace github.com/utrack/clay/transport/middlewares v0.0.0 => ../transport/middlewares

replace github.com/utrack/clay/transport/server v0.0.0 => ../transport/server

replace github.com/utrack/clay/server v0.0.0 => ../server

replace github.com/utrack/clay/transport/v2 v2.0.0 => ../transport

replace github.com/utrack/clay/log v0.0.0 => ../log

require (
	github.com/go-chi/chi v0.0.0-20171222161133-e83ac2304db3
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/utrack/clay/transport/v2 v2.0.0
)
