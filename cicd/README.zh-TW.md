# GitHub Actions 實戰工作坊

## 課程簡介

兩小時的 GitHub Actions 實戰課程，從 CI/CD 基礎概念出發，動手寫第一個 workflow，再為一個真正的 Go 專案建立完整的 CI pipeline。完成後你會具備在實際專案中導入 CI/CD 的基本能力。

## 適用對象

- 已具備 **GitHub** 使用經驗（熟悉 clone、push、pull、branch 等基本操作）
- 已具備 **Docker** 基礎知識（已上過 Docker 課程）
- 對自動化流程有興趣，想了解如何讓程式碼的測試與部署更有效率

## 課前準備

- [ ] **GitHub 帳號** — 如果還沒有，請到 [github.com](https://github.com) 註冊
- [ ] **安裝 Git** — 確認終端機可以執行 `git --version`
- [ ] **安裝 Go 1.24+** — 確認終端機可以執行 `go version`
- [ ] **安裝 Docker** — 確認終端機可以執行 `docker --version`
- [ ] **程式編輯器** — 推薦使用 [VS Code](https://code.visualstudio.com/)，並安裝 YAML 擴充套件

## 課程大綱與時間分配

| 時段 | 章節 | 內容 | 時間 |
|------|------|------|------|
| 0:00–0:05 | — | 開場、確認事前準備 | 5 分鐘 |
| 0:05–0:20 | **01** | CI/CD 概念介紹 | 15 分鐘 |
| 0:20–0:55 | **02** | GitHub Actions 基礎（hello.yml + 介面 + 結構元素） | 35 分鐘 |
| 0:55–1:10 | **02** | 手動觸發 + Pipeline 失敗 debug | 15 分鐘 |
| 1:10–1:25 | **02 練習** | 自訂 hello.yml | 15 分鐘 |
| 1:25–2:00 | **03** | Go 專案 CI Pipeline 講解 | 35 分鐘 |

> **總計：2 小時（120 分鐘）**
>
> 03 練習與 04 部署章節為課後延伸內容，建議學生回家自行完成。

### 各章節內容

**01 — CI/CD 概念介紹**
為什麼需要 CI/CD、它解決什麼問題、CI 和 CD 各自在做什麼、SDC 實際怎麼用。

**02 — GitHub Actions 基礎**
動手寫第一個 `hello.yml`、push 後觀察 GitHub Actions 介面、認識 workflow 的五個結構元素（Workflow / Event / Job / Step / Runner）、`uses` 引用 Action、`${{ }}` 讀 context、手動觸發、刻意製造失敗練習 debug。

**03 — Go 專案 CI Pipeline**
為一個 Go HTTP API 伺服器建立完整的 CI workflow：lint（golangci-lint）、test（含覆蓋率）、build，學會用 `needs` 控制 job 依賴關係，以及 `actions/upload-artifact` 在 job 之間傳遞檔案。最後介紹 PR 觸發的 CI 與相關的 permissions、fork 安全限制。

**04 — 部署到雲端平台（延伸）**
GitHub Environments 與 Secrets 的概念，以及實戰部署到 Fly.io 的範例。本章為課後延伸閱讀。

## 教材結構

```
cicd/
├── README.md                          # 教材總覽（英文）
├── README.zh-TW.md                    # 本檔案
├── 01-cicd-intro.md                   # CI/CD 概念介紹
├── 02-github-actions-basics.md        # GitHub Actions 基礎（含實作）
├── 03-go-ci-pipeline.md               # Go 專案 CI Pipeline
├── 04-deployment.md                   # 部署到雲端平台（延伸）
├── assets/                            # 截圖與圖表
├── examples/
│   ├── .github/workflows/
│   │   ├── hello.yml                  # 02 章 demo workflow
│   │   └── repo-info.yml              # 02 章失敗範例
│   └── sample-app/                    # Go HTTP API 範例專案
│       ├── main.go
│       ├── handler.go
│       ├── handler_test.go
│       ├── go.mod
│       ├── Dockerfile
│       ├── .golangci.yml
│       └── .github/workflows/
│           ├── ci.yml                 # 03 章 CI pipeline
│           └── cd.yml                 # 04 章部署到 Fly.io
└── exercises/
    ├── 01-basics.md                   # 02 章練習
    └── 02-ci-pipeline.md              # 03 章練習
```

## 範例專案說明

`examples/sample-app` 是一個簡單的 **Go HTTP API 伺服器**，提供基本的 RESTful endpoint。這個專案會在 03、04 章作為示範與練習的基礎，你將學會如何為它建立完整的 CI/CD pipeline，從自動化測試、程式碼風格檢查、Docker 映像檔建置，到自動部署到雲端平台。
