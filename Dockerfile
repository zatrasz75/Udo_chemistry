# Указываем базовый образ
FROM golang:latest as builder
LABEL authors="zatrasz"

# Создание рабочий директории
RUN mkdir -p /app/udo

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app/udo

# Копируем файлы проекта внутрь контейнера
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o udo ./cmd/main.go


# Второй этап: создание production образ
FROM ubuntu AS chemistry

WORKDIR /app/udo

RUN apt-get update

COPY --from=builder /app/udo/udo ./
COPY ./ ./

# Указываем переменную окружения
ENV APP_PORT=7654
ENV APP_HOST=0.0.0.0

CMD ["./udo"]


# docker build -t udo .

#docker run --rm -d -p 7655:7654 --name=udo_chemistry udo

# docker-compose up -d