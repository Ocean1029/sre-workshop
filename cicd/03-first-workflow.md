# 03 — 動手做：第一個 Workflow

> **時間：15 分鐘**

---

## 目錄

- [學習目標](#學習目標)
- [事前準備](#事前準備)
- [Step 1：建立 Workflow 檔案](#step-1建立-workflow-檔案)
- [Step 2：Push 並觀察結果](#step-2push-並觀察結果)
- [Step 3：手動觸發 Workflow](#step-3手動觸發-workflow)
- [Step 4：製造一個失敗](#step-4製造一個失敗)
- [GitHub Actions 的 Context 與表達式](#github-actions-的-context-與表達式)
- [環境變數](#環境變數)
- [常見問題排解](#常見問題排解)
- [小結與練習題](#小結與練習題)

---

## 學習目標

完成本章節後，你將能夠：

- 從零開始建立一個 GitHub Actions workflow 檔案
- 成功觸發並觀察 workflow 的執行結果
- 使用 `workflow_dispatch` 手動觸發 workflow
- 理解 workflow 成功與失敗的差異，並能排除基本錯誤
- 使用 `${{ }}` 表達式存取 GitHub Actions 的 context 資訊
- 在不同層級設定環境變數

---

## 事前準備

### 確認你的 GitHub 帳號

請確認你已經擁有一個 GitHub 帳號。如果還沒有，請前往 [github.com](https://github.com) 註冊一個。

### 建立一個新的 GitHub Repository

我們需要一個全新的 repository 來練習。請依照以下步驟操作：

1. 登入 GitHub 後，點擊右上角的 **「+」** 按鈕，選擇 **「New repository」**
2. 填寫 repository 資訊：
   - **Repository name**：輸入 `github-actions-lab`（或任何你喜歡的名稱）
   - **Description**（選填）：輸入 `My first GitHub Actions workflow`
   - **Visibility**：選擇 **Public**（公開 repository 的 GitHub Actions 完全免費）
   - **Initialize this repository with**：勾選 **「Add a README file」**
3. 點擊 **「Create repository」** 按鈕

> 💡 **講師提示：** 建議在投影幕上同步操作，讓學生跟著一步一步做。如果有學生還沒有 GitHub 帳號，可以讓他們先觀看示範，課後再自行練習。

### 將 Repository Clone 到本機

打開終端機，執行以下指令：

```bash
# Replace <your-username> with your GitHub username
git clone https://github.com/<your-username>/github-actions-lab.git
cd github-actions-lab
```

---

## Step 1：建立 Workflow 檔案

### 建立目錄結構

GitHub Actions 的 workflow 檔案必須放在 `.github/workflows/` 目錄下。我們先來建立這個目錄：

```bash
mkdir -p .github/workflows
```

> 💡 **講師提示：** 提醒學生 `.github` 是以「點」開頭的隱藏目錄，在檔案管理器中可能看不到。在終端機中可以用 `ls -a` 來查看。

### 建立 Workflow 檔案

在 `.github/workflows/` 目錄下建立一個名為 `hello.yml` 的檔案：

```bash
# You can use any editor you prefer
code .github/workflows/hello.yml
```

將以下內容貼入 `hello.yml`：

```yaml
name: Hello GitHub Actions

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  greeting:
    runs-on: ubuntu-latest
    steps:
      - name: Say Hello
        run: echo "Hello, GitHub Actions! 🚀"

      - name: Show System Info
        run: |
          echo "OS: $(uname -s)"
          echo "Architecture: $(uname -m)"
          echo "Go Version: $(go version)"
          echo "Docker Version: $(docker --version)"
          echo "Current Time: $(date)"

      - name: Show GitHub Context
        run: |
          echo "Repository: ${{ github.repository }}"
          echo "Branch: ${{ github.ref_name }}"
          echo "Commit SHA: ${{ github.sha }}"
          echo "Actor: ${{ github.actor }}"
          echo "Event: ${{ github.event_name }}"
```

### 逐段解說

讓我們來理解這個 workflow 的每個部分：

#### `name: Hello GitHub Actions`

這是 **workflow 的名稱**，會顯示在 GitHub 的 Actions 分頁中。取一個有意義的名稱可以幫助你快速辨識不同的 workflow。

#### `on:` — 觸發條件

```yaml
on:
  push:
    branches: [main]
  workflow_dispatch:
```

這段定義了 **什麼時候觸發** 這個 workflow：

- `push: branches: [main]` — 當有程式碼被 push 到 `main` 分支時自動觸發
- `workflow_dispatch` — 允許在 GitHub UI 上 **手動觸發**（我們在 Step 3 會用到）

#### `jobs:` — 工作定義

```yaml
jobs:
  greeting:
    runs-on: ubuntu-latest
```

- `jobs:` — 定義所有要執行的工作
- `greeting:` — 這個 job 的名稱（可以自訂，盡量用有意義的名稱）
- `runs-on: ubuntu-latest` — 指定在 **最新版 Ubuntu** 的 GitHub-hosted Runner 上執行

#### `steps:` — 步驟

**Step 1：Say Hello**

```yaml
- name: Say Hello
  run: echo "Hello, GitHub Actions! 🚀"
```

最簡單的步驟——用 `run` 執行一行 shell 指令。`name` 是步驟的名稱，會顯示在執行記錄中。

**Step 2：Show System Info**

```yaml
- name: Show System Info
  run: |
    echo "OS: $(uname -s)"
    echo "Architecture: $(uname -m)"
    echo "Go Version: $(go version)"
    echo "Docker Version: $(docker --version)"
    echo "Current Time: $(date)"
```

使用 `|`（literal block）語法可以 **執行多行指令**。這個步驟會印出 Runner 的系統資訊，讓你了解 GitHub-hosted Runner 的環境。你會發現 Go 和 Docker 都已經預裝好了！

**Step 3：Show GitHub Context**

```yaml
- name: Show GitHub Context
  run: |
    echo "Repository: ${{ github.repository }}"
    echo "Branch: ${{ github.ref_name }}"
    echo "Commit SHA: ${{ github.sha }}"
    echo "Actor: ${{ github.actor }}"
    echo "Event: ${{ github.event_name }}"
```

使用 `${{ }}` 表達式存取 **GitHub 提供的上下文資訊**。這些資訊在後面的章節會經常用到，例如自動在 commit message 中標示 SHA、在部署時判斷分支等。

> 💡 **講師提示：** 逐段解說時，可以讓學生先不看教材，嘗試猜測每個部分的意思。這樣可以加深印象。

---

## Step 2：Push 並觀察結果

### 推送程式碼到 GitHub

在終端機中執行以下指令，將 workflow 檔案推送到 GitHub：

```bash
git add .github/workflows/hello.yml
git commit -m "Add first GitHub Actions workflow"
git push origin main
```

### 在 GitHub 上觀察執行結果

1. 打開瀏覽器，進入你的 GitHub repository 頁面
2. 點擊上方的 **「Actions」** 分頁（如果看不到，請確認檔案已成功 push）
3. 你應該會看到一個正在執行（或已完成）的 workflow run，名稱為 **「Hello GitHub Actions」**
4. 點擊進入查看詳細資訊

### 執行結果頁面解說

進入 workflow run 頁面後，你會看到以下內容：

- **Workflow Run 列表**：顯示所有的執行記錄，每一條對應一次觸發
- **Job 列表**：左側或中央會顯示這次 run 包含的所有 job（我們只有一個 `greeting` job）
- **Step Logs**：點擊 job 可以展開每個 step 的執行記錄

點開 **greeting** job，你會看到：

```
✅ Set up job         — GitHub 自動進行的環境設定
✅ Say Hello          — 我們定義的第一個步驟
✅ Show System Info   — 印出系統資訊
✅ Show GitHub Context — 印出 GitHub context
✅ Complete job       — GitHub 自動進行的收尾工作
```

每個步驟旁邊都有一個 **綠色勾勾** ✅，代表該步驟 **執行成功**。

點擊任何一個步驟可以展開查看詳細的輸出內容。試著點開 **Show System Info** 看看 Runner 上裝了什麼版本的 Go 和 Docker。

> 💡 **講師提示：** 讓學生實際在自己的 repository 中操作，並比較彼此的 `github.actor` 和 `github.sha` 是否不同。這可以幫助他們理解 context 的動態特性。

---

## Step 3：手動觸發 Workflow

### 什麼是 `workflow_dispatch`？

還記得我們在 `on:` 中定義了 `workflow_dispatch` 嗎？這個觸發器允許你 **不需要 push 程式碼**，就能直接在 GitHub UI 上手動觸發 workflow。

這在以下場景非常實用：

- 手動重新執行失敗的 pipeline
- 測試 workflow 設定是否正確
- 觸發不需要程式碼變更的維運作業（例如：手動清理暫存檔案）

### 手動觸發步驟

1. 前往你的 repository 的 **Actions** 分頁
2. 在左側邊欄找到 **「Hello GitHub Actions」** workflow
3. 點擊它，你會看到右上方出現一個 **「Run workflow」** 按鈕
4. 點擊 **「Run workflow」**，選擇 branch 為 `main`，然後按下綠色的 **「Run workflow」** 按鈕
5. 等待幾秒鐘，頁面會出現新的 workflow run

觀察這次手動觸發的結果，點開 **Show GitHub Context** 步驟，你會發現 `Event` 欄位顯示的是 `workflow_dispatch`，而不是上次的 `push`。

---

## Step 4：製造一個失敗

### 為什麼要故意製造失敗？

在實際開發中，CI pipeline 失敗是常態。了解失敗時的表現可以幫助你快速排除問題。

### 修改 Workflow

編輯 `.github/workflows/hello.yml`，在 `steps` 最後面加入一個 **一定會失敗** 的步驟：

```yaml
      - name: This Step Will Fail
        run: |
          echo "About to fail..."
          exit 1
```

`exit 1` 代表以 **非零的退出碼** 結束，shell 會將其視為錯誤。

修改完的 `hello.yml` 完整內容如下：

```yaml
name: Hello GitHub Actions

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  greeting:
    runs-on: ubuntu-latest
    steps:
      - name: Say Hello
        run: echo "Hello, GitHub Actions! 🚀"

      - name: Show System Info
        run: |
          echo "OS: $(uname -s)"
          echo "Architecture: $(uname -m)"
          echo "Go Version: $(go version)"
          echo "Docker Version: $(docker --version)"
          echo "Current Time: $(date)"

      - name: Show GitHub Context
        run: |
          echo "Repository: ${{ github.repository }}"
          echo "Branch: ${{ github.ref_name }}"
          echo "Commit SHA: ${{ github.sha }}"
          echo "Actor: ${{ github.actor }}"
          echo "Event: ${{ github.event_name }}"

      - name: This Step Will Fail
        run: |
          echo "About to fail..."
          exit 1
```

### Push 並觀察失敗結果

```bash
git add .github/workflows/hello.yml
git commit -m "Add a failing step to test error handling"
git push origin main
```

前往 Actions 分頁，你會看到：

- Workflow run 旁邊會出現一個 **紅色叉叉** ❌
- 點進去查看，前面的步驟仍然是 ✅，但最後一個步驟會顯示 ❌
- 點開失敗的步驟，可以看到錯誤訊息：`Error: Process completed with exit code 1.`

**重要觀念**：當某個 step 失敗時，**後續的 step 預設不會執行**，整個 job 會被標記為失敗。

### 修復並重新 Push

了解失敗的表現後，讓我們把失敗的步驟移除（或註解掉），恢復正常：

```yaml
      # - name: This Step Will Fail
      #   run: |
      #     echo "About to fail..."
      #     exit 1
```

或者直接刪除那個步驟，然後 push：

```bash
git add .github/workflows/hello.yml
git commit -m "Remove failing step"
git push origin main
```

回到 Actions 分頁，確認新的 workflow run 又回到 **綠色勾勾** ✅。

> 💡 **講師提示：** 這是一個很好的教學點——讓學生親自體驗「紅轉綠」的過程，建立「失敗 → 排查 → 修復 → 成功」的 CI 心態。可以告訴學生：「在真實工作中，看到紅色不要慌張，CI 的目的就是在問題到達正式環境之前就攔截下來。」

---

## GitHub Actions 的 Context 與表達式

### `${{ }}` 表達式語法

在 GitHub Actions 中，`${{ }}` 是一種特殊語法，用來存取 **context 物件** 的資訊。GitHub 在 workflow 執行時會自動提供這些資訊。

```yaml
# Basic usage
run: echo "Hello, ${{ github.actor }}!"

# Use in conditionals
if: ${{ github.ref == 'refs/heads/main' }}

# Use in environment variables
env:
  REPO_NAME: ${{ github.repository }}
```

### 常用的 Context

#### `github` context

提供 workflow run 的相關資訊，是最常用的 context：

| 表達式 | 說明 | 範例值 |
|--------|------|--------|
| `github.repository` | Repository 全名 | `octocat/hello-world` |
| `github.ref` | 完整的 Git ref | `refs/heads/main` |
| `github.ref_name` | 分支或 tag 名稱 | `main` |
| `github.sha` | 觸發 commit 的完整 SHA | `abc123def456...` |
| `github.actor` | 觸發事件的使用者 | `octocat` |
| `github.event_name` | 觸發的事件名稱 | `push`、`pull_request` |
| `github.run_number` | Workflow 的執行序號 | `42` |
| `github.workspace` | Runner 上的工作目錄路徑 | `/home/runner/work/repo/repo` |

#### `env` context

用來存取在 workflow 中定義的環境變數（下一節會詳細介紹）：

```yaml
env:
  MY_VAR: hello

steps:
  - run: echo "${{ env.MY_VAR }}"
```

#### `secrets` context

用來存取在 repository 中設定的 **機密資訊**（如 API Key、Token 等）。Secrets 的內容在 log 中會自動被遮蔽，不會外洩。

```yaml
steps:
  - name: Deploy
    run: ./deploy.sh
    env:
      API_KEY: ${{ secrets.API_KEY }}
```

> 我們會在後面的章節詳細介紹如何設定和使用 secrets。現在只需要知道它的存在即可。

#### `runner` context

提供 Runner 環境的相關資訊：

| 表達式 | 說明 | 範例值 |
|--------|------|--------|
| `runner.os` | 作業系統 | `Linux`、`Windows`、`macOS` |
| `runner.arch` | 架構 | `X64`、`ARM64` |
| `runner.temp` | 暫存目錄路徑 | `/home/runner/work/_temp` |

---

## 環境變數

### 環境變數的三個層級

在 GitHub Actions 中，你可以在三個不同的層級設定環境變數，每個層級的 **作用範圍** 不同：

```yaml
name: Environment Variables Demo

on:
  push:
    branches: [main]

# Workflow-level: available to ALL jobs and ALL steps
env:
  APP_NAME: my-awesome-app
  ENVIRONMENT: production

jobs:
  demo:
    runs-on: ubuntu-latest

    # Job-level: available to ALL steps in THIS job only
    env:
      LOG_LEVEL: debug
      DATABASE_HOST: localhost

    steps:
      - name: Show All Variables
        # Step-level: available to THIS step only
        env:
          STEP_VAR: "I only exist in this step"
        run: |
          echo "App Name: $APP_NAME"
          echo "Environment: $ENVIRONMENT"
          echo "Log Level: $LOG_LEVEL"
          echo "Database Host: $DATABASE_HOST"
          echo "Step Var: $STEP_VAR"

      - name: Step Var Is Gone
        run: |
          echo "App Name: $APP_NAME"
          echo "Log Level: $LOG_LEVEL"
          echo "Step Var: $STEP_VAR"
          # STEP_VAR will be empty here because it was defined in the previous step
```

### 三個層級的作用範圍比較

| 層級 | 設定位置 | 作用範圍 |
|------|---------|---------|
| **Workflow level** | `on:` 同層的 `env:` | 所有 job 的所有 step 都能使用 |
| **Job level** | `jobs.<job_id>:` 底下的 `env:` | 該 job 內的所有 step 能使用 |
| **Step level** | `steps[*]:` 底下的 `env:` | 僅限該 step 能使用 |

**優先順序**：如果不同層級定義了相同名稱的環境變數，**越小範圍的優先**（Step > Job > Workflow）。

### 範例：環境變數的覆蓋

```yaml
env:
  MESSAGE: "I'm from workflow level"

jobs:
  demo:
    runs-on: ubuntu-latest
    env:
      MESSAGE: "I'm from job level"
    steps:
      - name: Check Message
        env:
          MESSAGE: "I'm from step level"
        run: echo "$MESSAGE"
        # Output: I'm from step level
```

> 💡 **講師提示：** 讓學生動手把「環境變數三層級」的範例加到自己的 workflow 中實際跑一次，觀察各層級變數的存取行為。這比單純看教材更容易理解。

---

## 常見問題排解

在第一次撰寫 workflow 時，以下是最常見的幾個問題：

### 1. YAML 縮排錯誤

YAML 對縮排非常敏感，必須使用 **空格（space）**，**不能使用 Tab**。

**常見錯誤訊息**：

```
Invalid workflow file: .github/workflows/hello.yml
  - yaml: line X: mapping values are not allowed in this context
```

**排解方式**：

- 確認你的編輯器設定為使用空格縮排（建議 2 個空格）
- 在 VS Code 中，右下角可以看到目前的縮排設定（`Spaces: 2`），點擊可以切換
- 使用線上 YAML 驗證工具（如 [yamllint.com](https://www.yamllint.com/)）檢查語法

```yaml
# Wrong — mixed indentation
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Test
       run: echo "wrong indent"    # ← only 1 space, should be 2

# Correct
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Test
        run: echo "correct indent"  # ← 2 spaces from the dash
```

### 2. 檔案路徑錯誤

Workflow 檔案 **必須** 放在 `.github/workflows/` 目錄下，路徑大小寫也要正確。

**常見錯誤**：

| 錯誤路徑 | 說明 |
|----------|------|
| `github/workflows/hello.yml` | 少了前面的「.」 |
| `.github/workflow/hello.yml` | `workflows` 少了 `s` |
| `.Github/Workflows/hello.yml` | 大小寫錯誤 |
| `.github/workflows/hello.yaml` | 副檔名 `.yaml` 其實可以，但要確認一致 |

### 3. 權限問題

如果你的 repository 是新建的或剛 fork 的，可能需要確認 Actions 權限：

1. 前往 repository 的 **Settings** → **Actions** → **General**
2. 在 **Actions permissions** 區塊，選擇 **「Allow all actions and reusable workflows」**
3. 在 **Workflow permissions** 區塊，確認權限設定為 **「Read and write permissions」**（後面章節會用到寫入權限）

### 4. 分支名稱不符

如果你的 repository 預設分支是 `master` 而非 `main`，workflow 的 `branches: [main]` 不會被觸發。

**排解方式**：確認你的預設分支名稱，並修改 workflow 中的 `branches` 設定：

```yaml
on:
  push:
    branches: [master]  # ← change to match your default branch
```

> 💡 **講師提示：** YAML 縮排錯誤是學生最常遇到的問題。建議請學生安裝 VS Code 的 YAML 擴充套件（如 Red Hat 的 YAML 套件），它能即時標示語法錯誤。另外也可以推薦安裝 GitHub Actions 擴充套件，它提供 workflow 檔案的語法高亮與自動補全。

---

## 小結與練習題

### 本章重點回顧

- GitHub Actions 的 workflow 檔案放在 **`.github/workflows/`** 目錄下
- 使用 `on:` 定義觸發條件，`jobs:` 定義工作，`steps:` 定義步驟
- `workflow_dispatch` 讓你可以 **手動觸發** workflow
- 當 step 失敗（exit code 非零）時，**後續步驟不會執行**，job 標記為失敗
- `${{ }}` 表達式可以存取各種 context 資訊（`github`、`env`、`secrets`、`runner`）
- 環境變數有三個層級：**Workflow level > Job level > Step level**，越小範圍的優先

### 練習題

**練習 1：自訂你的 Workflow**

修改 `hello.yml`，在 `greeting` job 中加入以下步驟：

1. 印出 `"My name is <your-name>"` 的訊息
2. 印出目前的工作目錄（提示：使用 `pwd` 指令）
3. 列出工作目錄中的所有檔案（提示：使用 `ls -la` 指令）

**練習 2：多 Job Workflow**

建立一個新的 workflow 檔案 `.github/workflows/multi-job.yml`，包含兩個 job：

1. 第一個 job 名為 `info`，印出 GitHub context 的各項資訊
2. 第二個 job 名為 `greet`，印出 `"Workflow completed successfully!"`
3. 讓 `greet` job **等待** `info` job 完成後才執行（提示：使用 `needs` 關鍵字）

> **接下來，我們要為一個真正的 Go 專案建立 CI Pipeline！**

---

[← 上一章：GitHub Actions 基礎](02-github-actions-basics.md) ｜ [下一章：Go 專案 CI Pipeline →](04-go-ci-pipeline.md)
