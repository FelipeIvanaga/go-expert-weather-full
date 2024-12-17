package weatherusecase

import (
	"context"

	"github.com/felipeivanaga/go-expert-weather-ms/internal/infra/provider/cep"
	"github.com/felipeivanaga/go-expert-weather-ms/internal/infra/provider/weather"
	"github.com/felipeivanaga/go-expert-weather-ms/internal/internal_error"
)

type WeatherUsecase struct {
	cepProvider     cep.CepProvider
	weatherProvider weather.WeatherProvider
}

func NewWeatherUsecase(cepProvider cep.CepProvider, weatherProvider weather.WeatherProvider) WeatherCaseInterface {
	return &WeatherUsecase{
		cepProvider:     cepProvider,
		weatherProvider: weatherProvider,
	}
}

type WeatherCaseInterface interface {
	GetWeather(
		ctx context.Context,
		cep string) (*WeatherOutputDTO, *internal_error.InternalError)
}

func (w *WeatherUsecase) GetWeather(ctx context.Context, cep string) (*WeatherOutputDTO, *internal_error.InternalError) {
	cidade, err := w.cepProvider.GetCityName(ctx, cep)
	if err != nil {
		return nil, err
	}

	now, err := w.weatherProvider.GetWeatherWithCityName(ctx, cidade)
	if err != nil {
		return nil, err
	}

	return &WeatherOutputDTO{
		TempC: now.Celsius,
		TempK: now.Celsius + 273,
		TempF: now.Fahrenheit,
	}, nil
}

func getWeather() float32 {
	return 0.0
}

type WeatherOutputDTO struct {
	TempC float64 `json:"temp_C"`
	TempK float64 `json:"temp_K"`
	TempF float64 `json:"temp_F"`
}
