# 練習二：CI Pipeline 實戰練習

> **難度：** 中級 | **對應章節：** 03-Go 專案 CI Pipeline

---

## 目錄

- [練習 2-1：擴充 CI Pipeline](#練習-2-1擴充-ci-pipeline)
- [練習 2-2：PR 檢查練習](#練習-2-2pr-檢查練習)
- [延伸思考](#延伸思考)

---

## 練習 2-1：擴充 CI Pipeline

### 目標

在課程教材中的 `ci.yml` 基礎上，加入更多程式碼品質檢查，建立更完整的 CI pipeline。

### 要求

在 `ci.yml` 的 `test` job 中，在 `Run tests` 步驟 **之前** 加入以下三個檢查：

1. **`go vet ./...`** — Go 內建的靜態分析工具，可以檢查常見的程式碼錯誤（例如 `fmt.Printf` 的格式化字串與參數不匹配）
2. **`go mod verify`** — 驗證 `go.sum` 中記錄的 module hash 是否正確，確保 dependencies 沒有被篡改
3. **程式碼格式檢查** — 使用 `gofmt -l .` 列出所有未格式化的檔案。如果有任何檔案未格式化，workflow 就應該 fail

### 提示

- `gofmt -l .` 會列出所有格式不正確的檔案。如果輸出為空，代表所有檔案都已正確格式化。
- 可以用 shell 的條件判斷來檢查 `gofmt` 的輸出是否為空。
- 這些檢查步驟都放在 `test` job 中，並且在 `Run tests` 步驟之前。

### 預期結果

- 如果有 `go vet` 發現的問題，CI 會在 "Run go vet" 步驟失敗。
- 如果 `go.sum` 不一致，CI 會在 "Verify dependencies" 步驟失敗。
- 如果有未格式化的檔案，CI 會在 "Check code formatting" 步驟失敗，並列出哪些檔案需要格式化。

<details>
<summary>點擊查看答案</summary>

只展示 `test` job 的改動，`lint` 和 `build` job 維持 03 章教材原樣即可：

```yaml
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      # === New checks start here ===
      - name: Run go vet
        run: go vet ./...

      - name: Verify dependencies
        run: go mod verify

      - name: Check code formatting
        run: |
          UNFORMATTED=$(gofmt -l .)
          if [ -n "$UNFORMATTED" ]; then
            echo "The following files are not properly formatted:"
            echo "$UNFORMATTED"
            echo ""
            echo "Please run 'gofmt -w .' to fix formatting."
            exit 1
          fi
          echo "All files are properly formatted."
      # === New checks end here ===

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

**重點說明：**

- `go vet` 是 Go 內建的靜態分析工具，與 `golangci-lint` 不同，它不需要額外安裝。
- `go mod verify` 會比對本地 module cache 與 `go.sum` 的 hash 值，確保沒有人篡改依賴套件。
- `gofmt -l .` 只列出檔案，不修改。如果輸出不為空，代表有檔案需要格式化。使用 `exit 1` 讓 step 失敗。
- 這三個檢查放在 `Run tests` 前面，如果基本品質不過關，就不用浪費時間跑測試。

</details>

---

## 練習 2-2：PR 檢查練習

### 目標

實際走一遍完整的 **PR 檢查流程**：建立分支、修改程式碼、撰寫測試、開 PR、觀察 CI 結果。

### 要求

1. 從 `main` 建立一個新的分支 `feature/add-info-endpoint`
2. 在範例專案中新增一個 `/info` endpoint，回傳應用程式的資訊
3. 為新的 endpoint 撰寫對應的測試
4. 開一個 Pull Request 到 `main`，並觀察 CI 檢查的結果

> 這個練習 **不提供完整答案**，而是提供 step by step 的引導。目的是讓你走一遍真實的開發流程。

### Step by Step 引導

#### Step 1：建立新分支

```bash
git checkout main
git pull origin main
git checkout -b feature/add-info-endpoint
```

#### Step 2：新增 `/info` endpoint

在 `handler.go` 中新增一個 handler 函式（維持與既有 handler 相同的命名風格與 method 路由風格）：

```go
// handleInfo returns application metadata.
func handleInfo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "app":     "sample-app",
        "version": version,
    })
}
```

然後在 `main.go` 中註冊這個 handler：

```go
mux.HandleFunc("GET /info", handleInfo)
```

#### Step 3：撰寫測試

在 `handler_test.go` 中新增測試：

```go
func TestHandleInfo(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/info", nil)
    w := httptest.NewRecorder()
    handleInfo(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", w.Code)
    }

    var result map[string]string
    if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
        t.Fatalf("failed to decode response: %v", err)
    }

    if result["app"] != "sample-app" {
        t.Errorf("expected app to be 'sample-app', got '%s'", result["app"])
    }
}
```

#### Step 4：在本地驗證

```bash
# Run tests locally (same flags as CI)
go test -v -race ./...

# Run go vet
go vet ./...

# Check formatting
gofmt -l .
```

確認所有檢查都通過後再繼續。

#### Step 5：Commit 並 Push

```bash
git add handler.go main.go handler_test.go
git commit -m "feat: add /info endpoint"
git push origin feature/add-info-endpoint
```

#### Step 6：開 Pull Request

1. 到 GitHub repository 頁面
2. 你應該會看到一個提示，讓你建立 PR
3. 填寫 PR 標題和描述
4. 點擊 **Create pull request**

#### Step 7：觀察 CI 結果

1. 在 PR 頁面下方，你會看到 CI 檢查的狀態
2. 點擊 **Details** 查看詳細的執行記錄
3. 確認所有檢查都通過（綠色勾勾）

### 思考問題

- 如果 CI 檢查失敗了，你會怎麼做？
- 如果你在 push 之後又發現了一個 bug，你會怎麼處理？（提示：直接在同一個分支上修復並 push，CI 會自動重跑）

---

## 延伸思考

以下問題沒有標準答案，供進階學生思考和討論：

1. **覆蓋率的迷思**：100% 的測試覆蓋率是否代表程式沒有 bug？為什麼？你認為合理的覆蓋率目標應該是多少？

2. **CI 效能最佳化**：如果你的 CI pipeline 跑了 15 分鐘，你會用什麼方法來加速？列出至少三種可能的最佳化策略。

---

[← 練習一：GitHub Actions 基礎練習](01-basics.md) ｜ [回到教材目錄 →](../README.md)
