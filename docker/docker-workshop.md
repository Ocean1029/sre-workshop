# Docker 入門指南

> **預計閱讀時間**：2 小時
> **適用對象**：具備 Linux 基礎，初次接觸 Docker 的開發者

---

## 目錄

- [Docker 入門指南](#docker-入門指南)
  - [目錄](#目錄)
  - [Part 1：Docker 核心概念與基本操作](#part-1docker-核心概念與基本操作)
    - [1.1 什麼是容器化？](#11-什麼是容器化)
      - [「ㄍㄢˋ，這在我的電腦上明明跑得起來。」](#ㄍㄢˋ這在我的電腦上明明跑得起來)
      - [所以 Docker 如何解決這個問題？？？](#所以-docker-如何解決這個問題)
      - [容器化的適用場景](#容器化的適用場景)
    - [1.2 容器 vs 虛擬機器](#12-容器-vs-虛擬機器)
      - [虛擬機器的運作方式](#虛擬機器的運作方式)
      - [容器的運作方式](#容器的運作方式)
      - [架構比較](#架構比較)
      - [選擇指引](#選擇指引)
      - [容器的底層技術](#容器的底層技術)
    - [1.3 第一次跑 Docker](#13-第一次跑-docker)
    - [1.4 Docker 架構](#14-docker-架構)
      - [Docker Client](#docker-client)
      - [Docker Daemon](#docker-daemon)
      - [containerd](#containerd)
      - [一個指令的完整旅程](#一個指令的完整旅程)
    - [1.5 映像檔與容器的關係](#15-映像檔與容器的關係)
      - [概念對比](#概念對比)
      - [映像檔的分層架構](#映像檔的分層架構)
      - [容器的可寫層](#容器的可寫層)
      - [多容器共用映像檔](#多容器共用映像檔)
      - [映像檔命名與 Tag](#映像檔命名與-tag)
      - [概念總結](#概念總結)
    - [1.6 映像檔基本操作](#16-映像檔基本操作)
      - [搜尋映像檔](#搜尋映像檔)
      - [下載映像檔](#下載映像檔)
      - [列出本機映像檔](#列出本機映像檔)
      - [刪除映像檔](#刪除映像檔)
      - [查看映像檔資訊](#查看映像檔資訊)
    - [1.7 容器基本操作](#17-容器基本操作)
      - [建立並啟動容器](#建立並啟動容器)
      - [常用 `docker run` 參數](#常用-docker-run-參數)
      - [列出容器](#列出容器)
      - [容器生命週期管理](#容器生命週期管理)
    - [1.8 Port Mapping](#18-port-mapping)
      - [Port 是什麼？](#port-是什麼)
      - [語法](#語法)
      - [常用變化](#常用變化)
    - [1.9 Volume](#19-volume)
    - [1.10 容器除錯](#110-容器除錯)
      - [查看日誌](#查看日誌)
      - [進入容器內部](#進入容器內部)
      - [查看資源使用狀況](#查看資源使用狀況)
    - [1.11 練習 1：執行你的第一個容器](#111-練習-1執行你的第一個容器)
  - [Part 2：Docker Compose 多容器編排](#part-2docker-compose-多容器編排)
    - [2.1 為什麼需要 Docker Compose？](#21-為什麼需要-docker-compose)
      - [services](#services)
      - [image](#image)
      - [ports](#ports)
      - [environment](#environment)
      - [volumes](#volumes)
      - [depends\_on](#depends_on)
      - [啟動與停止](#啟動與停止)
    - [2.2 服務定義詳解](#22-服務定義詳解)
      - [depends\_on 搭配 healthcheck](#depends_on-搭配-healthcheck)
      - [restart — 重啟策略](#restart--重啟策略)
    - [2.3 Network](#23-network)
    - [2.4 Volume](#24-volume)
    - [2.5 環境變數管理](#25-環境變數管理)
      - [做法：用 .env 檔案管理](#做法用-env-檔案管理)
    - [2.6 常用 Compose 指令](#26-常用-compose-指令)
      - [啟動與停止](#啟動與停止-1)
      - [查看狀態與日誌](#查看狀態與日誌)
      - [進入容器與執行命令](#進入容器與執行命令)
      - [速查表](#速查表)
  - [Part 3：Dockerfile 基礎](#part-3dockerfile-基礎)
    - [3.1 什麼是 Dockerfile？](#31-什麼是-dockerfile)
    - [3.2 Dockerfile 基礎指令](#32-dockerfile-基礎指令)
      - [指令速查表](#指令速查表)
    - [3.3 為 Go 應用撰寫 Dockerfile](#33-為-go-應用撰寫-dockerfile)
      - [基本版 Dockerfile（單階段建置）](#基本版-dockerfile單階段建置)
      - [Multi-stage Build](#multi-stage-build)
      - [建置與執行](#建置與執行)
  - [Part 4：綜合演練與延伸學習](#part-4綜合演練與延伸學習)
    - [4.1 SDC 的短網址服務 Shlink](#41-sdc-的短網址服務-shlink)
      - [什麼是 Shlink？](#什麼是-shlink)
      - [Shlink 的 Docker Compose 架構](#shlink-的-docker-compose-架構)
      - [解析 docker-compose.yaml](#解析-docker-composeyaml)
      - [整體運作流程](#整體運作流程)
    - [4.2 常見問題排查](#42-常見問題排查)
      - [容器啟動失敗](#容器啟動失敗)
      - [容器間無法連線](#容器間無法連線)
      - [磁碟空間不足](#磁碟空間不足)
      - [映像檔建置快取問題](#映像檔建置快取問題)
      - [常見錯誤訊息](#常見錯誤訊息)
    - [4.3 延伸學習資源](#43-延伸學習資源)
      - [官方文件](#官方文件)
      - [進階主題](#進階主題)
  - [附錄 A：Docker 指令速查表](#附錄-adocker-指令速查表)
    - [映像檔 (Image)](#映像檔-image)
    - [容器 (Container)](#容器-container)
    - [資料卷 (Volume)](#資料卷-volume)
    - [網路 (Network)](#網路-network)
    - [系統](#系統)
  - [附錄 B：Dockerfile 指令速查表](#附錄-bdockerfile-指令速查表)
  - [附錄 C：Docker Compose 速查表](#附錄-cdocker-compose-速查表)

---

## Part 1：Docker 核心概念與基本操作

### 1.1 什麼是容器化？

容器化（Containerization）是一種將「應用程式」及其「完整執行環境」包括程式語言、相依套件、設定檔與環境變數，封裝為方便轉移的元件的技術，而這個元件稱為容器（Container）。無論把容器部署到開發者的筆電、測試伺服器，還是雲端正式環境，執行結果都會保持一致。

這帶來幾個明顯的好處：團隊成員不再需要各自手動建置環境、不同專案的相依套件可以完全隔離、開發與正式環境的差異大幅縮小。容器化特別適合用在多人協作開發、CI/CD 流水線，以及需要快速複製環境的場景。

#### 「ㄍㄢˋ，這在我的電腦上明明跑得起來。」

想像一個情境，SDC 有很多個新鮮的肝正在協作開發一個由 Go 撰寫的網頁應用，專案依賴 PostgreSQL 資料庫與特定版本的函式庫。今天，身為 Project Team 的 Andrew 想要請 SRE Team 的 Ocean 把這個服務上線，兩個人的電腦環境差異很大：

| 開發者 | 作業系統 | Go 版本 | PostgreSQL 版本 | 函式庫版本 |
|--------|---------|---------|----------------|-----------|
| Andrew | macOS   | 1.21    | 15             | v2.3.1    |
| Ocean    | Ubuntu  | 1.19    | 14             | v2.1.0    |

Andrew 的程式原本在 macOS 上跑得順順的，但當程式轉移到 Ocean 的電腦之後就看到終端機噴錯了。Ocean 看著空空如也的 README 不知道要從哪裡改設定，在跟 Claude Code 奮戰兩個小時之後大喝了兩口 Whiskey，一不小心睡著了。

Andrew 因為答應了計中明天要上線，只好自己登入到系辦的機器上部署服務。因為正式環境又是另一套設定，程式剛一執行就跑出了 2000 行的 errors。可憐的 Andrew 只好一邊咒罵 SRE Team，一邊喝著自己的 Triple espresso with Vodka，跟這些 errors 奮鬥到天亮。

還好這些不是真實故事，~~因為 Ocean 會乖乖地部署完服務再睡覺~~，聰明的 Docker 也會避免這件事情發生。

#### 所以 Docker 如何解決這個問題？？？

開發者將執行環境的完整規格撰寫為 Dockerfile，定義以下資訊：

- 使用哪個版本的 Go
- 需要安裝哪些套件
- 需要哪些環境變數
- 程式如何啟動

任何人在任何機器上，都能透過 Docker 依照 Dockerfile 自動建置出完全相同的容器。新成員加入專案不再需要逐步照著手冊設定環境，執行一個指令即可開始工作。

#### 容器化的適用場景

- 多人協作開發：不同人使用不同作業系統與不同版本的工具，容器確保每個人面對完全相同的執行環境。

- 解決相依性衝突：不同專案可能依賴同一套件的不同版本。容器讓每個專案在隔離的空間中獨立運作，彼此不受干擾。

- 標準化的專案啟動程序：Dockerfile 是可執行的環境文件，明確記錄專案需求與啟動方式，取代容易過時的手動安裝說明。

- 一致的部署流程：開發、測試、正式環境使用同一個設定流程，大幅降低環境差異導致的問題。

- 快速實驗：想試 PostgreSQL 17 或 Redis 新版本？一行指令就可以快速啟動，並且用完直接回收，不需要安裝也不會汙染系統環境。

> 容器化將「環境」從各自的機器中抽離，轉變為可版本控制、可共享、可重現的成品。

### 1.2 容器 vs 虛擬機器

Ocean 在系上修過作業系統，對虛擬機器（VM）不陌生，計中機房裡就有好幾台 VM 在跑各種服務。在聽完 Snow 對 Docker 的介紹之後，Ocean 的第一個問題是：「這跟 VM 有什麼不一樣？不都是把東西隔離起來跑」

Snow：問得好，容器跟 VM 要解決的問題確實差不多，但做法完全不一樣。

#### 虛擬機器的運作方式

虛擬機器透過 Hypervisor（虛擬機管理程式）在硬體上模擬出多台獨立電腦。每台 VM 都有完整的作業系統，從核心到使用者空間全部包在裡面。

你可以把 VM 想成大樓裡的獨立公寓，各自有完整的水電管路跟門牌。隔離性很強，但每一間都要重複建設同樣的基礎設施，資源利用率偏低。

常見的 Hypervisor 有 VMware ESXi、KVM、Microsoft Hyper-V。

#### 容器的運作方式

容器不模擬硬體，而是直接利用 Linux 核心本身的功能（Namespace + Cgroup）來隔離程序。所有容器共用同一個宿主機核心，每個容器只包含應用程式跟它需要的函式庫。

容器比較像大樓裡的 co-working space，共用水電、網路、電梯這些基礎設施，各個團隊有自己獨立的工作區域，互不干擾。省空間、啟動快，但隔離強度不如獨立公寓。

#### 架構比較

```
┌──────────────────────────────┬──────────────────────────────┐
│        虛擬機器 (VM)          │        容器 (Container)       │
├──────────────────────────────┼──────────────────────────────┤
│                              │                              │
│  ┌────────┐  ┌────────┐      │  ┌────────┐  ┌────────┐      │
│  │ App A  │  │ App B  │      │  │ App A  │  │ App B  │      │
│  ├────────┤  ├────────┤      │  ├────────┤  ├────────┤      │
│  │ Bins/  │  │ Bins/  │      │  │ Bins/  │  │ Bins/  │      │
│  │ Libs   │  │ Libs   │      │  │ Libs   │  │ Libs   │      │
│  ├────────┤  ├────────┤      │  └────┬───┘  └───┬────┘      │
│  │Guest OS│  │Guest OS│      │       │          │           │
│  └────┬───┘  └───┬────┘      │  ┌────┴──────────┴────┐      │
│  ┌────┴──────────┴────┐      │  │   Docker Engine     │     │
│  │    Hypervisor      │      │  ├─────────────────────┤     │
│  ├────────────────────┤      │  │     Host OS         │     │
│  │     Host OS        │      │  ├─────────────────────┤     │
│  ├────────────────────┤      │  │    Infrastructure   │     │
│  │   Infrastructure   │      │  └─────────────────────┘     │
│  └────────────────────┘      │                              │
│                              │                              │
└──────────────────────────────┴──────────────────────────────┘
```

注意右邊，沒有 Guest OS 那一層。這就是容器比 VM 輕這麼多的原因。

| 比較項目 | 虛擬機器 (VM) | 容器 (Container) |
|---------|--------------|-----------------|
| **隔離方式** | 硬體層虛擬化 (Hypervisor) | 作業系統層虛擬化 (Namespace + Cgroup) |
| **Guest OS** | 每台 VM 都有完整的 OS | 共用宿主機核心，無 Guest OS |
| **啟動時間** | 分鐘級（需要開機） | 秒級（啟動一個程序） |
| **映像檔大小** | GB 級（含完整 OS） | MB 級（App + Libs） |
| **資源佔用** | 高（每台都需分配 CPU、記憶體給 OS） | 低（共用核心，按需分配） |
| **隔離強度** | 強（完整硬體隔離） | 較弱（共用核心，靠 Namespace 隔離） |
| **部署密度** | 一台伺服器通常跑數台至數十台 VM | 一台伺服器可跑數十至數百個容器 |
| **適用場景** | 強隔離需求、不同 OS、遺留系統 | 微服務、CI/CD、快速部署、開發環境 |

#### 選擇指引

實際上，這兩個東西經常搭配使用。像是在雲端平台上開一台 VM（EC2 / GCE），然後在 VM 裡面跑 Docker 容器。

**適合使用 VM 的場景：**
- 需要執行不同的作業系統（例如在 Linux 上執行 Windows）
- 需要硬體等級的隔離（多租戶環境、金融法規要求）
- 遺留系統無法容器化

**適合使用容器的場景：**
- CI/CD 流水線
- 開發與測試環境
- 需要快速橫向擴展（scale out）的應用

#### 容器的底層技術

容器底層靠的是 Linux 核心的三個機制來做隔離：

**Namespace — 隔離可見範圍**

Namespace 讓每個容器以為自己是系統上唯一的存在，只能看到屬於自己的資源。

| Namespace 類型 | 隔離內容 | 效果 |
|---------------|---------|------|
| **PID** | 程序 ID | 容器內的程序以為自己是 PID 1 |
| **Network** | 網路介面、IP、路由 | 每個容器有獨立的網路堆疊與 IP 位址 |
| **Mount** | 檔案系統掛載點 | 容器只能看到自己的檔案系統 |
| **UTS** | 主機名稱 | 每個容器可以有自己的 hostname |
| **IPC** | 跨程序通訊 | 容器間的 IPC 相互隔離 |
| **User** | 使用者/群組 ID | 容器內的 root 不等同於宿主機的 root |
| **Cgroup** | Cgroup 層級 | 容器只能看到自己的資源限制 |

簡單來說，你在宿主機上跑 `ps aux` 可以看到所有程序，但在容器裡面跑同樣的指令就只能看到自己的程序，這就是 PID Namespace 在做的事。

**Cgroup（Control Group）— 限制資源用量**

Cgroup 管的是每個容器可以用多少系統資源：

| 資源 | 說明 | 範例 |
|------|------|------|
| **CPU** | 限制 CPU 使用量 | 最多使用 1.5 個 CPU 核心 |
| **Memory** | 限制記憶體使用量 | 最多使用 512MB，超過則 OOM Kill |
| **I/O** | 限制磁碟讀寫速度 | 讀寫速度上限 100MB/s |
| **Network** | 限制網路頻寬 | 上行/下行各 100Mbps |

**Union Filesystem — 管理檔案系統**

Union Filesystem（像 OverlayFS）可以把多個目錄「疊加」在一起，看起來像一個統一的檔案系統。這是 Docker 映像檔分層架構的基礎（詳見 [1.5 節](#15-映像檔與容器的關係)）。

```
┌───────────────────────────────────────────────────────┐
│              容器的三大底層技術                          │
├───────────────────────────────────────────────────────┤
│                                                       │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────┐   │
│  │  Namespace  │  │   Cgroup    │  │ Union FS     │   │
│  │             │  │             │  │              │   │
│  │ 隔離可見性    │  │ 限制資源量   │  │ 管理檔案系統   │   │
│  │ 「看到什麼」  │  │ 「用多少」   │  │ 「存什麼」     │   │
│  └─────────────┘  └─────────────┘  └──────────────┘   │
│         │                │                │           │
│         └────────────────┼────────────────┘           │
│                          ▼                            │
│                 ┌─────────────────┐                   │
│                 │   Linux Kernel  │                   │
│                 └─────────────────┘                   │
│                                                       │
│  這些都是 Linux 核心的原生功能，Docker 的貢獻是將它們       │
│  封裝成簡單易用的工具。                                  │
│                                                       │
└───────────────────────────────────────────────────────┘
```

> **注意**：容器不是真正獨立的虛擬機器，它就是被核心功能隔離出來的程序。你在宿主機上跑 `ps aux` 其實看得到容器裡的程序。它們本質上就是普通的 Linux 程序，只是被 Namespace 跟 Cgroup 圍起來而已。

---

### 1.3 第一次跑 Docker

Ocean 聽完概念坐不住了，打開終端機想要親手試試看。幸好每次 Andrew 都會趁 Ocean 不在的時候，偷偷往他電腦裡面下載一些怪東西，而他前段時間下載的 OrbStack 剛好可以用來啟動 Docker。

安裝 docker:

```bash
brew install orbstack
```

因為 Andrew 歧視 Windows 系統，所以 Windows 用戶只好請你問你的 AI 夥伴怎麼裝 Docker 了。

```bash
# 驗證安裝
docker version

# 查看 Docker 系統資訊
docker info

# 執行測試容器
docker run hello-world
```

**預期輸出：**

```
Hello from Docker!
This message shows that your installation appears to be working correctly.
...
```

看到這段訊息就代表 Docker 裝好了。
接下來 Ocean 打算試點更有感覺的，跑一個真正的網頁伺服器：

```bash
# 跑一個 Nginx 網頁伺服器
docker run --name my-nginx -p 8080:80 nginx:1.27-alpine
```

一行指令就把一個網頁伺服器跑起來了。原來，`docker run` 是運行一個容器的意思，後面接了一些參數。`-p 8080:80` 把容器的 80 port 對映到你電腦的 8080 port，所以我們可以在 `http://localhost:8080` 找到這個應用，看到 404 Not Found 就代表服務啟動起來了。測試完成之後可以清理一下：

```bash
docker stop my-nginx && docker rm my-nginx
```

---

### 1.4 Docker 架構

Ocean 看著瀏覽器上出現的 Nginx 歡迎頁面，進入了自信之巔。但他馬上冒出下一個問題：「剛剛只打了一行指令，背後到底發生什麼事？」此時，Snow 的聲音從遠方冒了出來，而他手上拿著一張 Docker 架構圖，像是已經等著有人問這個問題很久了 ......

```
┌──────────────────────────────────────────────────────────────────┐
│                         Docker 架構                               │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐                                                │
│  │ Docker CLI   │                                                │
│  │              │                                                │
│  │ docker run   │                                                │
│  │ docker build │                                                │
│  │ docker pull  │                                                │
│  └──────┬───────┘                                                │
│         │  REST API（通常透過 Unix Socket）                        │
│         ▼                                                        │
│  ┌──────────────────────────────────────────┐                    │
│  │          Docker Daemon                   │                    │
│  │                                          │                    │
│  │  接收 Client 指令，協調所有 Docker 操作     │                    │
│  │                                          │                    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐  │                    │
│  │  │ Image    │ │ Network  │ │ Volume   │  │                    │
│  │  │ 管理      │ │ 管理     │ │ 管理      │  │                    │
│  │  └──────────┘ └──────────┘ └──────────┘  │                    │
│  └──────────────────┬───────────────────────┘                    │
│                     │                                            │
│                     ▼                                            │
│  ┌──────────────────────────────────────────┐                    │
│  │            containerd                    │                    │
│  │                                          │                    │
│  │          管理容器生命週期                  │                    │
│  └──────────────────┬───────────────────────┘                    │
│                     │                                            │
│                     ▲                                            │
│                     │ 拉取映像檔                                  │
│                     ▼                                            │
│  ┌──────────────────────────────────────────┐                    │
│  │          Docker Registry                 │                    │
│  │                                          │                    │
│  │  存放和分發映像檔的倉庫，例如                │                    │
│  │  nginx, postgres, redis, golang, ubuntu  │                    │
│  └──────────────────────────────────────────┘                    │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

#### Docker Client

就是你在終端機用的 `docker` 命令列工具。它本身不管理容器，只負責把你的指令轉成 REST API 請求丟給 Docker Daemon 處理。Client 跟 Daemon 預設透過 Unix Socket（`/var/run/docker.sock`）溝通，通常在同一台機器上，但也可以透過 TCP 連到遠端的 Daemon。

```bash
# 查看 Docker Client 與 Daemon 的版本資訊
docker version
```

#### Docker Daemon

Docker Daemon（程序名稱為 `dockerd`）是 Docker 的核心服務，一直在背景執行。它接收 Client 丟過來的 API 請求，負責「協調」映像檔管理、容器建立、網路設定跟儲存管理這些操作。

舉個例子，你跑 `docker run nginx` 的時候，Daemon 會依序做這些事：
1. 檢查本機是否有 `nginx` 映像檔
2. 若無，從 Registry 下載
3. 請 containerd 建立並啟動容器
4. 設定網路（分配 IP、建立 bridge）
5. 掛載 Volume
6. 回報結果給 Client

dockerd 是領導者（aka AK），他負責規劃與安排任務，不負責執行。實際的運作會交給 containerd 與更底層的工人去執行 Linux 指令。

#### containerd

containerd 是負責管理容器生命週期的高階執行環境（aka 高級主管）。包括拉取跟推送映像檔、建立和刪除容器、管理映像檔儲存。

**Docker Registry（映像檔倉庫）**

存放跟分發映像檔的服務，你可以把它想成映像檔的 GitHub。

| Registry | 說明 |
|----------|------------------|
| **Docker Hub** | 官方公開倉庫，有大量社群映像檔（預設） |
| **GitHub Container Registry (ghcr.io)** | GitHub 提供的映像檔倉庫 |
| **Harbor** | 開源的自建映像檔倉庫，SDC 的映像檔管理是用這個。 |

#### 一個指令的完整旅程

來看看當你打下 `docker run nginx` 之後，背後到底發生了什麼事：

```
┌─────────────────────────────────────────────────────────────┐
│  docker run nginx 的完整流程                                  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. 使用者輸入 docker run nginx                               │
│     │                                                       │
│     ▼                                                       │
│  2. Docker CLI 將指令轉換為 REST API 請求                      │
│     POST /containers/create + POST /containers/{id}/start   │
│     │                                                       │
│     ▼                                                       │
│  3. Docker Daemon 收到請求                                   │
│     ├─→ 本機有 nginx 映像檔嗎？                               │
│     │   ├─ 有 → 直接使用                                     │
│     │   └─ 沒有 → 從 Registry 下載                           │
│     │                                                       │
│     ▼                                                       │
│  4. Daemon 請 containerd 建立容器                            │
│     │                                                       │
│     ▼                                                       │
│  5. 容器啟動完成，nginx 開始運作                               │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

> **macOS 跟 Windows 的使用者注意**：Docker 容器需要 Linux 核心才能跑。在 macOS 跟 Windows 上，Docker Desktop 會在背景偷偷開一台輕量的 Linux VM，Docker Daemon 其實是跑在這台隱藏的 VM 裡面。

---

### 1.5 映像檔與容器的關係

回想前面學到的東西，Ocean 知道了自己的服務只要跑在容器裡面，就可以容易的轉移到別人的電腦上面。但前面幾節一直出現「映像檔」這個詞，讓他腦袋有點混亂。Ocean 舉手問 Snow：「映像檔跟容器不是同一個東西嗎？」Snow 笑了一下：「很多人一開始都這樣以為，但它們其實完全不同。」

映像檔（Image）是一個打包好的、唯讀的檔案，裡面定義並包含了運行一個應用程式所需的完整執行環境：作業系統的基礎檔案、程式語言的執行環境、你的程式碼、相依套件、設定檔。當你用 `docker run` 把映像檔跑起來，跑起來的那個東西就叫做容器（Container）。映像檔是一個不會變的藍圖，但你可以從同一個映像檔跑出很多個容器，每個容器各自獨立運作。

剛才跑 `docker run hello-world` 的時候，終端機上先後出現了 `Unable to find image 'hello-world:latest' locally`，以及 `Pulling from library/hello-world` 的訊息。他代表我們本地沒有那個藍圖，所以自動往 registry （預設是 Docker Hub）pull 映像檔，並繼續啟動容器。


#### 概念對比

| 比喻 | 映像檔（Image） | 容器（Container） |
|------|----------------|-------------------|
| **建築** | 建築藍圖 | 蓋好的房子 |
| **程式設計** | Class（類別定義） | Instance（實體物件） |
| **作業系統** | ISO 映像檔 | 執行中的系統 |

#### 映像檔的分層架構

映像檔不是一個大檔案，而是由好幾個唯讀的 Layer 疊起來的。以 `nginx:1.27-alpine` 為例，你可以用 `docker history` 看到它的分層：

```bash
docker history nginx:1.27-alpine
```

每一層都記錄了一個變更：最底層是 Alpine Linux 的基礎檔案，往上可能是安裝 Nginx、複製設定檔、設定啟動指令等等。這些層全部疊在一起，就組成了完整的映像檔。這個分層設計有幾個好處：

- 磁碟共用：假設 10 個映像檔都用 `alpine:3.21` 當基底，Alpine 那一層在磁碟上只需要存一份，省下空間。

- 傳輸效率：拉取映像檔的時候，只需要下載你本機還沒有的 Layer。如果你已經有 Alpine 的基礎層，拉 `nginx:1.27-alpine` 這張 Image（nginx built on Alpine layer），就只需要下載 Nginx 相關的層。

#### 容器的可寫層

映像檔的每一層都是唯讀的，不能改。那容器跑起來之後產生的資料寫在哪？

Docker 會在映像檔的所有唯讀層上面加一層可讀寫的容器層。容器執行期間的所有變更（新建檔案、修改設定、寫入日誌）都存在這一層。

```
┌─────────────────────────────────────────────┐
│  Container = Image Layers + 可寫層           │
├─────────────────────────────────────────────┤
│                                             │
│  ┌──────────────────────────────────┐       │
│  │  Container Layer（可讀寫）         │      │
│  │  新建的檔案、修改的設定、日誌等       │      │
│  │                                  │       │
│  │  ⚠ 容器刪除時，這一層就消失了！      │       │
│  ├──────────────────────────────────┤       │
│  │  Layer 3: Nginx 啟動設定（唯讀）   │       │
│  ├──────────────────────────────────┤       │
│  │  Layer 2: 安裝 Nginx（唯讀）       │Image  │
│  ├──────────────────────────────────┤Layers │
│  │  Layer 1: Alpine Linux（唯讀）    │       │
│  └──────────────────────────────────┘       │
│                                             │
└─────────────────────────────────────────────┘
```

#### 多容器共用映像檔

```
┌────────────────────────────────────────────────────────────┐
│                 多容器共用同一個映像檔                        │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  Container A          Container B          Container C     │
│  ┌──────────┐        ┌──────────┐        ┌──────────┐      │
│  │ 可寫層 A  │        │ 可寫層 B  │        │ 可寫層 C  │      │
│  │ (日誌、   │        │ (不同的   │        │ (各自獨   │      │
│  │  暫存等)  │        │  資料)    │        │  立的)    │     │
│  └─────┬────┘        └─────┬────┘        └─────┬────┘      │
│        │                   │                   │           │
│        └───────────────────┼───────────────────┘           │
│                            ▼                               │
│              ┌──────────────────────┐                      │
│              │  共用的映像檔 Layer    │                      │
│              │  (唯讀，只存一份)      │                      │
│              │                      │                      │
│              │  Layer 5: app        │                      │
│              │  Layer 4: deps       │                      │
│              │  Layer 3: runtime    │                      │
│              │  Layer 2: libs       │                      │
│              │  Layer 1: base OS    │                      │
│              └──────────────────────┘                      │
│                                                            │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

#### 映像檔命名與 Tag

每個映像檔可以有多個 Tag 來區分版本。命名格式長這樣：

```
┌────────────────────────────────────────────────────────────┐
│                    映像檔命名格式                            │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  完整格式：                                                  │
│  [registry/][username/]repository[:tag]                    │
│                                                            │
│  範例：                                                     │
│  docker.io/library/nginx:1.27-alpine                       │
│  ├────────┘ ├─────┘ ├───┘ ├─────────┘                      │
│  │          │       │     └─ Tag（版本標記）                 │
│  │          │       └─ Repository（映像檔名稱）              │
│  │          └─ Username（官方映像檔為 library，通常省略）      │
│  └─ Registry（預設 docker.io，通常省略）                      │
│                                                            │
│  常見 Tag 慣例：                                             │
│  nginx:latest        → 最新版（不建議正式環境使用）             │
│  nginx:1.27          → 主要版本                              │
│  nginx:1.27.3        → 精確版本（正式環境建議使用）             │
│  nginx:1.27-alpine   → 基於 Alpine 的精簡版                  │
│  nginx:1.27-bookworm → 基於 Debian Bookworm 的版本           │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

#### 概念總結

| 概念 | 映像檔（Image） | 容器（Container） |
|------|----------------|-------------------|
| **性質** | 唯讀模板，由多個 Layer 組成 | 可讀寫的執行實體，映像檔 + 可寫層 |
| **生命週期** | 建置後不可變（Immutable） | 可建立、啟動、暫停、停止、刪除 |
| **儲存** | 分層結構，可跨映像檔共用 | 可寫層存放執行期間的變更 |
| **可重複性** | 同一 Dockerfile 永遠建出相同的映像檔 | 每個容器的可寫層內容可能不同 |
| **數量關係** | 一個映像檔可產生多個容器 | 一個容器只對應一個映像檔 |
| **資料持久性** | 除非主動刪除，否則永久存在 | 可寫層隨容器刪除而消失（需 Volume 持久化） |

---

### 1.6 映像檔基本操作

知道映像檔是什麼之後，來看看怎麼找、下載、管理這些映像檔。

#### 搜尋映像檔

```bash
# 從 Docker Hub 搜尋映像檔
docker search nginx

# 建議直接前往 Docker Hub 網站搜尋，資訊更完整
# https://hub.docker.com
```

#### 下載映像檔

```bash
# 下載最新版（預設是 :latest）
docker pull nginx

# 下載特定版本
docker pull nginx:1.27

# 下載特定平台版本
docker pull --platform linux/amd64 nginx:1.27

# 下載 Alpine 精簡版（體積更小，適合正式環境）
docker pull nginx:1.27-alpine
```

> **注意**：Tag 就是映像檔的版本標記。正式環境一定要指定明確版本號，例如 `nginx:1.27.3`，否則當軟體升級時，容易出現不相容的問題。


#### 列出本機映像檔

```bash
docker images

# 輸出範例：
# REPOSITORY   TAG           IMAGE ID       CREATED       SIZE
# nginx        1.27-alpine   a2bd6dc6e5e6   2 weeks ago   43.3MB
# nginx        1.27          39286ab8a5e1   2 weeks ago   192MB
# hello-world  latest        d2c94e258dcb   9 months ago  13.3kB
```

`nginx:1.27-alpine`（43MB）跟 `nginx:1.27`（192MB）差了快五倍，這就是用 Alpine 當基底的好處。

#### 刪除映像檔

```bash
# 刪除指定映像檔
docker rmi nginx:1.27

# 刪除所有未使用的映像檔（dangling images）
docker image prune
```

#### 查看映像檔資訊

```bash
# 查看映像檔詳細資訊（層數、環境變數等）
docker inspect nginx:1.27-alpine

# 查看映像檔的建置歷史（每一層的指令）
docker history nginx:1.27-alpine
```

---

### 1.7 容器基本操作

有了映像檔，接下來就是把它跑起來。這一節整理容器從建立到刪除的完整操作。

#### 建立並啟動容器

```bash
# 基本語法
docker run [OPTIONS] IMAGE [COMMAND] [ARG...]

# 前景模式執行（Ctrl+C 可停止）
docker run nginx:1.27-alpine

# 指定容器名稱
docker run --name my-nginx nginx:1.27-alpine

# 執行後自動刪除容器（適合一次性任務，後面可以直接接參數是要執行的命令，nginx -v 是看 nginx 的版本。）
docker run --rm nginx:1.27-alpine nginx -v
```

#### 常用 `docker run` 參數

| 參數 | 說明 | 範例 |
|------|------|------|
| `-d` | 背景執行（detach） | `docker run -d nginx` |
| `--name` | 命名容器 | `docker run --name web nginx` |
| `-p` | Port mapping（host:container） | `docker run -p 8080:80 nginx` |
| `-v` | 掛載 Volume | `docker run -v /data:/app/data nginx` |
| `-e` | 設定環境變數 | `docker run -e DB_HOST=db nginx` |
| `--rm` | 停止後自動刪除 | `docker run --rm nginx` |
| `--network` | 指定網路 | `docker run --network=mynet nginx` |

#### 列出容器

```bash
# 列出執行中的容器
docker ps

# 列出所有容器（包含已停止的）
docker ps -a

# 只顯示容器 ID
docker ps -q

# 輸出範例：
# CONTAINER ID   IMAGE               STATUS          PORTS                  NAMES
# a1b2c3d4e5f6   nginx:1.27-alpine   Up 10 minutes   0.0.0.0:8080->80/tcp   my-nginx
```

#### 容器生命週期管理

```bash
# 停止容器（發送 SIGTERM，會）
docker stop my-nginx

# 啟動已停止的容器
docker start my-nginx

# 強制停止容器（發送 SIGKILL，不建議使用）
docker kill my-nginx

# 刪除已停止的容器
docker rm my-nginx

# 強制刪除執行中的容器
docker rm -f my-nginx

# 刪除所有已停止的容器
docker container prune
```

**容器生命週期圖：**

```
┌─────────────────────────────────────────────────────────────┐
│                    容器生命週期                               │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  docker create         docker start         docker stop     │
│  ┌──────────┐  ───→   ┌──────────┐  ───→   ┌──────────┐     │
│  │ Created  │         │ Running  │         │ Stopped  │     │
│  └──────────┘  ←───   └──────────┘  ←───   └──────────┘     │
│                                                             │
│  docker run = docker create + docker start                  │
│                                                             │
│              docker rm                                      │
│  任何狀態 ────────────────→ 刪除                              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

### 1.8 Port Mapping

Ocean 想 demo 給 Andrew 看自己學會用 Docker 了，很帥氣地打了 `docker run -d nginx:1.27-alpine`，打開瀏覽器想要連接 `http://localhost:3000`，結果發現 "This site can’t be reached"。Ocean 以為是容器沒有順利啟動，打了 `docker ps` 卻發現容器跑得好好的。

Andrew 看了一眼指令：「你沒加 `-p`，port 根本沒開出來啊。」

Ocean：「port 是什麼？」Snow：「6」

#### Port 是什麼？

一台電腦上同時會跑很多網路服務：Web Server、資料庫、SSH 等。它們共用同一個 IP 位址，作業系統透過 **Port**（連接埠，範圍 0–65535）區分流量該送給哪個程式。

每個容器都有獨立的網路空間（Network Namespace），容器內部的服務對宿主機而言是不可見的。**Port Mapping** 的作用就是在宿主機與容器之間建立一條轉發通道，將宿主機指定連接埠收到的流量轉發至容器內的連接埠。

| 名詞 | 說明 | 範例 |
|------|------|------|
| **Container Port** | 應用程式在容器內監聽的連接埠，僅存在於容器的網路空間中，宿主機無法直接存取 | Node.js `app.listen(3000)` → Container Port 為 3000 |
| **Host Port** | 宿主機對外開放的連接埠，Docker 監聽此埠並轉發流量至對應的 Container Port | `-p 8080:3000` → Host Port 為 8080 |

```
┌─────────────────────────────────────────────────┐
│                    Port Mapping                 │
├─────────────────────────────────────────────────┤
│                                                 │
│  宿主機 (Host)                                   │
│  ┌─────────────────────────────────────┐        │
│  │                                     │        │
│  │   瀏覽器 → http://localhost:8080     │        │
│  │                    │                │        │
│  │                    ▼                │        │
│  │          Host Port 8080             │        │
│  │                    │                │        │
│  │         ┌──────────┼──────────┐     │        │
│  │         │ Container│          │     │        │
│  │         │          ▼          │     │        │
│  │         │  Container Port 80  │     │        │
│  │         │          │          │     │        │
│  │         │     ┌────▼────┐     │     │        │
│  │         │     │  Nginx  │     │     │        │
│  │         │     └─────────┘     │     │        │
│  │         └─────────────────────┘     │        │
│  │                                     │        │
│  └─────────────────────────────────────┘        │
│                                                 │
└─────────────────────────────────────────────────┘
```

#### 語法

```bash
# -p <Host Port>:<Container Port>
docker run -d --name web -p 8080:80 nginx:1.27-alpine

# 測試
curl http://localhost:8080
```

#### 常用變化

| 用法 | 指令 | 說明 |
|------|------|------|
| 對映多個連接埠 | `docker run -d -p 8080:80 -p 8443:443 nginx` | 同時開放 HTTP 與 HTTPS |
| 限定綁定 IP | `docker run -d -p 127.0.0.1:8080:80 nginx` | 只允許本機連線，不對外開放 |
| 隨機分配 Host Port | `docker run -d -P nginx` | Docker 自動選一個可用的 port，用 `docker port` 查看 |
| 查看對映 | `docker port web` | 列出指定容器的所有 port mapping |

> **注意**：Container Port 不會衝突。每個容器有獨立的網路空間，十個容器都跑 port 3000 完全沒問題，只要 Host Port 不重複就好。

---

### 1.9 Volume

Ocean 學會 Port Mapping 之後，原本就已經在自信之巔的他再度信心大增，決定把社團的 PostgreSQL 資料庫也用 Docker 跑起來。他花了一個下午匯入資料、調整設定，一切運作良好。隔天早上 Ocean 發現容器不知道為什麼停了，他想說沒關係，砍掉重跑一個就好：

```bash
docker rm db && docker run -d --name db postgres:17-alpine
```

打開資料庫一看，裡面比新的還乾淨，喔不對他就是新的。Ocean 崩潰地跑去找 Snow：「我的資料呢？！」Snow 淡定地喝了一口咖啡：「容器砍掉，裡面的可寫層就沒了啊。還記得 1.5 節講的嗎？」

Ocean 這才想起來——容器的可寫層是跟著容器走的，`docker rm` 的瞬間，資料庫的紀錄、使用者上傳的檔案、辛辛苦苦產生的 log，全部歸零。

Snow：「噢不。」

Docker 的 **Volume** 把宿主機的目錄或 Docker 管理的儲存空間掛載到容器裡面，讓資料可以跨容器生命週期保留下來。

常見的掛載方式有兩種：

- **Bind Mount（綁定掛載）**：把宿主機的檔案目錄直接掛進容器。只需在 `docker run` 時添加 `-v /host/path:/container/path` 即為設定完成。因為改完本機的檔案後會直接影響容器內部，適合在開發時同步原始碼。
- **Named Volume（具名資料卷）**：由 Docker 管理的儲存空間，存放在 host 的 `/var/lib/docker/volumes/` 底下。只需在 `docker run` 時添加 `-v mydata:/container/path` 即為設定完成。適合資料庫這種需要持久化的資料，以及不需要從 host 直接存取內容的情況。


```bash
# Bind Mount
docker run -d --name web \
  -p 8080:80 \
  -v $(pwd)/html:/usr/share/nginx/html \
  nginx:1.27-alpine
# $(pwd) 在 Linux 當中是當前目錄的意思。

```

```bash
# Named Volume

docker volume create pgdata
# 在使用 Named Volume 時，你需要先像是變數宣告一樣 create 它。

docker run -d --name db \
  -v pgdata:/var/lib/postgresql/data \
  postgres:17-alpine
```

```bash
# 列出所有 Volume
docker volume ls

# 查看 Volume 詳細資訊
docker volume inspect pgdata

# 刪除 Voulem
docker volume rm pgdata

# 刪除未使用的 Volume，執行前請先問你家老大。
docker volume prune
```

> **注意**：Bind Mount 要用絕對路徑（不能用相對路徑），Named Volume 只要寫名稱就好。如果路徑是 `/` 開頭，Docker 就會當成 Bind Mount。

---

### 1.10 容器除錯

#### 查看日誌

```bash
docker logs my-nginx
```

| 參數 | 效果 | 範例 |
|------|------|------|
| `-f` | 持續追蹤，類似 `tail -f` | `docker logs -f my-nginx` |
| `--tail N` | 只看最後 N 行 | `docker logs --tail 100 my-nginx` |
| `-t` | 顯示時間戳記 | `docker logs -t my-nginx` |
| `--since` | 指定時間之後的日誌 | `docker logs --since 2024-01-01 my-nginx` |

參數可以組合使用，例如 `docker logs -f --tail 50 -t my-nginx` 就是持續追蹤最後 50 行並顯示時間。

#### 進入容器內部

```bash
# 在容器內開啟 shell (可以用 sh/bash)
docker exec -it my-nginx sh
# -i: 保持標準輸入開啟（interactive）
# -t: 分配偽終端（tty），兩個一起用是常見操作。

# 執行單一命令（不進入互動模式）
docker exec my-nginx cat /etc/nginx/nginx.conf

# 以 root 身份進入
docker exec -it -u root my-nginx sh
```

#### 查看資源使用狀況

```bash
# 即時顯示所有容器的資源使用
docker stats

# 只看特定容器
docker stats my-nginx

# 輸出範例：
# CONTAINER ID   NAME       CPU %   MEM USAGE / LIMIT     MEM %   NET I/O
# a1b2c3d4e5f6   my-nginx   0.00%   3.441MiB / 7.667GiB   0.04%   1.45kB / 0B
```

### 1.11 練習 1：執行你的第一個容器

Nginx 是一個常見的網頁伺服器，很多網站背後都用它來處理請求，這次的目標就是啟動一個 Nginx 容器，掛載自訂頁面，並透過瀏覽器存取。

  **Step 1：建立工作目錄和網頁**

  ```bash
  mkdir -p ~/docker-lab/lab1
  cd ~/docker-lab/lab1

  cat << 'EOF' > index.html
  <!DOCTYPE html>
  <html>
  <head>
      <title>Docker Lab 1</title>
  </head>
  <body>
      <h1>Hello from Docker!</h1>
      <p>This page is served from a Docker container.</p>
  </body>
  </html>
  EOF
  ```

  **Step 2：啟動 Nginx 容器**

  ```bash
  docker run -d \
    --name lab1-nginx \
    -p 8080:80 \
    -v $(pwd):/usr/share/nginx/html \
    nginx:1.27-alpine
  ```

  **Step 3：確認容器在跑，然後測試**

  ```bash
  docker ps
  curl http://localhost:8080
  ```

  打開瀏覽器連 `http://localhost:8080`，應該會看到你剛才寫的 HTML 頁面。

  **Step 4：試試除錯指令**

  ```bash
  # 看日誌
  docker logs lab1-nginx

  # 進去容器裡面逛逛
  docker exec -it lab1-nginx sh
  # 試試看在裡面：
  #   ls /usr/share/nginx/html/
  #   cat /etc/nginx/nginx.conf
  #   exit # 離開容器，回到自己的終端機。
  ```

  **Step 5：清理**

  ```bash
  docker stop lab1-nginx
  docker rm lab1-nginx
  ```

**延伸練習：**

1. 修改 `index.html` 的內容，不重啟容器，直接重新整理瀏覽器查看變化（思考：為什麼修改會立即生效？）

到這裡，Ocean 回頭看看自己在 Part 1 學了什麼：知道容器化是怎麼回事、容器跟 VM 差在哪、映像檔跟容器的關係、怎麼用 `docker run` 把服務跑起來、用 `-p` 開 port、用 `-v` 掛資料、容器壞了怎麼查 log 跟進去除錯。
Snow 看著 Ocean 的練習成果，點了點頭：「不錯嘛，基本操作都會了。剛好 SDC 缺了一個 SRE Team，你要不要來做做看？」

Ocean：「......那要做什麼？」

Snow：「用 Docker Compose 把這些服務的設定步驟寫成檔案。」

---

## Part 2：Docker Compose 多容器編排

### 2.1 為什麼需要 Docker Compose？

Ocean 現在會跑容器了，但每次要啟動服務都要手動打一串 `docker run`，參數又臭又長。更慘的是，上次好不容易把 Frontend + Backend + DB 三個容器都跑起來了，結果 Snow 問他：「指令記在哪？」Ocean：「...我記在腦袋裡。」Snow：「喔不。」

**Docker Compose** 是 Docker 官方提供的多容器編排工具，透過一個 YAML 設定檔定義所有容器的 configuration（包含映像檔、連接埠、Volume、環境變數、相依關係），再以單一指令 `docker compose up` 完成整個應用的部署。如果你使用 OrbStack，Docker Compose 已內建其中，並且可以透過 `docker compose version` 確認正在運行中的版本。

實際專案中，一個服務通常由多個容器組成。Docker Compose 將這些容器的設定與啟動流程集中定義在 `docker-compose.yml`（或 `compose.yml`）中，納入版本控制後，團隊成員只需 clone 專案並執行 `docker compose up`，即可在本機重現完整的服務環境。

```
┌──────────────────────────────────────────────────────┐
│                    典型的應用架構                      │
├──────────────────────────────────────────────────────┤
│                                                      │
│                   ┌──────────┐                       │
│                   │  Client  │                       │
│                   └────┬─────┘                       │
│                        │                             │
│                        ▼                             │
│              ┌──────────────────┐                    │
│              │  Frontend (React)│  ← Container 1     │
│              │    (Port 3000)   │                    │
│              └────────┬─────────┘                    │
│                       │                              │
│                       ▼                              │
│              ┌──────────────────┐                    │
│              │  Backend (Go API)│  ← Container 2     │
│              │    (Port 8080)   │                    │
│              └────────┬─────────┘                    │
│                       │                              │
│                       ▼                              │
│              ┌──────────────────┐                    │
│              │   PostgreSQL DB  │  ← Container 3     │
│              │    (Port 5432)   │                    │
│              └──────────────────┘                    │
│                                                      │
└──────────────────────────────────────────────────────┘

```

如果不用 Compose，你就得手動打一堆指令：

```bash
# 建立網路（後面會說明）
docker network create myapp

# 啟動 PostgreSQL
docker run -d --name db --network myapp \
  -e POSTGRES_PASSWORD=secret \
  -v pgdata:/var/lib/postgresql/data \
  postgres:17-alpine

# 啟動 Backend（Go API）
docker run -d --name api --network myapp \
  -p 8080:8080 \
  -e DATABASE_URL=postgres://postgres:secret@db:5432/mydb \
  my-go-app:v1

# 啟動 Frontend（React）
docker run -d --name frontend --network myapp \
  -p 3000:3000 \
  my-frontend:v1
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

  api:
    image: my-go-app:v1
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://postgres:secret@db:5432/mydb
    depends_on:
      - db

  frontend:
    image: my-frontend:v1
    ports:
      - "3000:3000"
    depends_on:
      - api

volumes:
  pgdata:
```

一個檔案就把三個容器的設定全部定義好了。來看看每個欄位在做什麼：

#### services

最上層的 `services` 區塊定義了所有要跑的容器。每個 key（`db`、`api`、`frontend`）就是一個服務名稱，同時也是這個容器在 Compose 網路中的 hostname，其他容器可以直接用這個名稱來連線。

#### image

```yaml
image: postgres:17-alpine
```

指定要用哪個映像檔。跟 `docker run postgres:17-alpine` 一樣的意思。

#### ports

```yaml
ports:
  - "8080:8080"
```

Port Mapping，格式是 `"宿主機:容器"`。等同於 `docker run -p 8080:8080`。

#### environment

```yaml
environment:
  POSTGRES_PASSWORD: secret
  DATABASE_URL: postgres://postgres:secret@db:5432/mydb
```

設定環境變數。等同於 `docker run -e POSTGRES_PASSWORD=secret`。注意 `DATABASE_URL` 裡面的 `@db` 是用服務名稱當 hostname，Compose 會自動幫你做 DNS 解析。

#### volumes

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

#### depends_on

```yaml
depends_on:
  - db
```

告訴 Compose 啟動順序：先啟動 `db`，再啟動 `api`。不過要注意，`depends_on` 只等容器啟動，不等服務真正準備好。如果需要等資料庫可以接受連線，要搭配 `healthcheck`（後面會講）。

#### 啟動與停止

```bash
# 啟動所有服務（背景執行）
docker compose up -d

# 看看服務狀態
docker compose ps

# 停止並移除所有容器
docker compose down
```

對比一下：原本要打三串 `docker run` 加一個 `docker network create`，現在一個 `docker compose up -d` 就搞定了。

---

### 2.2 服務定義詳解

2.1 介紹了 Compose 最基本的欄位。這一節補充幾個實務上常用的進階設定。

#### depends_on 搭配 healthcheck

在 2.1 裡我們用了簡單的 `depends_on`，但它只保證「容器有啟動」，不保證「服務可以連線」。實際上 PostgreSQL 容器啟動後還需要幾秒初始化，這段時間 API 連過去會直接噴錯。

解法是搭配 `healthcheck` 和 `condition: service_healthy`：

```yaml
services:
  api:
    image: my-go-app:v1
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

`pg_isready` 是 PostgreSQL 內建的檢查工具，Compose 每 5 秒問一次「資料庫好了沒」，連續失敗 5 次才判定不健康。API 會等到 db 的 healthcheck 通過才啟動。

#### restart — 重啟策略

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
  api:
    restart: unless-stopped
```

---

### 2.3 Network

在 1.8 我們學到，容器的網路是隔離的，外面連不進去要靠 Port Mapping。但反過來，容器跟容器之間呢？預設也是不通的。

要讓多個容器能互相溝通，就需要把它們放進同一個 **Docker Network**（虛擬網路）。你可以把它想成一條虛擬的網路線，只有接在同一條線上的容器才能互相找到對方。

2.1 的 `docker run` 範例裡，我們手動跑了 `docker network create myapp`，然後每個容器都加上 `--network myapp` 才能互連。但 Compose 版本完全沒寫到網路，為什麼容器之間還是能溝通？

因為 **Docker Compose 會自動建立一個預設網路**，把 `docker-compose.yml` 裡定義的所有服務都放進去。你不用寫任何網路設定，容器之間就能互相連線。

更方便的是，在這個網路裡面，**服務名稱就等於 hostname**。回頭看 2.1 的範例：

```yaml
environment:
  DATABASE_URL: postgres://postgres:secret@db:5432/mydb
```

這裡的 `@db` 不是什麼特殊語法，就是用 `docker-compose.yml` 裡定義的服務名稱 `db` 當作 hostname。Compose 內建的 DNS 會自動把 `db` 解析到 PostgreSQL 容器的 IP。如果你的服務改叫 `database`，這裡就要改成 `@database:5432`。

> **注意**：這個 DNS 只在 Compose 的虛擬網路裡面有效。你在宿主機上打 `ping db` 是不會通的。

---

### 2.4 Volume

1.9 介紹過 Volume 的概念：容器砍掉，可寫層就沒了，所以需要持久化的資料要掛出來。在 Compose 裡面，Volume 的寫法跟 `docker run -v` 一樣，只是搬進了 YAML。

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

除了 Named Volume，你也可以用 Bind Mount 把宿主機的檔案直接掛進去：

```yaml
volumes:
  - ./init.sql:/docker-entrypoint-initdb.d/init.sql:ro
```

這行把本機目前目錄下的 `init.sql` 掛進容器，`:ro` 代表容器只能讀不能改。PostgreSQL 的映像檔會在第一次啟動時自動執行 `/docker-entrypoint-initdb.d/` 裡的 SQL 檔案，很適合拿來做資料庫初始化。

> **注意**：`docker compose down` 不會刪除 Volume。要連 Volume 一起清掉要加 `-v`：`docker compose down -v`。

---

### 2.5 環境變數管理

2.1 的範例裡，我們把密碼直接寫在 `docker-compose.yml` 裡面：

```yaml
environment:
  POSTGRES_PASSWORD: secret
```

開發時這樣寫沒問題，但如果要把 `docker-compose.yml` 推進 Git，密碼就跟著推上去了。實際專案會把敏感資訊抽到一個 `.env` 檔案裡，然後在 `.gitignore` 排除它。

#### 做法：用 .env 檔案管理

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
  api:
    env_file:
      - .env
```

> **重要**：記得把 `.env` 加進 `.gitignore`，永遠不要把密碼推到 Git 上。

---

### 2.6 常用 Compose 指令

前面幾節學了怎麼寫 `docker-compose.yml`，這一節整理日常操作 Compose 時最常用的指令。所有指令都是 `docker compose` 開頭，後面接動作。

#### 啟動與停止

你最常打的大概就這兩個：

```bash
docker compose up -d      # 啟動所有服務（背景執行）
docker compose down        # 停止並移除所有容器和網路
```

`down` 不會刪除 Volume（2.4 有提到），所以資料庫資料還在。如果你想連 Volume 一起清掉：

```bash
docker compose down -v     # 停止、移除容器、刪除 Volume
```

#### 查看狀態與日誌

服務跑起來之後，要怎麼知道它是不是正常的？

```bash
docker compose ps          # 列出所有服務的狀態
docker compose logs api    # 看 api 服務的日誌
docker compose logs -f api # 持續追蹤日誌（Ctrl+C 停止）
```

#### 進入容器與執行命令

跟 Part 1 學的 `docker exec` 一樣，只是改成用服務名稱：

```bash
docker compose exec api sh           # 進入 api 容器的 shell
docker compose run --rm api go test ./...  # 開一個新容器跑一次性命令
```

`exec` 是對已經在跑的容器執行命令，`run` 是開一個全新的容器跑完就砍掉（`--rm`）。

#### 速查表

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

回顧一下 Part 2 學到的東西：

- 把多個 `docker run` 寫成一個 `docker-compose.yml`，一個指令啟動所有服務（2.1）
- 用 `depends_on` + `healthcheck` 控制啟動順序，確保資料庫準備好了 API 才啟動（2.2）
- Compose 自動建立虛擬網路，容器之間用服務名稱當 hostname 互連（2.3）
- Volume 讓資料跨容器生命週期保留，`docker compose down` 不會刪 Volume（2.4）
- 敏感資訊抽到 `.env`，不推進 Git（2.5）
- 日常操作靠 `up`、`down`、`ps`、`logs`、`exec` 這幾個指令就夠了（2.6）

Ocean 現在可以用一個 YAML 檔把 Frontend + Backend + DB 全部跑起來，而且任何人 clone 專案後都能用同樣的方式啟動。Andrew 終於不用擔心 Ocean 把指令記在腦袋裡了。

但 Andrew 馬上丟了一個新問題：「現在你跑的映像檔都是別人做好的，如果要把我們自己的 Go 程式打包成映像檔呢？」

---

## Part 3：Dockerfile 基礎

### 3.1 什麼是 Dockerfile？

Ocean 用 Compose 把服務跑得有聲有色，但有一天 Andrew 推了新的 Go 程式碼，Ocean 發現 Docker Hub 上根本沒有現成的映像檔可以用。Andrew：「別人的映像檔當然沒有我們的程式碼。要把自己的程式打包成映像檔，就要寫 Dockerfile。」

到目前為止，我們用的映像檔（`nginx`、`postgres`）都是別人做好放在 Docker Hub 上的。但如果你要跑的是自己寫的程式，就需要自己做一個映像檔。

做映像檔的方式就是寫一個 **Dockerfile**。它是一個純文字檔，放在專案的根目錄，裡面一行一行寫著建置映像檔的步驟，像是：用哪個基礎環境、複製哪些檔案進去、執行什麼安裝指令、程式怎麼啟動。

你可以把 Dockerfile 想成是一份「環境建置的 SOP」。手動裝環境的時候，你可能會：

1. 找一台乾淨的 Linux 機器
2. 安裝 Go 語言
3. 把程式碼複製過去
4. 編譯程式
5. 設定啟動指令

Dockerfile 做的就是把這些步驟寫下來，讓 Docker 自動執行。寫好之後跑一個 `docker build` 指令，Docker 就會照著步驟一步一步做，最後產出一個映像檔：

```
Dockerfile ─── docker build ──→ Docker Image ─── docker run ──→ Container
（建置腳本）      （建置過程）      （成品）          （執行）       （跑起來的容器）
```

這個映像檔可以推到 Docker Hub 或私有 Registry，讓其他人拉下來直接用，不需要重新建置。

---

### 3.2 Dockerfile 基礎指令

假設你有一個最簡單的 Go 專案，結構長這樣：

```
my-app/
├── main.go         # Go 程式碼
├── go.mod          # Go 模組定義
└── Dockerfile      # 建置映像檔的腳本
```

`main.go` 是一個簡單的 HTTP Server，`go.mod` 是 Go 的套件管理檔，而 `Dockerfile` 就是我們要寫的東西。來看看這個 Dockerfile 長什麼樣：

```dockerfile
FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go build -o server .
CMD ["./server"]
```

只有五行，但做了完整的事情。一行一行來看：

**`FROM golang:1.24-alpine`** — 起點。每個 Dockerfile 都從 `FROM` 開始，指定你要在什麼基礎環境上建置。這裡用的是已經裝好 Go 1.24 的 Alpine Linux。

**`WORKDIR /app`** — 設定映像檔內部的工作目錄。後面的 `COPY`、`RUN` 都會在 `/app` 這個目錄下執行。

**`COPY . .`** — 把你本機目前目錄下的所有檔案複製到映像檔的 `/app` 裡面。第一個 `.` 是本機路徑，第二個 `.` 是容器內的路徑（也就是 `WORKDIR` 設定的 `/app`）。

**`RUN go build -o server .`** — 在建置過程中執行命令。這裡是編譯 Go 程式，產出一個叫 `server` 的執行檔。`RUN` 的結果會被寫進映像檔的 Layer 裡。

**`CMD ["./server"]`** — 容器啟動時的預設命令。當你 `docker run` 這個映像檔的時候，它就會執行 `./server`。

建置並執行：

```bash
docker build -t my-app .     # 用目前目錄的 Dockerfile 建置映像檔，取名叫 my-app
docker run -p 8080:8080 my-app   # 跑起來
```

#### 指令速查表

上面的範例只用了五個指令，但 Dockerfile 還有其他常用指令。完整列表如下：

| 指令 | 說明 | 範例 |
|------|------|------|
| `FROM` | 指定基礎映像檔（必須是第一行） | `FROM golang:1.24-alpine` |
| `WORKDIR` | 設定工作目錄 | `WORKDIR /app` |
| `COPY` | 複製檔案至映像檔 | `COPY . .` |
| `RUN` | 建置時執行命令 | `RUN go build -o server .` |
| `CMD` | 容器啟動時的預設命令 | `CMD ["./server"]` |
| `ENTRYPOINT` | 容器啟動時的固定入口 | `ENTRYPOINT ["./server"]` |
| `EXPOSE` | 宣告容器監聽的連接埠（僅文件用途） | `EXPOSE 8080` |
| `ENV` | 設定環境變數 | `ENV GIN_MODE=release` |
| `USER` | 指定執行命令的使用者 | `USER nonroot` |
| `HEALTHCHECK` | 定義健康檢查 | `HEALTHCHECK CMD curl -f http://localhost/` |

---

### 3.3 為 Go 應用撰寫 Dockerfile

先來看一個簡單的 Go HTTP 伺服器範例：

**main.go**

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		fmt.Fprintf(w, "Hello from Go! Hostname: %s\n", hostname)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

#### 基本版 Dockerfile（單階段建置）

```dockerfile
FROM golang:1.24-alpine

WORKDIR /app

# 先複製 go.mod, go.sum 以利用 Layer Cache
COPY go.mod go.sum ./
RUN go mod download

# 複製原始碼並編譯
COPY . .
RUN go build -o server .

EXPOSE 8080
CMD ["./server"]
```

**問題在哪？** 最終映像檔裡面包含了整個 Go 編譯工具鏈（~300MB+），但跑起來的時候根本只需要那個編譯好的二進位檔。

#### Multi-stage Build

Ocean 照著上面的方式寫了第一個 Dockerfile，build 出來一看，映像檔 300 多 MB。Andrew 看到差點把咖啡噴出來：「你要拿這個上正式環境？裡面一整套 Go 編譯工具鏈都打包進去了啊。」Ocean：「那要怎麼辦？」Andrew：「用 Multi-stage Build，build 完只留下執行檔就好。」

Multi-stage Build（多階段建置）的概念很簡單：**建置階段**用完整的工具鏈編譯程式，**最終映像檔**只放跑起來需要的東西。

```dockerfile
# ==========================================
# Stage 1: Build（建置階段）
# ==========================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o server .

# ==========================================
# Stage 2: Runtime（執行階段）
# ==========================================
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

WORKDIR /app

# 從 builder 階段複製編譯好的二進位檔
COPY --from=builder /app/server .

USER appuser
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./server"]
```

**映像檔大小比較：**

| 方式 | 映像檔大小 |
|------|-----------|
| 單階段 `golang:1.24-alpine` | ~300MB |
| 多階段 `alpine:3.21` | ~15MB |
| 多階段 `scratch`（最精簡） | ~8MB |

#### 建置與執行

```bash
# 建置映像檔
docker build -t my-go-app:v1 .

# 查看映像檔大小
docker images my-go-app

# 執行容器
docker run -d --name my-app -p 8080:8080 my-go-app:v1

# 測試
curl http://localhost:8080
curl http://localhost:8080/health
```

---

## Part 4：綜合演練與延伸學習

### 4.1 SDC 的短網址服務 Shlink

前面三個 Part 我們學了容器怎麼跑、映像檔怎麼建、Compose 怎麼把多個服務串起來。現在來看一個真的跑在 SDC 基礎設施上的服務——**Shlink 短網址服務**

#### 什麼是 Shlink？

Shlink 是一個開源的短網址（URL Shortener）服務。簡單來說，就是把又臭又長的網址變短：

```
原始：https://docs.google.com/spreadsheets/d/1aBcDeFgHiJkLmNoPqRsTuVwXyZ/edit#gid=0
縮短：https://link.sdc.tw/abc123
```

為什麼我們自己架而不用 bit.ly？因為：

- **自訂網域**：用 `link.sdc.tw` 比較有辨識度
- **數據自己掌控**：點擊次數、來源分析都在自己手上
- **不受第三方限制**：免費方案有額度限制，自架沒有這個問題

#### Shlink 的 Docker Compose 架構

Shlink 整個服務由三個容器組成，架構很乾淨：

```
┌─────────────────────────────────────────────────────┐
│                    Traefik (反向代理)                  │
│              處理 HTTPS、憑證、路由分流                   │
└──────────┬──────────────────────┬────────────────────┘
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

#### 解析 docker-compose.yaml

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

- **`depends_on` + `condition: service_healthy`**：等資料庫的 healthcheck 通過才啟動。不是只等容器跑起來，而是等到資料庫「真的準備好了」。這就是 Part 3 學的 `depends_on` 進階用法
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

#### 整體運作流程

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

### 4.2 常見問題排查

某天半夜，Ocean 被 Andrew 的訊息吵醒：「服務掛了你快看一下！」Ocean 睡眼惺忪地打開電腦，容器起不來，log 一片紅字。這種時候不要慌，先照著以下的排查步驟走。

#### 容器啟動失敗

```bash
# 查看容器日誌
docker logs <container-name>

# 查看容器詳細資訊
docker inspect <container-name>

# 檢查容器退出碼
docker inspect --format='{{.State.ExitCode}}' <container-name>
# 0   = 正常退出
# 1   = 應用程式錯誤
# 137 = OOM Killed 或 SIGKILL
# 143 = SIGTERM（正常停止）
```

#### 容器間無法連線

```bash
# 確認容器在同一網路
docker network inspect <network-name>

# 進入容器測試連線
docker exec -it api sh
# 在容器內執行：
#   ping db
#   wget -qO- http://api:8080/health
```

#### 磁碟空間不足

```bash
# 查看 Docker 佔用空間
docker system df

# 清理所有未使用的資源（映像檔、容器、網路、Volume）
docker system prune -a --volumes
# ⚠ 注意：此操作將刪除所有未使用的資料，請謹慎執行
```

#### 映像檔建置快取問題

```bash
# 強制不使用快取重新建置
docker build --no-cache -t my-app .

# Compose 中強制重新建置
docker compose build --no-cache
docker compose up -d --build --force-recreate
```

#### 常見錯誤訊息

| 錯誤訊息 | 原因 | 解決方式 |
|---------|------|---------|
| `port is already allocated` | 宿主機連接埠被佔用 | 更換連接埠或停止佔用的程序 |
| `no space left on device` | 磁碟空間不足 | 執行 `docker system prune` |
| `network not found` | 網路不存在 | 檢查 `docker-compose.yml` 的 `networks` 設定 |
| `exec format error` | 映像檔平台不符 | 加上 `--platform linux/amd64` 重新建置 |
| `permission denied` | 權限不足 | 檢查 `USER` 指令或目錄權限 |
| `connection refused` | 服務尚未就緒 | 新增 `depends_on` + `healthcheck` |

---

### 4.3 延伸學習資源

#### 官方文件

- [Docker 官方文件](https://docs.docker.com/)
- [Dockerfile 最佳實踐](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [Docker Compose 規格](https://docs.docker.com/compose/compose-file/)

#### 進階主題

| 主題 | 說明 |
|------|------|
| **Docker Network 進階** | Bridge、Host、Overlay 等網路模式 |
| **Docker Security** | 映像檔掃描、Rootless 模式、Secret 管理 |
| **Container Registry** | Harbor、ECR、GCR 等私有映像檔倉庫 |
| **Container Orchestration** | Kubernetes、Docker Swarm |
| **CI/CD 整合** | GitHub Actions、GitLab CI 中使用 Docker |
| **多平台建置** | `docker buildx` 建置 ARM/AMD64 映像檔 |
| **日誌管理** | 日誌驅動程式、集中式日誌收集 |
| **監控** | Prometheus + Grafana 監控容器 |

---

## 附錄 A：Docker 指令速查表

### 映像檔 (Image)

```bash
docker pull <image>:<tag>           # 下載映像檔
docker images                        # 列出映像檔
docker rmi <image>                   # 刪除映像檔
docker build -t <name>:<tag> .       # 建置映像檔
docker tag <src> <dest>              # 標記映像檔
docker push <image>:<tag>            # 推送映像檔
docker save -o file.tar <image>      # 匯出映像檔
docker load -i file.tar              # 匯入映像檔
docker history <image>               # 查看建置歷史
docker inspect <image>               # 查看詳細資訊
docker image prune                   # 清理 dangling 映像檔
docker image prune -a                # 清理所有未使用映像檔
```

### 容器 (Container)

```bash
docker run [opts] <image> [cmd]      # 建立並啟動容器
docker create [opts] <image> [cmd]   # 建立容器（不啟動）
docker start <container>             # 啟動容器
docker stop <container>              # 停止容器
docker restart <container>           # 重啟容器
docker kill <container>              # 強制停止
docker rm <container>                # 刪除容器
docker ps                            # 列出執行中容器
docker ps -a                         # 列出所有容器
docker logs [-f] <container>         # 查看日誌
docker exec -it <container> <cmd>    # 在容器內執行命令
docker cp <src> <container>:<dest>   # 複製檔案
docker inspect <container>           # 查看詳細資訊
docker stats                         # 即時資源使用
docker top <container>               # 查看程序
docker port <container>              # 查看連接埠對映
docker container prune               # 清理已停止容器
```

### 資料卷 (Volume)

```bash
docker volume create <name>          # 建立資料卷
docker volume ls                     # 列出資料卷
docker volume inspect <name>         # 查看詳細資訊
docker volume rm <name>              # 刪除資料卷
docker volume prune                  # 清理未使用資料卷
```

### 網路 (Network)

```bash
docker network create <name>         # 建立網路
docker network ls                    # 列出網路
docker network inspect <name>        # 查看詳細資訊
docker network rm <name>             # 刪除網路
docker network connect <net> <ctn>   # 將容器加入網路
docker network disconnect <net> <c>  # 將容器移出網路
docker network prune                 # 清理未使用網路
```

### 系統

```bash
docker system df                     # 磁碟使用量
docker system prune                  # 清理未使用資源
docker system prune -a --volumes     # 深度清理（含映像檔和資料卷）
docker info                          # 系統資訊
docker version                       # 版本資訊
```

---

## 附錄 B：Dockerfile 指令速查表

```dockerfile
# Base image
FROM <image>:<tag> [AS <name>]

# Working directory
WORKDIR /path

# Copy files
COPY [--from=<stage>] <src> <dest>
ADD <src> <dest>

# Execute commands (build stage)
RUN <command>

# Environment variables
ENV KEY=value
ARG KEY=default

# Declare port
EXPOSE <port>

# Container startup command
CMD ["executable", "param1"]
ENTRYPOINT ["executable", "param1"]

# User
USER <user>[:<group>]

# Health check
HEALTHCHECK [OPTIONS] CMD <command>

# Metadata
LABEL key="value"

# Stop signal
STOPSIGNAL signal

# Shell
SHELL ["executable", "parameters"]
```

---

## 附錄 C：Docker Compose 速查表

```bash
docker compose up -d                 # 啟動服務（背景）
docker compose up -d --build         # 啟動並重新建置
docker compose down                  # 停止並移除
docker compose down -v               # 停止、移除、刪除 Volume
docker compose ps                    # 服務狀態
docker compose logs [-f] [service]   # 服務日誌
docker compose exec <svc> <cmd>      # 在服務內執行命令
docker compose run --rm <svc> <cmd>  # 一次性命令
docker compose build [--no-cache]    # 建置映像檔
docker compose pull                  # 拉取映像檔
docker compose restart [service]     # 重啟服務
docker compose config                # 驗證設定
docker compose top                   # 查看服務程序
```