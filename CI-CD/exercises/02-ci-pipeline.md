# 練習二：CI Pipeline 實戰練習

> **難度：** 中級 | **對應章節：** 03-Go 專案 CI Pipeline

---

## 目錄

- [練習 2-1：擴充 CI Pipeline](#練習-2-1擴充-ci-pipeline)
- [練習 2-2：PR 檢查練習](#練習-2-2pr-檢查練習)

---

## 練習 2-1：擴充 CI Pipeline

### 目標

把 03 章介紹的 `go mod verify` 加進 `ci.yml`，確認依賴套件沒被篡改。

### 要求

在 `ci.yml` 的 `test` job 中，在 `Run tests` 步驟 **之前** 加入：

- **`go mod verify`** — 驗證 `go.sum` 中記錄的 module hash 是否正確，確保 dependencies 沒有被篡改

### 提示

- 用一個獨立的 `- name: Verify dependencies` step

### 預期結果

- 如果 `go.sum` 跟本地 module cache 不一致，CI 會在 `Verify dependencies` 這一步失敗

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

      # === New check ===
      - name: Verify dependencies
        run: go mod verify

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

</details>

---

## 練習 2-2：PR 檢查練習

### 目標

實際走一遍完整的 **PR 檢查流程**：建立分支、做一個小改動、開 PR、觀察 CI 結果。重點不在程式碼有多複雜，而在感受「push → CI → review → merge」這條動線。

### 要求

1. 從 `main` 建立一個新的分支 `feature/update-greeting`
2. 把 `/` endpoint 的歡迎訊息改成你想要的字串，並同步更新對應的 test
3. 開一個 PR，觀察 CI 檢查跑起來、變綠之後合併

### Step by Step 引導

#### 前置：Fork workshop repo

學員對工作坊 repo 沒有寫入權限，所以要先 fork 一份到自己帳號下，後續所有 push、PR 都在你的 fork 內進行。

1. 到 SRE Workshop repo 頁面，右上角按 **Fork**

   ![GitHub Fork 按鈕位置](../assets/github-fork-button.png)

2. 把 fork 到自己帳號下的 repo Clone 下來：

   ```bash
   git clone git@github.com:<你的-username>/sre-workshop.git
   cd sre-workshop/CI-CD/examples/sample-app
   ```

#### Step 1：建立新分支

```bash
git checkout main
git pull origin main
git checkout -b feature/update-greeting
```

#### Step 2：改歡迎訊息

打開 `handler.go`，把 `handleRoot` 裡的字串換成你想要的（例如加上你的名字）：

```go
// before
if _, err := w.Write([]byte("Hello, GitHub Actions!")); err != nil {

// after
if _, err := w.Write([]byte("Hello from <你的名字>!")); err != nil {
```

#### Step 3：同步更新測試

打開 `handler_test.go`，把 `TestHandleRoot` 裡的 `expected` 改成一樣的字串：

```go
expected := "Hello from <你的名字>!"
```

> 故意先**不**改 test 也行 — push 上去 CI 會失敗，可以體驗 CI 擋住合併的感覺，再修正 test 後 push 一次就會變綠。

#### Step 4：在本地驗證

```bash
go test -v -race ./...
```

確認測試通過後再繼續。

#### Step 5：Commit 並 Push

```bash
git add handler.go handler_test.go
git commit -m "chore: personalize greeting message"
git push origin feature/update-greeting
```

#### Step 6：開 Pull Request（在你的 fork 內）

1. 到你 fork 的 GitHub 頁面
2. 點擊 **Pull requests** → **New pull request**
3. **確認 base 是你 fork 的 `main`**（GitHub 預設可能指向 upstream，要切回你自己的 fork）
4. compare 選擇 `feature/add-info-endpoint`
5. 填寫 PR 標題和描述，點擊 **Create pull request**

#### Step 7：觀察 CI 結果

1. 在 PR 頁面下方，你會看到 CI 檢查的狀態
2. 點擊 **Details** 查看詳細的執行記錄
3. 確認所有檢查都通過（綠色勾勾）

### 思考問題

- 如果 CI 檢查失敗了，你會怎麼做？
- 如果你在 push 之後又發現了一個 bug，你會怎麼處理？（提示：直接在同一個分支上修復並 push，CI 會自動重跑）

[← 練習一：GitHub Actions 基礎練習](01-basics.md) ｜ [回到教材目錄 →](../README.md)
