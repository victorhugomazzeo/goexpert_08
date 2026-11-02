# Leilões com fechamento automático

Sistema de leilões em Go (desafio Go Expert). Ao criar um leilão, uma goroutine aguarda o tempo definido em `AUCTION_INTERVAL` e altera o status do leilão para fechado (`Completed`) no MongoDB, sem bloquear a requisição.

Implementação: [`internal/infra/database/auction/create_auction.go`](internal/infra/database/auction/create_auction.go)

## Requisitos

- Docker e Docker Compose
- Go 1.20+ (apenas para rodar os testes)

## Como rodar

```bash
docker compose up --build
```

API disponível em `http://localhost:8080`.

## Configuração do tempo

No arquivo [`cmd/auction/.env`](cmd/auction/.env):

```env
AUCTION_INTERVAL=20s
```

- Aceita o formato de duração do Go: `30s`, `5m`, `1h`.
- Se estiver vazio ou inválido, usa 5 minutos.
- O mesmo valor define o fechamento automático e o bloqueio de novos lances.

Após alterar, suba novamente com `docker compose up --build`.

## Testes

O teste cria um leilão, confirma que está aberto, espera o intervalo (1s no teste) e confirma que foi fechado automaticamente.

```bash
docker compose up -d mongodb
go test ./internal/infra/database/auction/ -v
```

## Teste manual

Com a aplicação rodando:

```bash
# 1. Criar leilão
curl -X POST http://localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{"product_name":"Notebook","category":"Eletronicos","description":"Notebook usado em bom estado","condition":1}'

# 2. Listar leilões fechados: vazio logo após criar,
#    o leilão aparece depois de AUCTION_INTERVAL
curl "http://localhost:8080/auction?status=1"
```

Status: `0` = aberto (`Active`), `1` = fechado (`Completed`).

## Limitação

O agendamento do fechamento fica em memória. Se a aplicação reiniciar antes do fim do leilão, ele não é fechado automaticamente.
