# 07 — 部署到雲端平台

> **時間：15 分鐘**

---

## 目錄

- [學習目標](#學習目標)
- [Environments 與 Secrets](#environments-與-secrets)
- [部署策略概覽](#部署策略概覽)
- [實戰：部署到 GCP Cloud Run](#實戰部署到-gcp-cloud-run)
- [其他平台部署方式簡介](#其他平台部署方式簡介)
- [多環境部署（Staging → Production）](#多環境部署staging--production)
- [部署通知](#部署通知)
- [Rollback（回滾）策略](#rollback回滾策略)
- [安全最佳實踐](#安全最佳實踐)
- [常見問題排解](#常見問題排解)
- [小結與練習題](#小結與練習題)
- [完整 CI/CD Pipeline 全景圖](#完整-cicd-pipeline-全景圖)

Ocean 的 CI pipeline 跑得很順，Release 也自動化了。Andrew 興奮地說：「那接下來就是把程式部署到伺服器上讓大家用了吧？」Ocean 點點頭，但又有點緊張，畢竟部署到正式環境可不能出差錯。Snow 說：「別擔心，我們可以設計一套安全的部署流程，先部署到測試環境驗證，確認沒問題再上正式環境。」

---

## 學習目標

完成本章節後，你將能夠：

- 理解 **GitHub Environments** 與 **Secrets** 的概念和設定方式
- 區分 **環境變數** 與 **Secrets** 的使用時機
- 了解常見的 **部署策略**（Recreate、Rolling Update、Blue-Green、Canary）
- 使用 GitHub Actions 部署應用到 **GCP Cloud Run**
- 了解 **AWS、Vercel、Fly.io** 等平台的部署概念
- 設定 **多環境部署** 流程（Staging → Production）
- 實作 **部署通知** 與 **Rollback** 機制
- 遵循部署相關的 **安全最佳實踐**

---

## Environments 與 Secrets

### GitHub Environments 的概念

**GitHub Environments** 是 GitHub 提供的功能，讓你可以定義不同的部署環境（例如 `staging`、`production`），並為每個環境設定獨立的：

- **Secrets**：每個環境可以有不同的 secrets（例如不同環境的 API key）
- **Protection Rules**：例如 production 需要手動核准才能部署
- **等待時間**：部署前的等待期間

```
┌─────────────────────────────────────────────────────┐
│                GitHub Repository                    │
│                                                     │
│  ┌──────────────────┐  ┌──────────────────┐         │
│  │   staging        │  │   production     │         │
│  │                  │  │                  │         │
│  │  Secrets:        │  │  Secrets:        │         │
│  │  ├ DB_HOST       │  │  ├ DB_HOST       │         │
│  │  ├ API_KEY       │  │  ├ API_KEY       │         │
│  │  └ GCP_SA_KEY    │  │  └ GCP_SA_KEY    │         │
│  │                  │  │                  │         │
│  │  Protection:     │  │  Protection:     │         │
│  │  └ (none)        │  │  ├ Required      │         │
│  │                  │  │  │ reviewers     │         │
│  │                  │  │  └ Wait timer:   │         │
│  │                  │  │    5 minutes     │         │
│  └──────────────────┘  └──────────────────┘         │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### 如何設定 Secrets

**Step 1：進入 Settings**
1. 到你的 GitHub repository 頁面
2. 點擊 **Settings** 標籤

**Step 2：設定 Repository Secrets（全域）**
1. 左側選單點擊 **Secrets and variables** → **Actions**
2. 點擊 **New repository secret**
3. 輸入 **Name**（例如 `GCP_PROJECT_ID`）和 **Secret**（實際的值）
4. 點擊 **Add secret**

**Step 3：設定 Environment Secrets（環境專屬）**
1. 左側選單點擊 **Environments**
2. 點擊 **New environment**，輸入名稱（例如 `production`）
3. 在環境設定頁面中，點擊 **Add environment secret**
4. 輸入 Name 和 Secret

### Environment Protection Rules

在環境設定頁面中，你可以設定：

- **Required reviewers**：指定哪些人需要核准才能部署到此環境
- **Wait timer**：部署前等待指定時間（例如 5 分鐘），讓你有時間取消
- **Deployment branches**：限制只有特定分支能部署到此環境

### 環境變數 vs Secrets 的差異

| 面向 | 環境變數（Variables） | Secrets |
|------|---------------------|---------|
| **存放內容** | 非敏感資料 | 敏感資料 |
| **範例** | `REGION=asia-east1`、`APP_NAME=myapp` | API Key、密碼、憑證 |
| **可見性** | 可以在 log 中看到 | 自動遮蔽，不會顯示在 log 中 |
| **使用方式** | `${{ vars.REGION }}` | `${{ secrets.API_KEY }}` |
| **修改** | 可以直接查看和修改值 | 只能覆寫，無法查看原始值 |

### Secrets 的安全注意事項

- **不會印在 log 中**：GitHub Actions 會自動將 secrets 的值在 log 中替換為 `***`
- **Fork PR 不能存取**：從 fork repository 提交的 PR 無法讀取 secrets（防止惡意存取）
- **加密儲存**：Secrets 在 GitHub 端使用 Libsodium sealed box 加密
- **不要硬寫在程式碼中**：永遠不要在 workflow 檔案中直接寫入敏感資料

```yaml
# Good — use secrets
env:
  API_KEY: ${{ secrets.API_KEY }}

# Bad — hardcoded secret
env:
  API_KEY: "sk-abc123xyz789"  # NEVER do this!
```

<details>
<summary>💡 講師提示</summary>

> 可以提醒學生，即使是 private repository，也不應該把 secret 直接寫在程式碼中。因為 git history 會永久保留，一旦 commit 了就很難完全清除。

</details>

---

## 部署策略概覽

在部署應用程式時，有幾種常見的策略，各有不同的風險和複雜度：

### 1. 直接部署（Recreate）

先停掉舊版本，再啟動新版本。

```
時間軸 ──────────────────────────────────────────▶

舊版本 ████████████████████████████┐
                                   │ 停機時間
新版本                              └─ ███████████
```

| 優點 | 缺點 |
|------|------|
| 最簡單 | 有 **停機時間**（downtime） |
| 不需要額外資源 | 使用者會感受到服務中斷 |

### 2. 滾動更新（Rolling Update）

逐步將舊版本替換成新版本，一次更新一部分。

```
時間軸 ──────────────────────────────────────────▶

Instance 1  ████████ 舊 ████████┐
                                └── ████ 新 ████████
Instance 2  ████████████ 舊 ████████┐
                                    └── ████ 新 ████
Instance 3  ████████████████ 舊 ████████┐
                                        └── █ 新 ██
```

| 優點 | 缺點 |
|------|------|
| **零停機** | 短暫同時存在新舊版本 |
| 逐步降低風險 | 需要新舊版本相容 |

### 3. 藍綠部署（Blue-Green）

同時準備兩套環境，部署好新版本後一次切換流量。

```
┌──────────────────────────────────────────┐
│               Load Balancer              │
│               ┌─────────┐                │
│               │ Router  │                │
│               └────┬────┘                │
│         ┌──────────┴──────────┐          │
│         │                     │          │
│         ▼                     ▼          │
│  ┌─────────────┐      ┌─────────────┐   │
│  │  Blue (舊)  │      │ Green (新)  │   │
│  │  v1.0.0     │      │  v1.1.0     │   │
│  │  (運行中)   │      │  (待命)     │   │
│  └─────────────┘      └─────────────┘   │
│                                          │
│  切換：將 Router 從 Blue 指向 Green       │
│  回滾：將 Router 從 Green 指回 Blue       │
└──────────────────────────────────────────┘
```

| 優點 | 缺點 |
|------|------|
| **零停機** | 需要雙倍的基礎設施 |
| 瞬間切換，風險低 | 成本較高 |
| **回滾極快**（切回去就好） | |

### 4. 金絲雀部署（Canary）

先將一小部分流量導向新版本，觀察沒問題後再逐步擴大。

```
┌──────────────────────────────────────────────────┐
│                 Load Balancer                     │
│                                                  │
│  階段一：5% 流量 → 新版本                          │
│  ┌──────────────────┐  ┌─────┐                   │
│  │  舊版 v1.0 (95%) │  │ 新版│ (5%)              │
│  └──────────────────┘  └─────┘                   │
│                                                  │
│  階段二：觀察 metrics 正常，擴大到 25%              │
│  ┌───────────────┐  ┌────────┐                   │
│  │  舊版 v1.0    │  │ 新版   │ (25%)             │
│  │  (75%)        │  │ v1.1   │                   │
│  └───────────────┘  └────────┘                   │
│                                                  │
│  階段三：確認沒問題，全量切換                        │
│  ┌──────────────────────────┐                    │
│  │      新版 v1.1 (100%)    │                    │
│  └──────────────────────────┘                    │
└──────────────────────────────────────────────────┘
```

| 優點 | 缺點 |
|------|------|
| **風險最低**，問題只影響少數使用者 | 最複雜 |
| 可以基於 metrics 自動決策 | 需要良好的 monitoring |
| 適合高流量服務 | 需要流量控制能力 |

### 策略比較

| 策略 | 停機 | 風險 | 複雜度 | 適用場景 |
|------|------|------|--------|---------|
| Recreate | 有 | 高 | 低 | 開發/測試環境 |
| Rolling Update | 無 | 中 | 中 | 一般 production |
| Blue-Green | 無 | 低 | 中高 | 需要快速回滾 |
| Canary | 無 | 最低 | 高 | 高流量、高風險變更 |

Snow 用社團網站來跟 Ocean 解釋這四種策略的差別：

- **Recreate**：「把舊網站關掉，部署好新版再打開。簡單，但社員會看到一段時間的 502。」
- **Rolling Update**：「如果跑了 3 台伺服器，一台一台換。隨時都有服務在線，但短暫會有新舊版同時存在。」
- **Blue-Green**：「開一組全新的伺服器部署新版，確認沒問題後把流量切過去。有問題的話一秒切回來。」
- **Canary**：「先讓 5% 的流量跑新版，觀察一陣子沒問題再慢慢放大到 100%。最安全，但也最複雜。」

Andrew 聽完說：「我們社團網站用 Rolling Update 就夠了吧。」Snow 點點頭。

---

## 實戰：部署到 GCP Cloud Run

### Cloud Run 是什麼？

**Google Cloud Run** 是 Google Cloud Platform（GCP）提供的 **serverless 容器平台**。你只需要提供一個 Docker image，Cloud Run 就會自動：

- 啟動容器來處理 HTTP 請求
- 根據流量自動擴縮容（auto-scaling），沒有流量時可以縮到 0 個 instance
- 自動處理 HTTPS 憑證和域名
- 按實際使用量計費（pay-per-use）

對於我們的範例 Go HTTP API 伺服器來說，Cloud Run 是非常合適的部署平台。

### 完整的部署 Workflow

請在 `.github/workflows/deploy.yml` 中建立此檔案：

```yaml
name: Deploy to Cloud Run

on:
  push:
    branches: [main]
  workflow_dispatch:

env:
  PROJECT_ID: my-project
  REGION: asia-east1
  SERVICE_NAME: my-app

jobs:
  deploy:
    name: Deploy
    runs-on: ubuntu-latest
    environment: production
    permissions:
      contents: read
      id-token: write
    steps:
      - uses: actions/checkout@v4

      - name: Authenticate to GCP
        uses: google-github-actions/auth@v2
        with:
          workload_identity_provider: ${{ secrets.WIF_PROVIDER }}
          service_account: ${{ secrets.SA_EMAIL }}

      - name: Setup Cloud SDK
        uses: google-github-actions/setup-gcloud@v2

      - name: Build and push Docker image
        run: |
          gcloud builds submit --tag gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:${{ github.sha }}

      - name: Deploy to Cloud Run
        run: |
          gcloud run deploy ${{ env.SERVICE_NAME }} \
            --image gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:${{ github.sha }} \
            --region ${{ env.REGION }} \
            --platform managed \
            --allow-unauthenticated

      - name: Show deployment URL
        run: |
          URL=$(gcloud run services describe ${{ env.SERVICE_NAME }} --region ${{ env.REGION }} --format 'value(status.url)')
          echo "Deployed to: ${URL}"
```

### 逐段解說

#### 觸發條件

```yaml
on:
  push:
    branches: [main]
  workflow_dispatch:
```

- **push to main**：當程式碼合併到 `main` 時自動部署
- **workflow_dispatch**：允許在 GitHub UI 上 **手動觸發** 部署（在 Actions 頁面會出現 "Run workflow" 按鈕）

`workflow_dispatch` 非常實用，在需要重新部署但不想改程式碼時，可以手動觸發。

#### Environment 設定

```yaml
environment: production
```

這行告訴 GitHub 這個 job 使用 `production` 環境。如果 `production` 環境設定了 protection rules（例如需要核准），workflow 會在執行此 job 前暫停，等待核准。

#### Permissions

```yaml
permissions:
  contents: read
  id-token: write
```

`id-token: write` 是 **Workload Identity Federation** 所需的權限，用於讓 GitHub Actions 向 GCP 證明自己的身份。

#### 認證到 GCP

```yaml
- name: Authenticate to GCP
  uses: google-github-actions/auth@v2
  with:
    workload_identity_provider: ${{ secrets.WIF_PROVIDER }}
    service_account: ${{ secrets.SA_EMAIL }}
```

使用 **Workload Identity Federation (WIF)** 認證。WIF 是一種不需要管理長期金鑰的認證方式，GitHub Actions 向 GCP 出示自己的 OIDC token，GCP 驗證後授予臨時權限。

這比傳統的 service account key 更安全，因為：
- 不需要儲存長期的金鑰
- 臨時 token 會自動過期
- 可以限制特定 repository 和分支才能認證

#### 建置與部署

```yaml
- name: Build and push Docker image
  run: |
    gcloud builds submit --tag gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:${{ github.sha }}

- name: Deploy to Cloud Run
  run: |
    gcloud run deploy ${{ env.SERVICE_NAME }} \
      --image gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:${{ github.sha }} \
      --region ${{ env.REGION }} \
      --platform managed \
      --allow-unauthenticated
```

1. **`gcloud builds submit`**：使用 Cloud Build 在 GCP 端建置 Docker image，並推送到 Google Container Registry
2. **`gcloud run deploy`**：將 image 部署到 Cloud Run

使用 `${{ github.sha }}` 作為 image tag，確保每次部署都使用唯一的、可追溯的 tag。

<details>
<summary>💡 講師提示</summary>

> 如果沒有 GCP 帳號也沒關係，重點是讓學生理解部署 workflow 的結構和概念。實際的 GCP 設定可以在課後作為延伸練習。

</details>

---

## 其他平台部署方式簡介

### AWS — ECS/ECR

**Amazon Elastic Container Service (ECS)** 搭配 **Elastic Container Registry (ECR)** 是 AWS 上的主流容器部署方案。

```yaml
jobs:
  deploy-aws:
    name: Deploy to AWS ECS
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: ${{ secrets.AWS_ROLE_ARN }}
          aws-region: ap-northeast-1

      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v2

      - name: Build and push image to ECR
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          IMAGE_TAG: ${{ github.sha }}
        run: |
          docker build -t $ECR_REGISTRY/myapp:$IMAGE_TAG .
          docker push $ECR_REGISTRY/myapp:$IMAGE_TAG

      - name: Deploy to ECS
        uses: aws-actions/amazon-ecs-deploy-task-definition@v2
        with:
          task-definition: task-definition.json
          service: my-service
          cluster: my-cluster
```

概念流程：設定 AWS 認證 → 登入 ECR → 建置並推送 image → 更新 ECS service。

### Vercel

**Vercel** 主要用於前端專案和 serverless 函式，但也支援其他類型的應用。

```yaml
jobs:
  deploy-vercel:
    name: Deploy to Vercel
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Deploy to Vercel
        uses: amondnet/vercel-action@v25
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
          vercel-project-id: ${{ secrets.VERCEL_PROJECT_ID }}
          vercel-args: '--prod'
```

Vercel 的優勢在於設定極為簡單，幾乎零配置就能部署。

### Fly.io

**Fly.io** 是一個開發者友善的容器平台，特別適合部署後端應用。

```yaml
jobs:
  deploy-fly:
    name: Deploy to Fly.io
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - name: Deploy
        run: flyctl deploy --remote-only
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

Fly.io 的 `flyctl deploy` 會自動讀取專案中的 `fly.toml` 和 `Dockerfile`，一個指令就完成建置和部署。

### 各平台比較

| 平台 | 特色 | 免費額度 | 適合 |
|------|------|---------|------|
| **GCP Cloud Run** | Serverless 容器 | 每月 200 萬次請求免費 | 後端 API |
| **AWS ECS** | 企業級容器服務 | Free tier 有限 | 企業級應用 |
| **Vercel** | 零配置部署 | 個人免費 | 前端、Serverless |
| **Fly.io** | 開發者友善 | 免費 3 個小型 VM | 後端、Side project |

<details>
<summary>💡 講師提示</summary>

> 根據學生的背景選擇重點介紹的平台。如果學生未來會接觸 GCP，就重點講 Cloud Run；如果是做 side project，Fly.io 可能更實用。

</details>

---

## 多環境部署（Staging → Production）

### 為什麼需要多環境？

在正式的開發流程中，程式碼通常會先部署到 **staging 環境** 進行驗證，確認沒問題後再部署到 **production 環境**。

```
┌────────┐     ┌────────────┐     ┌──────────────┐
│  Code  │────▶│  Staging   │────▶│  Production  │
│  Merge │     │  (自動)     │     │  (手動核准)   │
└────────┘     └────────────┘     └──────────────┘
                    │                     │
                 自動部署              需要核准者
                 自動測試              點擊 Approve
```

### Staging + Production 的 Workflow 範例

```yaml
name: Deploy Pipeline

on:
  push:
    branches: [main]

env:
  PROJECT_ID: my-project
  SERVICE_NAME: my-app

jobs:
  # ──────────────────────────────────────
  # Stage 1: Deploy to Staging (automatic)
  # ──────────────────────────────────────
  deploy-staging:
    name: Deploy to Staging
    runs-on: ubuntu-latest
    environment: staging
    permissions:
      contents: read
      id-token: write
    steps:
      - uses: actions/checkout@v4

      - name: Authenticate to GCP
        uses: google-github-actions/auth@v2
        with:
          workload_identity_provider: ${{ secrets.WIF_PROVIDER }}
          service_account: ${{ secrets.SA_EMAIL }}

      - name: Setup Cloud SDK
        uses: google-github-actions/setup-gcloud@v2

      - name: Deploy to Staging
        run: |
          gcloud run deploy ${{ env.SERVICE_NAME }}-staging \
            --image gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:${{ github.sha }} \
            --region asia-east1 \
            --platform managed

      - name: Run smoke tests
        run: |
          STAGING_URL=$(gcloud run services describe ${{ env.SERVICE_NAME }}-staging \
            --region asia-east1 --format 'value(status.url)')
          # Verify health endpoint
          STATUS=$(curl -s -o /dev/null -w "%{http_code}" ${STAGING_URL}/health)
          if [ "$STATUS" != "200" ]; then
            echo "::error::Smoke test failed! Health check returned ${STATUS}"
            exit 1
          fi
          echo "Smoke test passed! Staging is healthy."

  # ──────────────────────────────────────────
  # Stage 2: Deploy to Production (manual approval)
  # ──────────────────────────────────────────
  deploy-production:
    name: Deploy to Production
    needs: deploy-staging
    runs-on: ubuntu-latest
    environment: production
    permissions:
      contents: read
      id-token: write
    steps:
      - uses: actions/checkout@v4

      - name: Authenticate to GCP
        uses: google-github-actions/auth@v2
        with:
          workload_identity_provider: ${{ secrets.WIF_PROVIDER }}
          service_account: ${{ secrets.SA_EMAIL }}

      - name: Setup Cloud SDK
        uses: google-github-actions/setup-gcloud@v2

      - name: Deploy to Production
        run: |
          gcloud run deploy ${{ env.SERVICE_NAME }} \
            --image gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:${{ github.sha }} \
            --region asia-east1 \
            --platform managed \
            --allow-unauthenticated

      - name: Verify deployment
        run: |
          PROD_URL=$(gcloud run services describe ${{ env.SERVICE_NAME }} \
            --region asia-east1 --format 'value(status.url)')
          echo "Production deployed to: ${PROD_URL}"
```

### 關鍵設計

1. **Staging 自動部署**：push 到 `main` 後自動部署到 staging
2. **Smoke Test**：在 staging 部署後自動執行基本的健康檢查
3. **Production 手動核准**：`environment: production` 設定了 protection rules，workflow 會暫停等待核准
4. **相同的 image**：staging 和 production 使用 **同一個 image**（同一個 `github.sha`），確保測試過的就是要部署的

### 設定 Environment Protection Rules

要讓 production 部署需要核准，在 GitHub 上設定：

1. **Settings** → **Environments** → **production**
2. 勾選 **Required reviewers**
3. 加入需要核准的人（例如 team lead）
4. 可選：設定 **Wait timer**（例如 5 分鐘的冷靜期）

設定完成後，當 workflow 執行到 `deploy-production` job 時，會出現一個等待核准的畫面，核准者需要到 GitHub Actions 頁面點擊 **Approve and deploy** 才會繼續。

<details>
<summary>💡 講師提示</summary>

> 可以現場展示等待核准的畫面。如果時間允許，讓兩位學生組隊，一個人 push code 觸發部署，另一個人去核准。

</details>

---

## 部署通知

### 使用 Slack 通知

部署完成（或失敗）後，發送通知到 Slack 讓團隊知道：

```yaml
  notify:
    name: Notify Deployment
    needs: [deploy-production]
    runs-on: ubuntu-latest
    if: always()
    steps:
      - name: Notify Slack on success
        if: needs.deploy-production.result == 'success'
        uses: slackapi/slack-github-action@v2
        with:
          webhook: ${{ secrets.SLACK_WEBHOOK_URL }}
          webhook-type: incoming-webhook
          payload: |
            {
              "text": "Deployment Successful",
              "blocks": [
                {
                  "type": "section",
                  "text": {
                    "type": "mrkdwn",
                    "text": "*Deployment Successful* :white_check_mark:\n*Service:* ${{ env.SERVICE_NAME }}\n*Version:* `${{ github.sha }}`\n*Deployed by:* ${{ github.actor }}\n*Commit:* ${{ github.event.head_commit.message }}"
                  }
                }
              ]
            }

      - name: Notify Slack on failure
        if: needs.deploy-production.result == 'failure'
        uses: slackapi/slack-github-action@v2
        with:
          webhook: ${{ secrets.SLACK_WEBHOOK_URL }}
          webhook-type: incoming-webhook
          payload: |
            {
              "text": "Deployment Failed",
              "blocks": [
                {
                  "type": "section",
                  "text": {
                    "type": "mrkdwn",
                    "text": "*Deployment Failed* :x:\n*Service:* ${{ env.SERVICE_NAME }}\n*Commit:* `${{ github.sha }}`\n*Action:* <${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}|View logs>"
                  }
                }
              ]
            }
```

#### 關鍵設計

- **`if: always()`**：不管前面的 job 成功或失敗，通知 job 都會執行
- **`needs.deploy-production.result`**：根據部署結果發送不同的通知
- **包含關鍵資訊**：服務名稱、版本（commit SHA）、部署者、commit 訊息
- **失敗時附上 log 連結**：方便快速排查問題

<details>
<summary>💡 講師提示</summary>

> Slack webhook 的設定步驟可以在課後補充說明。重點是讓學生理解「部署後通知團隊」這個概念。如果團隊不用 Slack，也可以用 Discord、Microsoft Teams 或 email 的 Action。

</details>

---

## Rollback（回滾）策略

### 使用 workflow_dispatch 手動觸發 Rollback

當部署出了問題，需要快速回到前一個版本。以下是一個 rollback workflow：

```yaml
name: Rollback

on:
  workflow_dispatch:
    inputs:
      target_sha:
        description: 'Commit SHA to rollback to'
        required: true
        type: string
      environment:
        description: 'Target environment'
        required: true
        type: choice
        options:
          - staging
          - production

env:
  PROJECT_ID: my-project
  REGION: asia-east1
  SERVICE_NAME: my-app

jobs:
  rollback:
    name: Rollback to ${{ github.event.inputs.target_sha }}
    runs-on: ubuntu-latest
    environment: ${{ github.event.inputs.environment }}
    permissions:
      contents: read
      id-token: write
    steps:
      - name: Authenticate to GCP
        uses: google-github-actions/auth@v2
        with:
          workload_identity_provider: ${{ secrets.WIF_PROVIDER }}
          service_account: ${{ secrets.SA_EMAIL }}

      - name: Setup Cloud SDK
        uses: google-github-actions/setup-gcloud@v2

      - name: Rollback deployment
        run: |
          SERVICE=${{ env.SERVICE_NAME }}
          if [ "${{ github.event.inputs.environment }}" = "staging" ]; then
            SERVICE="${SERVICE}-staging"
          fi
          gcloud run deploy ${SERVICE} \
            --image gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:${{ github.event.inputs.target_sha }} \
            --region ${{ env.REGION }} \
            --platform managed

      - name: Verify rollback
        run: |
          echo "Rolled back to commit: ${{ github.event.inputs.target_sha }}"
          echo "Environment: ${{ github.event.inputs.environment }}"
```

### 使用方式

1. 到 GitHub Actions 頁面
2. 選擇 **Rollback** workflow
3. 點擊 **Run workflow**
4. 輸入要回滾到的 commit SHA（可以從 git log 或上一次成功部署的紀錄中取得）
5. 選擇環境（staging 或 production）
6. 點擊 **Run workflow** 執行

### 為什麼用 Commit SHA 而不是 Tag？

因為我們在部署時是用 `github.sha` 作為 Docker image tag，所以每次部署的 image 都可以用 commit SHA 精確定位。這比用 `latest` tag 更可靠。

```
Image tags in registry:
  gcr.io/my-project/my-app:a1b2c3d  ← 前前一次部署
  gcr.io/my-project/my-app:e4f5g6h  ← 前一次部署（要回滾到這裡）
  gcr.io/my-project/my-app:i7j8k9l  ← 目前有問題的版本
```

<details>
<summary>💡 講師提示</summary>

> 可以提醒學生，回滾只是 **暫時的應急措施**。回滾後還是需要找出問題、修復、然後重新部署。重點是「先止血」。

</details>

---

## 安全最佳實踐

### 1. 最小權限原則（Principle of Least Privilege）

只給 workflow 需要的最小權限，不多給：

```yaml
# Good — only the permissions needed
permissions:
  contents: read
  id-token: write

# Bad — too many permissions
permissions: write-all
```

同樣的原則也適用於 GCP/AWS 的 service account，只給部署所需的權限，不要給 admin。

### 2. 定期輪換 Secrets

- 定期更換 API key、token 等 secrets
- 如果有人離開團隊，立即更換所有 secrets
- 可以設定提醒（例如每 90 天輪換一次）

### 3. 使用 OIDC 取代長期金鑰

**OIDC（OpenID Connect）** 讓 GitHub Actions 可以用臨時 token 向雲端平台認證，不需要儲存長期金鑰：

```
┌──────────────┐          ┌──────────────┐
│   GitHub     │  OIDC    │    GCP /     │
│   Actions    │  Token   │    AWS       │
│              │─────────▶│              │
│              │          │  驗證後發     │
│              │◀─────────│  臨時憑證     │
│              │  短期     │              │
│              │  Token    │              │
└──────────────┘          └──────────────┘
```

| 方式 | 安全性 | 管理成本 |
|------|--------|---------|
| **長期金鑰（Service Account Key）** | 低，金鑰可能洩漏 | 高，需要輪換 |
| **OIDC（Workload Identity Federation）** | 高，臨時 token 自動過期 | 低，不需要管理金鑰 |

我們在 Cloud Run 部署 workflow 中使用的 `google-github-actions/auth@v2` 就是用 OIDC 方式。

### 4. Audit Log

- 在 GitHub 上可以查看 **Audit Log**（Settings → Audit Log），追蹤誰做了什麼操作
- 記錄每次部署的 commit SHA、部署者、時間
- 部署通知（如 Slack 通知）也是一種 audit trail

### 安全檢查清單

| 項目 | 說明 |
|------|------|
| ✅ 使用 OIDC 認證 | 不儲存長期金鑰 |
| ✅ 設定 `permissions` | 每個 workflow 都明確設定最小權限 |
| ✅ 使用 Environment protection | Production 部署需要核准 |
| ✅ 不在程式碼中硬寫 secrets | 使用 GitHub Secrets |
| ✅ Pin action 版本 | 使用 `@v4` 而非 `@main` |
| ✅ 限制 fork PR 的權限 | Fork PR 不能存取 secrets |
| ✅ 定期輪換 secrets | 設定提醒，定期更換 |
| ✅ 保留部署紀錄 | 透過通知和 audit log 追蹤 |

<details>
<summary>💡 講師提示</summary>

> 安全是一個很大的主題，這裡只列出最基本的幾個重點。可以鼓勵有興趣的學生去看 GitHub 的 [Security hardening for GitHub Actions](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions) 官方文件。

</details>

---

## 常見問題排解

### 1. GCP 認證失敗

最常見的原因是 Workload Identity Federation 的設定不正確。確認以下幾點：
- `WIF_PROVIDER` 和 `SA_EMAIL` 兩個 secrets 都已正確設定
- Service account 有部署到 Cloud Run 的權限
- Workload Identity Pool 有允許你的 GitHub repository

如果是第一次設定，建議先參考 [Google 的官方教學](https://cloud.google.com/iam/docs/workload-identity-federation-with-deployment-pipelines)。

### 2. 部署成功但 Smoke Test 失敗

可能是 Cloud Run 的新版本還沒完全啟動就開始打 health check。可以在 curl 前加一個短暫的等待，或用 retry 機制：

```bash
for i in 1 2 3 4 5; do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" ${URL}/health)
  if [ "$STATUS" = "200" ]; then
    echo "Health check passed!"
    exit 0
  fi
  echo "Attempt $i: got ${STATUS}, retrying in 10s..."
  sleep 10
done
echo "Health check failed after 5 attempts"
exit 1
```

---

## 小結與練習題

### 本章重點回顧

- **GitHub Environments**：為不同環境（staging、production）設定獨立的 secrets 和 protection rules
- **Secrets vs 環境變數**：Secrets 存放敏感資料（自動遮蔽），環境變數存放非敏感設定
- **部署策略**：Recreate（簡單但有停機）、Rolling Update（漸進式）、Blue-Green（瞬間切換）、Canary（最安全）
- **Cloud Run 部署**：使用 Workload Identity Federation 認證，gcloud 指令建置與部署
- **多環境部署**：staging 自動部署 + smoke test → production 手動核准部署
- **部署通知**：使用 Slack webhook 通知團隊部署結果
- **Rollback**：使用 `workflow_dispatch` 搭配 commit SHA 手動回滾
- **安全最佳實踐**：最小權限、OIDC 認證、定期輪換 secrets、audit log

### 練習題

完成以下練習來鞏固本章所學：

👉 [練習三：進階練習](exercises/exercise-03-advanced.md)（練習 3-3）

---

### 完整 CI/CD Pipeline 全景圖

恭喜你完成了整個 CI/CD 工作坊！讓我們回顧一下，你現在已經有能力建立這樣一套完整的自動化流程：

```
開發者 push 程式碼
        │
        ▼
┌───────────────────────────────────────────────────────┐
│                   CI Pipeline (ch04)                   │
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
│              Release 自動化 (ch06)                     │
│                                                       │
│  git tag v1.0.0 ──▶ 多平台 Binary + Docker Image     │
│                 ──▶ GitHub Release + Changelog         │
└───────────────────────────────────────────────────────┘
        │
        ▼
┌───────────────────────────────────────────────────────┐
│                  CD Pipeline (ch07)                    │
│                                                       │
│  ┌──────────┐    ┌──────────┐    ┌──────────────┐    │
│  │ Deploy   │───▶│ Smoke    │───▶│ Deploy       │    │
│  │ Staging  │    │ Test     │    │ Production   │    │
│  │ (自動)   │    │          │    │ (手動核准)    │    │
│  └──────────┘    └──────────┘    └──────────────┘    │
└───────────────────────────────────────────────────────┘
```

從寫程式碼到部署上線，每個環節都有自動化的品質把關。這就是 CI/CD 的力量。

---

[← 上一章：Release 自動化](06-release-automation.md) ｜ [回到目錄 →](README.md)
