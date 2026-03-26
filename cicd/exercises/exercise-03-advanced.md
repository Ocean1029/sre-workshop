# 練習三：進階練習

> **難度：** 進階 | **對應章節：** 06-Release 自動化、07-部署到雲端平台

---

## 目錄

- [練習 3-1：Release 流程練習](#練習-3-1release-流程練習)
- [練習 3-2：Docker 建置與推送](#練習-3-2docker-建置與推送)
- [練習 3-3：Deployment Workflow 設計](#練習-3-3deployment-workflow-設計)
- [練習 3-4：自訂 Composite Action（選做）](#練習-3-4自訂-composite-action選做)
- [延伸思考](#延伸思考)

---

## 練習 3-1：Release 流程練習

### 目標

完成一次完整的 **release 流程**——從建立 Git tag 到 GitHub Release 頁面出現 binary 和 changelog。

### 要求

1. 確認 `.github/workflows/release.yml` 已經設定好（參考教材第六章）
2. 確認 `main.go` 中有 `var version = "dev"` 變數
3. 建立一個 Git tag `v0.1.0`
4. Push tag 到 GitHub
5. 觀察 Release workflow 的執行過程
6. 確認 GitHub Release 頁面有正確的 binary 和 changelog

### Step by Step 指令

#### Step 1：確認 release.yml 已設定

確認 `.github/workflows/release.yml` 已經存在，並且包含 tag 觸發的設定：

```yaml
on:
  push:
    tags:
      - 'v*'
```

如果還沒有建立，請參考教材第六章的完整 release workflow。

#### Step 2：確認程式碼已 push 到 main

```bash
# Make sure everything is committed and pushed
git status
git push origin main
```

#### Step 3：建立 Annotated Tag

```bash
# Create an annotated tag
git tag -a v0.1.0 -m "First release: basic HTTP API server"

# Verify the tag was created
git tag -l
git show v0.1.0
```

#### Step 4：Push Tag 到 GitHub

```bash
# Push the specific tag to remote
git push origin v0.1.0
```

> **注意：** 使用 `git push origin v0.1.0` 只推送特定的 tag。

#### Step 5：觀察 Release Workflow

1. 前往 GitHub repository 的 **Actions** 分頁
2. 你應該會看到一個名為 **"Release"** 的 workflow run 正在執行
3. 點擊進去觀察每個 job 的執行狀態：
   - **Test** — 執行測試，確保程式碼正確
   - **Build** — 為多個平台建置 binary（觀察 matrix 產生的多個 job）
   - **Release** — 建立 GitHub Release 並上傳 binary

#### Step 6：確認 GitHub Release

1. 前往 repository 的 **Releases** 頁面（在 Code 分頁的右側欄）
2. 確認以下內容：
   - Release 標題包含 `v0.1.0`
   - 有自動產生的 changelog
   - Assets 區域有多個平台的 binary 可供下載

### 驗證清單

- [ ] Tag `v0.1.0` 成功建立並推送
- [ ] Release workflow 所有 job 都通過（綠色勾勾）
- [ ] GitHub Release 頁面有 `v0.1.0` 的 release
- [ ] Release 中有可下載的 binary 檔案
- [ ] Changelog 中有 commit 訊息

---

## 練習 3-2：Docker 建置與推送

### 目標

在 CI workflow 中加入 **Docker image 建置**，並推送到 **GitHub Container Registry（GHCR）**。

### 要求

1. 在 `ci.yml`（或 `release.yml`）中加入一個 `docker` job
2. 使用範例專案中已有的 **多階段建置 Dockerfile**
3. 推送到 GitHub Container Registry（`ghcr.io`）
4. 為 image 加上兩個 tag：
   - **commit SHA**（如 `ghcr.io/<owner>/<repo>:abc123d`）
   - **latest**（如 `ghcr.io/<owner>/<repo>:latest`）

### 提示

- 使用 `docker/login-action@v3` 來登入 GHCR。
- 使用 `docker/build-push-action@v6` 來建置和推送 image。
- `${{ secrets.GITHUB_TOKEN }}` 是 GitHub 自動提供的 token，不需要額外設定。
- `${{ github.sha }}` 可以取得完整的 commit SHA，但 Docker tag 通常用前 7 碼就好（可以用 `${{ github.sha }}` 的完整值或自行截斷）。
- 記得在 workflow 的 `permissions` 中加入 `packages: write`。

### 預期結果

- Docker image 成功推送到 `ghcr.io/<your-username>/<your-repo>`。
- 在 GitHub repository 的右側欄 **Packages** 中可以看到你的 Docker image。
- Image 有 `latest` 和 commit SHA 兩個 tag。

<details>
<summary>點擊查看答案</summary>

```yaml
name: Go CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

permissions:
  contents: read
  packages: write

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest

  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Show coverage
        run: go tool cover -func=coverage.out

  build:
    name: Build
    needs: [lint, test]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Build binary
        run: go build -o bin/app ./...
      - name: Upload binary
        uses: actions/upload-artifact@v4
        with:
          name: app-binary
          path: bin/app

  docker:
    name: Build & Push Docker Image
    needs: [lint, test]
    runs-on: ubuntu-latest
    # Only push image on main branch (skip for PRs)
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4

      - name: Login to GitHub Container Registry
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata for Docker
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ghcr.io/${{ github.repository }}
          tags: |
            type=sha,prefix=
            type=raw,value=latest

      - name: Build and push Docker image
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          build-args: |
            VERSION=${{ github.sha }}
```

**重點說明：**

- **`permissions: packages: write`** 是推送到 GHCR 所必需的權限，放在 workflow 最上層，讓所有 job 都能使用。
- **`if: github.event_name == 'push' && github.ref == 'refs/heads/main'`** 確保只有在 push 到 main 時才推送 image。PR 中只做測試，不推送 image。
- **`docker/metadata-action@v5`** 是一個很好用的 action，可以自動產生 image tag 和 label。`type=sha` 會用 commit SHA 的前 7 碼作為 tag。
- **`docker/build-push-action@v6`** 負責建置和推送 image。`build-args` 可以傳遞 build argument 給 Dockerfile（例如版本資訊）。
- **`docker` job 與 `build` job 平行執行**：兩者都 `needs: [lint, test]`，所以 lint 和 test 通過後，docker 和 build 會同時開始。

</details>

---

## 練習 3-3：Deployment Workflow 設計

### 目標

設計一個完整的 **部署 workflow**，涵蓋 staging 和 production 兩個環境，包含自動部署、手動核准和部署通知。

> 這是一個 **設計題**，不提供完整答案。目的是讓你思考真實世界的部署流程應該如何設計。

### 要求

1. **Staging 環境**：push to `main` 時自動部署
2. **Production 環境**：需要手動核准才能部署
3. **部署後驗證**：部署完成後進行 smoke test（健康檢查）
4. **失敗通知**：部署失敗時發送通知

### 參考架構圖

```
┌──────────────────────────────────────────────────────────────┐
│                    Deploy Pipeline                           │
│                                                              │
│  ┌──────────┐                                                │
│  │  Push to  │                                               │
│  │   main    │                                               │
│  └────┬─────┘                                                │
│       │                                                      │
│       ▼                                                      │
│  ┌──────────────────┐                                        │
│  │  Deploy Staging   │  ← 自動執行                            │
│  │  (environment:    │                                       │
│  │   staging)        │                                       │
│  └────────┬─────────┘                                        │
│           │                                                  │
│           ▼                                                  │
│  ┌──────────────────┐                                        │
│  │  Smoke Test       │  ← curl /health, 確認回應 200          │
│  │  (staging)        │                                       │
│  └────────┬─────────┘                                        │
│           │                                                  │
│           │ (等待手動核准)                                     │
│           ▼                                                  │
│  ┌──────────────────┐                                        │
│  │  Deploy Prod      │  ← 需要 reviewer 核准                  │
│  │  (environment:    │                                       │
│  │   production)     │                                       │
│  └────────┬─────────┘                                        │
│           │                                                  │
│           ▼                                                  │
│  ┌──────────────────┐                                        │
│  │  Smoke Test       │  ← 確認 production 正常                │
│  │  (production)     │                                       │
│  └────────┬─────────┘                                        │
│           │                                                  │
│           ▼                                                  │
│  ┌──────────────────┐                                        │
│  │  Notify           │  ← 成功或失敗都要通知                   │
│  │  (always)         │                                       │
│  └──────────────────┘                                        │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### 設計提示

1. **Environment 設定**
   - 在 GitHub repository 的 Settings 中建立 `staging` 和 `production` 兩個 environment。
   - 為 `production` 設定 **Required reviewers**。

2. **Smoke Test 設計**
   - 部署完成後，用 `curl` 打 health endpoint。
   - 檢查 HTTP status code 是否為 200。
   - 如果不是 200，用 `exit 1` 讓 workflow 失敗。

3. **通知設計**
   - 使用 `if: always()` 確保不管前面的 job 成功或失敗，通知 job 都會執行。
   - 使用 `needs.<job_id>.result` 來判斷前面 job 的結果。

4. **安全考量**
   - 使用 `permissions` 限制最小權限。
   - Secrets 只在需要的 environment 中定義。
   - 使用 OIDC 認證，不使用長期金鑰。

### 討論問題

和同學一起討論以下問題：

1. 如果 staging 的 smoke test 失敗了，整個 workflow 會怎樣？production 會被部署嗎？
2. 如果 production 部署成功但 smoke test 失敗，代表什麼？你應該怎麼做？
3. 如果需要緊急回滾（rollback），你會怎麼設計？
4. 在這個 workflow 中，staging 和 production 使用的是同一個 Docker image 嗎？為什麼這很重要？

---

## 練習 3-4：自訂 Composite Action（選做）

### 目標

建立一個 **可重用的 composite action**，把 Go 專案中常用的設定步驟（checkout、setup-go、cache、install dependencies）封裝成一個 action，讓多個 workflow 可以共用。

### 要求

1. 建立 `.github/actions/setup-go-project/action.yml`
2. 這個 action 封裝以下步驟：
   - Checkout 程式碼
   - 安裝指定版本的 Go
   - 下載 Go module dependencies
3. 接受 `go-version` 作為 **input 參數**（預設值為 `1.22`）
4. 在 CI workflow 中使用這個自訂 action

### 提示

- Composite action 的檔案必須命名為 `action.yml`。
- 使用 `inputs` 定義參數，用 `${{ inputs.go-version }}` 在 step 中引用。
- Composite action 的 step 必須加上 `shell: bash`（或其他 shell）。
- 在 workflow 中使用時，路徑是相對於 repo 根目錄的，例如 `uses: ./.github/actions/setup-go-project`。

### 預期結果

- `ci.yml` 中的 lint、test、build job 都可以用一行 `uses` 取代原本的 checkout + setup-go 步驟。
- 修改 Go 版本時，只需要在一個地方修改就好。

<details>
<summary>點擊查看答案</summary>

#### 1. 建立 Composite Action

在 `.github/actions/setup-go-project/action.yml` 中：

```yaml
name: Setup Go Project
description: Checkout code, install Go, and download dependencies

inputs:
  go-version:
    description: 'Go version to install'
    required: false
    default: '1.22'

runs:
  using: composite
  steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Setup Go ${{ inputs.go-version }}
      uses: actions/setup-go@v5
      with:
        go-version: ${{ inputs.go-version }}

    - name: Download dependencies
      shell: bash
      run: go mod download

    - name: Verify dependencies
      shell: bash
      run: go mod verify
```

#### 2. 在 CI Workflow 中使用

修改 `ci.yml`，將原本的 checkout + setup-go 步驟替換成自訂 action：

```yaml
name: Go CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - name: Setup Go Project
        uses: ./.github/actions/setup-go-project
        with:
          go-version: '1.22'
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest

  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - name: Setup Go Project
        uses: ./.github/actions/setup-go-project
        with:
          go-version: '1.22'
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Show coverage
        run: go tool cover -func=coverage.out
      - name: Upload coverage
        uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: coverage.out

  build:
    name: Build
    needs: [lint, test]
    runs-on: ubuntu-latest
    steps:
      - name: Setup Go Project
        uses: ./.github/actions/setup-go-project
        with:
          go-version: '1.22'
      - name: Build binary
        run: go build -o bin/app ./...
      - name: Upload binary
        uses: actions/upload-artifact@v4
        with:
          name: app-binary
          path: bin/app
```

**重點說明：**

- **Composite action** 讓你可以把多個步驟封裝成一個可重用的 action。相較於直接在每個 job 中重複寫 checkout + setup-go，composite action 讓維護更輕鬆。
- **`uses: ./.github/actions/setup-go-project`** 是使用本地 action 的語法。路徑是相對於 repo 根目錄。
- **`shell: bash`** 在 composite action 的 `run` step 中是必須的（與一般 workflow 不同，一般 workflow 會自動推斷 shell）。
- **集中管理版本**：如果你想統一升級到 Go 1.23，只需要修改 `action.yml` 的 `default` 值，或在各個 workflow 中修改 `with: go-version`。
- **可擴充性**：你還可以在 composite action 中加入更多步驟，例如安裝 linter 工具、設定環境變數等。
- **注意**：`golangci-lint-action` 已經會自己做 checkout，所以在 lint job 中使用這個 composite action 時，checkout 實際上會執行兩次（一次在 composite action 中，一次在 golangci-lint-action 中）。這不會造成問題，但如果你想最佳化，可以考慮建立不含 checkout 的版本。

</details>

---

## 延伸思考

以下問題沒有標準答案，供進階學生思考和討論：

1. **Release 頻率**：你的團隊應該多久 release 一次？每天？每週？每個 PR merge 就 release？不同策略各有什麼優缺點？

2. **Docker Image 安全性**：在 CI 中建置的 Docker image，你會做哪些安全掃描？（提示：搜尋 Trivy、Snyk 等工具。）

3. **Composite Action vs Reusable Workflow**：GitHub Actions 提供了兩種程式碼重用機制——composite action 和 reusable workflow。它們的差異是什麼？各自適合什麼場景？

4. **GitOps**：如果你的部署不是直接在 CI 中執行 `gcloud run deploy`，而是修改一個 Git repository 中的 Kubernetes manifest，然後讓 ArgoCD 自動同步，這種模式叫什麼？它有什麼優勢？

5. **Multi-Arch Docker Image**：如果你想讓你的 Docker image 同時支援 `amd64` 和 `arm64` 架構，你會怎麼做？（提示：搜尋 `docker buildx` 和 `docker/build-push-action` 的 `platforms` 參數。）

---

[← 練習二：CI Pipeline 實戰練習](exercise-02-ci-pipeline.md) ｜ [回到教材目錄 →](../README.md)
