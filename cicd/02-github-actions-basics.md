# 02 — GitHub Actions 基礎

## Table of Contents

- [事前準備](#事前準備)
- [先跑一個 Workflow](#先跑一個-workflow)
- [GitHub Actions 介面](#github-actions-介面)
- [Workflow 結構元素](#workflow-結構元素)
- [手動觸發與失敗處理](#手動觸發與失敗處理)
- [小結](#小結)

## 事前準備

### 確認你的 GitHub 帳號

請確認你已經擁有一個 GitHub 帳號。如果還沒有，請前往 [github.com](https://github.com) 註冊一個。

### 建立一個新的 GitHub Repository

我們需要一個全新的 repository 來練習。請依照以下步驟操作：

1. 登入 GitHub 後，點擊右上角的 **「+」** 按鈕，選擇 **「New repository」**
2. 填寫 repository 資訊：
   - **Repository name**：輸入 `github-actions-lab`
   - **Description**（選填）：輸入 `My first GitHub Actions workflow`
   - **Visibility**：選擇 **Public**
3. 點擊 **「Create repository」** 按鈕

### 將 Repository Clone 到本機

打開終端機，執行以下指令：

```bash
# Replace <your-username> with your GitHub username
git clone https://github.com/<your-username>/github-actions-lab.git
cd github-actions-lab
```

## 先跑一個 Workflow

先不用管每一行在幹嘛，照著貼就好，我們等一下會回來解釋。

### 建立目錄結構

GitHub Actions 的 workflow 檔案必須放在 `.github/workflows/` 目錄下。我們先來建立這個目錄：

```bash
mkdir -p .github/workflows
```

### 建立 Workflow 檔案

在 `.github/workflows/` 目錄下建立一個名為 `hello.yml` 的檔案（用你熟悉的編輯器打開即可），將以下內容貼入並存檔：

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
        run: echo "Hello, GitHub Actions!"

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

> 完整檔案也可以在 `cicd/examples/.github/workflows/hello.yml` 找到，直接複製過去即可。（本章的 demo workflow 不綁定後面的 sample-app，所以放在 `cicd/examples/` 根目錄。）

### Push 到 GitHub

```bash
git add .github/workflows/hello.yml
git commit -m "Add first GitHub Actions workflow"
git push origin main
```

## GitHub Actions 介面

打開瀏覽器，進入你的 GitHub repository 頁面，點擊上方的 **「Actions」** 分頁（如果看不到，請確認檔案已成功 push）。你應該會看到一個正在執行（或已完成）的 workflow run，名稱就是你剛剛的 commit message：

![GitHub Actions 介面](assets/github-actions-tab.png)

點進去之後，你會看到這次 workflow run 的詳細資訊：

![Workflow Run 頁面](assets/github-actions-workflow-run.png)

- **Summary**：這次 run 的摘要，包括觸發方式、commit SHA、狀態、執行時間
- **Job 列表**：左側會顯示這次 run 包含的所有 job（我們只有一個 `greeting` job）

點開 **greeting** job，你會看到這個 job 裡每個 step 的執行紀錄：

![Job Steps 列表](assets/github-actions-job-steps.png)

每個步驟旁邊都有打勾符號代表執行成功。點擊任何一個步驟可以展開查看詳細的輸出內容：

![Step Logs 展開畫面](assets/github-actions-step-logs.png)

試著看一下 **Show System Info**，你會發現 Runner 上已經預裝了 Go 和 Docker，並可以從中得知它們的版本。

## Workflow 結構元素

OK，你已經成功跑完第一個 workflow 了。接下來我們回頭看 `hello.yml`，認識它的結構元素。

### 1. Workflow

**Workflow** 是一個定義自動化流程的 **YAML 檔案**。剛剛那個 `hello.yml` 就是一個 workflow。

- 檔案位置：必須放在 repository 的 **`.github/workflows/`** 目錄下
- 副檔名：`.yml` 或 `.yaml`
- 一個 repository 可以有 **多個** workflow

```text
my-repo/
├── .github/
│   └── workflows/
│       ├── ci.yml          # workflow 1
│       ├── pr-check.yml    # workflow 2
│       └── release.yml     # workflow 3
├── main.go
└── go.mod
```

### 2. Event

**Event** 是觸發 workflow 執行的條件。在剛剛的 `hello.yml` 裡，我們定義了兩個 event：

```yaml
on:
  push:
    branches: [main]
  workflow_dispatch:
```

- `push: branches: [main]` — 當有程式碼被 push 到 `main` 分支時自動觸發
- `workflow_dispatch` — 允許在 GitHub UI 上手動觸發

常見的 event 還有這些：

| Event | 觸發時機 |
|-------|---------|
| `push` | 推送程式碼到指定分支時 |
| `pull_request` | 建立或更新 PR 時 |
| `schedule` | 定時排程（cron 語法） |
| `workflow_dispatch` | 手動觸發 |
| `release` | 建立 Release 時 |
| `workflow_run` | 另一個 workflow 完成後觸發 |

每個 event 還可以加上篩選條件（例如只在特定 branch、特定檔案被改動時才觸發），後面章節用到時再詳細介紹。

### 3. Job

**Job** 是 workflow 中的**執行單元**。在剛剛的 `hello.yml` 裡，我們只有一個 job 叫做 `greeting`：

```yaml
jobs:
  greeting:
    runs-on: ubuntu-latest
```

一個 workflow 可以有多個 job，預設會平行執行（下一章會看到平行 lint / test / build 的例子）。

### 4. Step

**Step** 就是 job 裡面的執行步驟。在剛剛的 `hello.yml` 裡，`greeting` job 有三個 step：`Say Hello`、`Show System Info`、`Show GitHub Context`，它們會依序執行（同一個 job 內的 step 不能平行運作），並共用相同的檔案系統。

每個 step 有兩種寫法：`run` 直接執行 shell 指令，`uses` 則是引用別人寫好的 **Action**。很多共通步驟（例如把程式碼 checkout 下來），社群都有現成的 Action 可以直接用，不用自己從頭寫：

```yaml
# run — execute a shell command yourself
- run: echo "Hello"

# uses — use a pre-built Action from the community
- uses: actions/checkout@v4
```

`uses` 的格式是 `{owner}/{repo}@{version}`，例如 `actions/checkout@v4` 就是 GitHub 官方提供的 checkout Action 第 4 版。建議固定版本號（用 `@v4` 而非 `@main`），避免未預期的變更。常用的 Actions 可以在 [GitHub Marketplace](https://github.com/marketplace?type=actions) 找到。

### 5. 表達式與 Context

仔細看 `hello.yml` 的 `Show GitHub Context` 步驟，你會發現有一些 `${{ ... }}` 的寫法，例如 `${{ github.repository }}`、`${{ github.actor }}`。這個 `${{ }}` 是 GitHub Actions 的**表達式語法**，用來存取 GitHub 在 workflow 執行時自動提供的資訊（稱為 **context**）。最簡單的用法就是把 context 的值嵌入 shell 指令裡：

```yaml
- run: echo "Hello, ${{ github.actor }}!"
```

最常用的是 `github` context，包含這次 workflow run 的相關資訊：

| 表達式 | 說明 | 範例值 |
|--------|------|--------|
| `github.repository` | Repository 全名 | `octocat/hello-world` |
| `github.ref_name` | 分支或 tag 名稱 | `main` |
| `github.sha` | 觸發 commit 的 SHA（40 字元） | `abc123def4567890abcdef1234567890abcdef12` |
| `github.actor` | 觸發事件的使用者 | `octocat` |
| `github.event_name` | 觸發的事件名稱 | `push`、`pull_request` |

Context 就像是 GitHub Actions 提供的內建變數集合，除了 `github` 之外還有 `secrets` 用來取機密資訊（例如 API key、token）。這些後面章節用到時會詳細介紹，例如 04 章部署到 Fly.io 就會用到 `secrets.FLY_API_TOKEN`。

### 6. Runner

**Runner** 是實際執行 job 的**機器環境**。在剛剛的 `hello.yml` 裡，`runs-on: ubuntu-latest` 代表我們使用 GitHub 提供的最新版 Ubuntu 虛擬機。常用的標籤有：

| Runner 標籤 | 作業系統 | 備註 |
|-------------|---------|------|
| `ubuntu-latest` | Ubuntu（最新 LTS） | 最常用，建議預設選擇 |
| `windows-latest` | Windows Server | 需要 Windows 環境時使用 |
| `macos-latest` | macOS | 需要 macOS 環境時使用（例如 iOS 開發） |

還記得 `Show System Info` 步驟印出來的 Go 和 Docker 版本嗎？那是因為 GitHub-hosted Runner 已經預裝了大量常用工具（Go、Node.js、Python、Git、Docker、各種雲端 CLI 等等），所以你不用自己 `apt install` 就能直接用。

GitHub-hosted Runner 是 GitHub 提供的雲端虛擬機，用完就釋放，不需要自己管理伺服器。如果你需要特殊硬體、軟體或內網環境，也可以用 **Self-hosted Runner**，自己架機器、自己維護，但完全可控。本課程都使用 GitHub-hosted Runner。

### 概念關係圖

這五個元素（Workflow / Event / Job / Step / Runner）構成了 workflow 的結構骨架，而 `${{ }}` 則是讓你可以在裡面讀 GitHub 自動提供的資訊。用剛剛的 `hello.yml` 來看這些元素的階層關係：

```
Workflow (.github/workflows/hello.yml)
├── Event
│   ├── on: push (branches: [main])
│   └── on: workflow_dispatch
├── Job: greeting
│   ├── Runner: ubuntu-latest
│   ├── Step 1 — run: echo "Hello, GitHub Actions!"
│   ├── Step 2 — run: echo "OS: ..." (Show System Info)
│   └── Step 3 — run: echo "Repository: ..." (Show GitHub Context)
```

## 手動觸發與失敗處理

### 手動觸發 Workflow

我們剛剛在 event 中定義了 `on: workflow_dispatch`，這個觸發器允許我們在不 push 任何程式碼的情況下，直接到 GitHub UI 上手動觸發 workflow。它常用於測試 workflow 設定是否正確，或是觸發不需要程式碼變更的維運作業。

1. 前往你的 repository 的 **Actions** 分頁
2. 在左側邊欄找到 **「Hello GitHub Actions」** workflow
3. 點擊它，你會看到右上方出現一個 **「Run workflow」** 按鈕
4. 點擊 **「Run workflow」**，選擇分支為 `main`，然後按下綠色的 **「Run workflow」** 按鈕
5. 等待幾秒鐘，頁面會出現新的 workflow run

![手動觸發 workflow](assets/github-actions-run-workflow.png)

觀察這次手動觸發的結果，點開 **Show GitHub Context** 步驟，你會發現 `Event` 欄位顯示的是 `workflow_dispatch`，而不是上次的 `push`。

### Pipeline 失敗

在實際開發中，CI pipeline 失敗是常態，學會看失敗訊息、知道從哪裡 debug 是 SRE 的基本功。我們接下來要拿一個會失敗的 pipeline 來跑跑看。

這個範例會用 `curl` 向 GitHub API 抓 repository 的基本資訊，再用 `jq`（一個命令列 JSON 解析工具）從 JSON 結果裡挑欄位出來印。我們會故意把其中一個欄位名稱寫錯，讓你看到 pipeline 失敗時長什麼樣子、以及要怎麼從 log 找出問題。請在 `.github/workflows/repo-info.yml` 中建立這個檔案：

```yaml
name: Repo Info

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  fetch-info:
    runs-on: ubuntu-latest
    steps:
      - name: Fetch repository info
        run: curl -s "https://api.github.com/repos/${{ github.repository }}" > repo.json

      - name: Show repository name
        run: jq -r .full_name repo.json

      - name: Show star count
        run: |
          STARS=$(jq -e .stars repo.json)
          echo "This repo has $STARS stars"

      - name: Show description
        run: jq -r .description repo.json
```

> 完整檔案也可以在 `cicd/examples/.github/workflows/repo-info.yml` 找到，直接複製過去即可。

push 上去之後到 Actions 分頁看，這次的 workflow run 旁邊會變成紅色叉叉：

![失敗的 workflow run 列表](assets/github-actions-failed-runs.png)

點進去看 summary，會顯示這個 run 的狀態是 **Failure**，下面的 `fetch-info` job 也標示失敗：

![失敗 run 的 summary 頁面](assets/github-actions-failed-summary.png)

點開 `fetch-info` job，可以看到每個 step 的執行狀況。前兩個 step（`Fetch repository info` 和 `Show repository name`）都正常執行，但 `Show star count` 失敗了，更值得注意的是失敗那一步後面的 `Show description` 完全沒被執行：

![Job 內 step 列表，Show description 被跳過](assets/github-actions-failed-steps.png)

這是 GitHub Actions 的預設行為，**只要任何一個 step 失敗（exit code 不是 0），後面的 step 就不會跑**，整個 job 也會被標記為失敗。

點開失敗的 step 看詳細 log，可以看到實際的錯誤訊息：

![失敗 step 的詳細錯誤訊息](assets/github-actions-failed-step-detail.png)

錯誤訊息告訴我們 `jq` 在 JSON 裡找不到 `stars` 這個欄位，那正確的欄位叫什麼呢？最快的確認方式就是直接去看 GitHub API 實際回的 JSON：在 `Fetch repository info` 那個 step 的 log 裡找到展開後的 `curl` URL（類似 `https://api.github.com/repos/你的帳號/你的-repo`），複製到瀏覽器打開，就能看到這個 repo 的完整 JSON。翻一下就會發現，star 數的欄位其實叫 `stargazers_count`，不是 `stars`：

![GitHub API 回傳的 JSON，highlight stargazers_count 欄位](assets/github-api-json-stargazers.png)

把 `.stars` 改成 `.stargazers_count`，commit、push，回到 Actions 分頁等幾秒，workflow 就會回到綠色打勾的狀態，`Show star count` 這次也能正確印出 star 數：

![修正後的 workflow run，所有 step 都打勾](assets/github-actions-fixed-run.png)

這整個「看失敗訊息 → 回頭對照程式碼與原始資料 → 修正 → 重跑」的循環，就是你之後每次遇到 CI 失敗會重複做的事。

## 小結

到這裡你已經完成了第一個 GitHub Actions workflow 的旅程：用 YAML 寫一個 `hello.yml` 放到 `.github/workflows/` 底下、push 上去看它跑、認識 workflow 的五個結構元素（**Workflow → Event → Job → Step → Runner**）和兩個常用的功能（`uses` 引用 Action、`${{ }}` 讀 context），到手動觸發 workflow、刻意製造失敗來看 debug 流程。下一章我們會用這些觀念為一個真正的 Go 專案建立完整的 CI pipeline。

### 練習題

完成以下練習來鞏固本章所學：

[練習一：GitHub Actions 基礎練習](exercises/01-basics.md)

[← 上一章：CI/CD 概念介紹](01-cicd-intro.md) ｜ [下一章：Go 專案 CI Pipeline →](03-go-ci-pipeline.md)
