# Go Common Library - Modular & Reusable Golang Utilities 🚀

## 📌 Overview
The **Go Common Library** is a modular and reusable set of utilities designed to simplify and accelerate Golang application development. It provides a collection of well-structured, production-ready packages for logging, configuration management, database access, caching, messaging, security, and more.

This library ensures **code reusability, consistency, and maintainability** across multiple projects by offering standardized implementations of essential components required in backend services and microservices.

---

## 🚀 Features
✔️ **Structured Logging** – Logrus-based logging with built-in log rotation.  
✔️ **Configuration Management** – Viper-based loader for JSON, YAML, and environment variables.  
✔️ **Database Helper** – Easy-to-use wrappers for MySQL, PostgreSQL, and MongoDB.  
✔️ **HTTP Client** – Customizable HTTP client with retries, timeouts, and request logging.  
✔️ **Caching** – Redis-based caching and in-memory LRU cache with TTL support.  
✔️ **Messaging & Pub-Sub** – Kafka, NATS, and RabbitMQ integration with consumer/producer abstractions.  
✔️ **Task Queues** – Background job processing with Asynq and worker pools.  
✔️ **Rate Limiting** – Redis-based API rate limiting middleware.  
✔️ **Security & Auth** – JWT-based authentication, password hashing, and encryption utilities.  
✔️ **Observability** – OpenTelemetry-based tracing and Prometheus metrics exporter.  
✔️ **Notifications** – Email (SendGrid), SMS (Twilio), and Webhook-based notifications.  
✔️ **Feature Flags** – Dynamic feature flag management for A/B testing.  
✔️ **Cron Jobs** – Scheduled background job execution support.  

---

## 📂 Project Structure
```bash
go-common-lib/
│── internal/                      # Internal utilities
│── pkg/                           # Exported reusable packages
│   ├── cache/                     # Redis & in-memory caching
│   ├── config/                    # Configuration management (Viper)
│   ├── database/                   # Database helpers (MySQL, MongoDB, PostgreSQL)
│   ├── httpclient/                 # HTTP client with retries & logging
│   ├── kafka/                      # Kafka producer & consumer utilities
│   ├── logger/                     # Structured logging (Logrus + Lumberjack)
│   ├── metrics/                    # Prometheus metrics integration
│   ├── middleware/                 # CORS, Auth, Rate Limiting, Logging
│   ├── notifier/                   # Email, SMS, Webhooks
│   ├── pubsub/                     # Messaging queues (Kafka, NATS, RabbitMQ)
│   ├── rate_limiter/               # API Rate Limiting
│   ├── search/                     # OpenSearch & Elasticsearch wrapper
│   ├── security/                   # JWT, OAuth, encryption
│   ├── tracing/                    # OpenTelemetry distributed tracing
│   ├── worker/                     # Background task processing
│── examples/                        # Usage examples
│── tests/                           # Integration tests
│── go.mod
│── README.md
│── LICENSE
