package weather

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel/trace"

	"github.com/felipeivanaga/go-expert-weather-ms/internal/internal_error"
)

type wheaterapiResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type WeatherapiProvider struct {
	token      string
	OTELTracer trace.Tracer
}

func NewWeatherapiProvider(OTELTracer trace.Tracer, token string) WeatherProvider {
	return &WeatherapiProvider{
		token:      token,
		OTELTracer: OTELTracer,
	}
}

func (p *WeatherapiProvider) GetWeatherWithCityName(ctx context.Context, city string) (*GetWeatherResponseDTO, *internal_error.InternalError) {
	ctx, span := p.OTELTracer.Start(ctx, "weatherapi-request")
	defer span.End()

	params := url.Values{}
	params.Add("key", p.token)
	params.Add("q", city)
	params.Add("aqi", "no")
	resp, err := http.Get("https://api.weatherapi.com/v1/current.json?" + params.Encode())
	if err != nil {
		return nil, internal_error.NewInternalServerError("unable to reach weatherapi service")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, internal_error.NewInternalServerError("unable to reach weatherapi service")
	}

	var wheaterapiResp wheaterapiResponse
	json.Unmarshal(body, &wheaterapiResp)

	return &GetWeatherResponseDTO{
		Celsius:    wheaterapiResp.Current.TempC,
		Fahrenheit: wheaterapiResp.Current.TempF,
	}, nil
}
