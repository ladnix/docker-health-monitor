# ⬢ DHM - Docker Health Monitor (v1.0.0)

A lightweight terminal-based (TUI) tool written in **Go 1.24** for monitoring Docker containers, visualizing their dependencies, and managing their lifecycle.

![License](https://img.shields.io)
![Go Version](https://img.shields.io)

## 🖥️ Interface Preview

╔═════════════════════  Docker Health Monitor  ════════════════════╗
║⬢ DHM                                                             ║
║├──report-service                                                 ║
║│  └──  🔗 db-postgres                                            ║
║├──backend-api                                                    ║
║│  ├──  🔗 db-postgres                                            ║
║│  └──  🔗 cache-redis                                            ║
║├──cache-redis                                                    ║
║                                                                  ║
╚══════════════════════════════════════════════════════════════════╝

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
cd docker-health-monitor
go build -ldflags="-s -w -X 'main.Version=v1.0.0'" -o dhm .
```

### Run

```bash
./dhm
```

### ⌨️ Hotkeys

| Key | Action |
| :--- | :--- |
| **Arrows** | Navigate through container tree |
| **Enter** | Show container details (ID, IP, Status) |
| **L** | Open container logs |
| **R** | Restart selected container |
| **F1 / H** | Show help menu |
| **ESC** | Back / Exit |


## 🧬 Core Logic

DHM scans container environment variables. If it finds a variable string containing another container's name (e.g., `DB_HOST=postgres-db`), it automatically renders it as a visual dependency link in the tree.

## 🤝 Author

Created by **ladnix**.

## 📄 License

MIT License
