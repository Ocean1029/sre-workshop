# Docker Workshop — 從零開始的容器化之旅

## 課程簡介

**2 小時**的實戰工作坊，帶你從 Docker 基礎概念一路學到 Docker Compose 多容器編排。所有範例皆使用 Go 作為應用程式語言。

## 適用對象

- 熟悉基本 Linux 命令列操作
- Docker 初學者——不需要任何容器經驗

## 課前準備

- [ ] **Docker Engine** 已安裝 — `docker --version`
- [ ] **Docker Compose** 已安裝 — `docker compose version`
- [ ] **Go 1.22+** 已安裝 — `go version`
- [ ] **程式編輯器** — 推薦使用 [VS Code](https://code.visualstudio.com/)

## 課程大綱

| 單元 | 主題 | 時長 |
|------|------|------|
| 1 | Docker 核心概念與基本操作 | 50 分鐘 |
| 2 | Dockerfile 深入實作 | 35 分鐘 |
| ☕ | 休息 | 5 分鐘 |
| 3 | Docker Compose 多容器編排 | 30 分鐘 |
| 4 | 綜合演練與延伸學習 | 15 分鐘 |

> **總計：2 小時（120 分鐘，含休息）**

### Part 1：Docker 核心概念與基本操作（50 分鐘）

- 什麼是 Docker？容器化概念與發展歷史
- 容器 vs 虛擬機器
- Docker 架構（Client / Daemon / Registry）
- Image 與 Container 的關係
- Image 操作：`pull`、`list`、`remove`
- Container 操作：`run`、`stop`、`rm`、`exec`
- Port Mapping 與 Volume 掛載
- Container 除錯技巧
- **練習 1**：執行你的第一個容器

### Part 2：Dockerfile 深入實作（35 分鐘）

- Dockerfile 語法與核心指令（`FROM`、`COPY`、`RUN`、`CMD`、`EXPOSE` 等）
- `.dockerignore` 最佳實踐
- 為 Go 應用撰寫 Dockerfile
- Multi-stage Build 最佳化映像檔
- 映像檔最佳實踐（Layer 快取、最小化 Base Image）
- **練習 2**：建置你的 Go 應用映像檔

### Part 3：Docker Compose 多容器編排（30 分鐘）

- 為什麼需要 Docker Compose？
- `docker-compose.yml` 基礎語法
- 服務定義、Network 與 Volume
- 環境變數管理
- 常用 Compose 指令
- 實戰示範：Go API + PostgreSQL + Redis
- **練習 3**：用 Compose 部署完整服務

### Part 4：綜合演練（15 分鐘）

- 綜合練習
- 常見問題排查
- 延伸學習資源

## 教材說明

完整工作坊內容在 [docker-workshop.md](docker-workshop.md)，包含詳細解說、程式碼範例、架構圖，以及附錄（Docker / Dockerfile / Compose 指令速查表）。
