# 02 — GitHub Actions 基礎概念

> **時間：15 分鐘**

---

## 目錄

- [GitHub Actions 是什麼？](#github-actions-是什麼)
- [核心概念](#核心概念)
- [概念關係圖](#概念關係圖)
- [Workflow 檔案結構](#workflow-檔案結構)
- [常見 Events（觸發器）](#常見-events觸發器)
- [常見 Actions](#常見-actions)
- [GitHub-hosted Runners](#github-hosted-runners)
- [免費額度與限制](#免費額度與限制)
- [YAML 語法速查](#yaml-語法速查)
- [小結](#小結)

---

## GitHub Actions 是什麼？

**GitHub Actions** 是 GitHub 內建的 CI/CD 與自動化平台。它讓你可以：

- 用 **YAML 檔案** 定義自動化流程（稱為 workflow）
- 在各種 **GitHub 事件** 發生時自動觸發執行（例如：push 程式碼、建立 Pull Request、發布 Release 等）
- 利用社群提供的 **可重用 Actions** 快速組裝你的 pipeline

簡單來說：**你只要把一個 YAML 檔案放進 repository，GitHub 就會幫你自動執行裡面定義的任務。**

> 💡 **講師提示：** 可以先打開一個有 GitHub Actions 的 repository，讓學生看到 Actions 分頁的介面，建立初步印象。

---

## 核心概念

GitHub Actions 有六個核心概念，了解它們之間的關係是使用 GitHub Actions 的基礎。

### 1. Workflow（工作流程）

**Workflow** 是一個定義自動化流程的 **YAML 檔案**。

- 檔案位置：必須放在 repository 的 **`.github/workflows/`** 目錄下
- 副檔名：`.yml` 或 `.yaml`
- 一個 repository 可以有 **多個** workflow
- 每個 workflow 可以被不同的事件觸發

```
your-repo/
├── .github/
│   └── workflows/
│       ├── ci.yml           # CI workflow — push 時觸發
│       ├── pr-check.yml     # PR 檢查 — 開 PR 時觸發
│       └── release.yml      # Release — 打 tag 時觸發
├── main.go
└── go.mod
```

### 2. Event（事件）

**Event** 是觸發 workflow 執行的條件。當特定事件發生時，GitHub 會自動啟動對應的 workflow。

常見事件包括：

- 推送程式碼（`push`）
- 建立或更新 Pull Request（`pull_request`）
- 定時排程（`schedule`）
- 手動觸發（`workflow_dispatch`）

```yaml
# Example: trigger on push to main branch
on:
  push:
    branches: [main]
```

### 3. Job（工作）

**Job** 是 workflow 中的 **執行單元**。

- 一個 workflow 可以包含一個或多個 job
- 預設情況下，多個 job **平行執行**
- 可以透過 `needs` 關鍵字設定 job 之間的 **依賴關係**，使其依序執行
- 每個 job 在一個獨立的 **Runner** 環境中執行

```yaml
jobs:
  # These two jobs run in parallel
  test:
    runs-on: ubuntu-latest
    steps:
      - run: echo "Running tests..."

  lint:
    runs-on: ubuntu-latest
    steps:
      - run: echo "Running linter..."

  # This job waits for both test and lint to finish
  deploy:
    needs: [test, lint]
    runs-on: ubuntu-latest
    steps:
      - run: echo "Deploying..."
```

### 4. Step（步驟）

**Step** 是 job 中的 **個別操作**。

- 每個 step 依序執行（同一個 job 內的 step 不能平行）
- Step 有兩種形式：
  - **`uses`**：使用一個現成的 Action
  - **`run`**：執行一段 shell 指令
- 同一個 job 內的 step 共享 **相同的檔案系統**

```yaml
steps:
  # Step 1: use an existing action
  - uses: actions/checkout@v4

  # Step 2: run a shell command
  - name: Run tests
    run: go test ./...

  # Step 3: use another action
  - uses: actions/upload-artifact@v4
    with:
      name: test-results
      path: results/
```

### 5. Action（動作）

**Action** 是一個 **可重用的步驟**，可以想像成是 workflow 中的「積木」。

- 來源：
  - **GitHub Marketplace**：社群提供的現成 Actions（如 `actions/checkout`）
  - **自訂 Action**：自己撰寫的 Action（放在 repository 中）
- 引用格式：`{owner}/{repo}@{version}`
  - 例如：`actions/checkout@v4`
  - 建議固定版本（用 `@v4` 而非 `@main`），避免未預期的變更

> 💡 **講師提示：** 可以用「樂高積木」的比喻——每個 Action 都是一塊積木，你可以把不同的積木組合起來，拼出你要的 pipeline。

### 6. Runner（執行器）

**Runner** 是實際執行 job 的 **機器環境**。

- **GitHub-hosted Runner**：GitHub 提供的雲端虛擬機，用完即丟
  - 不需要自己管理伺服器
  - 提供 Ubuntu、Windows、macOS 等作業系統
- **Self-hosted Runner**：你自己管理的機器
  - 適用於需要特殊硬體、軟體或網路環境的場景
  - 需要自行維護與更新

本課程使用 **GitHub-hosted Runner**，不需要額外設定。

---

## 概念關係圖

以下是 GitHub Actions 各概念之間的階層關係：

```
┌─────────────────────────────────────────────────────────────┐
│                        Workflow                             │
│                  (.github/workflows/ci.yml)                 │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ Event (觸發條件)                                      │  │
│  │ on: push, pull_request, schedule ...                  │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─────────────────────────┐  ┌─────────────────────────┐  │
│  │ Job 1: build            │  │ Job 2: deploy           │  │
│  │ (runs-on: ubuntu-latest)│  │ (needs: build)          │  │
│  │                         │  │                         │  │
│  │  ├── Step 1             │  │  ├── Step 1             │  │
│  │  │   uses: action       │  │  │   uses: action       │  │
│  │  ├── Step 2             │  │  └── Step 2             │  │
│  │  │   run: command       │  │      run: command       │  │
│  │  └── Step 3             │  │                         │  │
│  │      uses: action       │  │                         │  │
│  └─────────────────────────┘  └─────────────────────────┘  │
│        │                             ▲                      │
│        └─────── needs ───────────────┘                      │
└─────────────────────────────────────────────────────────────┘
```

用樹狀結構來看更清楚：

```
Workflow (.github/workflows/ci.yml)
├── Event (觸發條件)
│   └── on: push (branches: [main])
├── Job 1: build
│   ├── runs-on: ubuntu-latest
│   ├── Step 1 — uses: actions/checkout@v4
│   ├── Step 2 — run: go build ./...
│   └── Step 3 — uses: actions/upload-artifact@v4
└── Job 2: deploy
    ├── needs: [build]
    ├── runs-on: ubuntu-latest
    ├── Step 1 — uses: actions/checkout@v4
    └── Step 2 — run: ./deploy.sh
```

> 💡 **講師提示：** 建議在白板上畫出這個階層關係，一邊畫一邊說明。強調 Workflow 包含 Job，Job 包含 Step，Step 使用 Action 或執行指令。

---

## Workflow 檔案結構

以下是一個完整的 workflow YAML 範例，我們逐行來解說：

```yaml
# Workflow name — shown in the Actions tab on GitHub
name: My First Workflow

# Event — when to trigger this workflow
on:
  push:
    branches: [main]

# Jobs — what to do
jobs:
  build:
    # Runner — where to run this job
    runs-on: ubuntu-latest

    # Steps — individual tasks
    steps:
      # Step 1: Check out the repository code
      - uses: actions/checkout@v4

      # Step 2: Run a shell command
      - name: Say Hello
        run: echo "Hello, GitHub Actions!"
```

### 逐行解說

| 行 | 程式碼 | 說明 |
|----|--------|------|
| 1 | `name: My First Workflow` | 為 workflow 取名，會顯示在 GitHub 的 Actions 分頁中 |
| 2 | `on:` | 定義觸發條件 |
| 3 | `push:` | 當有 push 事件時觸發 |
| 4 | `branches: [main]` | 只在 push 到 `main` 分支時觸發 |
| 5 | `jobs:` | 定義所有的 job |
| 6 | `build:` | 定義一個名為 `build` 的 job |
| 7 | `runs-on: ubuntu-latest` | 在最新版 Ubuntu 環境中執行 |
| 8 | `steps:` | 定義這個 job 的所有步驟 |
| 9 | `- uses: actions/checkout@v4` | 使用官方的 checkout action，將 repository 程式碼下載到 Runner |
| 10 | `- name: Say Hello` | 為這個步驟取一個名稱 |
| 11 | `run: echo "Hello, GitHub Actions!"` | 執行一段 shell 指令 |

> 💡 **講師提示：** 可以在 GitHub 上即時建立這個 workflow 檔案並 push，讓學生在 Actions 分頁中看到它被觸發並執行的過程。

---

## 常見 Events（觸發器）

| Event | 觸發時機 | 範例 |
|-------|---------|------|
| `push` | 推送程式碼到指定分支時 | `on: push` / `on: push: branches: [main]` |
| `pull_request` | 建立或更新 PR 時 | `on: pull_request: branches: [main]` |
| `schedule` | 定時排程（cron 語法） | `on: schedule: - cron: '0 0 * * *'` |
| `workflow_dispatch` | 手動觸發（在 Actions 頁面按按鈕） | `on: workflow_dispatch` |
| `release` | 建立 Release 時 | `on: release: types: [published]` |
| `issues` | Issue 被建立或修改時 | `on: issues: types: [opened]` |
| `workflow_run` | 另一個 workflow 完成後觸發 | `on: workflow_run: workflows: [CI]` |

### 事件篩選

你可以進一步篩選觸發條件，只在特定情況下才執行：

```yaml
on:
  push:
    # Only trigger on specific branches
    branches: [main, develop]
    # Only trigger when specific files change
    paths:
      - 'src/**'
      - '*.go'
    # Ignore specific paths
    paths-ignore:
      - 'docs/**'
      - '*.md'

  pull_request:
    # Only trigger on specific event types
    types: [opened, synchronize, reopened]
    branches: [main]
```

> 💡 **講師提示：** `paths` 篩選器非常實用——例如只修改了文件不需要重新跑測試。可以提醒學生善用這個功能來節省 CI 資源。

---

## 常見 Actions

以下是你會經常用到的官方 Actions：

### actions/checkout

**用途**：將 repository 的程式碼下載到 Runner 上。幾乎每個 workflow 的第一步都是這個。

```yaml
- uses: actions/checkout@v4
```

### actions/setup-go

**用途**：安裝指定版本的 Go 語言環境。

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.22'
```

### actions/cache

**用途**：快取依賴套件，加速後續的建置。

```yaml
- uses: actions/cache@v4
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

### actions/upload-artifact

**用途**：將建置產出的檔案上傳為 artifact，供下載或傳給其他 job 使用。

```yaml
- uses: actions/upload-artifact@v4
  with:
    name: my-binary
    path: ./build/app
```

### actions/download-artifact

**用途**：下載先前上傳的 artifact（通常在不同的 job 中使用）。

```yaml
- uses: actions/download-artifact@v4
  with:
    name: my-binary
    path: ./build/
```

### 如何尋找更多 Actions？

前往 [GitHub Marketplace](https://github.com/marketplace?type=actions) 搜尋你需要的功能。在搜尋時，注意以下幾點：

- 優先選擇有 **Verified Creator** 標章的 Action
- 查看 **星星數** 和 **使用人數**
- 確認 Action 的 **最後更新時間**，避免使用已停止維護的 Action
- 閱讀 Action 的 README，了解所有可用參數

---

## GitHub-hosted Runners

GitHub 提供免費的雲端 Runner，不需要自行架設伺服器。

### 可用的作業系統

| Runner 標籤 | 作業系統 | 說明 |
|-------------|---------|------|
| `ubuntu-latest` | Ubuntu 22.04 | 最常用，建議預設選擇 |
| `ubuntu-24.04` | Ubuntu 24.04 | 指定特定 Ubuntu 版本 |
| `windows-latest` | Windows Server 2022 | 需要 Windows 環境時使用 |
| `macos-latest` | macOS 14 (Sonoma) | 需要 macOS 環境時使用（例如 iOS 開發） |

### 預裝的軟體

GitHub-hosted Runner 已預裝大量常用軟體，包括：

- **語言**：Go、Node.js、Python、Java、Ruby、Rust 等
- **套件管理器**：npm、pip、gem 等
- **工具**：Git、Docker、docker compose、curl、wget、jq 等
- **雲端 CLI**：AWS CLI、Azure CLI、Google Cloud SDK 等

> 完整清單請參考 [GitHub 官方文件](https://github.com/actions/runner-images)。

### 資源限制

| 項目 | 限制 |
|------|------|
| vCPU | 2 核（Linux/Windows）/ 3 核（macOS） |
| 記憶體 | 7 GB（Linux/Windows）/ 14 GB（macOS） |
| 磁碟空間 | 14 GB（SSD） |

---

## 免費額度與限制

### 費用方案

| 方案 | 免費額度 |
|------|---------|
| **公開（Public）Repository** | 完全免費，無分鐘數限制 |
| **私人（Private）Repository — Free 方案** | 每月 2,000 分鐘 |
| **私人（Private）Repository — Pro 方案** | 每月 3,000 分鐘 |
| **私人（Private）Repository — Team 方案** | 每月 3,000 分鐘 |

> **注意**：不同作業系統的分鐘數消耗倍率不同。Linux = 1x、Windows = 2x、macOS = 10x。

### 執行時間限制

| 項目 | 限制 |
|------|------|
| 單次 job 最長執行時間 | **6 小時** |
| 單次 workflow 最長執行時間 | **35 天**（通常用於含有 approval 等待的 workflow） |
| API 請求頻率 | 每個 repository 每小時 1,000 次 |
| 同時執行的 job 數量 | Free 方案最多 20 個 |

> 💡 **講師提示：** 對於學生來說，使用公開 repository 基本上不需要擔心額度問題。建議學生在練習時使用公開 repository。

---

## YAML 語法速查

GitHub Actions 的 workflow 使用 YAML 格式。如果你不熟悉 YAML，以下是快速入門：

### 縮排規則

YAML 使用 **空格**（space）進行縮排，**不能用 Tab**。通常使用 **2 個空格** 為一層縮排。

```yaml
# Good — using spaces
parent:
  child:
    grandchild: value

# Bad — using tabs (will cause errors!)
parent:
	child:           # ← This is a tab, YAML will reject this
```

### Key-Value Pairs（鍵值對）

```yaml
name: My Workflow
version: 1.0
enabled: true
count: 42
```

### 列表（List）

```yaml
# Block style
fruits:
  - apple
  - banana
  - cherry

# Inline style
fruits: [apple, banana, cherry]
```

### 巢狀結構

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Build
        run: go build ./...
```

### 多行字串

```yaml
# Literal block (preserves newlines)
description: |
  This is line 1.
  This is line 2.
  This is line 3.

# Folded block (joins lines with spaces)
description: >
  This is a long sentence
  that will be joined
  into a single line.
```

### 變數引用

在 GitHub Actions 中，使用 `${{ }}` 語法引用變數：

```yaml
steps:
  - name: Print info
    run: |
      echo "Repository: ${{ github.repository }}"
      echo "Branch: ${{ github.ref_name }}"
      echo "Actor: ${{ github.actor }}"

  - name: Use secrets
    run: echo "Token is set"
    env:
      MY_TOKEN: ${{ secrets.MY_SECRET_TOKEN }}
```

常用的內建變數：

| 變數 | 說明 |
|------|------|
| `github.repository` | Repository 名稱（如 `owner/repo`） |
| `github.ref` | 完整的 ref（如 `refs/heads/main`） |
| `github.ref_name` | 分支或 tag 名稱（如 `main`） |
| `github.sha` | Commit 的 SHA 雜湊值 |
| `github.actor` | 觸發 workflow 的使用者名稱 |
| `github.event_name` | 觸發事件的名稱（如 `push`） |
| `runner.os` | Runner 的作業系統（如 `Linux`） |
| `secrets.GITHUB_TOKEN` | GitHub 自動提供的認證 token |

> 💡 **講師提示：** 不需要讓學生記住所有變數。重點是讓他們知道 `${{ }}` 語法，以及在需要時去哪裡查。可以引導他們到 [GitHub Actions 官方文件](https://docs.github.com/en/actions/learn-github-actions/contexts) 查閱完整變數列表。

---

## 小結

- **GitHub Actions** 是 GitHub 內建的 CI/CD 平台，用 YAML 檔案定義自動化流程
- 六個核心概念：**Workflow → Event → Job → Step → Action → Runner**
- Workflow 檔案放在 **`.github/workflows/`** 目錄下
- 透過 **Event** 決定何時觸發，透過 **Job** 與 **Step** 定義要做什麼
- 可以利用 **GitHub Marketplace** 找到社群提供的 Actions，像積木一樣組裝你的 pipeline
- 公開 repository 的 GitHub Actions 使用 **完全免費**

> **接下來，我們要動手建立第一個 Workflow！**

---

[← 上一章：CI/CD 概念介紹](01-cicd-intro.md) ｜ [下一章：動手做：第一個 Workflow →](03-first-workflow.md)
