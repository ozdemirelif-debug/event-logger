#Event Logger

Event Logger, RabbitMQ üzerinden gelen event mesajlarını dinleyip, PostgreSQL veritabanına kaydeden, GO ile yazılmış mikroservis tabanlı bir uygulamadır. Ayrıca Prometheus ve Grafana ile izleme ve görselleştirme imkanı sunar.

---

## İçerik

- RabbitMQ (Message Broker)
- PostgreSQL (Veri tabanı)
- API Servisi (Go)
- Consumer Servisi (Go)
- Prometheus (Monitoring)
- Grafana (Dashboard)

---

## Başlangıç

Projeyi çalıştırmak için Docker ve Docker Compose kurulu olmalıdır.


### Adımlar:

1. Proje klasörüne gel:
-  cd event-logger
2. Docker Compose ile tüm servisleri başlat:
-  docker-compose up --build
3. API ve Consumer servisleri RabbitMQ ve PostgreSQL'e bağlanarak çalışmaya başlayacaktır.
4. Projeye event göndermek için send.go dosyası kullanılabilir:
- go run send.go
5. Manuel olarak API'a event gönderilebilir:
Manuel olarak API'a event gönderilebilir:

```bash
curl -X POST http://localhost:8080/events \
-H "Content-Type: application/json" \
-d '{
  "eventId": "12345",
  "source": "test-service",
  "type": "user.signup",
  "payload": {
    "username": "elif",
    "email": "elif@example.com"
  },
  "timestamp": "2025-06-14T10:00:00Z"
}'
````
### İzleme ve Dashboard:

- Prometheus http://localhost:9090
- Grafana http://localhost:3000 (kullanıcı: admin, şifre: admin)
- RabbitMQ http://localhost:15672 (kullanıcı: guest, şifre: guest)

### Kullanılan Teknolojiler

- Go
- RabbitMQ
- PostgreSQL
- Prometheus
- Grafana
- Docker & Docker Compose

  ### Proje Yapısı

- event-logger/
- |----api/
- |----consumer/
- |----send.go
- |----docker-compose.yml
- |----Dockerfile.api
- |----Dockerfile.consumer
- |----prometheus.yml
