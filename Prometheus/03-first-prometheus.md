# 03 — 動手做：部署 Prometheus

> **時間：25 分鐘**

---

## 目錄

- [03 — 動手做：部署 Prometheus](#03--動手做部署-prometheus)
  - [目錄](#目錄)
  - [事前準備](#事前準備)
    - [確認工具已安裝](#確認工具已安裝)
  - [](#)
  - [Step 1：/examples 檔案結構](#step-1examples-檔案結構)
  - [Step 2：建立 Prometheus 設定檔](#step-2建立-prometheus-設定檔)
    - [為什麼要監測 Prometheus 自己？](#為什麼要監測-prometheus-自己)
  - [Step 3：建立 Docker Compose 檔案](#step-3建立-docker-compose-檔案)
    - [逐行解說](#逐行解說)
  - [Step 4：啟動 Prometheus](#step-4啟動-prometheus)
  - [Step 5：探索 Prometheus Web UI](#step-5探索-prometheus-web-ui)
    - [Targets 頁面](#targets-頁面)
    - [Graph 頁面](#graph-頁面)
    - [直接查看 /metrics endpoint](#直接查看-metrics-endpoint)
  - [Step 6：加入 Node Exporter](#step-6加入-node-exporter)
    - [更新 Docker Compose](#更新-docker-compose)
    - [更新 Prometheus 設定](#更新-prometheus-設定)
    - [重啟服務](#重啟服務)
    - [驗證 Node Exporter](#驗證-node-exporter)
  - [Step 7：用 PromQL 查詢 Node Exporter 的 Metrics](#step-7用-promql-查詢-node-exporter-的-metrics)
    - [查看所有 targets 是否正常](#查看所有-targets-是否正常)
    - [記憶體使用率（百分比）](#記憶體使用率百分比)
    - [CPU 使用率（百分比）](#cpu-使用率百分比)
    - [硬碟使用率（百分比）](#硬碟使用率百分比)
    - [系統開機時間](#系統開機時間)
  - [設定檔逐行解說](#設定檔逐行解說)
    - [prometheus.yml 完整解說](#prometheusyml-完整解說)
    - [關鍵概念](#關鍵概念)
  - [常見問題排解](#常見問題排解)
    - [1. Port 被占用](#1-port-被占用)
    - [2. Prometheus 設定檔語法錯誤](#2-prometheus-設定檔語法錯誤)
    - [3. Node Exporter 在 macOS 上 metrics 很少](#3-node-exporter-在-macos-上-metrics-很少)
    - [4. Target 顯示 DOWN](#4-target-顯示-down)
  - [小結與練習題](#小結與練習題)
    - [本章重點回顧](#本章重點回顧)
    - [練習題](#練習題)


---

## 事前準備

### 確認工具已安裝

在 Docker workshop 中你已經安裝了 Docker 和 Docker Compose。請確認它們都能正常使用：

```bash
docker --version          # Docker version 27.x.x 或以上
docker compose version    # Docker Compose version v2.x.x 或以上
```

下載 [Prometheus Formatter Extension](https://chromewebstore.google.com/detail/jhfbpphccndhifmpfbnpobpclhedckbb?utm_source=item-share-cb)

![alt](./assets/metricspage.png)
![alt text](./assets/formatter.png)
---

## Step 1：/examples 檔案結構
``


```
Prometheus/examples/monitoring-stack/
├── docker-compose.yml          
└── config/
    └── prometheus.yml          
```

---

## Step 2：建立 Prometheus 設定檔

Prometheus 的核心設定檔是一個 YAML 檔案，告訴 Prometheus：

1. **多久** scrape 一次 metrics
2. **去哪裡** scrape（也就是 scrape targets）

建立 `Prometheus/examples/monitoring-stack/config/prometheus.yml`：

```yaml
# Prometheus 主設定檔
global:
  scrape_interval: 15s      # 每 15 秒去 scrape 一次（預設是 1 分鐘）
  evaluation_interval: 15s  # 每 15 秒評估一次 alert rules

# Scrape 設定 — 定義要監測的 targets
scrape_configs:
  # Job 1：監測 Prometheus 自己
  - job_name: "prometheus"
    static_configs:
      - targets: ["localhost:9090"]
```

### 為什麼要監測 Prometheus 自己？

這叫做 **self-monitoring**。Prometheus 本身也是一個服務，它也可能出問題（記憶體不足、TSDB 損壞等）。讓 Prometheus scrape 自己的 `/metrics` endpoint，你可以知道：

- Prometheus 吃了多少記憶體
- TSDB 儲存了多少 time series
- scrape 有沒有失敗

> 💡 **講師提示：** 這裡可以順便介紹 `scrape_interval` 的意義。15 秒是教學用途，正式環境通常用 30 秒到 1 分鐘。間隔越短，資料越精細，但儲存和 CPU 消耗也越大。

---

## Step 3：建立 Docker Compose 檔案

建立 `Prometheus/examples/monitoring-stack/docker-compose.yml`：

```yaml
services:
  prometheus:
    image: prom/prometheus:v3.2.1
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./config/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus-data:/prometheus
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"
      - "--storage.tsdb.path=/prometheus"
      - "--storage.tsdb.retention.time=7d"
      - "--web.console.libraries=/etc/prometheus/console_libraries"
      - "--web.console.templates=/etc/prometheus/consoles"
    restart: unless-stopped

volumes:
  prometheus-data:
```

### 逐行解說

| 設定 | 說明 |
|------|------|
| `image: prom/prometheus:v3.2.1` | 使用 Prometheus 官方 Docker image，指定版本避免意外升級 |
| `ports: "9090:9090"` | 把 container 內的 9090 port 映射到主機的 9090 port |
| `volumes: ./config/prometheus.yml:...` | 把我們的設定檔掛載到 container 內，`:ro` 表示唯讀 |
| `volumes: prometheus-data:/prometheus` | 用 Docker volume 持久化 TSDB 資料，重啟不會遺失 |
| `--storage.tsdb.retention.time=7d` | 資料保留 7 天，超過自動清理 |
| `restart: unless-stopped` | 除非手動停止，否則 container 異常退出時自動重啟 |

> 💡 **講師提示：** 提醒學生在 Docker workshop 中學過的 volumes、ports 等概念。這裡是實際應用這些知識的好機會。

---

## Step 4：啟動 Prometheus

在 `Prometheus/examples/monitoring-stack/` 目錄下執行：

```bash
cd Prometheus/examples/monitoring-stack
docker compose up -d
```

確認 container 正常運行：

```bash
docker compose ps
```

你應該看到類似以下的輸出：

```
NAME         IMAGE                    STATUS         PORTS
prometheus   prom/prometheus:v3.2.1   Up 10 seconds  0.0.0.0:9090->9090/tcp
```

確認 Prometheus 有在回應：

```bash
curl -s http://localhost:9090/-/healthy
# 預期輸出：Prometheus Server is Healthy.
```

> 💡 **講師提示：** 如果有學生的 container 沒有起來，請他們跑 `docker compose logs prometheus` 查看錯誤訊息。最常見的問題是 YAML 縮排錯誤或 port 被占用。

---

## Step 5：探索 Prometheus Web UI

打開瀏覽器，進入 http://localhost:9090

### Targets 頁面

點擊上方導覽列的 **Status → Targets**：

- 你應該看到一個名為 **prometheus** 的 target
- 狀態應該是 **UP**（綠色）
- `Last Scrape` 欄位顯示上次 scrape 的時間
- `Scrape Duration` 顯示 scrape 花了多久

如果狀態是 **DOWN**（紅色），表示 scrape 失敗了——通常是設定檔的 target 地址寫錯。

### Graph 頁面

回到首頁（Graph 頁面），在查詢框中輸入以下 PromQL 查詢：

**查詢 1：所有 targets 的狀態**

```promql
up
```

你應該看到一筆結果：`up{instance="localhost:9090", job="prometheus"} = 1`。

`1` 代表正常，`0` 代表掛了。

**查詢 2：Prometheus 的記憶體使用量**

```promql
process_resident_memory_bytes{job="prometheus"}
```

點擊 **Graph** 分頁可以看到記憶體使用量隨時間變化的曲線圖。

**查詢 3：Prometheus 已經收錄了多少 time series**

```promql
prometheus_tsdb_head_series
```

這個數字代表 TSDB 中目前有多少條活躍的 time series。

### 直接查看 /metrics endpoint

你也可以在瀏覽器中直接打開 http://localhost:9090/metrics，看到 Prometheus 暴露的所有 metrics 原始資料。這就是 Prometheus 對自己做 scrape 時看到的內容。

> 💡 **講師提示：** 讓學生花 2-3 分鐘自由探索 Web UI。可以引導他們試試看在 Graph 頁面輸入 `prometheus_` 然後看自動完成提示，感受一下 Prometheus 自己產生了哪些 metrics。

---

## Step 6：加入 Node Exporter

目前我們只有 Prometheus 自己的 metrics。接下來加入 **Node Exporter**，讓我們能夠監測主機的硬體指標（CPU、記憶體、硬碟等）。

### 更新 Docker Compose

編輯 `docker-compose.yml`，在 `services` 區塊中加入 Node Exporter：

```yaml
services:
  prometheus:
    image: prom/prometheus:v3.2.1
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./config/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus-data:/prometheus
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"
      - "--storage.tsdb.path=/prometheus"
      - "--storage.tsdb.retention.time=7d"
      - "--web.console.libraries=/etc/prometheus/console_libraries"
      - "--web.console.templates=/etc/prometheus/consoles"
    restart: unless-stopped

  node-exporter:
    image: prom/node-exporter:v1.9.0
    container_name: node-exporter
    ports:
      - "9100:9100"
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs:ro
    command:
      - "--path.procfs=/host/proc"
      - "--path.sysfs=/host/sys"
      - "--path.rootfs=/rootfs"
      - "--collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)"
    restart: unless-stopped

volumes:
  prometheus-data:
```

> **macOS 使用者注意**：在 macOS 上 `/proc` 和 `/sys` 不存在，Node Exporter 的某些 collector 會無法使用。你可以改用以下精簡版，只掛載可用的路徑：
>
> ```yaml
>   node-exporter:
>     image: prom/node-exporter:v1.9.0
>     container_name: node-exporter
>     ports:
>       - "9100:9100"
>     restart: unless-stopped
> ```

### 更新 Prometheus 設定

編輯 `config/prometheus.yml`，加入 Node Exporter 的 scrape job：

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: "prometheus"
    static_configs:
      - targets: ["localhost:9090"]

  # Job 2：監測主機的硬體指標
  - job_name: "node-exporter"
    static_configs:
      - targets: ["node-exporter:9100"]
```

注意 target 地址用的是 `node-exporter:9100` 而不是 `localhost:9100`——因為在 Docker 網路中，container 之間用 **service name** 來通訊。

### 重啟服務

```bash
docker compose up -d
```

Prometheus 會自動偵測到設定檔的變更。你也可以用以下指令讓 Prometheus 重新讀取設定：

```bash
# 方式 1：重啟 container
docker compose restart prometheus

# 方式 2：發送 SIGHUP 信號（不需要重啟）
docker compose kill -s SIGHUP prometheus
```

### 驗證 Node Exporter

1. 打開 http://localhost:9100/metrics —— 你會看到大量的硬體 metrics
2. 回到 Prometheus 的 Targets 頁面 (http://localhost:9090/targets) —— 你應該看到兩個 target 都是 **UP**

> 💡 **講師提示：** 讓學生打開 Node Exporter 的 `/metrics` 頁面，看看原始的 metrics 長什麼樣子。搜尋 `node_memory` 或 `node_cpu` 來感受 metrics 的命名規則。

---

## Step 7：用 PromQL 查詢 Node Exporter 的 Metrics

回到 Prometheus Web UI (http://localhost:9090)，試試以下查詢：

### 查看所有 targets 是否正常

```promql
up
```

現在應該會看到兩筆結果——一筆是 `prometheus`，一筆是 `node-exporter`。

### 記憶體使用率（百分比）

```promql
(1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * 100
```

切換到 **Graph** 分頁可以看到記憶體使用率的趨勢圖。

### CPU 使用率（百分比）

```promql
(1 - avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m]))) * 100
```

這個查詢做了三件事：
1. `node_cpu_seconds_total{mode="idle"}` — 取得 CPU 閒置時間
2. `rate(...[5m])` — 計算過去 5 分鐘的每秒變化率
3. `1 - avg by(instance)(...)` — 用 1 減去閒置率 = CPU 使用率

### 硬碟使用率（百分比）

```promql
(1 - node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100
```

### 系統開機時間

```promql
node_boot_time_seconds
```

這是一個 Gauge metric，記錄了系統開機的 Unix timestamp。

> 💡 **講師提示：** 在這個步驟花足夠的時間讓學生嘗試不同的查詢。鼓勵他們在查詢框裡輸入 `node_` 然後看自動完成的建議，自己探索有哪些 metrics 可以查。

---

## 設定檔逐行解說

### prometheus.yml 完整解說

```yaml
# === 全域設定 ===
global:
  scrape_interval: 15s       # 所有 scrape jobs 的預設 scrape 間隔
                              # 可以在個別 job 中覆蓋這個設定
  evaluation_interval: 15s   # 評估 alert rules 的間隔
                              # 每 15 秒檢查一次有沒有 alert 需要觸發

# === Scrape 設定 ===
# 定義 Prometheus 要去哪裡蒐集 metrics
scrape_configs:
  # 每個 job_name 代表一組邏輯上相關的 targets
  - job_name: "prometheus"         # Job 的名稱（會自動加成 label）
    static_configs:                # 靜態設定 targets（vs. service discovery）
      - targets: ["localhost:9090"] # target 的地址列表

  - job_name: "node-exporter"
    static_configs:
      - targets: ["node-exporter:9100"]
    # 你也可以在這裡覆蓋 global 的 scrape_interval：
    # scrape_interval: 30s
```

### 關鍵概念

| 概念 | 說明 |
|------|------|
| **Job** | 一組目的相同的 targets。所有 targets 會自動被加上 `job="<job_name>"` label |
| **Target** | 一個 Prometheus 要 scrape 的 endpoint（`host:port`） |
| **static_configs** | 手動列出 targets。另一種方式是 service discovery（自動發現） |
| **scrape_interval** | 多久去 scrape 一次。global 的是預設值，可以在 job 層級覆蓋 |

---

## 常見問題排解

### 1. Port 被占用

```
Error starting userland proxy: listen tcp4 0.0.0.0:9090: bind: address already in use
```

**解決方式**：找出誰在用這個 port，停掉它，或改用其他 port：

```bash
# 找出誰在用 9090
lsof -i :9090

# 或改用其他 port（例如 9091）
# 修改 docker-compose.yml 的 ports 為 "9091:9090"
```

### 2. Prometheus 設定檔語法錯誤

**症狀**：Prometheus container 啟動後立刻退出。

```bash
# 查看錯誤訊息
docker compose logs prometheus
```

**常見錯誤**：YAML 縮排不正確、key 名稱拼錯。

**驗證設定檔語法**：

```bash
# 用 Prometheus 的內建工具檢查設定檔
docker run --rm -v $(pwd)/config:/etc/prometheus prom/prometheus:v3.2.1 \
  promtool check config /etc/prometheus/prometheus.yml
```

### 3. Node Exporter 在 macOS 上 metrics 很少

這是正常的。macOS 不是 Linux，所以 `/proc` 和 `/sys` 不存在，很多 Linux-specific 的 collector 無法使用。你仍然可以看到一些基本 metrics（如 Go runtime metrics）。完整的 Node Exporter 功能需要在 Linux 環境中才能體驗。

### 4. Target 顯示 DOWN

**排查步驟**：

```bash
# 1. 確認 container 有在跑
docker compose ps

# 2. 確認 endpoint 可以連到
docker compose exec prometheus wget -qO- http://node-exporter:9100/metrics | head

# 3. 確認 prometheus.yml 中的 target 地址正確
```

> 💡 **講師提示：** 如果多數學生遇到問題，建議統一排查。最常見的問題是 YAML 縮排錯誤和 Docker 網路問題（container name 寫錯）。

---

## 小結與練習題

### 本章重點回顧

- Prometheus 的核心設定檔是 `prometheus.yml`，定義了 **scrape_interval** 和 **scrape_configs**
- 用 Docker Compose 可以快速部署 Prometheus，設定檔用 volume 掛載進 container
- Prometheus Web UI 的 **Targets** 頁面可以查看所有 scrape target 的狀態
- **Node Exporter** 負責將主機的硬體指標（CPU、記憶體、硬碟等）轉換成 Prometheus 格式的 metrics
- 在 Docker 網路中，container 之間用 **service name** 通訊（不是 `localhost`）

### 練習題

**練習 1：修改 scrape interval**

把 `prometheus.yml` 的 `scrape_interval` 改成 `5s`，重啟 Prometheus，觀察 Targets 頁面的 `Last Scrape` 間隔是否變成 5 秒。改完後記得改回 `15s`。

**練習 2：探索 Prometheus 自己的 metrics**

在 Graph 頁面查詢以下 metrics 並觀察它們的變化：

1. `prometheus_tsdb_head_series` — TSDB 中有多少 time series？
2. `prometheus_tsdb_head_samples_appended_total` — 總共收錄了多少筆 sample？（這是一個 Counter，試試用 `rate()` 看看每秒收錄多少筆）
3. `prometheus_target_scrape_pool_targets` — 每個 scrape pool 有多少 targets？

**練習 3：查看 Node Exporter 的 raw metrics**

打開 http://localhost:9100/metrics，找到以下 metrics 並理解它們的含義：

1. 找到一個 **Counter** 類型的 metric（提示：看 `# TYPE` 那行）
2. 找到一個 **Gauge** 類型的 metric
3. 找到 `node_uname_info` metric，它的 labels 裡有什麼資訊？

> **接下來，我們要設定 Alert Rules 和 Alertmanager，讓系統在出問題時自動通知你！**

---

[← 上一章：Prometheus 核心概念](02-prometheus-fundamentals.md) ｜ [下一章：Alertmanager 告警系統 →](04-alerting.md)
