# Stage 1: Build application
FROM golang:alpine as builder 


WORKDIR /app

# Copy toàn bộ mã nguồn vào container
COPY . .

# Download và cài đặt các module cần thiết
RUN go mod tidy

# Tải CompileDaemon và build nó
RUN go get github.com/githubnemo/CompileDaemon \
    && go build -o /go/bin/CompileDaemon github.com/githubnemo/CompileDaemon

# Biên dịch ứng dụng
RUN go build -o main cmd/api/main.go


# Stage 2: Runner (Chạy ứng dụng)
FROM golang:alpine as runner


WORKDIR /app

# Copy binary từ builder stage
COPY --from=builder /app/main /app/main
COPY --from=builder /go/bin/CompileDaemon /usr/local/bin/CompileDaemon

# Thiết lập biến môi trường
ENV config=docker

# Mở các cổng ứng dụng sử dụng
EXPOSE 5000
EXPOSE 5555
EXPOSE 7070

# Chạy ứng dụng với CompileDaemon (tự động reload khi code thay đổi)
ENTRYPOINT ["/usr/local/bin/CompileDaemon", "--build=go build -o /app/main cmd/api/main.go", "--command=/app/main"]
