package cep

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/felipeivanaga/go-expert-weather-ms/internal/internal_error"
)

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Erro        string `json:"erro"`
}

type ViaCepProvider struct {
	OTELTracer trace.Tracer
}

func NewViaCepProvider(OTELTracer trace.Tracer) CepProvider {
	return &ViaCepProvider{
		OTELTracer: OTELTracer,
	}
}

func (p *ViaCepProvider) GetCityName(ctx context.Context, cep string) (string, *internal_error.InternalError) {
	ctx, span := p.OTELTracer.Start(ctx, "viacep-request")
	defer span.End()

	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return "", internal_error.NewInternalServerError("unable to reach viacep service")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", internal_error.NewInternalServerError("unable to read response body of viacep")
	}

	var viaCepResp ViaCEPResponse
	err = json.Unmarshal(body, &viaCepResp)
	if err != nil {
		return "", internal_error.NewInternalServerError("unable to parse response body")
	}

	if viaCepResp.Erro == "true" {
		return "", internal_error.NewNotFoundError("can not find zipcode")
	}

	return viaCepResp.Localidade, nil
}
