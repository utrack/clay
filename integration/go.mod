module github.com/utrack/clay/integration

replace github.com/utrack/clay v1.0.2 => github.com/utrack/clay/v2 v2.1.0

replace github.com/utrack/clay/v2 v2.1.0 => ../

require (
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/utrack/clay/v2 v2.1.0
)
