# 練習一：GitHub Actions 基礎練習

> **難度：** 入門 | **對應章節：** 02-GitHub Actions 基礎、03-動手做：第一個 Workflow

---

## 目錄

- [練習 1-1：自訂 Hello World](#練習-1-1自訂-hello-world)
- [練習 1-2：多個 Job](#練習-1-2多個-job)
- [練習 1-3：環境變數挑戰](#練習-1-3環境變數挑戰)
- [練習 1-4：Scheduled Workflow](#練習-1-4scheduled-workflow)
- [延伸思考](#延伸思考)

---

## 練習 1-1：自訂 Hello World

### 目標

修改 `.github/workflows/hello.yml`，加入你自己的自訂步驟，熟悉 GitHub Actions step 的基本寫法。

### 要求

1. 加入一個步驟，印出你自己的名字和學號
2. 加入一個步驟，使用 `actions/checkout@v4` 把 repo 的程式碼 checkout 下來，然後用 `ls -la` 列出工作目錄的檔案
3. 加入一個步驟，顯示目前的 Git 資訊（分支名稱、最近一次 commit 的訊息和作者）

### 提示

- 要能看到 repository 裡的檔案，**必須先使用 `actions/checkout@v4`** 來 checkout 程式碼。如果沒有 checkout，Runner 上的工作目錄是空的。
- 多行指令可以使用 `|`（literal block）語法。
- Git 資訊可以用 `git log -1` 來查看最近一次 commit。

### 預期結果

- 在 GitHub Actions 的 log 中可以看到你的名字和學號
- 可以看到 repo 中所有檔案的列表
- 可以看到 Git 分支名稱和 commit 資訊

<details>
<summary>點擊查看答案</summary>

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

      - name: Show My Info
        run: |
          echo "Name: Wang Xiao-Ming"
          echo "Student ID: B10901001"

      - name: Checkout Repository
        uses: actions/checkout@v4

      - name: List Files in Working Directory
        run: ls -la

      - name: Show Git Info
        run: |
          echo "Branch: $(git branch --show-current)"
          echo "Latest Commit:"
          git log -1 --pretty=format:"  SHA: %H%n  Author: %an <%ae>%n  Date: %ad%n  Message: %s"
```

**重點說明：**

- `actions/checkout@v4` 必須放在 `ls -la` 和 `git log` **之前**，否則工作目錄是空的，也沒有 `.git` 資料夾。
- `git log -1` 的 `--pretty=format:` 可以自訂輸出格式，`%H` 是完整 SHA、`%an` 是作者名稱、`%s` 是 commit 訊息。

</details>

---

## 練習 1-2：多個 Job

### 目標

建立一個有多個 Job 的 workflow，練習使用 `needs` 關鍵字來建立 Job 之間的依賴關係。

### 要求

1. 建立一個新的 workflow 檔案 `.github/workflows/multi-job.yml`
2. **Job 1（`preparation`）：** 印出 `"Step 1: Preparation"` 並顯示系統資訊（OS、架構、時間）
3. **Job 2（`processing`）：** 依賴 Job 1，印出 `"Step 2: Processing"` 並顯示 GitHub context 資訊
4. **Job 3（`cleanup`）：** 依賴 Job 2，印出 `"Step 3: Cleanup"` 並印出 `"All steps completed successfully!"`

### 提示

- 使用 `needs` 關鍵字來指定 job 的依賴關係。
- `needs: preparation` 表示這個 job 會等 `preparation` 完成後才開始。
- 沒有 `needs` 的 job 會平行執行，有 `needs` 的 job 會序列執行。

### 預期結果

在 GitHub Actions 頁面上，你應該會看到三個 job 像瀑布一樣 **依序執行**：

```
preparation → processing → cleanup
```

<details>
<summary>點擊查看答案</summary>

```yaml
name: Multi Job Workflow

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  preparation:
    name: Step 1 - Preparation
    runs-on: ubuntu-latest
    steps:
      - name: Preparation
        run: |
          echo "Step 1: Preparation"
          echo "========================"
          echo "OS: $(uname -s)"
          echo "Architecture: $(uname -m)"
          echo "Current Time: $(date)"
          echo "Runner Name: ${{ runner.name }}"

  processing:
    name: Step 2 - Processing
    needs: preparation
    runs-on: ubuntu-latest
    steps:
      - name: Processing
        run: |
          echo "Step 2: Processing"
          echo "========================"
          echo "Repository: ${{ github.repository }}"
          echo "Branch: ${{ github.ref_name }}"
          echo "Commit SHA: ${{ github.sha }}"
          echo "Triggered by: ${{ github.actor }}"

  cleanup:
    name: Step 3 - Cleanup
    needs: processing
    runs-on: ubuntu-latest
    steps:
      - name: Cleanup
        run: |
          echo "Step 3: Cleanup"
          echo "========================"
          echo "All steps completed successfully!"
          echo "Workflow run number: ${{ github.run_number }}"
```

**重點說明：**

- `processing` 設定 `needs: preparation`，所以會等 `preparation` 完成後才開始。
- `cleanup` 設定 `needs: processing`，所以會等 `processing` 完成後才開始。
- 如果前一個 job 失敗，後面依賴它的 job **不會執行**。
- 在 GitHub Actions 頁面上，你可以看到 job 之間的依賴關係圖。

</details>

---

## 練習 1-3：環境變數挑戰

### 目標

練習在 GitHub Actions 中使用 **Workflow level**、**Job level** 和 **Step level** 三個不同層級的環境變數。

### 要求

1. 設定 **Workflow level** 的環境變數：`APP_NAME`（值設為你的應用程式名稱，例如 `my-cool-app`）
2. 設定 **Job level** 的環境變數：`ENVIRONMENT`（值設為 `development`）
3. 設定 **Step level** 的環境變數：`MESSAGE`（值設為 `Hello from step level!`）
4. 在最後一個 step 中，**同時印出所有三個環境變數**
5. 再加一個 step，驗證 Step level 的 `MESSAGE` 變數在這個新的 step 中 **已經消失**

### 提示

- Workflow level 的 `env` 放在 `on:` 的同一層。
- Job level 的 `env` 放在 `jobs.<job_id>:` 底下。
- Step level 的 `env` 放在 `steps[*]:` 底下。
- 環境變數的優先順序：**Step > Job > Workflow**（越小範圍越優先）。

### 預期結果

- 在 "Show All Variables" step 中，三個變數都有值。
- 在 "Verify Step Scope" step 中，`MESSAGE` 變數是空的（因為它只在上一個 step 中定義）。

<details>
<summary>點擊查看答案</summary>

```yaml
name: Environment Variables Challenge

on:
  push:
    branches: [main]
  workflow_dispatch:

# Workflow-level environment variables
env:
  APP_NAME: my-cool-app

jobs:
  env-demo:
    runs-on: ubuntu-latest

    # Job-level environment variables
    env:
      ENVIRONMENT: development

    steps:
      - name: Show All Variables
        # Step-level environment variables
        env:
          MESSAGE: "Hello from step level!"
        run: |
          echo "=== All Environment Variables ==="
          echo "APP_NAME (workflow level): $APP_NAME"
          echo "ENVIRONMENT (job level): $ENVIRONMENT"
          echo "MESSAGE (step level): $MESSAGE"

      - name: Verify Step Scope
        run: |
          echo "=== Verify Variable Scopes ==="
          echo "APP_NAME (workflow level): $APP_NAME"
          echo "ENVIRONMENT (job level): $ENVIRONMENT"
          echo "MESSAGE (step level): '$MESSAGE'"
          echo ""
          if [ -z "$MESSAGE" ]; then
            echo "MESSAGE is empty — step-level variables are NOT accessible in other steps."
          else
            echo "MESSAGE is still set — this should not happen!"
          fi
```

**重點說明：**

- `APP_NAME` 是 workflow level，所以 **所有 job 的所有 step** 都能存取。
- `ENVIRONMENT` 是 job level，所以 **同一個 job 內的所有 step** 都能存取。
- `MESSAGE` 是 step level，所以 **只有定義它的那個 step** 能存取。到了下一個 step，`MESSAGE` 就不見了。
- 這個範例可以幫助你理解三個層級的作用範圍差異。

</details>

---

## 練習 1-4：Scheduled Workflow

### 目標

建立一個 **定時執行** 的 workflow，練習使用 `schedule` 觸發器和 cron 語法。

### 要求

1. 建立一個新的 workflow 檔案 `.github/workflows/scheduled.yml`
2. 使用 cron 語法設定為 **每天台灣時間早上 9 點**（UTC+8）執行
3. 在 workflow 中印出當前的 UTC 時間和台灣時間
4. 同時支援 **手動觸發**（`workflow_dispatch`）

### Cron 語法說明

GitHub Actions 使用標準的 cron 語法，所有時間都是 **UTC 時區**：

```
┌───────────── minute (0–59)
│ ┌───────────── hour (0–23)
│ │ ┌───────────── day of the month (1–31)
│ │ │ ┌───────────── month (1–12)
│ │ │ │ ┌───────────── day of the week (0–6, Sunday=0)
│ │ │ │ │
│ │ │ │ │
* * * * *
```

常用範例：

| Cron 表達式 | 說明 |
|------------|------|
| `0 * * * *` | 每小時的第 0 分鐘 |
| `0 9 * * *` | 每天 UTC 09:00 |
| `0 1 * * *` | 每天 UTC 01:00（台灣時間 09:00） |
| `30 8 * * 1-5` | 週一至週五 UTC 08:30 |
| `0 0 * * 0` | 每週日 UTC 00:00 |

### 提示

- 台灣時間（UTC+8）早上 9 點 = UTC 時間凌晨 1 點。
- 所以 cron 表達式應該是 `0 1 * * *`。
- 可以使用 `TZ` 環境變數搭配 `date` 指令來顯示特定時區的時間。
- **注意**：GitHub Actions 的 scheduled workflow 可能會有幾分鐘的延遲，不會精確地在指定時間執行。

### 預期結果

- Workflow 可以透過手動觸發來測試（不用等到隔天早上 9 點）。
- 在 log 中可以看到 UTC 時間和台灣時間。

<details>
<summary>點擊查看答案</summary>

```yaml
name: Scheduled Workflow

on:
  schedule:
    # Run at 01:00 UTC every day (= 09:00 Taiwan Time, UTC+8)
    - cron: '0 1 * * *'
  workflow_dispatch:

jobs:
  daily-check:
    runs-on: ubuntu-latest
    steps:
      - name: Show Current Time
        run: |
          echo "=== Current Time ==="
          echo "UTC Time: $(date -u '+%Y-%m-%d %H:%M:%S %Z')"
          echo "Taiwan Time (UTC+8): $(TZ='Asia/Taipei' date '+%Y-%m-%d %H:%M:%S %Z')"

      - name: Show Trigger Info
        run: |
          echo "=== Trigger Info ==="
          echo "Event: ${{ github.event_name }}"
          echo "Triggered by: ${{ github.actor }}"
          if [ "${{ github.event_name }}" = "schedule" ]; then
            echo "This was a scheduled run."
          else
            echo "This was a manual run."
          fi

      - name: Daily Report
        run: |
          echo "=== Daily Report ==="
          echo "Repository: ${{ github.repository }}"
          echo "Default Branch: ${{ github.event.repository.default_branch }}"
          echo "Workflow Run Number: ${{ github.run_number }}"
          echo "Daily check completed!"
```

**重點說明：**

- `cron: '0 1 * * *'` 代表 UTC 01:00，也就是台灣時間 09:00。
- `TZ='Asia/Taipei' date` 可以在 Ubuntu Runner 上顯示台灣時區的時間。
- 使用 `${{ github.event_name }}` 可以判斷 workflow 是被 schedule 觸發還是手動觸發。
- **建議先用 `workflow_dispatch` 手動觸發來測試**，確認 workflow 沒有語法錯誤，不用等到隔天凌晨。
- GitHub 的 scheduled workflow 只會在 **預設分支**（通常是 `main`）上執行。

</details>

---

## 延伸思考

以下問題沒有標準答案，供進階學生思考和討論：

1. **Job 平行 vs 序列**：在練習 1-2 中，如果 `processing` 和 `cleanup` 都只依賴 `preparation`（而不是 `cleanup` 依賴 `processing`），執行流程會有什麼改變？畫出 job 之間的依賴關係圖。

2. **環境變數安全性**：如果你需要在 workflow 中使用一個 API Key，應該用環境變數還是 GitHub Secrets？為什麼？兩者的差異是什麼？

3. **Cron 的限制**：GitHub Actions 的 scheduled workflow 有什麼限制？（提示：查一下 GitHub 官方文件關於 schedule 觸發器的說明，了解最短間隔和延遲問題。）

4. **Workflow 除錯**：如果你的 workflow 一直失敗，你會用什麼方法來除錯？列出至少三種方法。

---

[回到教材目錄 →](../README.md)
