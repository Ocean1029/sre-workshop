# 04 — 部署到雲端平台

Ocean 的 CI pipeline 跑得很順。Andrew 興奮地說：「那接下來就是把程式部署到伺服器上讓大家用了吧？」Ocean 點點頭，但又有點緊張，畢竟部署到正式環境可不能出差錯。Snow 說：「別擔心，我們可以設計一套安全的部署流程——先部署到測試環境驗證，確認沒問題再上正式環境。」

## Table of Contents

- [學習目標](#學習目標)
- [Environments 與 Secrets](#environments-與-secrets)
- [實戰：部署到 Fly.io](#實戰部署到-flyio)
- [小結與練習題](#小結與練習題)
- [補充：多環境部署（Staging → Production）](#補充多環境部署staging--production)
- [完整 CI/CD Pipeline 全景圖](#完整-cicd-pipeline-全景圖)


## 學習目標

完成本章節後，你將能夠：

- 理解 **GitHub Environments** 與 **Secrets** 的概念和設定方式
- 區分 **環境變數** 與 **Secrets** 的使用時機
- 使用 GitHub Actions 部署應用到 **Fly.io**


## Environments 與 Secrets

### GitHub Environments 的概念

**GitHub Environments** 是 GitHub 提供的功能，讓你可以定義不同的部署環境（例如 `staging`、`production`），並為每個環境設定獨立的：

- **Secrets**：每個環境可以有不同的 secrets（例如不同環境的 API key）
- **Protection Rules**：例如 production 需要手動核准才能部署
- **Deployment branches**：限制只有特定分支能部署到此環境

### 如何設定 Secrets

1. 到你的 GitHub repository 頁面，點擊 **Settings**
2. 左側選單點擊 **Environments** → **New environment**，輸入名稱（例如 `production`）
3. 在環境設定頁面中，點擊 **Add environment secret**，輸入 Name 和 Secret
4. 如果要設定手動核准，勾選 **Required reviewers** 並加入核准者

在 workflow 中透過 `${{ secrets.XXX }}` 讀取。

### 環境變數 vs Secrets 的差異

| 面向 | 環境變數（Variables） | Secrets |
|------|---------------------|---------|
| **存放內容** | 非敏感資料 | 敏感資料 |
| **範例** | `REGION=asia-east1`、`APP_NAME=myapp` | API Key、密碼、憑證 |
| **可見性** | 可以在 log 中看到 | 自動遮蔽，顯示為 `***` |
| **使用方式** | `${{ vars.REGION }}` | `${{ secrets.API_KEY }}` |
| **修改** | 可以直接查看和修改值 | 只能覆寫，無法查看原始值 |

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


## 實戰：部署到 Fly.io

### 為什麼選 Fly.io？

[Fly.io](https://fly.io) 是一個開發者友善的容器平台，特別適合部署後端應用：

- **讀你的 Dockerfile 就好**：不需要自己管 registry，`flyctl deploy` 會自己建置和推送
- **一個指令完成部署**：設定檔 `fly.toml` 產生後，後續每次部署就是 `flyctl deploy`
- **免費額度足夠實驗**：可以跑幾個小型 VM
- **全球邊緣節點**：自動把應用部署到離使用者最近的區域

對我們這個 Go HTTP API 範例來說，Fly.io 是最輕量的選擇。

### 前置準備

在設定 workflow 之前，先在本機完成一次性的 Fly.io 初始化：

1. 安裝 `flyctl`：[官方安裝指南](https://fly.io/docs/hands-on/install-flyctl/)
2. 登入：`flyctl auth login`
3. 在專案根目錄執行 `flyctl launch`——這會產生 `fly.toml` 設定檔（記得把它 commit 到 repo）
4. 產生 CI 用的 API token：`flyctl tokens create deploy`
5. 把 token 存到 GitHub Secret，名稱 `FLY_API_TOKEN`

### 部署 Workflow

請在你的專案中建立 `.github/workflows/cd.yml`：

```yaml
name: Deploy to Fly.io

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  deploy:
    name: Deploy
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v4

      - name: Setup flyctl
        uses: superfly/flyctl-actions/setup-flyctl@master

      - name: Deploy to Fly.io
        run: flyctl deploy --remote-only
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

> 完整檔案也可以在 `cicd/examples/sample-app/.github/workflows/cd.yml` 找到。

### 逐段解說

**觸發條件**

```yaml
on:
  push:
    branches: [main]
  workflow_dispatch:
```

- `push to main`：程式碼合併到 `main` 時自動部署
- `workflow_dispatch`：也允許在 GitHub UI 手動觸發——需要重新部署但不想改程式碼時很實用

**Environment**

```yaml
environment: production
```

告訴 GitHub 這個 job 使用 `production` 環境。如果 `production` 環境有設定 protection rules（例如需要核准），workflow 會在執行此 job 前暫停，等待核准。

**Deploy 步驟**

```yaml
- name: Setup flyctl
  uses: superfly/flyctl-actions/setup-flyctl@master

- name: Deploy to Fly.io
  run: flyctl deploy --remote-only
  env:
    FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

這兩個 step 就是全部了：

1. `setup-flyctl` 把 `flyctl` 指令裝到 runner 上
2. `flyctl deploy --remote-only` 讀取 repo 根目錄的 `fly.toml` 和 `Dockerfile`，在 Fly.io 的 builder 上建置 image 並部署。`--remote-only` 代表不在 runner 本地建置，省時也省資源

`FLY_API_TOKEN` 透過環境變數傳遞給 `flyctl`，它會自動識別並用來認證。


## 小結與練習題

### 本章重點回顧

- **GitHub Environments**：為不同環境（staging、production）設定獨立的 secrets 和 protection rules
- **Secrets vs 環境變數**：Secrets 存放敏感資料（自動遮蔽），環境變數存放非敏感設定
- **Fly.io 部署**：只要 `fly.toml` + `FLY_API_TOKEN`，一個 `flyctl deploy` 就完成

## 補充：多環境部署（Staging → Production）

> 本節是進階補充。前面 Fly.io 單環境部署的流程已經足夠完成一個基本 CD pipeline；這裡示範在真實團隊中常見的「先 staging、再 production」兩階段流程，有興趣再看即可。

### 為什麼需要多環境？

在正式的開發流程中，程式碼通常會先部署到 **staging 環境** 進行驗證，確認沒問題後再部署到 **production 環境**。

```
┌────────┐     ┌────────────┐     ┌──────────────┐
│  Code  │────▶│  Staging   │────▶│  Production  │
│  Merge │     │  (自動)     │     │  (手動核准)   │
└────────┘     └────────────┘     └──────────────┘
```

Fly.io 的做法是為每個環境建立獨立的 app（例如 `myapp-staging` 和 `myapp-prod`），並使用不同的 `fly.toml` 或用 `--config` / `--app` 參數切換。

### Staging + Production 的 Workflow 範例

```yaml
name: Deploy Pipeline

on:
  push:
    branches: [main]

jobs:
  # ──────────────────────────────────────
  # Stage 1: Deploy to Staging (automatic)
  # ──────────────────────────────────────
  deploy-staging:
    name: Deploy to Staging
    runs-on: ubuntu-latest
    environment: staging
    steps:
      - uses: actions/checkout@v4
      - uses: superfly/flyctl-actions/setup-flyctl@master

      - name: Deploy to Staging
        run: flyctl deploy --remote-only --app myapp-staging
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}

      - name: Smoke test
        run: |
          for i in 1 2 3 4 5; do
            STATUS=$(curl -s -o /dev/null -w "%{http_code}" https://myapp-staging.fly.dev/health)
            if [ "$STATUS" = "200" ]; then
              echo "Health check passed!"
              exit 0
            fi
            echo "Attempt $i: got ${STATUS}, retrying in 10s..."
            sleep 10
          done
          echo "::error::Smoke test failed"
          exit 1

  # ──────────────────────────────────────────
  # Stage 2: Deploy to Production (manual approval)
  # ──────────────────────────────────────────
  deploy-production:
    name: Deploy to Production
    needs: deploy-staging
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v4
      - uses: superfly/flyctl-actions/setup-flyctl@master

      - name: Deploy to Production
        run: flyctl deploy --remote-only --app myapp-prod
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

### 關鍵設計

1. **Staging 自動部署**：push 到 `main` 後自動部署到 staging
2. **Smoke Test**：staging 部署完成後做一次健康檢查；用 retry 迴圈避免因為新 instance 還沒完全啟動就誤判失敗
3. **Production 手動核准**：`environment: production` 設定了 protection rules，workflow 會暫停等待核准
4. **`needs: deploy-staging`**：production 一定要在 staging 成功後才會執行

### 設定 Environment Protection Rules

要讓 production 部署需要核准，在 GitHub 上設定：

1. **Settings** → **Environments** → 點擊 **production**
2. 勾選 **Required reviewers**，加入需要核准的人
3. 可選：設定 **Wait timer**（例如 5 分鐘冷靜期）

設定完成後，當 workflow 執行到 `deploy-production` 時會暫停，核准者需要到 Actions 頁面點擊 **Approve and deploy** 才會繼續。


## 完整 CI/CD Pipeline 全景圖

恭喜你完成了整個 CI/CD 工作坊！讓我們回顧一下，你現在已經有能力建立這樣一套完整的自動化流程：

```
開發者 push 程式碼
        │
        ▼
┌───────────────────────────────────────────────────────┐
│                   CI Pipeline (ch03)                   │
│                                                       │
│  ┌──────┐  ┌──────┐                                  │
│  │ Lint │  │ Test │  ← 平行執行                       │
│  └──┬───┘  └──┬───┘                                  │
│     └────┬────┘                                      │
│          ▼                                           │
│     ┌─────────┐                                      │
│     │  Build  │                                      │
│     └─────────┘                                      │
└───────────────────────────────────────────────────────┘
        │
        │ PR 合併到 main
        ▼
┌───────────────────────────────────────────────────────┐
│                  CD Pipeline (ch04)                   │
│                                                       │
│  ┌──────────┐    ┌──────────┐    ┌──────────────┐    │
│  │ Deploy   │───▶│ Smoke    │───▶│ Deploy       │    │
│  │ Staging  │    │ Test     │    │ Production   │    │
│  │ (自動)   │    │          │    │ (手動核准)    │    │
│  └──────────┘    └──────────┘    └──────────────┘    │
└───────────────────────────────────────────────────────┘
```

從寫程式碼到部署上線，每個環節都有自動化的品質把關。這就是 CI/CD 的力量。

## 結語

四章走下來，你做了這些事：

- **第二章**：寫了第一個 workflow，push 上去看它跑，認識 workflow / job / step / runner 的關係，也刻意製造一次失敗看 debug 流程
- **第三章**：把 lint、test、build 三個 job 串成真正的 Go CI。`lint` 和 `test` 透過 `needs` 平行執行、`build` 等兩者都過才啟動；`golangci-lint` 守程式碼品質，`go test -race -coverprofile` 抓 race condition 並收 coverage，`upload-artifact` 把產物交給下一個 job。同時也看到 `on: pull_request` 怎麼在 merge 前先攔一次
- **第四章**：把建好的東西部署到 Fly.io，用 GitHub Environments 管 secrets，設計 staging 自動部署加 smoke test、production 手動核准的兩階段流程

回頭看第一章 Ocean 和 Andrew 那些抱怨：

| 原本的痛點 | 現在怎麼解 |
|-----------|-----------|
| 手動測試每次都花好久 | push 後 CI 自動跑完 lint + test + build |
| 合併之後才發現壞掉 | PR 階段就先跑一次 CI，壞的改不了合進來 |
| 部署又忘了步驟 | `flyctl deploy` 一個指令，前置步驟寫進 workflow 裡跑 |
| 誰都可以亂部署 production | Environment protection 強制手動核准 |
| API key 散落在各處 | 統一放進 GitHub Secrets，log 自動遮蔽 |

沒講到的東西還很多，Release 自動化、image 掃描、GitOps、監控告警之類的，工作上真的碰到再查就好。

### 延伸資源

- [GitHub Actions 官方文件](https://docs.github.com/en/actions)
- [Awesome Actions](https://github.com/sdras/awesome-actions)，社群整理的 Actions 清單
- [Fly.io Docs](https://fly.io/docs/)
- [The Twelve-Factor App](https://12factor.net/)
- [Google SRE Book](https://sre.google/sre-book/table-of-contents/)


[← 上一章：Go 專案 CI Pipeline](03-go-ci-pipeline.md) ｜ [回到目錄 →](README.md)
