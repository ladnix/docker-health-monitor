# ⬢ DHM - Docker Health Monitor (v1.1.0)

A high-performance terminal dashboard (TUI) for Docker, featuring smart
dependency mapping and real-time resource tracking.

![License](https://img.shields.io)
![Go Version](https://img.shields.io)

## 🖥️ Interface Preview

```text
 ╔═══════════  Docker Health Monitor  ═══════════╗
 ║                                               ║
 ║ DHM                                           ║
 ║ ├──[red]backend-api[-]                        ║
 ║ │  └──  [gray]🔗 db-postgres[-]               ║
 ║ ├──[green]frontend-ui [ 4.2 % | 120.5 MB ][-] ║
 ║ │  ├──  [blue]🔗 backend-api[-]               ║
 ╚═══════════════════════════════════════════════╝
```
---

## 🚀 Key Features (v1.1.0 Update)

Dual Monitoring Modes:
    Lite Mode: High-speed updates (1s) for large-scale environments.
    Full Mode: Deep inspection (3s) including real-time CPU and RAM metrics.
    Smart Dependency Graph: Automatically visualizes links between containers
    by analyzing environment variables.
    Intelligent Status Analysis: Distinguishes between manual stops (Exit 0/137)
    and actual application crashes (Exit != 0).
    Interactive Details: Lock monitoring on a specific container to track its vital
    signs while navigating the tree.
    Health-First Sorting: Critical errors and failed services are automatically
    pinned to the top of the list.
    Log Streamer: Clean view of container logs with severity highlighting
    (ERROR/WARN).

---

## 🛠 Installation

### Prerequisites

- Docker Engine installed and running.
- Go 1.24+ (to build from source).

### Build from source

```bash
git clone https://github.com/ladnix/docker-health-monitor
cd docker-health-monitor
go build -ldflags="-s -w -X 'main.Version=v1.1.0'" -o dhm .
```

### Install Man Page (Linux)

```bash
sudo cp dhm.1 /usr/local/share/man/man1/
sudo mandb
```

### Run

```bash
./dhm
```
---

## ⌨️ Hotkeys

| Key | Action |
| :--- | :--- |
| **Arrows** | Navigate through the container tree |
| **Enter** | **Lock Focus**: Select container for active monitoring in Detailsl |
| **I** | **Toggle Info**: Switch between Lite and Full (CPU/RAM) modes |
| **L** | **Logs**: Open container log viewer |
| **R** | **Restart**: Re-launch the selected container |
| **F1 / H** | **Help**: Show interactive help menu |
| **ESC** | Back / Exit |

---

## 🧬 Core Logic

- **Dependencies**: DHM scans environment variables. If a variable contains
another container's name (e.g., `DB_HOST=postgres-db`), it renders a visual
link `🔗`. Links change color based on the target's health.
- **Exit Codes**: The UI interprets Docker exit codes. Exit **137** (SIGKILL)
is treated as a manual stop (Gray), while Exit **1** is treated as a Failure (Red).
- **Resource Highlighting**: CPU and RAM metrics change color
(White → Yellow → Red) as they hit 40% and 80% thresholds.

---

## 🤝 Author
Created by **ladnix**.

## 📄 License
MIT License