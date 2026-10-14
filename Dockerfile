# Estágio 1: Build do Frontend Vue 3 SPA
FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# Estágio 2: Build da aplicação Go
FROM golang:alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Garante que os assets compilados do frontend estejam no web/dist para o go:embed
COPY --from=frontend-builder /app/web/dist ./web/dist

# Compilar binário estático ultraleve e stripped
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# Estágio 3: Imagem final mínima e segura
FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

COPY --from=backend-builder /app/server /app/server

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

ENTRYPOINT ["/app/server"]
