# SRE Workshop

涵蓋 SRE（Site Reliability Engineering）核心技能的實戰工作坊系列——Docker 容器化與 GitHub Actions CI/CD 自動化。

## 工作坊列表

| 工作坊 | 時長 | 說明 |
|--------|------|------|
| [Docker](docker/) | 2 小時 | Docker 基礎、Dockerfile 撰寫、Docker Compose 編排 |
| [CI/CD](cicd/) | 2 小時 | GitHub Actions、CI pipeline、Release 與部署自動化 |

> **建議上課順序**：先上 Docker，再上 CI/CD——CI/CD 課程預設學員已具備 Docker 基礎知識。

## 課前準備

- [ ] **Git** — `git --version`
- [ ] **Go 1.22+** — `go version`
- [ ] **Docker** — `docker --version`
- [ ] **GitHub 帳號** — [github.com](https://github.com)
- [ ] **程式編輯器** — 推薦使用 [VS Code](https://code.visualstudio.com/)（搭配 YAML 擴充套件）

## 專案結構

```
sre-workshop/
├── README.md              # 英文版 README
├── README.zh-TW.md        # 本檔案
├── docker/                # Docker 工作坊
│   ├── README.md          # 工作坊總覽（英文）
│   ├── README.zh-TW.md    # 工作坊總覽（中文）
│   └── docker-workshop.md # 完整教材
└── cicd/                  # CI/CD 工作坊
    ├── README.md          # 工作坊總覽
    ├── README.zh-TW.md    # 工作坊總覽（中文）
    ├── 01-cicd-intro.md   # CI/CD 概念介紹
    ├── 02-github-actions-basics.md
    ├── 03-first-workflow.md
    ├── 04-go-ci-pipeline.md
    ├── 05-pr-checks.md
    ├── 06-release-automation.md
    ├── 07-deployment.md
    ├── examples/          # 範例 Go 應用程式
    └── exercises/         # 實作練習
```

## 授權

本教材僅供 SDC 工作坊教學使用。
