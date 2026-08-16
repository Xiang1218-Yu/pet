# Benzhi 构建说明

本仓库提供 `benzhi.Dockerfile` 和 `build_benzhi_docker.sh`，用于构建包含 Go 1.26.5 工具链的评测镜像。

## 构建镜像

```bash
./build_benzhi_docker.sh
```

默认生成镜像 `pet-benzhi:latest`。也可以指定镜像名和标签：

```bash
./build_benzhi_docker.sh <image-name> <image-tag>
```

等价 Docker 命令：

```bash
docker build -f benzhi.Dockerfile -t pet-benzhi:latest .
```

## 容器内验证

```bash
docker run --rm -it pet-benzhi:latest bash
go test ./...
```

镜像构建阶段会预下载 Go 模块依赖并执行 `go build ./...`；容器工作目录为 `/app`。
