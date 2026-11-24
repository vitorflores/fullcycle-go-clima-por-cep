# Clima por CEP

Serviço em Go que recebe um CEP, descobre a cidade pelo [ViaCEP](https://viacep.com.br/) e devolve a temperatura atual em Celsius, Fahrenheit e Kelvin usando a [Open-Meteo](https://open-meteo.com/), que é gratuita e não precisa de chave de API. O projeto foi feito para o desafio de deploy no Google Cloud Run da pós Go Expert.

## Acesso no Cloud Run

URL do serviço: **https://clima-por-cep-326680495202.us-central1.run.app**

Exemplo de chamada:

```bash
curl "https://clima-por-cep-326680495202.us-central1.run.app/?cep=01001000"
```

## Como funciona a API

Existe uma única rota, `GET /`, que recebe o CEP pelo parâmetro `cep`.

| Situação | Status | Resposta |
| --- | --- | --- |
| Deu tudo certo | 200 | `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}` |
| CEP sem 8 dígitos ou com letras, traço etc. | 422 | `invalid zipcode` |
| CEP com formato certo, mas que não existe | 404 | `can not find zipcode` |
| Alguma das APIs externas falhou | 500 | `internal server error` |

Por baixo dos panos são três chamadas: o ViaCEP devolve a cidade e o estado, a geocodificação da Open-Meteo transforma isso em latitude e longitude, e a previsão da Open-Meteo devolve a temperatura atual nesse ponto. Como existem várias cidades com o mesmo nome no Brasil (tem Bom Jesus no Piauí, no Rio Grande do Sul, em Santa Catarina...), o serviço escolhe o resultado cujo estado bate com o que veio do ViaCEP.

As conversões seguem as fórmulas pedidas no desafio: `F = C * 1.8 + 32` e `K = C + 273`. Os valores saem arredondados em duas casas decimais.

## Pré-requisitos

- Docker com o Compose (ou Go 1.27 se preferir rodar sem container)

Não é preciso criar conta nem chave em lugar nenhum.

## Configuração

A única configuração é a porta, definida pela variável `PORT`. Se ela não for informada, a porta padrão é a `8080`. O Cloud Run preenche essa variável sozinho. Para mudar a porta localmente, dá para copiar o `.env.example` para `.env` e ajustar o valor.

## Rodando com Docker

Suba a aplicação:

```bash
docker compose up --build app
```

Depois é só testar:

```bash
curl "http://localhost:8080/?cep=01001000"
curl "http://localhost:8080/?cep=123"
curl "http://localhost:8080/?cep=99999999"
```

Se preferir usar o Docker direto, sem o Compose:

```bash
docker build -t clima-por-cep .
docker run --rm -p 8080:8080 clima-por-cep
```

## Rodando os testes com Docker

Os testes rodam dentro do estágio de build da imagem, então não é preciso ter Go instalado:

```bash
docker compose run --rm testes
```

Eles não dependem de internet, porque as chamadas para o ViaCEP e para a Open-Meteo são simuladas com `httptest`. A suíte cobre as conversões de temperatura, a validação do CEP, o caso de uso (com mocks do testify), os gateways HTTP (incluindo a escolha da cidade certa quando há nomes repetidos) e o handler com os cenários 200, 404, 422 e 500.

## Rodando sem Docker

```bash
go test ./...
go run ./cmd/server
```

Se existir um arquivo `.env`, o `go run` lê ele da pasta em que o comando é executado, então rode a partir da raiz do projeto.
