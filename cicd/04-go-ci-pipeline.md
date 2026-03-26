# 04 — Go 專案 CI Pipeline

> **時間：20 分鐘**

---

## 目錄

- [學習目標](#學習目標)
- [範例專案介紹](#範例專案介紹)
- [完整 CI Workflow](#完整-ci-workflow)
- [逐段解說](#逐段解說)
- [Job 依賴關係圖](#job-依賴關係圖)
- [Caching 快取](#caching-快取)
- [Matrix Strategy](#matrix-strategy)
- [Artifacts 深入](#artifacts-深入)
- [實用技巧](#實用技巧)
- [小結與練習題](#小結與練習題)

---

## 學習目標

完成本章節後，你將能夠：

- 為 Go 專案建立一套包含 **lint、test、build** 的完整 CI pipeline
- 理解 **job 之間的依賴關係** 以及平行與序列執行的差異
- 使用 **caching** 加速 CI 執行速度
- 使用 **matrix strategy** 同時在多個環境中測試
- 使用 **artifacts** 在 job 之間傳遞檔案

---

## 範例專案介紹

### 專案結構

本課程使用一個簡單的 Go HTTP API 伺服器作為範例專案，它的結構如下：

```
examples/sample-app/
├── main.go            # HTTP server entry point
├── handler.go         # HTTP handler functions
├── handler_test.go    # Unit tests
├── go.mod             # Go module definition
├── Dockerfile         # Container image definition
└── .golangci.yml      # golangci-lint configuration
```

### 各檔案的功能

#### `main.go` — 主程式進入點

這是伺服器的啟動點，負責設定路由並啟動 HTTP 伺服器：

```go
package main

import (
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", HomeHandler)
    mux.HandleFunc("/health", HealthHandler)
    mux.HandleFunc("/version", VersionHandler)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

#### `handler.go` — Handler 函式

定義了三個 API endpoint 的處理邏輯：

```go
package main

import (
    "encoding/json"
    "net/http"
)

// HomeHandler returns a welcome message
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Hello, World!",
    })
}

// HealthHandler returns the health status
func HealthHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{
        "status": "ok",
    })
}

// VersionHandler returns the application version
func VersionHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{
        "version": "1.0.0",
    })
}
```

#### 三個 Endpoint 的功能

| Endpoint | HTTP Method | 功能 | 回傳範例 |
|----------|-------------|------|---------|
| `/` | GET | 首頁歡迎訊息 | `{"message": "Hello, World!"}` |
| `/health` | GET | 健康檢查 | `{"status": "ok"}` |
| `/version` | GET | 版本資訊 | `{"version": "1.0.0"}` |

#### `handler_test.go` — 單元測試

使用 Go 標準函式庫的 `testing` 和 `net/http/httptest` 套件撰寫測試：

```go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHomeHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    w := httptest.NewRecorder()
    HomeHandler(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", w.Code)
    }
}
```

> 💡 **講師提示：** 如果學生之前上過 Docker 課程並接觸過類似的範例專案，可以快速帶過。重點放在 CI workflow 本身，而非 Go 程式碼的細節。

---

## 完整 CI Workflow

以下是我們要為範例專案建立的完整 CI workflow。請在 `.github/workflows/ci.yml` 中建立這個檔案：

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
```

> 💡 **講師提示：** 建議先展示完整的 workflow 檔案，讓學生有個全貌，接著再逐段解說。可以問學生：「看到這個檔案，你能猜出它做了哪些事情嗎？」

---

## 逐段解說

### 觸發條件

```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```

這個 workflow 在兩種情況下會觸發：

1. **Push 到 `main` 分支**：當程式碼直接 push 到 `main` 時執行
2. **Pull Request 指向 `main` 分支**：當有 PR 要合併進 `main` 時執行

為什麼兩種都要？

- **PR 觸發**：在合併前就檢查程式碼品質，確保不會有問題的程式碼進入 `main`
- **Push 觸發**：作為最後的安全網，確認合併後的程式碼仍然正常

---

### Lint Job — 程式碼品質檢查

```yaml
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
```

#### 什麼是 Linter？

**Linter** 是一種靜態分析工具，它不會執行你的程式碼，而是 **檢查程式碼的風格、潛在錯誤和不良寫法**。就像文章的拼字檢查器，但針對的是程式碼。

#### 什麼是 golangci-lint？

[golangci-lint](https://golangci-lint.run/) 是 Go 生態系中最受歡迎的 linter 工具，它整合了數十個 linter（如 `govet`、`errcheck`、`staticcheck` 等），只需要一個指令就能執行所有檢查。

#### 各步驟說明

| 步驟 | 做什麼 | 為什麼需要 |
|------|--------|-----------|
| `actions/checkout@v4` | 下載 repository 程式碼 | Linter 需要讀取程式碼才能分析 |
| `actions/setup-go@v5` | 安裝 Go 1.22 | golangci-lint 需要 Go 環境 |
| `golangci-lint-action@v6` | 執行 golangci-lint | 檢查程式碼品質 |

#### golangci-lint-action 的使用

[golangci/golangci-lint-action](https://github.com/golangci/golangci-lint-action) 是一個社群維護的 Action，它會自動：

- 下載並安裝 golangci-lint
- 快取 golangci-lint 的執行結果，加速後續 run
- 如果 repository 中有 `.golangci.yml` 設定檔，自動套用其中的設定

---

### Test Job — 自動化測試

```yaml
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
    - name: Upload coverage
      uses: actions/upload-artifact@v4
      with:
        name: coverage-report
        path: coverage.out
```

#### `go test` 的各個 Flag

| Flag | 用途 | 說明 |
|------|------|------|
| `-v` | **Verbose** 詳細輸出 | 印出每個測試函式的名稱與結果，而不只是通過/失敗的摘要 |
| `-race` | **Race Detector** 競爭偵測 | 在測試時偵測 goroutine 之間的 data race，這是 Go 並行程式中常見的 bug |
| `-coverprofile=coverage.out` | **Coverage Profile** 覆蓋率報告 | 將測試覆蓋率資料輸出到 `coverage.out` 檔案 |
| `./...` | **所有套件** | 遞迴執行所有子目錄中的測試 |

#### 測試覆蓋率（Test Coverage）

**測試覆蓋率** 是衡量「你的測試涵蓋了多少程式碼」的指標。

```bash
# Example output of go tool cover -func
github.com/user/app/handler.go:10:  HomeHandler     100.0%
github.com/user/app/handler.go:17:  HealthHandler   100.0%
github.com/user/app/handler.go:24:  VersionHandler  80.0%
total:                              (statements)    93.3%
```

- **100%** 表示該函式的每一行都被測試到了
- **80%** 表示有 20% 的程式碼路徑沒有被測試覆蓋
- 一般來說，核心業務邏輯建議達到 **80% 以上** 的覆蓋率

> 💡 **講師提示：** 提醒學生覆蓋率不是越高越好，重點是測試 **有意義的邏輯**，而非追求數字。100% 覆蓋率不代表沒有 bug。

#### Artifact 上傳

```yaml
- name: Upload coverage
  uses: actions/upload-artifact@v4
  with:
    name: coverage-report
    path: coverage.out
```

**Artifact** 是 workflow 執行過程中產生的檔案，上傳後可以：

- 在 GitHub Actions 頁面 **下載查看**
- 在後續的 job 中 **下載使用**（我們在 Build Job 會看到類似的用法）

---

### Build Job — 建置可執行檔

```yaml
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
```

#### `needs` 關鍵字 — Job 依賴關係

```yaml
needs: [lint, test]
```

`needs` 是這個 workflow 中最重要的概念之一。它定義了 **job 之間的執行順序**：

- **沒有 `needs`** 的 job（如 `lint` 和 `test`）會 **平行執行**，同時開始
- **有 `needs`** 的 job（如 `build`）會 **等待** 指定的 job 全部完成後才開始
- 如果 `needs` 中的任何一個 job **失敗**，這個 job **不會執行**

#### 為什麼 Build 要等 Lint 和 Test 通過？

想想看：如果程式碼有 lint 錯誤或測試失敗，我們還需要花時間去 build 嗎？

答案是 **不需要**。既然程式碼品質有問題，就應該先修好再建置。這樣可以：

1. **節省 CI 資源**：不浪費時間在必然無用的建置上
2. **明確的失敗訊號**：開發者知道問題出在 lint 或 test，而不是 build
3. **邏輯上的正確性**：只有品質通過檢查的程式碼才值得建置

#### Artifact 的用途

Build Job 產出的 binary 透過 `upload-artifact` 上傳後，可以用於：

- **手動下載測試**：在 Actions 頁面下載 binary 來手動驗證
- **後續 job 使用**：例如後面的 deploy job 可以下載這個 binary 來部署
- **Release 發佈**：在 Release workflow 中將 binary 附加到 GitHub Release

---

## Job 依賴關係圖

這三個 job 的執行流程如下：

```
┌──────┐     ┌──────┐
│ Lint │     │ Test │
└──┬───┘     └──┬───┘
   │            │
   │  (平行執行) │
   │            │
   └─────┬──────┘
         │
         │ (等待兩者都完成)
         │
   ┌─────▼─────┐
   │   Build   │
   └───────────┘
```

- **Lint** 和 **Test** 同時開始，平行執行（因為它們之間沒有依賴關係）
- **Build** 會等到 Lint 和 Test **都通過** 後才開始
- 如果 Lint 或 Test 任何一個失敗，Build **不會執行**

> 💡 **講師提示：** 可以在 GitHub Actions 頁面上展示這個依賴關係圖。GitHub 會自動根據 `needs` 設定，在 workflow run 頁面上以圖形化的方式呈現 job 之間的關係。

---

## Caching 快取

### 為什麼需要快取？

每次 CI 執行時，Runner 都是一個 **全新的環境**。這意味著每次都要重新下載所有的依賴套件（如 Go modules）。對於大型專案來說，這可能需要好幾分鐘。

**快取** 可以把這些下載過的依賴儲存起來，下次執行時直接取用，大幅 **縮短 CI 的執行時間**。

### actions/setup-go 內建的快取

好消息是 `actions/setup-go@v5` **已經內建了快取功能**！當它偵測到 `go.sum` 檔案存在時，會自動快取 Go module 的依賴。

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.22'
    # Cache is enabled by default when go.sum exists
```

你不需要額外設定任何東西，`setup-go` 會自動：

1. 計算 `go.sum` 的 hash 值作為 cache key
2. 第一次執行時儲存快取
3. 後續執行時，如果 `go.sum` 沒有變化，直接使用快取

### 手動使用 actions/cache

如果你需要更精細的快取控制（例如快取其他工具或自訂路徑），可以手動使用 `actions/cache`：

```yaml
- name: Cache Go modules
  uses: actions/cache@v4
  with:
    path: |
      ~/go/pkg/mod
      ~/.cache/go-build
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-

- uses: actions/setup-go@v5
  with:
    go-version: '1.22'
    cache: false  # Disable built-in cache since we manage it manually
```

#### Cache Key 的設計

| 部分 | 說明 |
|------|------|
| `${{ runner.os }}` | 作業系統，確保不同 OS 的快取不會混用 |
| `go` | 用來標識這是 Go 的快取 |
| `${{ hashFiles('**/go.sum') }}` | 根據 `go.sum` 的內容產生 hash，依賴有變就更新快取 |

`restore-keys` 是備用 key——如果完全匹配的快取不存在，會嘗試用前綴匹配找到最近的快取。這比從零開始下載好很多。

### 快取效果

| 情境 | 下載 Go modules 時間 | 說明 |
|------|---------------------|------|
| **無快取（第一次）** | ~30–60 秒 | 需要從網路下載所有依賴 |
| **有快取（後續執行）** | ~2–5 秒 | 直接從快取還原 |

> 💡 **講師提示：** 可以實際展示兩次 CI run 的時間差異。第一次 run 後，在不改動 `go.mod` 的情況下再觸發一次，讓學生看到快取帶來的速度提升。

---

## Matrix Strategy

### 什麼是 Matrix？

**Matrix strategy** 讓你可以用一組變數的 **所有組合** 自動產生多個 job。最常見的用途是 **同時測試多個版本或多個作業系統**。

### 範例：多 Go 版本 + 多作業系統測試

```yaml
jobs:
  test:
    name: Test (Go ${{ matrix.go-version }}, ${{ matrix.os }})
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        go-version: ['1.21', '1.22']
        os: [ubuntu-latest, macos-latest]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
      - name: Run tests
        run: go test -v -race ./...
```

### Matrix 產生的 Job 組合

上面的設定會自動產生 **4 個 job**（2 個 Go 版本 x 2 個作業系統）：

| Job | Go Version | OS |
|-----|-----------|-----|
| Test (Go 1.21, ubuntu-latest) | 1.21 | ubuntu-latest |
| Test (Go 1.21, macos-latest) | 1.21 | macos-latest |
| Test (Go 1.22, ubuntu-latest) | 1.22 | ubuntu-latest |
| Test (Go 1.22, macos-latest) | 1.22 | macos-latest |

這 4 個 job 會 **同時平行執行**，大幅縮短測試時間。

### 進階用法：fail-fast

```yaml
strategy:
  fail-fast: false  # Don't cancel other jobs if one fails
  matrix:
    go-version: ['1.21', '1.22']
    os: [ubuntu-latest, macos-latest]
```

- **`fail-fast: true`**（預設）：一旦任何一個 matrix job 失敗，就 **取消其他所有 matrix job**
- **`fail-fast: false`**：即使某個 job 失敗，**其他 job 繼續執行**，讓你看到完整的測試結果

### 什麼時候該用 Matrix？

| 適合用 Matrix 的場景 | 不需要 Matrix 的場景 |
|---------------------|---------------------|
| 開源函式庫，需要支援多個 Go 版本 | 應用程式，只需要一個目標版本 |
| 跨平台工具，需要在多個 OS 上測試 | 只部署到 Linux 伺服器 |
| 需要測試不同資料庫版本的相容性 | 只使用一種資料庫 |

> 💡 **講師提示：** 對於學生的練習專案，通常不需要 matrix。但了解這個功能在開源專案和企業專案中非常重要。可以展示一些知名 Go 開源專案的 CI 設定，讓學生看到 matrix 的實際應用。

---

## Artifacts 深入

### Upload vs Download

在 GitHub Actions 中，不同的 job 在 **不同的 Runner** 上執行，它們的檔案系統 **互不相通**。如果你需要在不同 job 之間傳遞檔案，就需要使用 Artifact。

```
┌──────────────┐                          ┌──────────────┐
│   Test Job   │                          │  Build Job   │
│ (Runner A)   │                          │ (Runner B)   │
│              │                          │              │
│ coverage.out │──── upload-artifact ────▶│              │
│              │     (存到 GitHub)         │  binary      │
└──────────────┘                          └──────────────┘
```

#### upload-artifact — 上傳 Artifact

```yaml
- uses: actions/upload-artifact@v4
  with:
    name: coverage-report      # Artifact name (for identification)
    path: coverage.out         # File or directory to upload
    retention-days: 7          # How long to keep (default: 90 days)
```

#### download-artifact — 下載 Artifact

```yaml
- uses: actions/download-artifact@v4
  with:
    name: coverage-report      # Must match the upload name
    path: ./downloaded/        # Where to save the downloaded files
```

### Artifact 的保存期限

| 方案 | 預設保存期限 | 最長保存期限 |
|------|------------|------------|
| Public Repository | 90 天 | 90 天 |
| Private Repository (Free) | 90 天 | 90 天 |
| Private Repository (Pro/Team) | 90 天 | 400 天 |

你可以透過 `retention-days` 參數自訂保存期限（最短 1 天）。

### 在不同 Job 之間傳遞檔案

以下是一個完整的範例，展示如何在 build job 中產生檔案，並在 deploy job 中使用：

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Build
        run: go build -o bin/app ./...
      - name: Upload binary
        uses: actions/upload-artifact@v4
        with:
          name: app-binary
          path: bin/app

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Download binary
        uses: actions/download-artifact@v4
        with:
          name: app-binary
          path: ./bin/
      - name: Deploy
        run: |
          chmod +x ./bin/app
          echo "Deploying ./bin/app ..."
          # Actual deployment steps would go here
```

---

## 實用技巧

### 1. 使用 `if` 條件控制 Step 執行

你可以用 `if` 來決定某個 step 是否要執行：

```yaml
steps:
  - name: Run only on main branch
    if: github.ref == 'refs/heads/main'
    run: echo "This runs only on main"

  - name: Run only on pull requests
    if: github.event_name == 'pull_request'
    run: echo "This runs only on PRs"

  - name: Run even if previous steps failed
    if: always()
    run: echo "This always runs (cleanup, notifications, etc.)"

  - name: Run only if previous steps succeeded
    if: success()
    run: echo "Everything is fine!"

  - name: Run only if previous steps failed
    if: failure()
    run: echo "Something went wrong! Sending notification..."
```

常用的條件函式：

| 函式 | 說明 |
|------|------|
| `success()` | 前面的步驟都成功時為 true（預設行為） |
| `failure()` | 前面有任何步驟失敗時為 true |
| `always()` | 無論成功或失敗，一律執行 |
| `cancelled()` | workflow 被取消時為 true |

### 2. `continue-on-error` — 容許失敗

有時候你希望某個 step 即使失敗，也不要影響整個 job 的結果：

```yaml
steps:
  - name: Optional lint check
    continue-on-error: true
    run: golangci-lint run ./...

  - name: This step will still run
    run: echo "Previous step might have failed, but we continue"
```

適用場景：

- 非關鍵的 lint 規則檢查
- 實驗性的測試
- 通知服務（即使通知失敗也不影響 CI 結果）

### 3. `timeout-minutes` — 設定超時

防止 job 因為無限迴圈或 hang 住而佔用資源：

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    timeout-minutes: 10  # Kill the job if it runs longer than 10 minutes
    steps:
      - uses: actions/checkout@v4
      - name: Run tests
        run: go test -v ./...
```

建議根據你的專案正常 CI 時間的 **2-3 倍** 來設定超時，既不會太短導致誤殺，也不會讓異常的 job 跑太久。

### 4. `concurrency` — 避免重複執行

當你快速 push 多個 commit 時，可能會觸發多個 workflow run。`concurrency` 可以確保同一時間只有一個 run 在執行：

```yaml
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true
```

| 參數 | 說明 |
|------|------|
| `group` | 同一個 group 的 run 會互斥 |
| `cancel-in-progress: true` | 新的 run 啟動時，**取消** 正在執行的舊 run |

這在 PR workflow 中特別有用——當你在短時間內推了多個 commit 到 PR，只需要跑 **最新那次** 的 CI 就好，前面的可以取消。

> 💡 **講師提示：** `concurrency` 是一個進階功能，如果時間不夠可以先跳過。但建議讓學生知道它的存在，因為在實際專案中很常用。

---

## 小結與練習題

### 本章重點回顧

- 一個完整的 Go CI pipeline 通常包含 **Lint → Test → Build** 三個階段
- 使用 `needs` 定義 job 之間的依賴關係，`lint` 和 `test` 平行執行，`build` 等它們都通過才執行
- **golangci-lint** 是 Go 最受歡迎的 linter 工具，可透過 `golangci-lint-action` 在 CI 中使用
- `go test -v -race -coverprofile=coverage.out ./...` 是標準的 CI 測試指令
- **Caching** 可以大幅加速 CI，`actions/setup-go` 已內建快取功能
- **Matrix strategy** 讓你同時在多個版本/平台上測試
- **Artifacts** 用來在不同 job 之間傳遞檔案
- `if`、`continue-on-error`、`timeout-minutes`、`concurrency` 是常用的進階控制技巧

### 練習題

**練習 1：加入 `go vet` 檢查**

在 CI workflow 的 test job 中，在 `Run tests` 步驟之前，加入一個新的步驟來執行 `go vet`：

```yaml
- name: Run go vet
  run: go vet ./...
```

`go vet` 是 Go 內建的靜態分析工具，可以檢查常見的程式錯誤（如 `fmt.Printf` 的格式化字串與參數不匹配等）。

**練習 2：使用 Matrix 測試多個 Go 版本**

修改 test job，使用 matrix strategy 同時在 Go 1.21 和 Go 1.22 上執行測試。完成後，觀察 GitHub Actions 頁面上會出現幾個 test job。

提示：

```yaml
strategy:
  matrix:
    go-version: ['1.21', '1.22']
```

記得把 `setup-go` 的 `go-version` 改為 `${{ matrix.go-version }}`。

> **接下來，我們將學習如何在 Pull Request 中自動進行程式碼檢查！**

---

[← 上一章：動手做：第一個 Workflow](03-first-workflow.md) ｜ [下一章：PR 檢查自動化 →](05-pr-checks.md)
