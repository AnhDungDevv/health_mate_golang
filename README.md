# 📌 Health_mate

A brief description of the project, its main functionality, and objectives.

## 🚀 Technologies Used

- [Gin Gonic](https://gin-gonic.com/) - Web framework for Go
- [Kafka](https://kafka.apache.org/) - Message Queue system
- [Redis](https://redis.io/) - Caching & Pub/Sub
- [PostgreSQL](https://www.postgresql.org/) - Primary database
- [WebSocket](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API) - Real-time communication
- [Swagger](https://swagger.io/) - API documentation
- [Prometheus](https://prometheus.io/) - Monitoring & Metrics
- [Docker](https://www.docker.com/) - Containerization
- [Docker Compose](https://docs.docker.com/compose/) - Service management

---

## ⚙️ Installation & Running the Project

### 1️⃣ Clone the Repository

\`\`\`sh
git clone https://github.com/username/repository.git
cd repository
\`\`\`

### 2️⃣ Run with Docker Compose

\`\`\`sh
docker-compose up --build
\`\`\`

### 4️⃣ Run Manually (Without Docker)

1. **Start Required Services**: Redis, PostgreSQL, Kafka
2. **Run the Go Application**
   \`\`\`sh
   go run main.go
   \`\`\`

---

## 📌 Key API Endpoints

| Method | Endpoint         | Description                 |
| ------ | ---------------- | --------------------------- |
| GET    | `/health`        | Check server status         |
| POST   | `/auth/login`    | User login                  |
| POST   | `/auth/register` | User registration           |
| GET    | `/ws`            | WebSocket connection        |
| GET    | `/metrics`       | Prometheus metrics endpoint |
| GET    | `/swagger/*`     | API documentation           |

---

## 🛠️ Database Migrations

## 📜 Swagger API Documentation

Swagger documentation is available at:
\`\`\`
http://localhost:5000/swagger/index.html
\`\`\`

---

## 📊 Monitoring with Prometheus

Application metrics can be accessed at:
\`\`\`
http://localhost:7070/metrics
\`\`\`

---

## 🐳 Docker & Docker Compose

### 📌 Run with Docker Compose

\`\`\`sh
docker-compose up --build
\`\`\`

---

## 🔗 Useful Resources

- [Gin Gonic](https://gin-gonic.com/)
- [Kafka Documentation](https://kafka.apache.org/documentation/)
- [Redis Guide](https://redis.io/documentation)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [Swagger Docs](https://swagger.io/docs/)
- [Prometheus Docs](https://prometheus.io/docs/)
- [Docker Docs](https://docs.docker.com/)

---
