# Go Kafka Test

Dự án này minh họa việc sử dụng Apache Kafka với Go để xây dựng hệ thống producer-consumer đơn giản. Producer nhận tin nhắn từ API endpoint và gửi đến Kafka, và consumer đọc và xử lý các tin nhắn này.

## Yêu cầu

- Go 1.18 trở lên
- Docker và Docker Compose
- cURL hoặc công cụ API client tương tự (như Postman)

## Cấu trúc dự án

```
go-kafka-test/
  ├── consumer/            # Ứng dụng consumer
  │   └── main.go          # Code của consumer
  ├── producer/            # Ứng dụng producer
  │   └── main.go          # Code của producer
  ├── docker-compose.yaml  # Cấu hình Docker cho Kafka, ZooKeeper và Kafka UI
  ├── go.mod               # Quản lý dependencies của Go
  ├── go.sum               # Checksums cho Go modules
  └── README.md            # Tài liệu hướng dẫn
```

## Cài đặt và chạy

### 1. Khởi động môi trường Kafka với Docker

```bash
docker-compose up -d
```

Lệnh này sẽ khởi động:

- ZooKeeper trên cổng 2181
- Kafka broker trên cổng 9092
- Kafka UI trên cổng 8080

### 2. Chạy producer

```bash
cd producer
go run main.go
```

Producer sẽ khởi động một HTTP server trên cổng 8999, cung cấp endpoint `/action` để nhận tin nhắn và gửi đến Kafka.

### 3. Chạy consumer

```bash
cd consumer
go run main.go
```

Consumer sẽ kết nối đến Kafka và bắt đầu lắng nghe tin nhắn từ topic "order-topic".

## Sử dụng

### Gửi tin nhắn tới producer

Sử dụng cURL để gửi một tin nhắn đến producer:

```bash
curl -X POST http://localhost:8999/action \
  -H "Content-Type: application/json" \
  -d '{"content":"Đây là tin nhắn test"}'
```

Sau khi gửi tin nhắn, bạn sẽ thấy:

- Producer log thông tin về tin nhắn đã gửi
- Consumer log thông tin về tin nhắn đã nhận

### Xem dữ liệu qua Kafka UI

Truy cập Kafka UI tại địa chỉ [http://localhost:8080](http://localhost:8080) để:

- Theo dõi các topics
- Kiểm tra số lượng tin nhắn
- Theo dõi các nhóm consumer
- Xem nội dung tin nhắn

## Cấu hình

### Producer

Producer được cấu hình để kết nối đến Kafka broker tại `localhost:9092` và gửi tin nhắn đến topic "order-topic".

### Consumer

Consumer được cấu hình với các thông số sau:

- Kết nối đến Kafka broker tại `localhost:9092`
- Lắng nghe topic "order-topic"
- Sử dụng consumer group ID "khach-hang-vip"
- Các tham số hiệu suất đã được tối ưu để giảm độ trễ:
  - MinBytes: 1 (nhận tin nhắn ngay khi có)
  - MaxBytes: 1MB
  - MaxWait: 100ms
  - ReadBatchTimeout: 100ms
  - CommitInterval: 0 (commit ngay lập tức)

## Xử lý lỗi thường gặp

### Không thể kết nối đến Kafka

Nếu bạn gặp lỗi "no such host" hoặc "connection refused", hãy kiểm tra:

- Docker containers có đang chạy không
- Địa chỉ kết nối trong code có đúng không
- Cổng Kafka có mở không (9092)

### Topic không tồn tại

Nếu gặp lỗi "Unknown Topic Or Partition", hãy kiểm tra:

- Topic đã được tạo trong Kafka chưa
- Xác nhận rằng KAFKA_CREATE_TOPICS trong docker-compose.yaml đúng với topic bạn đang sử dụng

### Độ trễ tin nhắn

Nếu bạn thấy có độ trễ giữa việc gửi và nhận tin nhắn, hãy xem lại các tham số cấu hình trong consumer và producer. Consumer có thể cần điều chỉnh thêm các tham số như `MinBytes`, `MaxBytes`, và `MaxWait`.

## Tham khảo thêm

- [Kafka Go Client Documentation](https://pkg.go.dev/github.com/segmentio/kafka-go)
- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
