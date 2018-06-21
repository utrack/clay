package transport

import (
	"github.com/go-openapi/spec"
)

type SwaggerOption func(swagger *spec.Swagger)

func WithHost(host string) SwaggerOption {
	return func(swagger *spec.Swagger) {
		swagger.Host = host
	}
}

func WithVersion(version string) SwaggerOption {
	return func(swagger *spec.Swagger) {
		if swagger.Info == nil {
			swagger.Info = &spec.Info{}
		}
		swagger.Info.Version = version
	}
}

func WithTitle(title string) SwaggerOption {
	return func(swagger *spec.Swagger) {
		if swagger.Info == nil {
			swagger.Info = &spec.Info{}
		}
		swagger.Info.Title = title
	}
}

func WithDescription(desc string) SwaggerOption {
	return func(swagger *spec.Swagger) {
		if swagger.Info == nil {
			swagger.Info = &spec.Info{}
		}
		swagger.Info.Description = desc
	}
}

func WithSecurityDefinitions(secDef spec.SecurityDefinitions) SwaggerOption {
	return func(swagger *spec.Swagger) {
		swagger.SecurityDefinitions = secDef
	}
}