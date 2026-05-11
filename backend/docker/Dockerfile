FROM registry-vpc.cn-shenzhen.aliyuncs.com/dreamdev/golang:1.22 AS builder
# 构建参数
ARG  APP_ROOT
ARG APP_NAME
# 设置工作目录
WORKDIR ${APP_ROOT}
COPY .. ./
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${APP_NAME} main.go

FROM registry-vpc.cn-shenzhen.aliyuncs.com/ebag/dream-ebag-new-alpine:v3
# 构建参数
ARG APP_ROOT
ARG EXPOSE_PORT
ARG APP_NAME
ARG ENV_STRING
ARG NAMESPACE
ARG POD_NAME
# 设置环境变量
ENV ENV_STRING=${ENV_STRING}
ENV APPNAMEENV ${APP_NAME}
ENV NAMESPACE=${NAMESPACE}
ENV POD_NAME=${POD_NAME}
# 工作目录和端口暴露
WORKDIR /app
EXPOSE ${EXPOSE_PORT}
# 复制上层构建产物
COPY --from=builder ${APP_ROOT}/${APP_NAME} ./${APP_NAME}
COPY --from=builder ${APP_ROOT}/docker/run.sh ./run.sh
COPY --from=builder ${APP_ROOT}/application.${ENV_STRING}.yml ./application.${ENV_STRING}.yml

# 设置执行权限
RUN chmod +x ./run.sh && \
    chmod +x ./${APP_NAME}

# 入口点：运行脚本。不需要功能或语法时直接执行，性能更高
ENTRYPOINT ["./run.sh"]
#ENTRYPOINT ["/bin/sh","-c","/app/./run.sh"]