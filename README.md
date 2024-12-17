# Desafio técnico Go Expert FullCycle

## Configuração

Copiar o arquivo `.env.exemplo` com o nome `.env`

## Rodando localmente

- Rode a imagem `docker compose up --build`
- Acessar o Zipkin no `http://localhost:9411`

## Exemplo de funcionamento

```bash
curl --location 'http://localhost:8080/weather' \
      --header 'Content-Type: application/json' \
      --data '{
          "cep": "81900550"
      }'
```
