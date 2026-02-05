# ⬢ DHM - Docker Health Monitor (v1.0.0)


A lightweight terminal-based (TUI) tool written in **Go 1.24** for monitoring Docker containers, visualizing their dependencies, and managing their lifecycle.

![License](https://img.shields.io)
![Go Version](https://img.shields.io)

## 🚀 Features

- **Dependency Graph**: Automatically detects relations between containers via environment variables.
- **Real-time Monitoring**: Instant status updates (Running/Exited/Restarting).
- **Log Streamer**: View container logs with timestamps and keyword highlighting (ERROR/WARN).
- **Quick Actions**: One-key container restart ('R').
- **Health-First Sorting**: Failed containers are automatically moved to the top.
- **Minimalist UI**: Clean, single-binary tool with intuitive navigation.

## 🛠 Installation

### Prerequisites
- Docker installed and running.
- Go 1.24+ (only for building from source).

### Build from source
```bash
git clone github.com/ladnix/docker-health-monitor
cd dhm
go build -ldflags="-s -w -X 'main.Version=v1.0.0'" -o dhm .

### Run
./dhm

### ⌨️ Hotkeys

| Key | Action |
| :--- | :--- |
| **Arrows** | Navigate through container tree |
| **Enter** | Show container details (ID, IP, Status) |
| **L** | Open container logs |
| **R** | Restart selected container |
| **F1 / H** | Show help menu |
| **ESC** | Back / Exit |

### 🧬 How it works (Dependencies)
DHM scans container environment variables. If a variable contains the name of another container (e.g., `DB_HOST=postgres-db`), it automatically creates a visual link in the tree.

## 🤝 Author
Created by **ladnix**.

## 📄 License
MIT License