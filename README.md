\# Server Stats — Go Monitoring Tool

A fast, single-binary server monitoring tool written in Go.

Zero dependencies · Works on any Linux server · ~2.5 MB

\### Features

\- CPU / Memory / Disk usage

\- System info + uptime

\- Top 5 processes (CPU \& Memory)

\- Load average (1/5/15 min)

\### Build (Linux binary from any OS)

```bash

git clone https://github.com/YOUR\_USERNAME/server-stats-go.git

cd server-stats-go

GOOS=linux GOARCH=amd64 CGO\_ENABLED=0 go build -ldflags="-s -w" -o server-stats

```

\### Run (on any Linux server)

```bash

chmod +x server-stats

./server-stats

```
