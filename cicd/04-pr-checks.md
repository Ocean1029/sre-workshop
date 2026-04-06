# 05 — PR 檢查自動化

> **時間：15 分鐘**

---

Ocean 把 CI pipeline 設好之後很得意，直接在 main 分支上改了一段程式碼就 push 上去。結果 Andrew 隔天一拉程式碼，發現 build 壞了。Snow 搖搖頭說：「你不能讓每個人都直接 push 到 main 啊，我來教你怎麼用 PR 檢查保護主分支。」

---

## 目錄

- [學習目標](#學習目標)
- [為什麼需要 PR 檢查？](#為什麼需要-pr-檢查)
- [PR 觸發的 Workflow](#pr-觸發的-workflow)
- [完整 PR 檢查 Workflow](#完整-pr-檢查-workflow)
- [使用 Reviewdog 進行 PR 註解](#使用-reviewdog-進行-pr-註解)
- [自動標籤（Auto Labeling）](#自動標籤auto-labeling)
- [PR 範本（Pull Request Template）](#pr-範本pull-request-template)
- [Branch Protection Rules](#branch-protection-rules)
- [Status Checks 的運作方式](#status-checks-的運作方式)
- [常見問題排解](#常見問題排解)
- [小結與練習題](#小結與練習題)

---

## 學習目標

完成本章節後，你將能夠：

- 理解 **PR 檢查** 在團隊協作中的重要性
- 建立 **PR 觸發的 CI workflow**，自動執行 lint、test、build 檢查
- 使用 **Reviewdog** 將 lint 結果直接顯示在 PR 的 diff 上
- 設定 **自動標籤（Auto Labeling）** 讓 PR 根據修改的檔案自動加上標籤
- 建立 **PR 範本**，幫助團隊成員撰寫高品質的 PR 描述
- 了解 **Branch Protection Rules** 與 **Status Checks** 的運作方式

---

## 為什麼需要 PR 檢查？

### Code Review 的重要性

在團隊開發中，**Pull Request（PR）** 是程式碼進入主分支前的最後一道關卡。Code Review 可以幫助團隊：

- **發現潛在的 bug**：一雙以上的眼睛看程式碼，更容易發現邏輯錯誤
- **維持程式碼風格一致性**：確保團隊的程式碼遵循相同的規範
- **知識分享**：讓團隊成員了解彼此的程式碼，降低 bus factor（只有一個人懂某段程式碼的風險）

### 自動化檢查讓 Reviewer 專注在邏輯

但如果 reviewer 需要花大量時間在格式檢查、變數命名等機械式的問題上，就沒有精力去看更重要的 **商業邏輯** 和 **架構設計**。

自動化 PR 檢查可以幫我們處理這些 **機械式的檢查**：

```
┌──────────────────────────────────────────────────────┐
│                 PR 提交後自動執行                      │
├──────────────────────────────────────────────────────┤
│                                                      │
│  🤖 自動化檢查                 👤 人工 Review          │
│  ├── 程式碼格式（Lint）         ├── 商業邏輯正確性       │
│  ├── 單元測試                  ├── 架構設計合理性       │
│  ├── 建置成功                  ├── 安全性考量          │
│  ├── 測試覆蓋率門檻            ├── 效能影響            │
│  └── 依賴整潔度                └── 可維護性            │
│                                                      │
└──────────────────────────────────────────────────────┘
```

### 確保品質門檻

透過 **自動化 + Branch Protection**，我們可以設定「品質門檻」：

- 所有自動檢查都通過才能合併
- 至少一位 reviewer 核准才能合併
- 禁止直接 push 到 `main` 分支

這些機制確保每一行進入 `main` 的程式碼都經過了品質把關。

<details>
<summary>💡 講師提示</summary>

> 可以問學生：「你有沒有遇過改了一個地方，結果其他功能壞掉的經驗？」用實際情境引導他們理解自動化檢查的價值。

</details>

---

## PR 觸發的 Workflow

### `on: pull_request` 的各種觸發類型

PR workflow 的觸發設定比 push 更豐富。`pull_request` 事件可以指定多種 **activity types**：

```yaml
on:
  pull_request:
    branches: [main]
    types: [opened, synchronize, reopened]
```

常用的 `types`：

| Type | 觸發時機 | 說明 |
|------|---------|------|
| `opened` | PR 剛建立時 | 第一次提交 PR |
| `synchronize` | 有新 commit push 到 PR 時 | 更新程式碼後重新檢查 |
| `reopened` | 關閉後重新開啟 PR 時 | 重新開啟需要再次檢查 |
| `edited` | PR 的標題或描述被修改時 | 通常不需要重跑 CI |
| `ready_for_review` | 從 Draft 變成 Ready 時 | 適合只在正式 review 時才跑的檢查 |
| `closed` | PR 被關閉或合併時 | 可用於清理資源 |

<details>
<summary>💡 講師提示</summary>

> 如果不指定 `types`，GitHub 預設只會在 `opened`、`synchronize`、`reopened` 三種情況下觸發，這通常就夠用了。

</details>

### PR Workflow 與 Push Workflow 的差異

| 差異項目 | Push Workflow | PR Workflow |
|---------|--------------|-------------|
| 觸發時機 | 程式碼 push 到指定分支 | PR 建立或更新時 |
| 檢查的程式碼 | push 後的最新 commit | **模擬合併後的結果** |
| 目的 | 確認分支上的程式碼正確 | 確認合併後不會壞掉 |
| 適用場景 | 持續整合 | 合併前的品質把關 |

一個重要的細節：PR workflow 會在 **merge commit** 上執行，測試的不只是你的 PR 程式碼，而是「你的修改合併進目標分支後的結果」。這樣可以及早發現合併衝突或不相容的問題。

### PR 的安全限制

當有人從 **fork** 的 repository 提交 PR 時，GitHub 會有一些安全限制：

- ❌ **無法存取 secrets**：fork PR 的 workflow 無法讀取 repository 的 secrets
- ❌ **寫入權限受限**：fork PR 的 workflow 預設只有 read 權限
- ✅ **程式碼檢查仍然可以執行**：lint、test、build 等不需要 secrets 的檢查不受影響

這是為了防止惡意的 fork PR 竊取你的 secrets 或修改你的 repository。

---

## 完整 PR 檢查 Workflow

以下是一個完整的 PR 檢查 workflow，請在 `.github/workflows/pr-check.yml` 中建立此檔案：

```yaml
name: PR Checks

on:
  pull_request:
    branches: [main]
    types: [opened, synchronize, reopened]

permissions:
  contents: read
  pull-requests: write
  checks: write

jobs:
  lint:
    name: Code Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest

  test:
    name: Unit Tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Run tests with coverage
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Total coverage: ${COVERAGE}%"
          if (( $(echo "$COVERAGE < 60" | bc -l) )); then
            echo "::error::Coverage ${COVERAGE}% is below threshold 60%"
            exit 1
          fi

  build-check:
    name: Build Check
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Verify build
        run: go build ./...
      - name: Check go mod tidy
        run: |
          go mod tidy
          git diff --exit-code go.mod go.sum
```

### 逐段解說

#### Permissions 設定

```yaml
permissions:
  contents: read
  pull-requests: write
  checks: write
```

| 權限 | 用途 |
|------|------|
| `contents: read` | 讀取 repository 的程式碼 |
| `pull-requests: write` | 允許在 PR 上留言或更新 status |
| `checks: write` | 允許建立 check run 結果 |

遵循 **最小權限原則**，只給 workflow 需要的權限，不多給。

#### Lint Job — 程式碼風格檢查

```yaml
lint:
  name: Code Lint
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.22'
    - name: golangci-lint
      uses: golangci/golangci-lint-action@v6
      with:
        version: latest
```

使用 `golangci-lint` 進行靜態分析。如果有任何 lint 違規，這個 job 就會失敗，PR 上會出現紅色的 ✗ 標記。

#### Test Job — 單元測試與覆蓋率門檻

這個 job 除了跑測試，還多了一個 **覆蓋率門檻檢查**：

```yaml
- name: Check coverage threshold
  run: |
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    echo "Total coverage: ${COVERAGE}%"
    if (( $(echo "$COVERAGE < 60" | bc -l) )); then
      echo "::error::Coverage ${COVERAGE}% is below threshold 60%"
      exit 1
    fi
```

逐行解說：

1. **擷取覆蓋率數值**：從 `go tool cover` 的輸出中取出總覆蓋率百分比
2. **印出覆蓋率**：方便在 log 中查看
3. **門檻檢查**：如果覆蓋率低於 60%，就用 `::error::` 語法在 GitHub Actions 中顯示錯誤訊息，並以 `exit 1` 讓 job 失敗

<details>
<summary>💡 講師提示</summary>

> `::error::` 是 GitHub Actions 的 **Workflow Command**，可以在 Actions UI 中顯示錯誤訊息。類似的還有 `::warning::` 和 `::notice::`。60% 只是範例門檻，實際專案中可以根據團隊情況調整。

</details>

#### Build Check Job — 建置與依賴檢查

```yaml
- name: Check go mod tidy
  run: |
    go mod tidy
    git diff --exit-code go.mod go.sum
```

這段檢查非常實用，它確保 `go.mod` 和 `go.sum` 是「乾淨」的：

1. 先執行 `go mod tidy`（清理不需要的依賴、補齊缺少的依賴）
2. 再用 `git diff --exit-code` 檢查是否有差異

如果有人忘記跑 `go mod tidy` 就提交 PR，這個檢查就會失敗，提醒開發者補上。

---

## 使用 Reviewdog 進行 PR 註解

### Reviewdog 是什麼？

[Reviewdog](https://github.com/reviewdog/reviewdog) 是一個 **自動 code review 工具**，它可以將 linter 的檢查結果 **直接顯示在 PR 的 diff 上**，就像人類 reviewer 留下的 comment 一樣。

```
┌──────────────────────────────────────────────────────┐
│  PR Diff View                                        │
├──────────────────────────────────────────────────────┤
│                                                      │
│  handler.go                                          │
│  ─────────────                                       │
│  10 │ func HomeHandler(w http.ResponseWriter, ...    │
│  11 │     data := map[string]string{                 │
│  12 │         "msg": "hello",                        │
│     │  🐶 reviewdog [golangci-lint]                   │
│     │  exported function HomeHandler should have     │
│     │  comment or be unexported (golint)             │
│  13 │     }                                          │
│                                                      │
└──────────────────────────────────────────────────────┘
```

### 使用 Reviewdog 的 Workflow 範例

```yaml
name: Reviewdog

on:
  pull_request:
    branches: [main]

permissions:
  contents: read
  pull-requests: write

jobs:
  reviewdog-lint:
    name: Reviewdog Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Run golangci-lint with Reviewdog
        uses: reviewdog/action-golangci-lint@v2
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          reporter: github-pr-review
          level: warning
          golangci_lint_flags: "--config=.golangci.yml"

      - name: Run staticcheck with Reviewdog
        uses: reviewdog/action-staticcheck@v1
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          reporter: github-pr-review
          level: warning
```

#### 參數說明

| 參數 | 說明 |
|------|------|
| `reporter: github-pr-review` | 將結果以 PR review comment 的形式顯示在 diff 上 |
| `reporter: github-check` | 將結果以 check annotation 的形式顯示 |
| `reporter: github-pr-check` | 將結果以 PR check 的形式顯示 |
| `level` | 報告的嚴重程度：`info`、`warning`、`error` |

使用 `github-pr-review` 的好處是 reviewer 不需要切到 CI log 去看錯誤，**直接在 diff 的對應行就能看到問題**，大幅提升 review 效率。

<details>
<summary>💡 講師提示</summary>

> 如果有 GitHub 帳號可以直接展示 Reviewdog 在 PR 上的效果，讓學生看到自動 comment 出現在 diff 上，會非常直覺。

</details>

---

## 自動標籤（Auto Labeling）

### 使用 actions/labeler 自動加標籤

當 PR 修改了特定目錄的檔案時，可以 **自動加上對應的標籤**。例如修改了 `docs/` 下的檔案就加上 `documentation` 標籤，修改了 `internal/` 就加上 `backend` 標籤。

#### Workflow 設定

```yaml
name: Auto Labeler

on:
  pull_request:
    types: [opened, synchronize, reopened]

permissions:
  contents: read
  pull-requests: write

jobs:
  label:
    name: Auto Label
    runs-on: ubuntu-latest
    steps:
      - uses: actions/labeler@v5
        with:
          repo-token: ${{ secrets.GITHUB_TOKEN }}
```

#### `.github/labeler.yml` 設定範例

在 repository 根目錄的 `.github/labeler.yml` 中定義標籤規則：

```yaml
# Add 'documentation' label when docs files are changed
documentation:
  - changed-files:
      - any-glob-to-any-file: 'docs/**'

# Add 'ci' label when workflow files are changed
ci:
  - changed-files:
      - any-glob-to-any-file: '.github/**'

# Add 'go' label when Go source files are changed
go:
  - changed-files:
      - any-glob-to-any-file: '**/*.go'

# Add 'docker' label when Dockerfile is changed
docker:
  - changed-files:
      - any-glob-to-any-file: 'Dockerfile'
      - any-glob-to-any-file: 'docker-compose*.yml'

# Add 'dependencies' label when go.mod or go.sum are changed
dependencies:
  - changed-files:
      - any-glob-to-any-file:
          - 'go.mod'
          - 'go.sum'

# Add 'tests' label when test files are changed
tests:
  - changed-files:
      - any-glob-to-any-file: '**/*_test.go'
```

自動標籤的好處：

- **快速分類 PR**：一眼就知道這個 PR 改了哪些類型的東西
- **方便篩選**：在 PR 列表中可以用標籤篩選
- **輔助 review 分配**：例如有 `database` 標籤的 PR 可以指派給 DBA review

---

## PR 範本（Pull Request Template）

### 建立 PR 範本

在 `.github/pull_request_template.md` 建立 PR 範本，當團隊成員建立新的 PR 時，描述欄位會 **自動填入範本內容**，引導他們提供完整的資訊：

```markdown
## Description

<!-- Briefly describe what this PR does -->

## Type of Change

- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Refactoring (no functional changes)
- [ ] Documentation update
- [ ] CI/CD changes

## How Has This Been Tested?

<!-- Describe the tests you ran to verify your changes -->

- [ ] Unit tests
- [ ] Integration tests
- [ ] Manual testing

## Checklist

- [ ] My code follows the project's coding style
- [ ] I have added tests that prove my fix/feature works
- [ ] All new and existing tests pass
- [ ] I have updated the documentation (if applicable)
- [ ] I have run `go mod tidy`
- [ ] I have run `golangci-lint run`

## Screenshots (if applicable)

<!-- Add screenshots to help explain your changes -->

## Related Issues

<!-- Link any related issues: Fixes #123, Closes #456 -->
```

### 為什麼需要 PR 範本？

| 沒有範本 | 有範本 |
|---------|--------|
| PR 描述常常是空的 | 引導撰寫完整描述 |
| Reviewer 不知道改了什麼、為什麼 | 有清楚的上下文 |
| 忘記跑測試或 lint | Checklist 提醒 |
| 不知道如何測試這個變更 | 明確的測試說明 |

<details>
<summary>💡 講師提示</summary>

> 可以展示一些知名開源專案（如 Kubernetes、Go 本身）的 PR 範本，讓學生知道這在業界是標準做法。

</details>

---

## Branch Protection Rules

### 什麼是 Branch Protection？

**Branch Protection Rules** 是 GitHub 提供的功能，可以保護重要的分支（通常是 `main`），防止不符合規則的變更被合併或直接 push 進去。

```
┌─────────────────────────────────────────────────────────┐
│                   main branch                           │
│                                                         │
│  ┌────────────────────────────────────────────────────┐  │
│  │              Branch Protection Rules               │  │
│  │                                                    │  │
│  │  ✓ 要求 PR review（至少 1 位 reviewer 核准）        │  │
│  │  ✓ 要求 status checks 通過（lint, test, build）    │  │
│  │  ✓ 禁止直接 push                                  │  │
│  │  ✓ 禁止 force push                                │  │
│  │  ✓ 要求分支是最新的（up to date）                   │  │
│  │                                                    │  │
│  └────────────────────────────────────────────────────┘  │
│                                                         │
│  只有通過所有規則的 PR 才能合併進來                        │
└─────────────────────────────────────────────────────────┘
```

### 如何設定 Branch Protection Rules

按照以下步驟在 GitHub 上設定：

**Step 1：進入設定頁面**
1. 到你的 GitHub repository 頁面
2. 點擊 **Settings** 標籤
3. 左側選單點擊 **Branches**
4. 在 **Branch protection rules** 區域點擊 **Add branch protection rule**

**Step 2：設定分支名稱**
1. 在 **Branch name pattern** 輸入 `main`
2. 這表示規則會套用到 `main` 分支

**Step 3：設定 PR Review 要求**
1. 勾選 **Require a pull request before merging**
2. 勾選 **Require approvals**
3. 設定 **Required number of approvals before merging** 為 `1`（至少 1 位 reviewer 核准）
4. 建議也勾選 **Dismiss stale pull request approvals when new commits are pushed**（新的 commit push 後，之前的核准會被撤銷，強制 reviewer 重新 review）

**Step 4：設定 Status Checks**
1. 勾選 **Require status checks to pass before merging**
2. 勾選 **Require branches to be up to date before merging**
3. 在搜尋框中搜尋你的 workflow job 名稱（例如 `Code Lint`、`Unit Tests`、`Build Check`）
4. 將它們加入 required status checks

**Step 5：禁止直接 Push**
1. 勾選 **Do not allow bypassing the above settings**（連管理員都必須遵守）
2. 確保沒有勾選 **Allow force pushes**

**Step 6：儲存**
1. 點擊 **Create** 按鈕儲存規則

<details>
<summary>💡 講師提示</summary>

> 建議在課堂上實際操作一次設定流程，讓學生看到完整的畫面。如果是免費的 GitHub 帳號，某些功能（如 required reviewers 數量 > 1）可能需要升級方案才能使用。

</details>

---

## Status Checks 的運作方式

### Workflow 如何與 PR Status Checks 連結

當你在 PR 上觸發一個 workflow 時，GitHub 會自動將每個 **job** 變成一個 **status check**。在 PR 頁面的底部，你可以看到所有 status checks 的結果：

```
┌──────────────────────────────────────────────────┐
│  PR #42: Add user authentication                 │
├──────────────────────────────────────────────────┤
│                                                  │
│  Status Checks                                   │
│  ─────────────                                   │
│  ✅ Code Lint          — passed (45s)            │
│  ✅ Unit Tests         — passed (1m 12s)         │
│  ✅ Build Check        — passed (32s)            │
│  ✅ Auto Label         — passed (8s)             │
│                                                  │
│  ┌──────────────────────────────────────────┐    │
│  │  ✅ All checks have passed               │    │
│  │  🟢 Ready to merge                       │    │
│  └──────────────────────────────────────────┘    │
│                                                  │
└──────────────────────────────────────────────────┘
```

如果有任何 check 失敗：

```
┌──────────────────────────────────────────────────┐
│  PR #42: Add user authentication                 │
├──────────────────────────────────────────────────┤
│                                                  │
│  Status Checks                                   │
│  ─────────────                                   │
│  ❌ Code Lint          — failed (23s)            │
│  ✅ Unit Tests         — passed (1m 12s)         │
│  ✅ Build Check        — passed (32s)            │
│                                                  │
│  ┌──────────────────────────────────────────┐    │
│  │  ❌ Some checks have failed              │    │
│  │  🔴 Merging is blocked                   │    │
│  └──────────────────────────────────────────┘    │
│                                                  │
└──────────────────────────────────────────────────┘
```

### Required vs Optional Checks

在 Branch Protection Rules 中，你可以設定哪些 checks 是 **必要的（Required）**，哪些是 **選擇性的（Optional）**：

| 類型 | 失敗時 | 適用場景 |
|------|--------|---------|
| **Required** | 阻止合併，PR 無法 merge | 核心品質檢查：lint、test、build |
| **Optional** | 不阻止合併，只是顯示警告 | 非關鍵檢查：code coverage 報告、performance benchmark |

建議的設定：

| Check | 建議設為 | 理由 |
|-------|---------|------|
| Code Lint | Required | 確保程式碼風格一致 |
| Unit Tests | Required | 確保功能正確 |
| Build Check | Required | 確保程式碼可以成功編譯 |
| Auto Label | Optional | 標籤只是輔助分類，失敗不影響品質 |
| Coverage Report | Optional | 覆蓋率報告是參考用，不一定要阻止合併 |

### Status Check 的命名

Status check 的名稱對應的是 workflow 中 job 的 `name` 屬性：

```yaml
jobs:
  lint:
    name: Code Lint       # ← This becomes the status check name
    runs-on: ubuntu-latest
    # ...
```

所以在 Branch Protection 中搜尋 `Code Lint` 就能找到這個 check。**請確保 job 的 name 有意義且穩定**，因為改了 name 就需要重新設定 Branch Protection。

---

## 常見問題排解

### 1. 設定 Branch Protection 後自己也不能 push

這是正常的！Branch Protection 設計上就是要阻止所有人（包括管理員）直接 push 到受保護的分支。你需要透過 PR 來合併程式碼。

如果你勾選了 **Do not allow bypassing the above settings**，連 repository 管理員都必須遵守規則。如果你想在緊急情況下保留繞過的能力，可以不勾這個選項。

### 2. Status Check 搜尋不到 workflow 的 job 名稱

Status check 只有在 workflow **至少執行過一次** 之後才會出現在搜尋清單中。如果你剛建立 workflow 還沒觸發過，先 push 一次讓它跑完，再回來設定 Branch Protection。

### 3. Fork PR 的 CI 檢查失敗

從 fork 提交的 PR 無法存取 repository 的 secrets，所以如果你的 CI 步驟需要 secrets（例如推送 Docker image），這些步驟會失敗。解法是把需要 secrets 的步驟用 `if: github.event.pull_request.head.repo.full_name == github.repository` 條件跳過。

---

## 小結與練習題

### 本章重點回顧

- **PR 檢查** 讓自動化工具處理機械式的檢查，讓 reviewer 專注在邏輯
- `on: pull_request` 支援多種 activity types，預設觸發 `opened`、`synchronize`、`reopened`
- **PR workflow 與 push workflow 的差異**：PR workflow 在模擬合併後的結果上執行
- **Fork PR 的安全限制**：無法存取 secrets，寫入權限受限
- **Permissions** 應遵循最小權限原則
- **覆蓋率門檻** 可以用 `::error::` workflow command 在 CI 中顯示錯誤
- **Reviewdog** 可以將 lint 結果顯示在 PR diff 的對應行上
- **Auto Labeling** 根據修改的檔案自動分類 PR
- **PR 範本** 引導開發者撰寫完整的 PR 描述
- **Branch Protection Rules** 保護主分支，確保只有通過檢查的 PR 才能合併
- **Status Checks** 分為 Required 和 Optional，Required checks 失敗會阻止合併

### 練習題

完成以下練習來鞏固本章所學：

👉 [練習二：CI Pipeline 實戰練習](exercises/exercise-02-ci-pipeline.md)（練習 2-3 至 2-4）

> **接下來，我們將學習如何在打 Git Tag 後自動建立 Release 並發佈建置產物！**

---

[← 上一章：Go 專案 CI Pipeline](04-go-ci-pipeline.md) ｜ [下一章：Release 自動化 →](06-release-automation.md)
