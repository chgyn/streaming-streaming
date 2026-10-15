# 🎬 Local Stream Hub

Plataforma unificada e autocontida para catálogo, streaming e canais de TV ao vivo, construída com arquitetura BFF em **Go (1.25)**, **Vue 3 (Composition API)**, **Pinia**, **Tailwind CSS v4** e embarcada com `//go:embed` para consumo mínimo de memória (< 15 MB RAM).

---

## ✨ Funcionalidades

- **Catálogo de Filmes & Séries:** Integração direta com a API do TMDB com proxy reverso e cache em memória (TTL configurável).
- **Canais de TV Ao Vivo:** Catálogo categorizado e player de canais da SuperFlixAPI com resiliência a rate limit e fallback automático.
- **Player Unificado:** Suporte para reprodução de filmes, séries (por temporada e episódio) e canais ao vivo com suporte a transmissão / Smart TV via HTTPS direto.
- **Eficiência Extrema:** Executa em um único binário de ~20 MB consumindo menos de 15 MB de RAM no Docker.

---

## 🚀 Como Executar

### Com Docker (Recomendado)

```bash
# 1. Configurar variáveis no .env
cp .env.example .env

# 2. Construir imagem
docker build -t local-stream-hub .

# 3. Executar o container
docker run -d -p 8080:8080 --env-file .env --name stream-hub local-stream-hub
```

Acesse no navegador: `http://localhost:8080` (ou pelo IP local da máquina na sua rede Wi-Fi).

### Desenvolvimento Local

```bash
# Frontend
cd frontend
npm install
npm run dev

# Backend (em outro terminal)
go run ./cmd/server/main.go
```

---

## 📖 Documentação & Histórico

- [Módulo de Canais Ao Vivo](docs/MODULO_CANAIS_AO_VIVO.md)
- [Histórico Técnico de Commits (Changelog)](docs/HISTORICO_COMMITS.md)

