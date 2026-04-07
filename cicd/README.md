# GitHub Actions Workshop

## Overview

A two-hour hands-on workshop on GitHub Actions. We start from CI/CD fundamentals, build your first workflow, then construct a full CI pipeline for a real Go project. By the end you'll have the basics needed to introduce CI/CD into your own projects.

## Target Audience

- Familiar with **GitHub** basics (clone, push, pull, branch)
- Has **Docker** foundational knowledge (completed the Docker workshop)
- Interested in automation for testing and deployment

## Prerequisites

- [ ] **GitHub account** — sign up at [github.com](https://github.com) if needed
- [ ] **Git** installed — `git --version`
- [ ] **Go 1.24+** installed — `go version`
- [ ] **Docker** installed — `docker --version`
- [ ] **Editor** — [VS Code](https://code.visualstudio.com/) recommended (with YAML extension)

## Schedule

| Time | Chapter | Content | Duration |
|------|---------|---------|----------|
| 0:00–0:05 | — | Intro and setup check | 5 min |
| 0:05–0:20 | **01** | CI/CD concepts | 15 min |
| 0:20–0:55 | **02** | GitHub Actions basics (hello.yml + UI + structure) | 35 min |
| 0:55–1:10 | **02** | Manual trigger + failure debugging | 15 min |
| 1:10–1:25 | **02 exercise** | Customize hello.yml | 15 min |
| 1:25–2:00 | **03** | Go CI pipeline walkthrough | 35 min |

> **Total: 2 hours (120 minutes)**
>
> The 03 exercise and chapter 04 (deployment) are extension material for self-study after the workshop.

### Chapter outline

**01 — CI/CD Concepts**
Why CI/CD exists, what problems it solves, what CI and CD each do, and how SDC actually uses it in production projects.

**02 — GitHub Actions Basics**
Write your first `hello.yml`, push it, explore the GitHub Actions UI, and learn the five structural elements of a workflow (Workflow / Event / Job / Step / Runner). Cover `uses` for Actions, `${{ }}` for contexts, manual triggers, and a deliberate failure walkthrough for debugging practice.

**03 — Go CI Pipeline**
Build a complete CI workflow for a Go HTTP API: lint with `golangci-lint`, test with coverage, build, and understand how to use `needs` for job dependencies and `actions/upload-artifact` to pass files between jobs. Wraps up with PR-triggered CI, permissions, and fork safety.

**04 — Deploy to Cloud (extension)**
GitHub Environments and Secrets, plus a hands-on deployment example to Fly.io. Self-study material.

## Project Structure

```
cicd/
├── README.md                          # This file
├── README.zh-TW.md                    # Chinese version
├── 01-cicd-intro.md                   # CI/CD concepts
├── 02-github-actions-basics.md        # GitHub Actions fundamentals + hands-on
├── 03-go-ci-pipeline.md               # Go CI pipeline
├── 04-deployment.md                   # Cloud deployment (extension)
├── assets/                            # Screenshots and diagrams
├── examples/
│   ├── .github/workflows/
│   │   ├── hello.yml                  # Ch02 demo workflow
│   │   └── repo-info.yml              # Ch02 failure example
│   └── sample-app/                    # Sample Go HTTP API
│       ├── main.go
│       ├── handler.go
│       ├── handler_test.go
│       ├── go.mod
│       ├── Dockerfile
│       ├── .golangci.yml
│       └── .github/workflows/
│           ├── ci.yml                 # Ch03 CI pipeline
│           └── cd.yml                 # Ch04 deploy to Fly.io
└── exercises/
    ├── 01-basics.md                   # Ch02 exercise
    └── 02-ci-pipeline.md              # Ch03 exercise
```

## Sample Application

`examples/sample-app` is a simple **Go HTTP API server** with basic RESTful endpoints. It serves as the foundation for chapters 03 and 04 — you'll build a complete CI/CD pipeline for it, covering automated testing, code style checks, Docker image builds, and deployment.
