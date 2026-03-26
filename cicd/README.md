# GitHub Actions Workshop

## Overview

A **2-hour** hands-on workshop on GitHub Actions, starting from CI/CD fundamentals and progressively covering workflow authoring, Go project CI pipelines, PR automation, release automation, and cloud deployment. By the end, you'll be able to integrate CI/CD into real-world projects.

## Target Audience

- Familiar with **GitHub** basics (clone, push, pull, branch)
- Has **Docker** foundational knowledge (completed the Docker workshop)
- Interested in automation for testing and deployment

## Prerequisites

- [ ] **GitHub account** — sign up at [github.com](https://github.com) if needed
- [ ] **Git** installed — `git --version`
- [ ] **Go 1.22+** installed — `go version`
- [ ] **Docker** installed — `docker --version`
- [ ] **Editor** — [VS Code](https://code.visualstudio.com/) recommended (with YAML extension)
- [ ] (Optional) **[act](https://github.com/nektos/act)** — test GitHub Actions workflows locally

## Curriculum

| Chapter | Topic | Duration |
|---------|-------|----------|
| 01 | CI/CD Concepts | 10 min |
| 02 | GitHub Actions Fundamentals | 15 min |
| 03 | Hands-on: Your First Workflow | 15 min |
| 04 | Go Project CI Pipeline | 20 min |
| 05 | PR Check Automation | 15 min |
| 06 | Release Automation | 20 min |
| 07 | Deploy to Cloud | 15 min |
| -- | Wrap-up & Q&A | 10 min |

> **Total: 2 hours (120 minutes)**

## Project Structure

```
cicd/
├── README.md                          # This file
├── README.zh-TW.md                    # Chinese version
├── 01-cicd-intro.md                   # CI/CD concepts
├── 02-github-actions-basics.md        # GitHub Actions fundamentals
├── 03-first-workflow.md               # Hands-on: first workflow
├── 04-go-ci-pipeline.md               # Go CI pipeline
├── 05-pr-checks.md                    # PR check automation
├── 06-release-automation.md           # Release automation
├── 07-deployment.md                   # Cloud deployment
├── examples/                          # Example code
│   └── sample-app/                    # Sample Go application
│       ├── main.go                    # Entry point
│       ├── handler.go                 # HTTP handler
│       ├── handler_test.go            # Tests
│       ├── go.mod                     # Go module definition
│       └── Dockerfile                 # Docker image definition
└── exercises/                         # Exercises
    ├── exercise-01-basics.md          # Chapters 01-03
    ├── exercise-02-ci-pipeline.md     # Chapters 04-05
    └── exercise-03-advanced.md        # Chapters 06-07
```

## Sample Application

`examples/sample-app` is a simple **Go HTTP API server** with basic RESTful endpoints. It serves as the foundation for all demos and exercises throughout the workshop — you'll build a complete CI/CD pipeline for it, covering automated testing, code quality checks, Docker image builds, and deployment.
