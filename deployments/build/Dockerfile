FROM docker:27-dind
# Подключаем edge репозитории
RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
    echo "https://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
# Устанавливаем Go, Git и Bash
RUN apk add --no-cache go git bash tini htop
# Устанавливаем рабочую директорию
WORKDIR /go/src/app
# Копируем только файлы зависимостей
COPY go.mod .
COPY go.sum .
RUN go mod download -x
COPY . .
