# SRE Workshop

A hands-on workshop series covering essential SRE (Site Reliability Engineering) skills — containerization with Docker and CI/CD automation with GitHub Actions.

## Workshops

| Workshop | Duration | Description |
|----------|----------|-------------|
| [Docker](docker/) | 2 hours | Docker fundamentals, Dockerfile, and Docker Compose |
| [CI/CD](cicd/) | 2 hours | GitHub Actions, CI pipeline, release and deployment automation |

> **Recommended order:** Docker first, then CI/CD — the CI/CD workshop assumes Docker knowledge.

## Prerequisites

- [ ] **Git** — `git --version`
- [ ] **Go 1.22+** — `go version`
- [ ] **Docker** — `docker --version`
- [ ] **GitHub account** — [github.com](https://github.com)
- [ ] **Editor** — [VS Code](https://code.visualstudio.com/) recommended (with YAML extension)

## Repository Structure

```
sre-workshop/
├── README.md              # This file
├── docker/                # Docker workshop
│   ├── README.md          # Workshop overview
│   └── docker-workshop.md # Full teaching material
└── cicd/                  # CI/CD workshop
    ├── README.md          # Workshop overview
    ├── 01-cicd-intro.md   # CI/CD concepts
    ├── 02-github-actions-basics.md
    ├── 03-first-workflow.md
    ├── 04-go-ci-pipeline.md
    ├── 05-pr-checks.md
    ├── 06-release-automation.md
    ├── 07-deployment.md
    ├── examples/          # Sample Go application
    └── exercises/         # Hands-on exercises
```

## License

This material is intended for educational use within SDC workshops.
