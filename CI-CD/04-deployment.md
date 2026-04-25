# 04 — 自動化部署

CI 確保程式碼品質之後，這章把程式部署到實際的伺服器讓使用者看到新版本。採兩段式做法：先用 GitHub-hosted runner build image 推到 GHCR，再用 self-hosted runner 把 image pull 下來重啟容器。

## 目錄

- [學習目標](#學習目標)
- [什麼是 CD？](#什麼是-cd)
- [Environments 與 Secrets](#environments-與-secrets)
- [實戰：用 self-hosted runner 部署](#實戰用-self-hosted-runner-部署)
- [結語：完整 CI/CD 全景與工作坊回顧](#結語完整-cicd-全景與工作坊回顧)


## 學習目標

完成本章節後，你將能夠：

- 理解 **GitHub Environments** 與 **Secrets** 的概念和設定方式
- 區分 **環境變數** 與 **Secrets** 的使用時機
- 用 GitHub Actions 把 image push 到 **GHCR**（GitHub Container Registry）
- 用 **self-hosted runner** 把 image pull 下來部署到自己的伺服器


## 什麼是 CD？

**部署（Deployment）** 就是把寫好的程式搬到實際運行的環境上，讓使用者可以用到新版本。CI 確保程式碼品質之後，CD 負責自動完成這件事：把通過檢查的程式碼打包成可以在伺服器上跑的格式（像是 Docker image），部署到對應的環境上，再執行後續的處理（像是通知團隊、更新監控）。從 push 程式碼到使用者看到新功能，整個過程全自動，不需要人工介入。

這種模式要求你對自動化測試有非常高的信心，因為沒有人會在中間攔一下，測試就是你唯一的安全網。


## Environments 與 Secrets

### GitHub Environments 的概念

**GitHub Environments** 是 GitHub 提供的功能，讓你可以定義不同的部署環境（例如 `staging`、`production`），並為每個環境設定獨立的：

- **Secrets**：每個環境可以有不同的 secrets（例如不同環境的 API key）
- **Protection Rules**：例如 production 需要手動核准才能部署
- **Deployment branches**：限制只有特定分支能部署到此環境

### Secret 放哪裡？Repository vs Environment

GitHub 在 Settings 裡有**兩個**可以放 secret 的地方，對新手常造成混淆：

| 放法 | 設定位置 | 適用情境 |
|------|---------|---------|
| **Repository secret** | Settings → Secrets and variables → Actions | 所有 workflow、所有 job 都可以讀到。最簡單的預設 |
| **Environment secret** | Settings → Environments → `<env>` → Secrets | 只有明確宣告 `environment: <name>` 的 job 才能讀到。可以搭配 protection rules |

兩種在 workflow 裡都是用 `${{ secrets.XXX }}` 讀取。本章 push image 到 GHCR 用 GitHub 自動產生的 `GITHUB_TOKEN`，不用自己設 secret；但實務上若部署需要外部 API key、資料庫密碼，就會放在這裡。如果你想讓 production 需要手動核准再部署，可以改用 environment secret 並設 protection rules。

### 如何設定 Repository Secret

1. 到你的 GitHub repository 頁面，點擊 **Settings**
2. 左側選單點擊 **Secrets and variables** → **Actions**
3. 點擊 **New repository secret**，輸入 Name 和 Secret 值
4. 儲存

### Secrets 的安全守則

- **永遠不要硬寫在程式碼中**，一律用 `${{ secrets.XXX }}` 讀取
- **Fork PR 不能存取 secrets**（GitHub 自動防護，避免惡意 PR 偷取）
- 離開團隊的成員，相關的 secrets 要立即輪換

```yaml
# Good
env:
  API_KEY: ${{ secrets.API_KEY }}

# Bad — NEVER do this
env:
  API_KEY: "sk-abc123xyz789"
```


## 實戰：用 self-hosted runner 部署

### 為什麼用 self-hosted runner？

**Self-hosted runner** 是你自己準備一台機器、安裝 GitHub Actions runner 程式並註冊到 repo，之後 workflow 就可以指定跑在這台機器上。跟 GitHub-hosted runner 對照：

| | GitHub-hosted | Self-hosted |
|---|---|---|
| 機器來源 | GitHub 提供，每次都是新的 | 自己的伺服器（VM、實體機都行） |
| 網路 | 公網 IP，連不到內部資源 | 在你的網段內，可直接接觸內部服務 |
| 費用 | 公開 repo 免費 | 機器費用自己出 |

CD 用 self-hosted 的最大理由是「**部署目標就在 runner 同一台（或同內網）**」。要部署到系辦的伺服器，從 GitHub-hosted runner 連回來要開 SSH 或 VPN，麻煩又有資安風險；直接在那台伺服器上跑 runner，部署就是執行本機的 docker 指令。

### 前置準備

本章假設：

1. 你的 repo 已經有一個 self-hosted runner，標籤為 `self-hosted`（工作坊現場會提供共用 runner）
2. 那台 runner 機器上已經裝好 Docker
3. 你的 GHCR image 設為 public（GitHub 上 Settings → Packages，這樣 pull 不用認證）

### 部署 Workflow

下面的 CD workflow 沒辦法在本機執行成功，這節只會觀察這個 workflow 的結構：

```yaml
name: Build and Deploy

on:
  push:
    branches: [main]
  workflow_dispatch:

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build-and-push:
    name: Build & Push Image
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v4

      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: |
            ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
            ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}

  deploy:
    name: Deploy to Server
    needs: build-and-push
    runs-on: [self-hosted]
    environment: production
    steps:
      - name: Pull latest image
        run: docker pull ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest

      - name: Stop and remove old container
        run: |
          docker stop sample-app || true
          docker rm sample-app || true

      - name: Run new container
        run: |
          docker run -d \
            --name sample-app \
            -p 8080:80 \
            --restart unless-stopped \
            ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
```

### 逐段解說

**兩段式設計**

`build-and-push` 用 GitHub-hosted runner（規格穩定、不佔自己機器資源）負責建置 image 並推到 GHCR；`deploy` 用 self-hosted runner（直接在部署機器上）pull image、重啟容器。

**Build & Push Job**

```yaml
permissions:
  contents: read
  packages: write
```

要 push image 到 GHCR 必須明確要 `packages: write`。預設權限不夠用。

```yaml
- uses: docker/login-action@v3
  with:
    registry: ${{ env.REGISTRY }}
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}
```

`GITHUB_TOKEN` 是 GitHub Actions 為每個 workflow run 自動產生的 token，搭配前面的 `permissions` 取得對應權限。push image 不用自己生 PAT。

```yaml
- uses: docker/build-push-action@v5
  with:
    context: .
    push: true
    tags: |
      ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
      ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
```

讀 repo 根目錄的 `Dockerfile` build 出 image，同時打兩個 tag：

- `latest`：方便 deploy 直接用
- `${{ github.sha }}`：留下這次 commit 的版本，未來想 rollback 才有依據

**Deploy Job**

```yaml
needs: build-and-push
runs-on: [self-hosted]
environment: production
```

- `needs: build-and-push`：image 沒 push 完就不用 deploy
- `runs-on: [self-hosted]`：用標籤指定要跑在 self-hosted runner，這台 runner 就在部署目標機器上
- `environment: production`：可以接 protection rules（例如手動核准），跟前面 [Environments 與 Secrets](#environments-與-secrets) 提到的設定方式一樣

```yaml
- name: Pull latest image
  run: docker pull ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest

- name: Stop and remove old container
  run: |
    docker stop sample-app || true
    docker rm sample-app || true

- name: Run new container
  run: |
    docker run -d \
      --name sample-app \
      -p 8080:80 \
      --restart unless-stopped \
      ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
```

三個動作：拉新 image → 砍舊 container → 跑新的。

- `|| true`：第一次部署沒有舊 container 時，`stop` / `rm` 會失敗，但不該擋住流程
- `--restart unless-stopped`：伺服器重開機後容器自動起來


## 結語：完整 CI/CD 全景與工作坊回顧

四章走下來，你做了這些事：

- **第二章**：寫了第一個 workflow，push 上去看它跑，認識 workflow / job / step / runner 的關係，也刻意製造一次失敗看 debug 流程
- **第三章**：把 lint、test、build 三個 job 串成真正的 Go CI。`lint` 和 `test` 透過 `needs` 平行執行、`build` 等兩者都過才啟動；`golangci-lint` 守程式碼品質，`go test -race -coverprofile` 抓 race condition 並收 coverage，`upload-artifact` 把產物上傳方便後續取用。同時也看到 `on: pull_request` 怎麼在 merge 前先攔一次
- **第四章**：用兩段式 workflow 把建好的東西部署上線。GitHub-hosted runner 負責 build image 並 push 到 GHCR；self-hosted runner 負責 pull image、重啟容器。`environment: production` 還可以接手動核准等 protection rules

拼起來，整條從寫程式碼到部署上線的自動化長這樣：

```
開發者 push 程式碼
        │
        ▼
┌───────────────────────────────────────────────────────┐
│                   CI Pipeline (ch03)                  │
│                                                       │
│  ┌──────┐  ┌──────┐                                   │
│  │ Lint │  │ Test │  ← 平行執行                       │
│  └──┬───┘  └──┬───┘                                   │
│     └────┬────┘                                       │
│          ▼                                            │
│     ┌─────────┐                                       │
│     │  Build  │                                       │
│     └─────────┘                                       │
└───────────────────────────────────────────────────────┘
        │
        │ PR 合併到 main
        ▼
┌───────────────────────────────────────────────────────┐
│                   CD Pipeline (ch04)                  │
│                                                       │
│       ┌──────────────────┐    ┌──────────────────┐    │
│       │   Build & Push   │───▶│  Pull & Restart  │    │
│       │      to GHCR     │    │  (self-hosted)   │    │
│       └──────────────────┘    └──────────────────┘    │
│                                                       │
└───────────────────────────────────────────────────────┘
```

回頭看第一章那些痛點：

| 原本的痛點 | 現在怎麼解 |
|-----------|-----------|
| 手動測試每次都花好久 | push 後 CI 自動跑完 lint + test + build |
| 合併之後才發現壞掉 | PR 階段就先跑一次 CI，壞掉的改動合不進來 |
| 部署又忘了步驟 | self-hosted runner 上跑 `docker pull` + `docker run`，前置步驟寫進 workflow 裡 |
| 誰都可以亂部署 production | Environment protection 可以強制手動核准 |
| API key 散落在各處 | 統一放進 GitHub Secrets，log 自動遮蔽 |

沒講到的東西還很多，Release 自動化、image 掃描、GitOps、監控告警之類的，工作上真的碰到再查就好。

### 延伸資源

- [GitHub Actions 官方文件](https://docs.github.com/en/actions)
- [Awesome Actions](https://github.com/sdras/awesome-actions)，社群整理的 Actions 清單
- [Working with the Container registry (GitHub Docs)](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [About self-hosted runners (GitHub Docs)](https://docs.github.com/en/actions/hosting-your-own-runners/managing-self-hosted-runners/about-self-hosted-runners)
- [The Twelve-Factor App](https://12factor.net/)
- [Google SRE Book](https://sre.google/sre-book/table-of-contents/)


[← 上一章：Go 專案 CI Pipeline](03-go-ci-pipeline.md) ｜ [回到目錄 →](README.md)
