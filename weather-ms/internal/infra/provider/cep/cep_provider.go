package cep

import (
	"context"

	"github.com/felipeivanaga/go-expert-weather-ms/internal/internal_error"
)

type CepProvider interface {
	GetCityName(ctx context.Context, cep string) (string, *internal_error.InternalError)
}
