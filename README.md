# Weather by CEP

Este projeto é um servidor HTTP em Go que fornece informações de temperatura com base no CEP fornecido.

## Links Google Cloud Run

- [Link da API](https://weather-by-cep-767789551196.us-central1.run.app)
- [Link da API com exemplo](https://weather-by-cep-767789551196.us-central1.run.app/cep/01310930)

## Pré-requisitos

- [Go](https://golang.org/doc/install) 1.22 ou superior
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [WeatherAPI](https://www.weatherapi.com/) API key

## Configuração

1. Clone o repositório:

```sh
git clone https://github.com/seu-usuario/weather-by-cep.git
cd weather-by-cep
```

2. Modifique o arquivo `docker-compose.yml` para adicionar a chave da API:

```yaml
    environment:
      - WEATHERAPI_KEY={YOUR_WEATHERAPI_KEY}
```

Substitua `YOUR_WEATHERAPI_KEY` pela sua chave da API do WeatherAPI.

3. Modifique o arquivo `main_test.go` para adicionar a chave da API:

```go
	// Set up environment variables for testing
	os.Setenv("WEATHERAPI_KEY", "YOUR_WEATHERAPI_KEY")
```

Substitua `YOUR_WEATHERAPI_KEY` pela sua chave da API do WeatherAPI.

## Executando o Projeto com Docker Compose

1. Construa e inicie o contêiner:

```sh
docker-compose up --build
```

2. O servidor estará disponível em: http://localhost:8080

## Testando o Projeto

1. Execute os testes diretamente:

```sh
go test -v
```

## Endpoints

### Obter Temperatura por CEP

- **URL:** `/cep/{cep}`
- **Método:** `GET`

#### Exemplo de Requisição

```sh
curl -X GET http://localhost:8080/cep/01310930
```

#### Exemplo de Resposta

```json
{
  "temp_C": 25.0,
  "temp_F": 77.0,
  "temp_K": 298.0
}
```

## Arquivo HTTP para Testes

Você pode usar o arquivo [cep.http](./test/cep.http) para testar os endpoints diretamente no VS Code com a extensão [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client).

## Estrutura do Projeto

- `main.go`: Código principal do servidor.
- `main_test.go`: Testes automatizados para o servidor.
- `Dockerfile`: Dockerfile para construir a imagem do servidor.
- `docker-compose.yml`: Arquivo Docker Compose para configurar e executar o contêiner.
- `cep.http`: Arquivo HTTP para testar os endpoints.