# 06 — Release 自動化

> **時間：20 分鐘**

---

## 目錄

- [學習目標](#學習目標)
- [Semantic Versioning（語意化版號）](#semantic-versioning語意化版號)
- [Git Tag 基礎](#git-tag-基礎)
- [Tag 觸發的 Release Workflow](#tag-觸發的-release-workflow)
- [ldflags 注入版本資訊](#ldflags-注入版本資訊)
- [Cross Compilation](#cross-compilation)
- [Changelog 自動產生](#changelog-自動產生)
- [Docker Image 建置與推送](#docker-image-建置與推送)
- [進階：使用 GoReleaser](#進階使用-goreleaser)
- [小結與練習題](#小結與練習題)

---

## 學習目標

完成本章節後，你將能夠：

- 理解 **Semantic Versioning** 的規則，知道何時該 bump 哪個版號
- 使用 **Git Tag** 標記版本，並推送到 remote
- 建立由 **tag push 觸發** 的 Release workflow，自動建置並發佈
- 使用 **ldflags** 在編譯時注入版本資訊到 Go 程式中
- 使用 **matrix strategy** 進行 cross compilation，產出多平台的 binary
- 自動產生 **changelog**，讓使用者了解每個版本的變更
- 在 release workflow 中加入 **Docker image 建置與推送**
- 了解 **GoReleaser** 工具的基本用法

---

## Semantic Versioning（語意化版號）

### 什麼是 SemVer？

**Semantic Versioning（語意化版號）** 是一套版號命名規則，格式為：

```
vMAJOR.MINOR.PATCH
```

例如：`v1.2.3`

每一個數字都有明確的意義：

```
v 1  .  2  .  3
│ │     │     │
│ │     │     └── PATCH：修正 bug，不影響功能
│ │     └──────── MINOR：新增功能，向下相容
│ └────────────── MAJOR：破壞性變更，不向下相容
└──────────────── prefix（Git tag 慣例）
```

### 各版號的意義

| 版號 | 何時遞增 | 向下相容？ | 範例 |
|------|---------|-----------|------|
| **MAJOR** | 有破壞性變更（Breaking Change） | ❌ 不相容 | `v1.0.0` → `v2.0.0` |
| **MINOR** | 新增功能，但不影響既有功能 | ✅ 相容 | `v1.0.0` → `v1.1.0` |
| **PATCH** | 修正 bug，不新增功能 | ✅ 相容 | `v1.0.0` → `v1.0.1` |

### 版號遞增範例

以我們的範例專案為例：

| 變更 | 版號變化 | 理由 |
|------|---------|------|
| 專案初始發佈 | `v1.0.0` | 第一個穩定版本 |
| 新增 `/users` endpoint | `v1.0.0` → `v1.1.0` | 新增功能，不影響 `/`, `/health`, `/version` |
| 修正 `/health` 回傳格式 bug | `v1.1.0` → `v1.1.1` | 修正 bug |
| 將 `/version` 的 JSON 結構改變 | `v1.1.1` → `v2.0.0` | 原本的使用者程式碼會壞掉，是 breaking change |

> 💡 **講師提示：** 可以用「手機 App 更新」來比喻——PATCH 像是修 bug 的小更新、MINOR 像是新增功能的更新、MAJOR 像是整個 UI 大改版到你找不到原本功能的更新。

---

## Git Tag 基礎

### 什麼是 Git Tag？

**Git Tag** 是 Git 中用來標記特定 commit 的「書籤」。它通常用來標記 **release 版本**，讓你可以快速找到每個版本對應的程式碼。

```
commit A ─── commit B ─── commit C ─── commit D ─── commit E
                           │                         │
                           v1.0.0                    v1.1.0
                           (tag)                     (tag)
```

### 建立 Tag

```bash
# Create a lightweight tag
git tag v1.0.0

# Create an annotated tag (recommended)
git tag -a v1.0.0 -m "Release version 1.0.0"
```

### Lightweight Tag vs Annotated Tag

| 類型 | 指令 | 內容 | 建議用途 |
|------|------|------|---------|
| **Lightweight** | `git tag v1.0.0` | 只是一個指向 commit 的指標 | 臨時標記 |
| **Annotated** | `git tag -a v1.0.0 -m "..."` | 包含建立者、日期、訊息 | Release 版本（推薦） |

Annotated tag 會記錄 **誰** 在 **什麼時候** 建立了這個 tag，以及附帶的訊息，資訊更完整，是標記 release 版本的推薦做法。

### 推送 Tag 到 Remote

建立 tag 後，它只存在於你的本地端。要讓 GitHub 知道這個 tag，需要推送：

```bash
# Push a specific tag
git push origin v1.0.0
```

> **注意：** `git push` 預設不會推送 tag，你需要明確指定要推送的 tag。

### 常用的 Tag 操作

```bash
# List all tags
git tag -l

# List tags matching a pattern
git tag -l "v1.*"

# Show tag details
git show v1.0.0

# Delete a local tag
git tag -d v1.0.0

# Delete a remote tag
git push origin --delete v1.0.0
```

> 💡 **講師提示：** 建議現場操作一次 tag 的建立和推送流程，讓學生看到 tag 推送後 GitHub 上的 Tags 頁面會出現新的 tag。

---

## Tag 觸發的 Release Workflow

### 完整的 Release Workflow

以下是由 tag push 觸發的完整 release workflow。當你推送一個 `v*` 格式的 tag 時，它會自動：

1. 執行測試，確保程式碼正確
2. 為多個平台建置 binary
3. 建立 GitHub Release 並上傳所有 binary

請在 `.github/workflows/release.yml` 中建立此檔案：

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write
  packages: write

jobs:
  test:
    name: Test Before Release
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - run: go test -v -race ./...

  build:
    name: Build Binaries
    needs: test
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, darwin, windows]
        goarch: [amd64, arm64]
        exclude:
          - goos: windows
            goarch: arm64
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Build
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
        run: |
          BINARY_NAME=myapp-${{ matrix.goos }}-${{ matrix.goarch }}
          if [ "${{ matrix.goos }}" = "windows" ]; then
            BINARY_NAME="${BINARY_NAME}.exe"
          fi
          go build -ldflags "-X main.version=${{ github.ref_name }}" -o ${BINARY_NAME} .
      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: myapp-${{ matrix.goos }}-${{ matrix.goarch }}
          path: myapp-*

  release:
    name: Create Release
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Download all artifacts
        uses: actions/download-artifact@v4
        with:
          path: artifacts
          merge-multiple: true
      - name: Generate changelog
        id: changelog
        run: |
          PREVIOUS_TAG=$(git tag --sort=-version:refname | head -2 | tail -1)
          if [ -z "$PREVIOUS_TAG" ]; then
            CHANGELOG=$(git log --pretty=format:"- %s (%h)" HEAD)
          else
            CHANGELOG=$(git log --pretty=format:"- %s (%h)" ${PREVIOUS_TAG}..HEAD)
          fi
          echo "CHANGELOG<<EOF" >> $GITHUB_OUTPUT
          echo "$CHANGELOG" >> $GITHUB_OUTPUT
          echo "EOF" >> $GITHUB_OUTPUT
      - name: Create GitHub Release
        uses: softprops/action-gh-release@v2
        with:
          body: |
            ## What's Changed
            ${{ steps.changelog.outputs.CHANGELOG }}
          files: artifacts/*
          generate_release_notes: true
```

### 逐段解說

#### 觸發條件：Tag Push

```yaml
on:
  push:
    tags:
      - 'v*'
```

這個 workflow **只在推送 tag 時觸發**，而且 tag 必須以 `v` 開頭（例如 `v1.0.0`、`v2.1.3`）。這確保只有正式的版本 tag 才會觸發 release 流程。

`${{ github.ref_name }}` 會自動取得 tag 的名稱（例如 `v1.0.0`），我們稍後會用它來注入版本資訊。

#### Permissions

```yaml
permissions:
  contents: write    # Required to create GitHub Release
  packages: write    # Required to push Docker image to GHCR
```

`contents: write` 是建立 GitHub Release 所必需的權限。

#### Test Job — Release 前的最後測試

```yaml
test:
  name: Test Before Release
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.22'
    - run: go test -v -race ./...
```

即使 PR 合併前已經通過了測試，在 release 前 **再跑一次** 是一個好習慣。這是你的最後一道安全網，確保要發佈的程式碼確實是正確的。

#### Build Job — 多平台建置

```yaml
build:
  name: Build Binaries
  needs: test
  strategy:
    matrix:
      goos: [linux, darwin, windows]
      goarch: [amd64, arm64]
      exclude:
        - goos: windows
          goarch: arm64
```

使用 matrix strategy 為多個平台建置 binary。這段會在後面的 [Cross Compilation](#cross-compilation) 章節詳細說明。

#### Release Job — 建立 GitHub Release

```yaml
release:
  name: Create Release
  needs: build
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0   # Fetch all history for changelog generation
```

`fetch-depth: 0` 表示取得完整的 git history，因為我們需要從 git log 產生 changelog。預設的 `fetch-depth: 1` 只會取得最新的一個 commit。

#### 流程圖

```
┌────────────────┐
│   git push     │
│   origin v1.0.0│
└───────┬────────┘
        │
        ▼
┌────────────────┐
│   Test Job     │
│  go test ...   │
└───────┬────────┘
        │ (通過)
        ▼
┌────────────────────────────────────────────────┐
│              Build Job (Matrix)                │
│                                                │
│  ┌──────────────┐  ┌──────────────┐            │
│  │ linux/amd64  │  │ linux/arm64  │            │
│  └──────────────┘  └──────────────┘            │
│  ┌──────────────┐  ┌──────────────┐            │
│  │ darwin/amd64 │  │ darwin/arm64 │            │
│  └──────────────┘  └──────────────┘            │
│  ┌──────────────┐                              │
│  │ windows/amd64│  (windows/arm64 excluded)    │
│  └──────────────┘                              │
│                                                │
└───────────────────────┬────────────────────────┘
                        │ (全部完成)
                        ▼
              ┌──────────────────┐
              │  Release Job     │
              │  ├ Download      │
              │  ├ Changelog     │
              │  └ GitHub Release│
              └──────────────────┘
```

---

## ldflags 注入版本資訊

### 原理

Go 的 `go build` 指令支援 `-ldflags` 參數，可以在 **編譯時** 覆寫程式碼中的變數值。`-X` flag 可以設定指定 package 中的 `string` 變數：

```bash
go build -ldflags "-X main.version=v1.0.0" -o myapp .
```

這行指令會在編譯時將 `main` package 中的 `version` 變數設為 `v1.0.0`。

### 在 Go 程式中讀取版本

首先，在 `main.go` 中宣告一個 `version` 變數：

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
)

// version is set at build time via ldflags
var version = "dev"

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", HomeHandler)
    mux.HandleFunc("/health", HealthHandler)
    mux.HandleFunc("/version", VersionHandler)

    log.Printf("Starting server v%s on :8080", version)
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

然後在 `handler.go` 的 `VersionHandler` 中使用它：

```go
// VersionHandler returns the application version
func VersionHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{
        "version": version,
    })
}
```

### 效果

| 建置方式 | version 的值 | 說明 |
|---------|-------------|------|
| `go build .` | `dev` | 預設值，開發時使用 |
| `go build -ldflags "-X main.version=v1.0.0" .` | `v1.0.0` | 編譯時注入版本 |
| CI: `go build -ldflags "-X main.version=${{ github.ref_name }}" .` | `v1.2.3`（tag 名稱） | 自動從 tag 取得版本 |

```bash
# After building with ldflags
$ curl http://localhost:8080/version
{"version":"v1.0.0"}
```

> 💡 **講師提示：** 可以現場示範 `go build` 和 `go build -ldflags ...` 的差異，然後分別 curl `/version` endpoint，讓學生看到版本的變化。這是一個非常實用的技巧，在業界廣泛使用。

---

## Cross Compilation

### GOOS 和 GOARCH 的組合

Go 的一大優勢是 **內建跨平台編譯**——不需要安裝任何額外工具，只要設定 `GOOS` 和 `GOARCH` 環境變數，就能在一台機器上編譯出其他平台的 binary。

| GOOS | GOARCH | 對應平台 |
|------|--------|---------|
| `linux` | `amd64` | Linux 64-bit（最常見的伺服器） |
| `linux` | `arm64` | Linux ARM（如 AWS Graviton、Raspberry Pi 4） |
| `darwin` | `amd64` | macOS Intel |
| `darwin` | `arm64` | macOS Apple Silicon (M1/M2/M3) |
| `windows` | `amd64` | Windows 64-bit |

### Matrix 的 exclude 用法

在我們的 workflow 中：

```yaml
strategy:
  matrix:
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    exclude:
      - goos: windows
        goarch: arm64
```

原本 3 x 2 = 6 種組合，但我們用 `exclude` 排除了 `windows/arm64`，最終產生 **5 個** build job。

### 為什麼排除某些組合？

| 排除的組合 | 理由 |
|-----------|------|
| `windows/arm64` | Windows ARM 裝置（如 Surface Pro X）目前佔比極小，且 Windows ARM 通常可以透過模擬層執行 amd64 binary |

你也可以根據專案需求排除其他組合。例如，如果你的應用只會部署在 Linux 伺服器上，就可以只保留 `linux/amd64` 和 `linux/arm64`。

### Build Step 的 Binary 命名

```yaml
- name: Build
  env:
    GOOS: ${{ matrix.goos }}
    GOARCH: ${{ matrix.goarch }}
  run: |
    BINARY_NAME=myapp-${{ matrix.goos }}-${{ matrix.goarch }}
    if [ "${{ matrix.goos }}" = "windows" ]; then
      BINARY_NAME="${BINARY_NAME}.exe"
    fi
    go build -ldflags "-X main.version=${{ github.ref_name }}" -o ${BINARY_NAME} .
```

每個平台的 binary 會用 `myapp-{os}-{arch}` 的格式命名：

| Binary 名稱 | 平台 |
|-------------|------|
| `myapp-linux-amd64` | Linux 64-bit |
| `myapp-linux-arm64` | Linux ARM |
| `myapp-darwin-amd64` | macOS Intel |
| `myapp-darwin-arm64` | macOS Apple Silicon |
| `myapp-windows-amd64.exe` | Windows（加上 `.exe` 副檔名） |

---

## Changelog 自動產生

### 從 Git Log 產生 Changelog

Release workflow 中的 changelog 產生步驟會從 git log 中提取兩個 tag 之間的所有 commit 訊息：

```yaml
- name: Generate changelog
  id: changelog
  run: |
    PREVIOUS_TAG=$(git tag --sort=-version:refname | head -2 | tail -1)
    if [ -z "$PREVIOUS_TAG" ]; then
      CHANGELOG=$(git log --pretty=format:"- %s (%h)" HEAD)
    else
      CHANGELOG=$(git log --pretty=format:"- %s (%h)" ${PREVIOUS_TAG}..HEAD)
    fi
    echo "CHANGELOG<<EOF" >> $GITHUB_OUTPUT
    echo "$CHANGELOG" >> $GITHUB_OUTPUT
    echo "EOF" >> $GITHUB_OUTPUT
```

逐行解說：

1. **找到前一個 tag**：`git tag --sort=-version:refname` 按版號降序排列，`head -2 | tail -1` 取第二個（也就是前一個版本）
2. **處理第一個 release**：如果沒有前一個 tag，就列出所有 commit
3. **產生 changelog**：用 `git log` 列出兩個 tag 之間的所有 commit，格式為 `- commit 訊息 (短 hash)`
4. **輸出到 GITHUB_OUTPUT**：使用 heredoc 語法將多行的 changelog 存入 output 變數

產生的 changelog 範例：

```markdown
## What's Changed
- feat: add /users endpoint (a1b2c3d)
- fix: handle nil pointer in health check (e4f5g6h)
- docs: update API documentation (i7j8k9l)
```

### GitHub 的 generate_release_notes 功能

```yaml
- name: Create GitHub Release
  uses: softprops/action-gh-release@v2
  with:
    body: |
      ## What's Changed
      ${{ steps.changelog.outputs.CHANGELOG }}
    files: artifacts/*
    generate_release_notes: true
```

`generate_release_notes: true` 會讓 GitHub **自動產生** release notes，包含：

- 新的貢獻者
- PR 列表（標題 + 連結）
- 與前一個 release 的 diff 連結

它會附加在你自訂的 `body` 之後，兩者結合提供完整的 release 資訊。

### 進階：使用 Conventional Commits 自動分類

如果你的團隊使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式（例如 `feat:`, `fix:`, `docs:`），可以進一步自動分類 changelog：

```
## Features
- feat: add /users endpoint (a1b2c3d)
- feat: support pagination (m1n2o3p)

## Bug Fixes
- fix: handle nil pointer in health check (e4f5g6h)
- fix: correct response status code (q4r5s6t)

## Documentation
- docs: update API documentation (i7j8k9l)
```

這可以透過工具如 [git-cliff](https://github.com/orhun/git-cliff) 或 GoReleaser 的 changelog 功能來實現。

> 💡 **講師提示：** 可以連結回前面課程提到的 Conventional Commits 概念。好的 commit message 不只是讓 git log 好看，還能自動化產生有意義的 changelog。

---

## Docker Image 建置與推送

由於學生們已經學過 Docker，我們可以在 release workflow 中加入 Docker image 的建置與推送。

### 推送到 GitHub Container Registry (GHCR)

**GHCR** 是 GitHub 內建的 container registry，與 GitHub 帳號整合，免費額度對公開 repository 非常友善。

在 release workflow 中新增一個 `docker` job：

```yaml
docker:
  name: Build Docker Image
  needs: test
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - name: Login to GHCR
      uses: docker/login-action@v3
      with:
        registry: ghcr.io
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    - name: Build and push
      uses: docker/build-push-action@v6
      with:
        context: .
        push: true
        tags: |
          ghcr.io/${{ github.repository }}:${{ github.ref_name }}
          ghcr.io/${{ github.repository }}:latest
```

### 逐段解說

#### 登入 GHCR

```yaml
- name: Login to GHCR
  uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}
```

- `ghcr.io` 是 GitHub Container Registry 的域名
- `github.actor` 是觸發 workflow 的使用者
- `secrets.GITHUB_TOKEN` 是 GitHub 自動提供的 token，**不需要額外設定**

#### 建置與推送

```yaml
- name: Build and push
  uses: docker/build-push-action@v6
  with:
    context: .
    push: true
    tags: |
      ghcr.io/${{ github.repository }}:${{ github.ref_name }}
      ghcr.io/${{ github.repository }}:latest
```

每個 image 會打上兩個 tag：

| Tag | 範例 | 用途 |
|-----|------|------|
| `${{ github.ref_name }}` | `ghcr.io/user/myapp:v1.0.0` | 指向特定版本 |
| `latest` | `ghcr.io/user/myapp:latest` | 永遠指向最新版本 |

### 完整的 Release 流程

加入 Docker 建置後，完整的 release 流程如下：

```
                  git push origin v1.0.0
                          │
                          ▼
                  ┌───────────────┐
                  │   Test Job    │
                  └───────┬───────┘
                          │
                ┌─────────┴─────────┐
                │                   │
                ▼                   ▼
      ┌──────────────────┐  ┌──────────────┐
      │   Build Job      │  │  Docker Job  │
      │   (5 binaries)   │  │  (image)     │
      └────────┬─────────┘  └──────┬───────┘
               │                   │
               │                   ▼
               │           ┌──────────────┐
               │           │    GHCR      │
               │           │  push image  │
               │           └──────────────┘
               ▼
      ┌──────────────────┐
      │   Release Job    │
      │   GitHub Release │
      │   + binaries     │
      └──────────────────┘
```

> 💡 **講師提示：** 可以提醒學生，因為 Docker 課程已經學過 Dockerfile 的撰寫，這裡的重點是「如何在 CI 中自動建置並推送 image」，而不是 Dockerfile 本身。

---

## 進階：使用 GoReleaser

### GoReleaser 是什麼？

[GoReleaser](https://goreleaser.com/) 是一個專門為 Go 專案設計的 release 自動化工具。它可以用一個指令完成：

- Cross compilation（多平台建置）
- Changelog 產生
- GitHub/GitLab Release 建立
- Docker image 建置與推送
- Homebrew formula 發佈
- 以及更多...

簡單來說，GoReleaser 把我們前面手動設定的所有步驟 **整合成一個工具**。

### `.goreleaser.yml` 設定範例

在專案根目錄建立 `.goreleaser.yml`：

```yaml
# GoReleaser configuration
version: 2

builds:
  - id: myapp
    main: .
    binary: myapp
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64
    ldflags:
      - -s -w
      - -X main.version={{.Version}}

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
      - '^ci:'

dockers:
  - image_templates:
      - "ghcr.io/{{ .Env.GITHUB_REPOSITORY }}:{{ .Tag }}"
      - "ghcr.io/{{ .Env.GITHUB_REPOSITORY }}:latest"
```

### 使用 goreleaser-action

在 GitHub Actions 中使用 GoReleaser：

```yaml
name: Release with GoReleaser

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write
  packages: write

jobs:
  goreleaser:
    name: GoReleaser
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

只需要 **一個 step** 就完成了所有 release 工作——cross compilation、changelog、GitHub Release、Docker image，全部由 GoReleaser 處理。

### GoReleaser vs 手動設定

| 面向 | 手動設定 | GoReleaser |
|------|---------|-----------|
| **設定複雜度** | 需要自己寫每個步驟 | 統一用 `.goreleaser.yml` 管理 |
| **學習曲線** | 了解每個 Action 的用法 | 學習 GoReleaser 的設定格式 |
| **靈活性** | 完全自訂 | 受限於 GoReleaser 支援的功能 |
| **維護成本** | 需要維護多個 workflow 步驟 | 只需維護一個設定檔 |
| **推薦場景** | 學習、簡單專案 | 正式專案、需要多功能 release |

> 💡 **講師提示：** 建議學生先用手動設定的方式學習，了解每個步驟的原理後，再進階使用 GoReleaser。這樣遇到問題時才知道要從哪裡 debug。

---

## 小結與練習題

### 本章重點回顧

- **Semantic Versioning**：`vMAJOR.MINOR.PATCH`，MAJOR 是破壞性變更，MINOR 是新功能，PATCH 是 bug 修正
- **Git Tag**：用來標記版本，推薦使用 annotated tag，需要明確推送到 remote
- **Tag 觸發 workflow**：使用 `on: push: tags: ['v*']` 在 tag push 時自動執行 release 流程
- **ldflags**：在編譯時注入版本資訊，讓程式知道自己的版本號
- **Cross Compilation**：Go 內建跨平台編譯，搭配 matrix strategy 可以同時建置多個平台
- **Changelog**：可以從 git log 自動產生，搭配 `generate_release_notes` 提供完整的 release notes
- **Docker image**：在 release 時同步建置並推送到 GHCR
- **GoReleaser**：一站式的 Go release 工具，適合正式專案使用

### 練習題

**練習 1：建立你的第一個 Release**

1. 將前面章節建立的 workflow 全部 push 到 GitHub
2. 把本章的 `release.yml` 加入 `.github/workflows/`
3. 修改 `main.go`，加入 `var version = "dev"` 變數，並在 `VersionHandler` 中使用它
4. Commit 並 push 到 `main`
5. 建立一個 annotated tag：`git tag -a v1.0.0 -m "First release"`
6. 推送 tag：`git push origin v1.0.0`
7. 到 GitHub Actions 頁面觀察 Release workflow 的執行過程
8. 到 Releases 頁面查看自動建立的 Release 和上傳的 binary

**練習 2：修改版號並發佈新版本**

1. 在 `handler.go` 中新增一個 `/ping` endpoint，回傳 `{"pong": true}`
2. 新增對應的測試
3. Commit 並 push
4. 思考：這個變更應該 bump 哪個版號？（答案：MINOR，因為是新增功能）
5. 建立 tag `v1.1.0` 並推送
6. 觀察新的 Release 中的 changelog 是否包含了新增 `/ping` endpoint 的 commit

**練習 3（進階）：加入 Docker Image 建置**

將本章的 Docker build & push job 加入你的 release workflow 中。push tag 後，到 GitHub 的 Packages 頁面確認 Docker image 是否成功推送。

> **接下來，我們將學習如何將應用程式部署到雲端平台，讓你的程式真正上線運行！**

---

[← 上一章：PR 檢查自動化](05-pr-checks.md) ｜ [下一章：部署到雲端平台 →](07-deployment.md)
