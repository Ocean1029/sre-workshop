# Final Exercise: End-to-End CI/CD + Monitoring

## Overview

A 25-minute capstone that ties together Docker, CI/CD, and Prometheus. You'll deploy a service via GitHub Actions to the SDC machine, expose metrics for Prometheus, and trigger a Discord alert when the service crashes.

> **The exercise lives in a separate repo:** [Ocean1029/sre-workshop-capstone](https://github.com/Ocean1029/sre-workshop-capstone).
>
> That repo is the one you clone, edit, and push. This directory only contains the briefing for instructors and anyone reading the main workshop repo.

**What the capstone repo provides:**

- Service source code and Dockerfile
- Most of the GitHub Actions workflow (CI complete, CD partially filled in)
- Prometheus scrape config and alert rule
- `docker-compose.yml` for the Part 1 local stack

**What's already running on the SDC machine:**

- GitHub self-hosted runner
- Prometheus and Alertmanager (Alertmanager pre-wired to a Discord channel)

## Target Audience

- Has completed the Docker, CI/CD, and Prometheus workshops earlier in the day
- Ready to put the three pieces together end to end

## Flow

### Part 0 — Walk through the exercise (5 min)

Target architecture: repo → GitHub Actions → self-hosted runner → `docker run` on SDC machine → Prometheus scrape → Alertmanager → Discord.

### Part 1 — Run it locally (10 min)

1. `docker compose up` to build and start the service locally.
2. Confirm Prometheus (also in the compose) is scraping the service.
3. Open the Prometheus UI and find the service's metrics.
4. Hit the `/crash` endpoint; watch the metric counter go up.

### Part 2 — Wire up CD and trigger a real alert (10 min)

1. Following the CI/CD workshop's ch04, fill in the missing CD job (self-hosted runner + `docker run`).
2. Push your assigned `student-<ID>` branch; watch the pipeline turn green on GitHub Actions.
3. Hit `/crash` on the deployed service.
4. Wait for the alert to fire; check Discord for the notification.

## Getting Started

```bash
git clone https://github.com/Ocean1029/sre-workshop-capstone.git
cd sre-workshop-capstone
docker compose up --build
```

Full step-by-step instructions are in the capstone repo's [README](https://github.com/Ocean1029/sre-workshop-capstone#flow).
## Thanks for Attending
![alt](./assets/qrcode.svg)