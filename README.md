# Weather by CEP (OTel)

Este projeto é um servidor HTTP em Go que fornece informações de temperatura com base no CEP fornecido utilizando OpenTelemetry

## Pré-requisitos

- [Go](https://golang.org/doc/install) 1.22 ou superior
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [WeatherAPI](https://www.weatherapi.com/) API key

## Configuração

1. Clone o repositório:

```sh
git clone https://github.com/seu-usuario/weather-by-cep-otel.git
cd weather-by-cep-otel
```

2. Modifique o arquivo `service-b/docker-compose.yml` para adicionar a chave da API:

```yaml
    environment:
      - WEATHERAPI_KEY={YOUR_WEATHERAPI_KEY}
```

Substitua `YOUR_WEATHERAPI_KEY` pela sua chave da API do WeatherAPI.

## Executando o Projeto com Docker Compose

1. Construa e inicie o contêiner:

```sh
docker-compose up --build
```

2. Endereço dos serviços
- O servidor-a estará disponível em: http://localhost:8080
- O servidor-b estará disponível em: http://localhost:8081

## Endpoints

### Obter Temperatura por CEP

- **URL:** `/`
- **Método:** `POST`
- **Body:** `{"cep": "your_cep"}`

#### Exemplo de Requisição

```sh
curl -X POST http://localhost:8080 -H "Content-Type: application/json" -d '{"cep": "01310930"}'
```

#### Exemplo de Resposta

```json
{
  "city": "São Paulo",
  "temp_C": 25.0,
  "temp_F": 77.0,
  "temp_K": 298.0
}
```
