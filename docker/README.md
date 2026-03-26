# Docker Workshop — From Zero to Containerization

## Overview

A **2-hour** hands-on workshop that takes you from Docker fundamentals to multi-container orchestration with Docker Compose. All examples use Go as the application language.

## Target Audience

- Familiar with basic Linux command-line operations
- Docker beginners — no prior container experience required

## Prerequisites

- [ ] **Docker Engine** installed — `docker --version`
- [ ] **Docker Compose** installed — `docker compose version`
- [ ] **Go 1.22+** installed — `go version`
- [ ] **Editor** — [VS Code](https://code.visualstudio.com/) recommended

## Curriculum

| Part | Topic | Duration |
|------|-------|----------|
| 1 | Docker Core Concepts & Basic Operations | 50 min |
| 2 | Dockerfile Deep Dive | 35 min |
| ☕ | Break | 5 min |
| 3 | Docker Compose Multi-Container Orchestration | 30 min |
| 4 | Comprehensive Exercise & Further Learning | 15 min |

> **Total: 2 hours (120 minutes, including break)**

### Part 1: Docker Core Concepts & Basic Operations (50 min)

- What is Docker? Containerization concepts and history
- Container vs Virtual Machine
- Docker architecture (Client / Daemon / Registry)
- Image & Container relationship
- Image operations: `pull`, `list`, `remove`
- Container operations: `run`, `stop`, `rm`, `exec`
- Port Mapping & Volume mounting
- Container debugging techniques
- **Lab 1:** Run your first container

### Part 2: Dockerfile Deep Dive (35 min)

- Dockerfile syntax and core instructions (`FROM`, `COPY`, `RUN`, `CMD`, `EXPOSE`, etc.)
- `.dockerignore` best practices
- Writing a Dockerfile for a Go application
- Multi-stage Build for optimized images
- Image best practices (layer caching, minimal base images)
- **Lab 2:** Build your Go application image

### Part 3: Docker Compose (30 min)

- Why Docker Compose?
- `docker-compose.yml` syntax
- Service definitions, networks, and volumes
- Environment variable management
- Common Compose commands
- Live demo: Go API + PostgreSQL + Redis
- **Lab 3:** Deploy a full-stack service with Compose

### Part 4: Wrap-up (15 min)

- Comprehensive exercise
- Common troubleshooting tips
- Further learning resources

## Teaching Material

The full workshop content is in [docker-workshop.md](docker-workshop.md), which includes detailed explanations, code examples, architecture diagrams, and appendices (Docker / Dockerfile / Compose cheat sheets).
