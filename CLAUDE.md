# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Build complet (frontend + backend)
make all

# Backend uniquement
make build-go

# Frontend dev server (proxy vers localhost:5915)
cd web && bun run dev

# Lint frontend
cd web && bun run lint

# Tous les tests Go
go test ./...

# Un package spécifique
go test ./internal/utils/...

# Nettoyer les artefacts
make clean
```

## Architecture

Le binaire Go embarque le SPA React (`internal/web/dist/`). Au démarrage il initialise la base JSON, monte les adapters, puis sert le frontend en statique et l'API REST sur le port `5915`.

La couche centrale est `DeployService` (`internal/api/service/deploy.go`) — un struct qui agrège tous les adapters et est injecté dans chaque handler HTTP.

Les couches backend :
- `internal/api/` — handlers, routes, DTOs, middleware
- `internal/application/` — logique métier (un fichier par opération)
- `internal/adapter/` — abstraction des services externes (Docker, Git, GitHub, SSH, filesystem, réseau, base JSON)
- `internal/domain/` — modèles de données

## Authentification

Premier lancement → `POST /api/setup` (public) crée l'admin, génère un JWT secret aléatoire et retourne un JWT.

Connexions suivantes → `POST /api/login` → retourne un JWT valable 30 jours.

Toutes les routes protégées attendent `Authorization: Bearer <token>`. Le token est stocké dans le `localStorage` du navigateur.

Endpoints publics (pas d'auth) : `POST /api/github/events`, `GET /api/info`, `POST /api/setup`, `POST /api/login`.

## Types de services & déploiement

Trois types de services coexistent dans `[]domain.Service` :
- `github_repo` — buildé via Dockerfile (si présent) ou Nixpacks, exposé via Traefik avec sous-domaine optionnel
- `database` — Postgres, Mongo, Redis, Minio — container Docker direct, non exposé publiquement
- `llm` — modèles Ollama, même logique que `database`

Le déploiement (`DeployApplication`) vérifie d'abord que Traefik tourne (le pull et démarre si besoin), puis boucle sur tous les services. Le serveur doit avoir un domaine configuré pour déployer des `github_repo`.

## Persistance

Tout l'état est persisté dans un seul fichier JSON : `~/.config/justdeploy/database.json`.  
Il contient trois clés : `Server`, `Services` (tableau), et `Settings` (email admin, hash mot de passe, JWT secret, tokens GitHub).

Les logs de build sont stockés séparément dans `~/.config/justdeploy/<serviceId>.log`.

Toutes les lectures/écritures passent par `DatabaseAdapter` (`internal/adapter/database.go`) — pas d'ORM, pas de migrations.

> Une migration vers SQLite est prévue prochainement.

## Temps réel (SSE)

Le frontend utilise `EventSource` (Server-Sent Events) pour suivre les événements de déploiement en temps réel.

Le hook `useSubEvent<T>(path)` (`web/src/hooks/useSubEvent.tsx`) encapsule la connexion SSE et expose les données typées au composant.

## Conventions

- Fichiers Go en kebab-case (`deploy-application.go`)
- `internal/application/` : un fichier par opération métier
- `internal/api/dto/` : tous les DTOs de l'API
- `web/src/services/` : un fichier par appel API côté frontend
- Le frontend utilise Zustand pour le state global (`web/src/contexts/`) et React local state pour l'UI éphémère
