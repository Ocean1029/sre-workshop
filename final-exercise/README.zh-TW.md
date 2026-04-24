# 最後練習：CI/CD + 監控整條串起來

## 課程簡介

25 分鐘的綜合練習，把 Docker、CI/CD、Prometheus 三項技能串在一起。你會透過 GitHub Actions 把服務部署到社辦機器，把 metric 吐給 Prometheus，並在服務壞掉時收到 Discord 告警。

> **練習本體在另一個 repo：** [Ocean1029/sre-workshop-capstone](https://github.com/Ocean1029/sre-workshop-capstone)。
>
> 那邊才是你要 clone、改、push 的 repo。這個資料夾只是給講師和看主 repo 的人一份說明。

**capstone repo 已經包含：**

- 服務的 source code 和 Dockerfile
- 大部分的 GitHub Actions workflow（CI 完整，CD 還缺一塊）
- Prometheus scrape 設定和 alert rule
- Part 1 本地環境用的 `docker-compose.yml`

**社辦機器上已經準備好：**

- GitHub self-hosted runner
- Prometheus 和 Alertmanager（Alertmanager 已經設定好對應的 Discord 頻道）

## 適用對象

- 已經上完當天的 Docker、CI/CD、Prometheus 三個工作坊
- 準備好把三個部分合在一起實際跑一次

## 練習流程

### Part 0 — 介紹練習（5 分鐘）

目標架構：repo → GitHub Actions → self-hosted runner → 在社辦機器上 `docker run` → Prometheus scrape → Alertmanager → Discord。

### Part 1 — 在本地跑起來（10 分鐘）

1. `docker compose up` 在本地 build 並啟動服務。
2. 確認 compose 裡面的 Prometheus 有成功 scrape 到服務。
3. 打開 Prometheus UI 找到服務的 metric。
4. 打 `/crash` endpoint，觀察 metric 次數變化。

### Part 2 — 補完 CD，觸發真正的告警（10 分鐘）

1. 參考 CI/CD 工作坊 ch04 的做法，把 CD job 缺的那一塊補齊（self-hosted runner + `docker run`）。
2. Push 上去，在 GitHub Actions 看到綠勾勾就代表部署成功。
3. 打已部署服務的 `/crash` endpoint。
4. 等告警觸發，到 Discord 頻道確認有收到通知。

## 開始練習

```bash
git clone https://github.com/Ocean1029/sre-workshop-capstone.git
cd sre-workshop-capstone
docker compose up --build
```

完整的步驟指引在 capstone repo 的 [README](https://github.com/Ocean1029/sre-workshop-capstone/blob/main/README.zh-TW.md#%E7%B7%B4%E7%BF%92%E6%B5%81%E7%A8%8B)。
