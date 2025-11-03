# Hướng dẫn cài đặt công cụ và chạy test với Ginkgo

## Mục lục
1. [Cài đặt Go](#cài-đặt-go)
2. [Cài đặt Ginkgo và Gomega](#cài-đặt-ginkgo-và-gomega)
3. [Cài đặt các công cụ hỗ trợ](#cài-đặt-các-công-cụ-hỗ-trợ)
4. [Chạy test](#chạy-test)
5. [Khắc phục sự cố](#khắc-phục-sự-cố)

## Cài đặt Go

### macOS

#### Cách 1: Sử dụng Homebrew (Khuyên dùng)
```bash
# Cài đặt Homebrew nếu chưa có
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Cài đặt Go
brew install go
```

#### Cách 2: Tải từ trang chủ
1. Truy cập https://golang.org/dl/
2. Tải file `.pkg` cho macOS
3. Chạy file installer và làm theo hướng dẫn

### Windows

#### Cách 1: Sử dụng Chocolatey
```cmd
# Cài đặt Chocolatey nếu chưa có (chạy PowerShell với quyền Administrator)
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Cài đặt Go
choco install golang
```

#### Cách 2: Tải từ trang chủ
1. Truy cập https://golang.org/dl/
2. Tải file `.msi` cho Windows
3. Chạy file installer và làm theo hướng dẫn

### Kiểm tra cài đặt Go
```bash
go version
```

## Cài đặt Ginkgo và Gomega

### Cài đặt Ginkgo CLI

```bash
# Cài đặt Ginkgo CLI tool
go install -a github.com/onsi/ginkgo/v2/ginkgo@latest
```

### Thêm Ginkgo và Gomega vào project

```bash
# Di chuyển đến thư mục project
cd /path/to/your/project

# Thêm dependencies
go get github.com/onsi/ginkgo/v2
go get github.com/onsi/gomega
```

### Kiểm tra cài đặt Ginkgo
```bash
ginkgo version
```

## Cài đặt các công cụ hỗ trợ

## Thiết lập PATH

### macOS (Bash/Zsh)

Thêm vào file `~/.bashrc` hoặc `~/.zshrc`:

```bash
export GOPATH=$HOME/go
export GOROOT=/usr/local/go  # hoặc path cài đặt Go của bạn
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

Reload shell:
```bash
source ~/.bashrc  # hoặc source ~/.zshrc
```

### Windows

#### PowerShell
Thêm vào PowerShell profile (`$PROFILE`):

```powershell
$env:GOPATH = "$env:USERPROFILE\go"
$env:PATH += ";$env:GOPATH\bin"
```

#### Command Prompt
Thiết lập biến môi trường trong System Properties > Environment Variables:
- `GOPATH`: `%USERPROFILE%\go`
- Thêm `%GOPATH%\bin` vào `PATH`

## Chạy test

### Chạy tất cả test
```bash
# Chạy tất cả test trong project
ginkgo -r

# Chạy với verbose output
ginkgo -r -v

# Chạy với progress reporting
ginkgo -r --progress
```

### Chạy test với coverage
```bash
# Chạy test và tạo coverage report
ginkgo -r --cover --coverprofile=coverage.out

# Xem coverage report dạng HTML
go tool cover -html=coverage.out -o coverage.html

# Mở file HTML (macOS)
open coverage.html

# Mở file HTML (Windows)
start coverage.html
```

### Chạy test trong package cụ thể
```bash
# Chạy test trong thư mục hiện tại
ginkgo

# Chạy test trong package service
ginkgo ./internal/application/service

```

## Cấu hình IDE

### VS Code
Cài đặt extension:
- Go (by Google)
- Go Test Explorer

Thêm vào `settings.json`:
```json
{
    "go.testFlags": ["-v"],
    "go.coverOnTestPackage": true,
    "go.testTimeout": "30s"
}
```

### GoLand/IntelliJ
- Cài đặt Go plugin nếu sử dụng IntelliJ
- Cấu hình Go SDK trong Project Settings
- Enable test coverage trong Run Configuration

## Khắc phục sự cố

### Lỗi "ginkgo command not found"

#### macOS/Linux
```bash
# Kiểm tra GOPATH và PATH
echo $GOPATH
echo $PATH

# Thêm GOPATH/bin vào PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

#### Windows
```cmd
# Kiểm tra biến môi trường
echo %GOPATH%
echo %PATH%

# Cài đặt lại Ginkgo
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

### Lỗi module dependency
```bash
# Làm sạch module cache
go clean -modcache

# Tải lại dependencies
go mod download

# Sync dependencies
go mod tidy
```

### Lỗi permission (macOS/Linux)
```bash
# Cấp quyền execute cho file binary
chmod +x $(go env GOPATH)/bin/ginkgo
```