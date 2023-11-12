# 第一阶段：构建Go程序
FROM golang:alpine AS build

ARG VERSION="v1.0.0"
ENV RootDir=/src
ENV ServerDir=${RootDir}/app/harmoni

# 将当前目录下的所有文件复制到工作目录
COPY ./app/harmoni ${ServerDir}
COPY ./internal ${RootDir}/internal

WORKDIR ${RootDir}

# 安装依赖
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod init harmoni
RUN go mod tidy

# 构建Go程序
RUN go build -ldflags "-X main.Version=$VERSION" -o ${ServerDir}/server ${ServerDir}/cmd/.

# 第二阶段：运行Go程序
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从第一阶段中复制构建好的Go程序
COPY --from=build /src/app/harmoni/server .
COPY ./app/harmoni/config ./configs/
# 设置环境变量
ENV PORT=8080
ENV PORT=8082
# 暴露端口
EXPOSE $PORT

# 运行Go程序
CMD ["./server"]
