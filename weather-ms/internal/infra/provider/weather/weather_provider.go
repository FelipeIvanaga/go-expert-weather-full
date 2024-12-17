package weather

import (
	"context"

	"github.com/felipeivanaga/go-expert-weather-ms/internal/internal_error"
)

type GetWeatherResponseDTO struct {
	Celsius    float64
	Fahrenheit float64
}

type WeatherProvider interface {
	GetWeatherWithCityName(ctx context.Context, city string) (*GetWeatherResponseDTO, *internal_error.InternalError)
}
