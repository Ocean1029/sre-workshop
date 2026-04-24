# 02 — Docker Compose 多容器編排

## 2.1 為什麼需要 Docker Compose？

Ocean 現在會跑容器了，但每次要啟動服務都要手動打一串 `docker run`，參數又臭又長。更慘的是，上次好不容易把資料庫跟一個網頁 Admin 介面兩個容器兜起來了，結果 Snow 問他：「指令記在哪？」Ocean：「...我記在腦袋裡。」Snow：「喔不。」

**Docker Compose** 是 Docker 官方提供的多容器編排工具，透過一個 YAML 設定檔定義所有容器的狀態（映像檔、連接埠、Volume、環境變數、相依關係），再以單一指令 `docker compose up` 完成整個應用的部署。如果你使用 OrbStack，Docker Compose 已內建其中，並且可以透過 `docker compose version` 確認正在運行中的版本。

實際專案中，一個服務通常由多個容器組成。Docker Compose 將這些容器的設定與啟動流程集中定義在 `docker-compose.yml`（或 `compose.yml`）中，納入版本控制後，團隊成員只需 clone 專案並執行 `docker compose up`，即可在本機重現完整的服務環境。

下面用「PostgreSQL + Adminer」這個組合來示範：`postgres` 是資料庫，`adminer` 是一個輕量的網頁資料庫管理介面（Adminer 支援 Postgres、MySQL、SQLite 等多種資料庫，常被 SRE 拿來查看正在跑的 DB）。使用者打開瀏覽器連 Adminer，Adminer 再透過內部網路連到 Postgres。

```
┌──────────────────────────────────────────────────────┐
│                  PostgreSQL + Adminer                │
├──────────────────────────────────────────────────────┤
│                                                      │
│                   ┌──────────┐                       │
│                   │  Client  │                       │
│                   └────┬─────┘                       │
│                        │ http://localhost:8080       │
│                        ▼                             │
│              ┌──────────────────┐                    │
│              │     Adminer      │  ← Container 1     │
│              │  (DB Admin UI)   │                    │
│              │   (Port 8080)    │                    │
│              └────────┬─────────┘                    │
│                       │  db:5432                     │
│                       ▼                              │
│              ┌──────────────────┐                    │
│              │   PostgreSQL     │  ← Container 2     │
│              │   (Port 5432)    │                    │
│              └──────────────────┘                    │
│                                                      │
└──────────────────────────────────────────────────────┘

```

如果不用 Compose，就得手動打這些指令：

```bash
# 建立網路（後面會說明）
docker network create myapp

# 啟動 PostgreSQL
docker run -d --name db --network myapp \
  -e POSTGRES_PASSWORD=secret \
  -v pgdata:/var/lib/postgresql/data \
  postgres:17-alpine

# 啟動 Adminer
docker run -d --name adminer --network myapp \
  -p 8080:8080 \
  -e ADMINER_DEFAULT_SERVER=db \
  adminer:4
```

現在把這些 `docker run` 寫成一個 `docker-compose.yml`：

```yaml
services:
  db:
    image: postgres:17-alpine
    environment:
      POSTGRES_PASSWORD: secret
    volumes:
      - pgdata:/var/lib/postgresql/data

  adminer:
    image: adminer:4
    ports:
      - "8080:8080"
    environment:
      ADMINER_DEFAULT_SERVER: db
    depends_on:
      - db

volumes:
  pgdata:
```

一個檔案就把兩個容器的設定全部定義好了。把它存成 `docker-compose.yml`，執行 `docker compose up -d`，然後打開 `http://localhost:8080`，系統名稱選 PostgreSQL、使用者填 `postgres`、密碼填 `secret`，就能連到資料庫了。接著來看每個欄位在做什麼：

### services

最上層的 `services` 區塊定義了所有要跑的容器。每個 key（`db`、`adminer`）就是一個服務名稱，同時也是這個容器在 Compose 網路中的 hostname，其他容器可以直接用這個名稱來連線。

### image

```yaml
image: postgres:17-alpine
```

指定要用哪個映像檔。跟 `docker run postgres:17-alpine` 一樣的意思。

### ports

```yaml
ports:
  - "8080:8080"
```

Port Mapping，格式是 `"主機:容器"`。等同於 `docker run -p 8080:8080`。

### environment

```yaml
environment:
  POSTGRES_PASSWORD: secret
  ADMINER_DEFAULT_SERVER: db
```

設定環境變數。等同於 `docker run -e POSTGRES_PASSWORD=secret`。注意 `ADMINER_DEFAULT_SERVER` 的值是 `db`——這個 `db` 就是上面 services 裡的服務名稱，Adminer 會把它當 hostname 連過去，Compose 會自動幫忙做 DNS 解析（2.3 會詳細說明）。

### volumes

```yaml
volumes:
  - pgdata:/var/lib/postgresql/data
```

掛載 Volume。等同於 `docker run -v pgdata:/var/lib/postgresql/data`。這裡用的是 Named Volume，資料由 Docker 管理，容器砍掉資料還在。

檔案最下面的 `volumes:` 區塊是宣告這個 Named Volume：

```yaml
volumes:
  pgdata:
```

### depends_on

```yaml
depends_on:
  - db
```

告訴 Compose 啟動順序：先啟動 `db`，再啟動 `adminer`。不過要注意，`depends_on` 只等容器啟動，不等服務真正準備好。如果需要等資料庫可以接受連線，要搭配 `healthcheck`（後面會講）。

### 啟動與停止

```bash
# 啟動所有服務（背景執行）
docker compose up -d

# 看看服務狀態
docker compose ps

# 停止並移除所有容器
docker compose down
```

對比一下：原本要打兩串 `docker run` 加一個 `docker network create`，現在一個 `docker compose up -d` 就搞定了。

---

## 2.2 服務定義詳解

Ocean 興高采烈地把 2.1 的 compose 跑起來，Adminer log 卻滿滿都是 `connection refused: db:5432`。他重跑了三次，現象一模一樣。Snow 看了一眼：「DB 還在初始化啦，Adminer 容器跑起來不代表 Postgres 已經可以接連線。」

2.1 介紹了 Compose 最基本的欄位，這一節補充幾個實務上常用的進階設定。

### depends_on 搭配 healthcheck

在 2.1 裡我們用了簡單的 `depends_on`，但它只保證「容器有啟動」，不保證「服務可以連線」。實際上 PostgreSQL 容器啟動後還需要幾秒初始化，這段時間只要有後端連過去就會直接噴錯（正是 Ocean 遇到的情況）。

解法是搭配 `healthcheck` 和 `condition: service_healthy`：

```yaml
services:
  adminer:
    image: adminer:4
    depends_on:
      db:
        condition: service_healthy  # 等 db 的 healthcheck 通過才啟動

  db:
    image: postgres:17-alpine
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5
```

`pg_isready` 是 PostgreSQL 內建的檢查工具，Compose 每 5 秒問一次「資料庫好了沒」，連續失敗 5 次才判定不健康。`adminer` 會等到 db 的 healthcheck 通過才啟動。

### restart — 重啟策略

容器掛了怎麼辦？`restart` 設定決定 Docker 要不要自動把它拉起來：

| 選項 | 行為 |
|------|------|
| `no` | 不自動重啟（預設） |
| `always` | 無論什麼原因停止都會重啟 |
| `on-failure` | 只在非正常退出（exit code ≠ 0）時重啟 |
| `unless-stopped` | 除非你手動 `docker stop`，否則都會重啟 |

正式環境通常用 `unless-stopped` 或 `always`：

```yaml
services:
  adminer:
    restart: unless-stopped
```

---

## 2.3 Network

Ocean：「等等，2.1 的 Adminer 直接填 `db` 就能連到資料庫，我完全沒寫網路設定，怎麼會通？」Snow：「Compose 幫你做了。」

在 1.6 我們學到，容器的網路是隔離的，外面連不進去要靠 Port Mapping。但反過來，容器跟容器之間呢？預設也是不通的。

要讓多個容器能互相溝通，就需要把它們放進同一個 **Docker Network**（虛擬網路）。你可以把它想成一條虛擬的網路線，只有接在同一條線上的容器才能互相找到對方。

2.1 的 `docker run` 範例裡，我們手動跑了 `docker network create myapp`，然後每個容器都加上 `--network myapp` 才能互連。但 Compose 版本完全沒寫到網路，為什麼容器之間還是能溝通？

因為 **Docker Compose 會自動建立一個預設網路**，把 `docker-compose.yml` 裡定義的所有服務都放進去。你不用寫任何網路設定，容器之間就能互相連線。

更方便的是，在這個網路裡面，**服務名稱就等於 hostname**。回頭看 2.1 的範例：

```yaml
environment:
  ADMINER_DEFAULT_SERVER: db
```

這裡的值 `db` 不是什麼特殊語法，就是 `docker-compose.yml` 裡定義的服務名稱。Adminer 啟動後會用 `db` 當 hostname 去連資料庫，Compose 內建的 DNS 會自動把 `db` 解析到 PostgreSQL 容器的 IP。換個角度，如果有個 API 服務想連同一個資料庫，環境變數會寫成類似 `DATABASE_URL: postgres://postgres:secret@db:5432/mydb`——裡面的 `@db` 也是同樣的道理。服務改叫 `database`，這裡就要跟著改成 `database`。

> **注意**：這個 DNS 只在 Compose 的虛擬網路裡面有效。在主機上打 `ping db` 是不會通的。

---

## 2.4 Volume

1.7 的 Postgres 事件讓 Ocean 學乖了——容器砍掉，可寫層就沒了，所以需要持久化的資料一定要掛出來。在 Compose 裡面，Volume 的寫法跟 `docker run -v` 一樣，只是搬進了 YAML。

回頭看 2.1 的範例，`db` 服務掛了一個 Named Volume：

```yaml
services:
  db:
    image: postgres:17-alpine
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

`pgdata:/var/lib/postgresql/data` 的意思是：建立一個叫 `pgdata` 的 Named Volume，掛載到容器內的 `/var/lib/postgresql/data`（PostgreSQL 存資料的地方）。最下面的 `volumes: pgdata:` 是宣告這個 Volume 的存在。

這樣就算你 `docker compose down` 把容器全砍了，下次 `docker compose up` 的時候資料還在，因為 `pgdata` 這個 Volume 是獨立於容器的。

除了 Named Volume，也可以用 Bind Mount 把主機的檔案直接掛進去：

```yaml
volumes:
  - ./init.sql:/docker-entrypoint-initdb.d/init.sql:ro
```

這行把本機目前目錄下的 `init.sql` 掛進容器，`:ro` 代表容器只能讀不能改。PostgreSQL 的映像檔會在第一次啟動時自動執行 `/docker-entrypoint-initdb.d/` 裡的 SQL 檔案，很適合拿來做資料庫初始化。

> **注意**：`docker compose down` 不會刪除 Volume。要連 Volume 一起清掉要加 `-v`：`docker compose down -v`。

---

## 2.5 環境變數管理

Ocean 把 2.1 的 compose 整份 commit 上 public repo，三分鐘後收到 GitHub Secret Scanning 的警告信。Snow：「噢不。」

2.1 的範例裡，我們把密碼直接寫在 `docker-compose.yml` 裡面：

```yaml
environment:
  POSTGRES_PASSWORD: secret
```

開發時這樣寫沒問題，但如果要把 `docker-compose.yml` 推進 Git，密碼就跟著推上去了。實際專案會把敏感資訊抽到一個 `.env` 檔案裡，然後在 `.gitignore` 排除它。

### 做法：用 .env 檔案管理

先建立一個 `.env` 檔案，跟 `docker-compose.yml` 放同一個目錄：

```bash
# .env
POSTGRES_USER=myuser
POSTGRES_PASSWORD=secret
POSTGRES_DB=mydb
```

然後在 `docker-compose.yml` 裡用 `${變數名}` 引用：

```yaml
services:
  db:
    image: postgres:17-alpine
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
```

Compose 會自動讀取同目錄下的 `.env` 檔案，把變數帶進去。這樣 `docker-compose.yml` 裡面就不會有密碼，可以安心推進版本控制。

如果覺得一個一個寫 `${...}` 太麻煩，也可以用 `env_file` 一口氣把整個 `.env` 載入容器：

```yaml
services:
  adminer:
    env_file:
      - .env
```

> **重要**：記得把 `.env` 加進 `.gitignore`，永遠不要把密碼推到 Git 上。

---

## 2.6 常用 Compose 指令

前面幾節學了怎麼寫 `docker-compose.yml`，這一節整理日常操作 Compose 時最常用的指令。所有指令都是 `docker compose` 開頭，後面接動作。

### 啟動與停止

最常打的大概就這兩個：

```bash
docker compose up -d      # 啟動所有服務（背景執行）
docker compose down        # 停止並移除所有容器和網路
```

`down` 不會刪除 Volume（2.4 有提到），所以資料庫資料還在。如果要連 Volume 一起清掉：

```bash
docker compose down -v     # 停止、移除容器、刪除 Volume
```

### 查看狀態與日誌

服務跑起來之後，要怎麼知道它是不是正常的？

```bash
docker compose ps          # 列出所有服務的狀態
docker compose logs api    # 看 api 服務的日誌
docker compose logs -f api # 持續追蹤日誌（Ctrl+C 停止）
```

### 進入容器與執行命令

跟 1.8 學的 `docker exec` 一樣，只是改成用服務名稱：

```bash
docker compose exec api sh           # 進入 api 容器的 shell
docker compose run --rm api go test ./...  # 開一個新容器跑一次性命令
```

`exec` 是對已經在跑的容器執行命令，`run` 是開一個全新的容器跑完就砍掉（`--rm`）。

### 速查表

| 指令 | 說明 |
|------|------|
| `up -d` | 啟動所有服務（背景） |
| `down` | 停止並移除容器、網路 |
| `down -v` | 同上，並刪除 Volume |
| `ps` | 列出服務狀態 |
| `logs -f [service]` | 追蹤日誌 |
| `exec [service] [cmd]` | 在執行中的容器內執行命令 |
| `run --rm [service] [cmd]` | 開新容器跑一次性命令 |
| `restart [service]` | 重啟特定服務 |
| `config` | 驗證設定（展開所有變數） |

---

## 2.7 SDC 的短網址服務 Shlink

Andrew 把 Ocean 抓過來：「你自己看看 `link.sdc.tw` 是怎麼架的，看得懂就代表 Compose 學會了。」Ocean 打開 repo，發現前幾節學過的東西全部都在裡面。

前面幾節我們學了容器怎麼跑、Compose 怎麼把多個服務串起來。現在來看一個真的跑在 SDC 基礎設施上的服務——**Shlink 短網址服務**

### 什麼是 Shlink？

Shlink 是一個開源的短網址服務。簡單來說，就是把又臭又長的網址變短：

```
原始：https://docs.google.com/spreadsheets/d/1aBcDeFgHiJkLmNoPqRsTuVwXyZ/edit#gid=0
縮短：https://link.sdc.tw/abc123
```

為什麼我們自己架而不用 bit.ly？因為：

- **自訂網域**：用 `link.sdc.tw` 比較有辨識度
- **數據自己掌控**：點擊次數、來源分析都在自己手上
- **不受第三方限制**：免費方案有額度限制，自架沒有這個問題

### Shlink 的 Docker Compose 架構

Shlink 整個服務由三個容器組成，架構很乾淨：

```
┌─────────────────────────────────────────────────────┐
│                    Traefik (反向代理)                │
│              處理 HTTPS、憑證、路由分流                │
└──────────┬──────────────────────┬───────────────────┘
           │                      │
           ▼                      ▼
┌─────────────────┐    ┌─────────────────────┐
│  shlink-web-    │    │  shlink-backend     │
│  client         │    │  (API 伺服器)        │
│  管理介面 UI     │    │  處理短網址轉址       │
│  port 8080      │    │  port 8080          │
└─────────────────┘    └──────────┬──────────┘
                                  │ shlink-net
                                  │ (內部網路)
                       ┌──────────▼──────────┐
                       │  shlink-db          │
                       │  PostgreSQL 15      │
                       │  儲存短網址資料       │
                       │  Volume 持久化       │
                       └─────────────────────┘
```

三個角色分工明確：

| 服務 | image | 負責什麼 |
|------|--------|---------|
| `shlink-db` | `postgres:15-alpine` | 資料庫，存所有短網址對應、點擊紀錄 |
| `shlink-backend` | `shlinkio/shlink:stable` | 核心 API 伺服器，負責建立短網址、處理轉址 |
| `shlink-web-client` | `shlinkio/shlink-web-client:stable` | 管理介面，讓你在網頁上操作短網址 |

### 解析 docker-compose.yaml

以下是 SDC 實際使用的 `docker-compose.yaml`，我們一段一段來看：

**1. 資料庫 — `shlink-db`**

```yaml
services:
  shlink-db:
    image: postgres:15-alpine
    container_name: shlink-db
    environment:
      - POSTGRES_DB=shlink
      - POSTGRES_USER=${DB_USER}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - shlink-db-data:/var/lib/postgresql/data
    networks:
      - shlink-net
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U shlink"]
      interval: 10s
      timeout: 5s
      retries: 5
```

- **`postgres:15-alpine`**：用 Alpine 版本的 PostgreSQL 15，映像檔比較小
- **`environment`**：用 `${DB_USER}` 這種寫法，代表值從 `.env` 檔案讀進來。密碼絕對不會寫死在 YAML 裡面
- **`volumes`**：用 Named Volume `shlink-db-data` 把資料庫的資料掛出來。這樣容器重啟或重建，資料都不會消失
- **`networks: shlink-net`**：只加入內部網路，外部完全碰不到資料庫
- **`healthcheck`**：用 `pg_isready` 指令定期檢查資料庫有沒有活著，每 10 秒查一次，連續失敗 5 次才算不健康

**2. 後端 API — `shlink-backend`**

```yaml
  shlink-backend:
    image: shlinkio/shlink:stable
    container_name: shlink-backend
    depends_on:
      shlink-db:
        condition: service_healthy
    environment:
      - DB_DRIVER=postgres
      - DB_HOST=shlink-db
      - DB_PORT=5432
      - DB_NAME=shlink
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DEFAULT_DOMAIN=${SHLINK_DOMAIN}
      - IS_HTTPS_ENABLED=true
      - INITIAL_API_KEY=${SHLINK_API_KEY}
      - GEOLITE_LICENSE_KEY=${GEOLITE_KEY}
      - ALLOWED_ORIGINS=https://${SHLINK_WEB_HOST}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.shlink-api.rule=Host(`${SHLINK_DOMAIN}`)"
      - "traefik.http.routers.shlink-api.entrypoints=websecure"
      - "traefik.http.routers.shlink-api.tls=true"
      - "traefik.http.routers.shlink-api.tls.certresolver=cloudflare"
      - "traefik.http.services.shlink-api.loadbalancer.server.port=8080"
    networks:
      - traefik
      - shlink-net
```

- **`depends_on` + `condition: service_healthy`**：等資料庫的 healthcheck 通過才啟動。不是只等容器跑起來，而是等到資料庫「真的準備好了」。這就是 2.2 學的 `depends_on` 進階用法
- **`DB_HOST=shlink-db`**：注意這裡填的是服務名稱，不是 IP。Docker 的內部 DNS 會自動把 `shlink-db` 解析成正確的容器 IP
- **`labels`（Traefik 設定）**：這些 label 是給 Traefik 反向代理看的。Traefik 會自動讀取這些 label 來設定路由規則和 HTTPS 憑證。你現在不需要完全理解 Traefik，只要知道它的作用是：「讓外部透過 `https://link.sdc.tw` 連到這個容器的 8080 port」
- **雙網路 `traefik` + `shlink-net`**：backend 同時接兩個網路。`traefik` 網路讓它接收外部流量，`shlink-net` 讓它連到資料庫。這是常見的安全設計

**3. 管理介面 — `shlink-web-client`**

```yaml
  shlink-web-client:
    image: shlinkio/shlink-web-client:stable
    container_name: shlink-web-client
    environment:
      - SHLINK_SERVER_URL=https://${SHLINK_DOMAIN}
      - SHLINK_SERVER_API_KEY=${SHLINK_API_KEY}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.shlink-web.rule=Host(`${SHLINK_WEB_HOST}`)"
      - "traefik.http.routers.shlink-web.entrypoints=websecure"
      - "traefik.http.routers.shlink-web.tls=true"
      - "traefik.http.routers.shlink-web.tls.certresolver=cloudflare"
      - "traefik.http.services.shlink-web.loadbalancer.server.port=8080"
    networks:
      - traefik
```

重點拆解：

- **只連 `traefik` 網路**：Web Client 不需要直接碰資料庫，它透過 `SHLINK_SERVER_URL` 呼叫 backend 的 API 來操作資料
- **`SHLINK_SERVER_API_KEY`**：Web Client 用 API Key 跟 backend 溝通，這也是從 `.env` 讀取的

**4. Network 與 Volume**

```yaml
networks:
  traefik:
    external: true
  shlink-net:
    internal: true

volumes:
  shlink-db-data:
```

重點拆解：

- **`traefik: external: true`**：這個網路不是 Shlink 自己建的，是整個伺服器共用的。其他服務也掛在同一個 Traefik 網路上
- **`shlink-net: internal: true`**：標記為 `internal`，代表這個網路**沒有對外連線能力**。只有 `shlink-db` 和 `shlink-backend` 在裡面，資料庫完全隔離，外部無法直接存取
- **`shlink-db-data`**：Named Volume，PostgreSQL 的資料持久化就靠它

### 整體運作流程

當有人點了 `https://link.sdc.tw/abc123`，會發生什麼事？

```
1. 使用者點擊 https://link.sdc.tw/abc123
         │
2. Traefik 收到請求，根據 Host 規則轉給 shlink-backend
         │
3. shlink-backend 查詢 shlink-db：「abc123 對應到哪個網址？」
         │
4. shlink-db 回傳原始網址
         │
5. shlink-backend 回應 HTTP 302 重新導向到原始網址
         │
6. 使用者的瀏覽器跳轉到目標網站
```

---

## 2.8 小結

這一章跟著 Ocean 走過一遍：從把一串又臭又長的 `docker run` 收進 `docker-compose.yml`（2.1），到被啟動順序搞到懷疑人生而學會 `depends_on` + `healthcheck`（2.2）；接著搞懂 Compose 為什麼不用自己建網路容器之間就能互連（2.3）、把資料庫資料用 Volume 保下來（2.4）、把密碼從 YAML 裡抽到 `.env` 避免再被 GitHub 寄信警告（2.5），最後整理了日常會用到的 `up`、`down`、`ps`、`logs`、`exec` 幾個指令（2.6），並在 2.7 對著 Shlink 的 compose 把前面學到的東西全部對照過一次。

Ocean 現在可以用一個 YAML 檔把多個容器串起來跑，從簡單的 DB + Admin 介面到像 Shlink 這種三層架構都不是問題，任何人 clone 專案後都能用同樣的方式啟動。Andrew 終於不用擔心 Ocean 把指令記在腦袋裡了。接下來第 3 章會教怎麼自己打包映像檔。

→ 下一章：[03-dockerfile.md](03-dockerfile.md)
