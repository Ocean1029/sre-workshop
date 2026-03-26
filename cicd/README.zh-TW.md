# GitHub Actions 實戰工作坊

## 課程簡介

**2 小時**的 GitHub Actions 實戰教學課程，從 CI/CD 基礎概念出發，一步步帶你認識 GitHub Actions 的核心觀念與實務操作。課程涵蓋 workflow 撰寫、Go 專案的 CI pipeline 建置、PR 自動化檢查、Release 自動化，以及部署到雲端平台等主題。透過動手實作，你將具備在真實專案中導入 CI/CD 的能力。

## 適用對象

- 已具備 **GitHub** 使用經驗（熟悉 clone、push、pull、branch 等基本操作）
- 已具備 **Docker** 基礎知識（已上過 Docker 課程）
- 對自動化流程有興趣，想了解如何讓程式碼的測試與部署更有效率

## 課前準備

- [ ] **GitHub 帳號** — 如果還沒有，請到 [github.com](https://github.com) 註冊
- [ ] **安裝 Git** — 確認終端機可以執行 `git --version`
- [ ] **安裝 Go 1.22+** — 確認終端機可以執行 `go version`，版本需為 1.22 以上
- [ ] **安裝 Docker** — 確認終端機可以執行 `docker --version`
- [ ] **程式編輯器** — 推薦使用 [VS Code](https://code.visualstudio.com/)，並安裝 YAML 擴充套件
- [ ] （選用）**安裝 [act](https://github.com/nektos/act)** — 可以在本地測試 GitHub Actions workflow，不必每次都 push 到 GitHub

## 課程大綱與時間分配

| 章節 | 主題 | 時間 |
|------|------|------|
| 01 | CI/CD 概念介紹 | 10 分鐘 |
| 02 | GitHub Actions 基礎 | 15 分鐘 |
| 03 | 動手做：第一個 Workflow | 15 分鐘 |
| 04 | Go 專案 CI Pipeline | 20 分鐘 |
| 05 | PR 檢查自動化 | 15 分鐘 |
| 06 | Release 自動化 | 20 分鐘 |
| 07 | 部署到雲端平台 | 15 分鐘 |
| -- | 總結與 Q&A | 10 分鐘 |

> **總計：2 小時（120 分鐘）**

### 01 — CI/CD 概念介紹（10 分鐘）

- 什麼是 CI/CD？
- 為什麼需要自動化？
- 常見 CI/CD 工具比較

### 02 — GitHub Actions 基礎（15 分鐘）

- Workflow / Job / Step / Action 核心概念
- YAML 語法與 workflow 檔案結構
- 觸發條件（Events）與執行環境（Runners）

### 03 — 動手做：第一個 Workflow（15 分鐘）

- 建立 `.github/workflows/` 目錄
- 撰寫 Hello World workflow
- Push 觸發並觀察執行結果
- **練習 1**：基礎 workflow 實作

### 04 — Go 專案 CI Pipeline（20 分鐘）

- Lint（程式碼風格檢查）
- Test（自動化測試）
- Build（編譯建置）
- Coverage（測試覆蓋率）
- **練習 2**：完整 CI pipeline 實作

### 05 — PR 檢查自動化（15 分鐘）

- Branch protection rules 設定
- Status checks 整合
- 自動標記與通知

### 06 — Release 自動化（20 分鐘）

- 語意化版本（Semantic Versioning）
- 自動建置 Docker image
- GitHub Release 發布
- **練習 3**：進階自動化實作

### 07 — 部署到雲端平台（15 分鐘）

- CD 流程設計
- 部署 workflow 撰寫
- 環境變數與 Secrets 管理

### 總結與 Q&A（10 分鐘）

- 課程重點回顧
- 延伸學習資源
- 問答時間

## 教材結構

```
cicd/
├── README.md                          # 教材總覽（英文）
├── README.zh-TW.md                    # 本檔案
├── 01-cicd-intro.md                   # CI/CD 概念介紹
├── 02-github-actions-basics.md        # GitHub Actions 基礎概念
├── 03-first-workflow.md               # 動手做：第一個 Workflow
├── 04-go-ci-pipeline.md               # Go 專案 CI Pipeline
├── 05-pr-checks.md                    # PR 檢查自動化
├── 06-release-automation.md           # Release 自動化
├── 07-deployment.md                   # 部署到雲端平台
├── examples/                          # 範例程式碼
│   └── sample-app/                    # 範例 Go 應用程式
│       ├── main.go                    # 主程式進入點
│       ├── handler.go                 # HTTP handler
│       ├── handler_test.go            # 測試檔案
│       ├── go.mod                     # Go module 定義
│       └── Dockerfile                 # Docker 映像檔定義
└── exercises/                         # 練習題
    ├── exercise-01-basics.md          # 章節 01-03 練習
    ├── exercise-02-ci-pipeline.md     # 章節 04-05 練習
    └── exercise-03-advanced.md        # 章節 06-07 練習
```

## 範例專案說明

`examples/sample-app` 是一個簡單的 **Go HTTP API 伺服器**，提供基本的 RESTful API 端點。這個專案會在整個工作坊中作為示範與練習的基礎，你將學會如何為它建立完整的 CI/CD pipeline——從自動化測試、程式碼品質檢查、Docker 映像檔建置，到自動部署。
