package weathercontroller

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/felipeivanaga/go-expert-gateway-ms/internal/configuration/rest_err"
)

type WeatherRequestBody struct {
	Cep string `json:"cep"`
}

type WeatherController struct {
	OTELTracer trace.Tracer
}

func NewWeatherController(OTELTracer trace.Tracer) *WeatherController {
	return &WeatherController{
		OTELTracer: OTELTracer,
	}
}

func (w *WeatherController) GetWeather(ctx *gin.Context) {
	carrier := propagation.HeaderCarrier(ctx.Request.Header)
	otelCtx := otel.GetTextMapPropagator().Extract(ctx, carrier)
	otelCtx, span := w.OTELTracer.Start(otelCtx, "gateway-app-request")
	defer span.End()

	var requestBody WeatherRequestBody
	if err := ctx.BindJSON(&requestBody); err != nil {
		errRest := rest_err.NewInternalServerError("unable to parse request body")
		ctx.JSON(errRest.Code, errRest)
		return
	}

	re := regexp.MustCompile("^[0-9]{8}$")

	if !re.MatchString(requestBody.Cep) {
		errRest := rest_err.NewUnprocessableEntityError("invalid zipcode")
		ctx.JSON(errRest.Code, errRest)
		return
	}

	req, _ := http.NewRequestWithContext(otelCtx, "GET", "http://weather-ms:8081/weather?CEP="+requestBody.Cep, nil)

	otel.GetTextMapPropagator().Inject(otelCtx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errRest := rest_err.NewInternalServerError("unable to reach weather-ms")
		ctx.JSON(errRest.Code, errRest)
		return
	}

	defer resp.Body.Close()
	ctx.DataFromReader(resp.StatusCode, resp.ContentLength, "application/json", resp.Body, nil)
}
