# 第一阶段：构建Go程序
FROM golang:alpine AS build

ARG VERSION="v1.0.0"
ENV RootDir=/src
ENV AppDir=app/notification
ENV ServerDir=${RootDir}/${AppDir}

# 将当前目录下的所有文件复制到工作目录
COPY ./app ${RootDir}/app
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
COPY --from=build /src/app/notification/server .
COPY ./app/notification/config /data/conf/
# 设置环境变量
ENV PORT=8181

# 暴露端口
EXPOSE $PORT

# 运行Go程序
CMD ["./server", "-conf", "/data/conf"]
