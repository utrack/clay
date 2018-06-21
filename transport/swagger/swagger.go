package swagger

import (
	"github.com/go-openapi/spec"
)

type Option func(swagger *spec.Swagger)

func WithHost(host string) Option {
	return func(swagger *spec.Swagger) {
		swagger.Host = host
	}
}

func WithVersion(version string) Option {
	return func(swagger *spec.Swagger) {
		if swagger.Info == nil {
			swagger.Info = &spec.Info{}
		}
		swagger.Info.Version = version
	}
}

func WithTitle(title string) Option {
	return func(swagger *spec.Swagger) {
		if swagger.Info == nil {
			swagger.Info = &spec.Info{}
		}
		swagger.Info.Title = title
	}
}

func WithDescription(desc string) Option {
	return func(swagger *spec.Swagger) {
		if swagger.Info == nil {
			swagger.Info = &spec.Info{}
		}
		swagger.Info.Description = desc
	}
}

func WithSecurityDefinitions(secDef spec.SecurityDefinitions) Option {
	return func(swagger *spec.Swagger) {
		swagger.SecurityDefinitions = secDef
	}
}
