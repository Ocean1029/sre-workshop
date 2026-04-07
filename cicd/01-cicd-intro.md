# 01 — CI/CD 概念介紹

角色設定：

- **Ocean** — SDC SRE Team 新進成員，負責維護社團的基礎設施和部署流程，還在摸索中，常常搞出各種狀況。
- **Andrew** — 整天在喝酒的 SDC 社長，熟悉社團大小事。
- **Snow** — SDC 開發部長，明明很聰明但常常把事情丟給 Ocean 做，口頭禪是「歐不」。

## Table of Contents

- [軟體開發的多人協作](#軟體開發的多人協作)
- [什麼是 CI/CD？](#什麼是-cicd)
- [CI/CD Pipeline 流程圖](#cicd-pipeline-流程圖)
- [常見的 CI/CD 工具](#常見的-cicd-工具)
- [小結](#小結)

## 軟體開發的多人協作

Ocean 最近在跟 Andrew 開發一個新專案「Special Ground」，兩個人一起在同一個 GitHub repository 上面做軟體開發。兩個人都直接在 `main` 分支（branch）上推程式碼，有一天 Ocean 推了一版上去，結果把 Andrew 前一天寫好的功能整個蓋掉了，上線中的服務也跟著壞掉。兩個人想把版本 revert 回去，結果發現 git history 已經亂成一鍋粥，用威士忌熬的那種，搞了半天也恢復不了原本的版本。

Snow 看到他們的窘境，決定教他們怎麼好好在 GitHub 上面用版本控制做多人協作，順便介紹一個叫做 CI/CD 的好東西。

Snow 第一件事就是叫他們不要再把程式直接推到 `main` 上面，而是各自開一條分支開發，寫到一個段落之後發 **[Pull Request（PR）](https://shoujhengduan.medium.com/%E4%BB%80%E9%BA%BC%E6%98%AF-pull-request-b476ee3e0217)**。PR 是 GitHub 上的一個功能，用途是告訴其他人「我這邊改好了，請來看一下」。其他人可以在 PR 上看到你改了哪些程式碼、留下討論、提出修改建議。如果需要調整，就繼續在同一條分支上 push 新的 commit，PR 則會自動更新。等大家都確認沒問題了，才由維護者把這條分支合併進 `main`。Snow 自己則負責盯著 `main` 的狀態，確保上面的程式碼是乾淨的。

![Pull Request 頁面](assets/github-pull-request.png)

上圖是一個實際的 PR 頁面。標題寫了這次改了什麼，下面標示著要從 `fix/update-infisical-env-default` 這條分支合併到 `main`。右邊可以指定 **Reviewer**（請誰來看）和 **Assignee**（誰負責這個 PR）。中間的 tab 可以切換看 **Conversation**（討論）、**Commits**（這個 PR 包含哪些 commit）、**Checks**（CI 自動檢查的結果，後面會詳細講）、**Files changed**（改了哪些檔案的哪些地方）。

這樣做之後好多了，但新的問題又出現了：分支要怎麼命名？什麼時候該從哪裡開？合併的順序是什麼？三個人各自有各自的習慣，分支開得越來越亂。Snow 說：「我們來了解一下 **Git Flow** 吧。」

Git Flow 是 Vincent Driessen 在 2010 年提出的分支管理規範，它定義了每條分支的角色與合併規則，讓多人協作時有明確的流程可以遵循，而不是各自亂開分支、合併時再來收拾殘局。

![Git Flow 分支模型示意圖](assets/git-flow-diagram.png)

**主幹分支（Long-lived branches）**

- `main`（同 `master`）永遠反映正式環境（production）的狀態。原則上，`main` 上的任何一個 commit 都代表一個可以安全部署的版本，並且會打上版本 tag（例如 `v1.0.0`）。開發者不應直接在 `main` 上提交程式碼。

- `develop` 是所有開發工作的整合中心，代表「下一個版本」的最新進度。開發者不直接在這條分支上寫功能，而是透過其他短命分支合併進來。這條分支的穩定程度介於功能分支與 `main` 之間。

**短命分支（Short-lived branches）**

- `feature/*` 從 `develop` 開出，用來開發一個獨立的功能，例如 `feature/user-login` 或 `feature/payment-api`。功能完成、通過 code review 後，合併回 `develop` 並刪除該分支。

- `release/*` 從 `develop` 開出，用來準備新版本的發布。這個階段只做最後的修正（例如修小 bug、更新版本號），而不加入新的功能。完成後同時合併到 `main` 與 `develop` 上面，確保修正也回流到開發線上。

- 另外還有 `hotfix/*` 從 `main` 開出，用來修正正式環境緊急的 bug，修完後同時合併回 `main` 和 `develop`。這裡不展開細節。

順帶一提，SDC 目前的專案（例如 [Core System backend](https://github.com/NYCU-SDC/core-system-backend)）實務上採用的是更簡化的模式：`main` 分支直接作為長期穩定分支，開發者每次增加新功能就從 `main` 開一條 `feature/*`，完成後發 PR 合回 `main`，沒有獨立的 `develop` 分支。這種模式一般稱為 **GitHub Flow**，比 Git Flow 輕量，適合持續部署的小團隊專案。

有了分支管理和 PR 之後，Ocean 和 Andrew 終於不會再互相蓋掉對方的程式碼了。但 Snow 發現另一個問題：每次有人發 PR，都得靠人記得去跑測試、檢查程式碼風格、確認編譯有沒有過。有時候 Ocean 忘了跑測試就按了合併，結果壞掉的程式碼又進到 `main` 裡面。Snow 說：「協作的流程有了，但這些檢查不應該靠人記得去做，我們來導入 **CI/CD** 吧。」

## 什麼是 CI/CD？

**CI/CD** 的全名是 Continuous Integration / Continuous Deployment（持續整合 / 持續部署）。做法是把編譯、測試、部署這些原本靠人手動做的步驟，定義成一條自動化的管線（pipeline）：程式碼一推上去，CI 自動跑檢查確保品質，通過之後 CD 自動部署到對應的環境上。

CI/CD 的核心目標是「維持程式品質」和縮短「寫完程式碼」到「使用者看到新功能」的時間。

品質的部分，每次程式碼變更都會自動跑完整的測試和檢查，每個人的程式碼都經過同一套標準，不會因為本機環境不一樣而出現「在我電腦上可以跑啊」的問題。部署流程也由 pipeline 定義好，每次步驟都一樣，不會因為人為操作漏掉關鍵步驟。

速度的部分，重複性的工作全部交給機器做，工程師把時間和腦力留給產品迭代本身。整個流程自動化之後，團隊也可以更頻繁地發版，每次變更更小、風險更低、出問題也更容易定位。

接下來分別看 CI 和 CD 各自在做什麼。

### CI（Continuous Integration）

CI 自動觸發的檢查通常包含這三項：

- **建置（Build）**：嘗試把程式碼編譯起來。如果編譯失敗，代表程式碼有語法錯誤或是缺少相依套件，連跑都跑不起來。
- **測試（Test）**：執行事先寫好的自動化測試（例如單元測試），確認功能邏輯是正確的，不會因為改了 A 功能結果 B 功能壞了。有些團隊還會設定覆蓋率門檻，比方說測試覆蓋率不能低於 80%。
- **風格檢查（Lint）**：確保大家寫出來的程式碼風格一致、沒有常見的 bad practice。Go 常用的是 golangci-lint，除了風格之外也能抓出潛在的 bug，例如未使用的變數、不安全的型別轉換等等。

這些檢查全部自動執行，開發者不需要記得「提交前要先跑哪些步驟」。如果任何一步失敗，系統會直接擋住合併，開發者收到通知後去修就好。

### CD（Continuous Deployment）

**部署（Deployment）** 就是把寫好的程式搬到實際運行的環境上，讓使用者可以用到新版本。CI 確保程式碼品質之後，CD 負責自動完成這件事：把通過檢查的程式碼打包成可以在伺服器上跑的格式（像是 Docker image），部署到對應的環境上，再執行後續的處理（像是通知團隊、更新監控）。從 push 程式碼到使用者看到新功能，整個過程全自動，不需要人工介入。

這種模式要求你對自動化測試有非常高的信心，因為沒有人會在中間攔一下，測試就是你唯一的安全網。

### SDC 實際怎麼做？

SDC 現行的專案（例如 [Core System backend](https://github.com/NYCU-SDC/core-system-backend)、[Clustron backend](https://github.com/NYCU-SDC/clustron-backend)）都具備完整的 CI/CD pipeline，並且他們在不同的環境（正式版 / 測試版）會經過不同的 pipeline 流程。

以 [Core System backend](https://github.com/NYCU-SDC/core-system-backend) 為例，CI 的部分使用 **GitHub Actions**，在每次有程式變更時都會跑一次 lint、test、build 等流程。開發者在處理手上的任務時，當流程全部正常通過，這個 PR 才可以被 merge。

CD 的部分則使用由 SRE Team 維護的 n8n deployment workflow 做自動化部署。這個 workflow 會在 CI 通過後自動部署一個 Snapshot 版本的服務，並建立一個域名（例如 [https://pr-153.core-system.sdc.nycu.club/api/healthz](https://pr-153.core-system.sdc.nycu.club/api/healthz)）讓人直接測試（你需要開啟 SDC VPN 才可以連結進去，點進去後會看到伺服器回 ok）。

回頭對照一下，Ocean 遇到的那些問題，在導入 CI/CD 之後是怎麼被解決的：

| 沒有 CI/CD 的時候 | 有了 CI/CD 之後 |
|------------------|----------------|
| Ocean 忘記跑測試就合併，壞掉的程式碼進到 `main` | 每次合併自動執行測試，沒通過就無法合併 |
| 本機測試環境不一樣，常出現「在我電腦上可以跑啊」的問題 | 每個人的程式碼都在同一套 runner 環境下跑過檢查，品質標準一致 |
| revert 不回去，git history 一團亂 | 每次變更都很小、都經過驗證，出問題也容易定位和復原 |

## CI/CD Pipeline 流程圖

一個典型的 CI/CD pipeline 長這樣：

```text
           ┌──────┐
        ┌─▶│ Lint │─┐
┌──────┐│  └──────┘ │  ┌───────┐   ┌─────────┐   ┌──────────┐   ┌─────────────┐
│ Code ├┤           ├─▶│ Build │──▶│ Package │──▶│  Deploy  │──▶│   Deploy    │
│  PR  ││  ┌──────┐ │  └───────┘   └─────────┘   │ (Staging)│   │(Production) │
└──────┘└─▶│ Test │─┘                            └──────────┘   └─────────────┘
           └──────┘
```

### 各階段說明

| 階段 | 做什麼 | 目的 |
|------|--------|------|
| **Code** | 開發者撰寫程式碼並推送到 Git | 觸發整個 pipeline |
| **Lint** | 靜態程式碼分析、風格檢查 | 確保程式碼品質與一致性 |
| **Test** | 執行單元測試、整合測試 | 確保功能正確、沒有 regression |
| **Build** | 編譯程式碼 | 確保程式碼可以成功編譯 |
| **Package** | 打包成可部署的成品（如 Docker image） | 產出可部署的 artifact |
| **Deploy (Staging)** | 部署到測試環境 | 讓 QA 或團隊成員驗證 |
| **Deploy (Production)** | 部署到正式環境 | 提供給終端使用者 |

## 常見的 CI/CD 工具

市面上有許多 CI/CD 工具，以下是最常見的幾個：

- **GitHub Actions** — 本課程的重點！與 GitHub 深度整合，免費且功能強大
- **GitLab CI/CD** — GitLab 內建的 CI/CD 平台，與 GitLab 生態系緊密結合
- **Jenkins** — 老牌的開源 CI/CD 工具，高度可客製化，但設定較複雜
- **CircleCI** — 雲端 CI/CD 服務，以速度與易用性著稱
- **Travis CI** — 早期非常流行的 CI 服務，對開源專案友善

## 小結

- **CI/CD** 是一套自動化實踐，讓程式碼從開發到部署的每個階段都能自動執行，**CI** 著重在頻繁合併與自動測試，**CD** 著重在自動化部署流程。
- CI/CD 能解決手動流程的痛點：忘記測試、環境不一致、部署出錯。
- **GitHub Actions** 是我們今天要學的 CI/CD 工具，與 GitHub 無縫整合。

> 下一章我們就來動手，寫出你的第一個 GitHub Actions workflow。

[下一章：GitHub Actions 基礎 →](02-github-actions-basics.md)
