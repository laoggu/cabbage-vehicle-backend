# make.ps1  —— 等价于 Makefile 的关键目标
param(
    [string]$Target = "help"
)

$ErrorActionPreference = "Stop"

function Proto {
    Write-Host ">>> buf generate" -ForegroundColor Cyan
    buf generate
}

function Gqlgen {
    Write-Host ">>> gqlgen generate" -ForegroundColor Cyan
    go run github.com/99designs/gqlgen generate
}

function Build {
    Write-Host ">>> go build" -ForegroundColor Cyan
    go build -ldflags="-w -s" -o bin/gateway.exe ./cmd/gateway
}

function Up {
    Write-Host ">>> docker-compose up -d" -ForegroundColor Cyan
    docker-compose -f build/docker-compose.dev.yml up -d
}

function Down {
    Write-Host ">>> docker-compose down" -ForegroundColor Cyan
    docker-compose -f build/docker-compose.dev.yml down
}

function Bench {
    Write-Host ">>> go bench" -ForegroundColor Cyan
    $dir = "./test/bench/gateway"
    New-Item -ItemType Directory -Force -Path $dir
    go test -bench=. -benchmem $dir
}

function Run {
    Build
    Up
    Write-Host ">>> starting gateway..." -ForegroundColor Green
    .\bin\gateway.exe
}

function Test {
    Write-Host ">>> running M1 tests" -ForegroundColor Green
    .\scripts\test_m1.ps1
}

function Help {
    Write-Host "Usage: .\make.ps1 <target>"
    Write-Host "  proto  - generate protobuf"
    Write-Host "  gqlgen - generate graphql"
    Write-Host "  build  - compile gateway.exe"
    Write-Host "  up     - start deps (MySQL/Redis/Kafka)"
    Write-Host "  down   - stop deps"
    Write-Host "  bench  - run benchmark"
    Write-Host "  run    - build+up+run gateway"
    Write-Host "  test   - run M1 test script"
}

&$Target