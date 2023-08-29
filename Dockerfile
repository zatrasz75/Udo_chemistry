# Указываем базовый образ
FROM golang:latest
LABEL authors="zatrasz"

# Создание рабочий директории
RUN mkdir -p /app/udo

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app/udo

# Копируем файлы проекта внутрь контейнера
COPY ./ ./

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd/main.go ./

# Указываем переменную окружения
ENV APP_PORT=7654
ENV APP_HOST=0.0.0.0

# Собираем приложение
RUN go build -o udo .

# Указываем, что контейнер будет слушать порт 7654
EXPOSE 7654

CMD ["./udo"]

# docker build -t udo .

#docker run --rm -d -p 7655:7654 --name=udo_chemistry udo