
# Server Stats Go — Advanced Server Monitoring Tool

A fast, beautiful, single-binary server monitoring tool written in Go.  
Perfect for DevOps, SRE, and quick server health checks.

- **Static binary** (~2.6 MB)
- **Zero dependencies** — just copy and run
- Colorful and clean output
- `--json` flag for automation
- Shows CPU, RAM, Disk, Swap, Network, Top processes, Load average

### Sample Output

```bash
$ ./server-stats
╔══════════════════════════════════════════════════════════╗
║            SERVER PERFORMANCE STATISTICS              ║
║          2025-12-04 22:10:15 • prod-server              ║
╚══════════════════════════════════════════════════════════╝

System Information
   Hostname   : prod-server
   OS         : Ubuntu 22.04 LTS
   Kernel     : 5.15.0-124-generic
   Arch       : linux/amd64
   Uptime     : 18d 14h 22m (18.6 days)
   Boot Time  : 2025-11-16 07:48
   Users      : 3 logged in

CPU Usage
   Cores      : 12
   Usage      : 18.7%

Memory & Swap
   RAM Total  : 31.2GB
   RAM Used   : 19.4GB (62.2%)
   RAM Free   : 4.8GB
   Swap Used  : 124MB / 2.0GB (6.1%)

Disk Usage
   Mount        Total    Used    Free   Use%
   ────────────────────────────────────────────────────────
   /            200GB   87.3GB 102.7GB  46.0%
   /var         100GB   62.1GB  32.9GB  65.4%

Network
   Local IP   : 10.10.20.50
   Public IP  : 185.123.45.67

Top 5 Processes by CPU
   PID    %CPU   %MEM  Command
   ──────────────────────────────────────────────
   8421   38.2   11.5  node
   1293   21.1    7.8  python3
   ...

Load Average (1/5/15 min)
   1.42  1.28  1.15

### Build

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o server-stats
```

### Run

```bash
chmod +x server-stats
./server-stats          # Pretty output
./server-stats --json   # JSON output for monitoring tools
```
