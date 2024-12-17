package weathercontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/felipeivanaga/go-expert-weather-ms/internal/configuration/rest_err"
	weatherusecase "github.com/felipeivanaga/go-expert-weather-ms/internal/usecase/weather_usecase"
)

type WeatherController struct {
	weatherUsecase weatherusecase.WeatherCaseInterface
	OTELTracer     trace.Tracer
}

func NewWeatherController(weatherUsecase weatherusecase.WeatherCaseInterface, OTELTracer trace.Tracer) *WeatherController {
	return &WeatherController{
		weatherUsecase: weatherUsecase,
		OTELTracer:     OTELTracer,
	}
}

func (w *WeatherController) GetWeather(ctx *gin.Context) {
	carrier := propagation.HeaderCarrier(ctx.Request.Header)
	otelCtx := ctx.Request.Context()
	otelCtx = otel.GetTextMapPropagator().Extract(otelCtx, carrier)
	otelCtx, span := w.OTELTracer.Start(otelCtx, "weather-app-request")
	defer span.End()

	cep := ctx.Query("CEP")

	weatherData, err := w.weatherUsecase.GetWeather(otelCtx, cep)
	if err != nil {
		restErr := rest_err.ConvertError(err)

		ctx.JSON(restErr.Code, restErr)
		return
	}

	ctx.JSON(http.StatusOK, weatherData)
}
