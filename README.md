# AP2 Assignment 4 - Performance Optimization & External Integrations

 **Name:** Adilbek Zhetpyspayev
 **Group:** SE-2406

## Contract-First Flow & Repositories
As per the course requirements, Protocol Buffers are managed remotely using a Contract-First approach. The automated generation pipeline is set up via GitHub Actions.

 **Protos Repository (Contracts):** https://github.com/Adilbek2006/protos
 **Generated Code Repository (CI/CD):** https://github.com/Adilbek2006/grpc-generated
 **Main Services Repository:** https://github.com/Adilbek2006/Assignment2_AP2

## Project Structure
 `order-service/`: Exposes REST API. Now optimized with **Redis Cache-aside** pattern and a **Rate Limiter**.
 `payment-service/`: Acts as a gRPC Server processing payments and publishes `payment.completed` events to RabbitMQ.
 `notification-service/`: Refactored into a robust **Background Worker**. Uses the Adapter Pattern to simulate external email providers and handles exponential backoff retries.
 `docker-compose.yml`: Orchestrates the entire infrastructure, now including **Redis** and persistent Docker Volumes for PostgreSQL.

## Architecture Highlights & Strategies

### 1. Caching & Invalidation Strategy (Order Service)
 **Cache-aside Pattern:** Implemented for `GET /orders/:id`. The system first checks Redis (Read Path). If a Cache Miss occurs, it queries PostgreSQL and sets the cache with a **5-minute TTL**.
 **Atomic Invalidation:** To prevent serving stale data, the cache is explicitly invalidated (`RedisClient.Del`) whenever an order's status changes in the database (e.g., during `CreateOrder` or `CancelOrder` operations).

### 2. Reliable Background Worker & Retry Logic (Notification Service)
 **Adapter Pattern:** Notification logic is decoupled using an `EmailSender` interface. Controlled via the `PROVIDER_MODE` environment variable, it switches between a Real SMTP implementation and a Simulated Provider.
 **Simulated Environment:** The mock provider introduces artificial network latency (`time.Sleep`) and a 30% chance of `503 Service Unavailable` failures to test resilience.
 **Exponential Backoff Strategy:** If the simulated provider fails, the worker does not drop the message. It retries the job up to 3 times, exponentially increasing the delay between attempts (**2s -> 4s -> 8s**) before either succeeding or negatively acknowledging (NACK) the message.

### 3. Distributed Idempotency
 The previous in-memory map was replaced with a distributed Redis store. The worker uses Redis `SetNX` to lock processed `payment_id`s, ensuring that duplicate events from RabbitMQ never result in duplicate email sends.

### 4. API Rate Limiter (Bonus +10%)
 A custom middleware was built for the `order-service`. It tracks client requests by IP in Redis and restricts usage to **10 requests per minute**. Exceeding this limit returns an `HTTP 429 Too Many Requests` error.

---

## Proof of Execution

**1. Cache Hit Performance (< 1ms) & Rate Limiting (429 Error):**
```text
order-service-1  | [GIN] 2026/05/11 - 16:37:24 | 200 |   1.78ms | POST  "/orders"
order-service-1  | [GIN] 2026/05/11 - 16:37:34 | 200 | 611.879µs| GET   "/orders/3303c..." (Cache Hit)
order-service-1  | [GIN] 2026/05/11 - 16:37:38 | 200 |   486.7µs| GET   "/orders/3303c..." (Cache Hit)
order-service-1  | [GIN] 2026/05/11 - 16:37:39 | 429 | 347.103µs| GET   "/orders/3303c..." (Rate Limited)
```
**2. Exponential Backoff in action:**
```text
notification-service-1  | [Worker Warning] Provider failed: simulated 503 Service Unavailable. Retrying in 2s (Attempt 1/3)...
notification-service-1  | [Worker Warning] Provider failed: simulated 503 Service Unavailable. Retrying in 4s (Attempt 2/3)...
notification-service-1  | [SIMULATED PROVIDER] Successfully sent email to user@example.com for Order #a71b4050...
```

## How to Run the Project
1. Prerequisites
Docker Desktop installed and running locally.

2. Start Infrastructure and Services
The entire application environment is containerized. Open a terminal in the root directory and run:
```bash
docker-compose up --build
```


## Architecture Diagram
<img width="4386" height="5498" alt="diagram" src="https://github.com/user-attachments/assets/a1dea939-8555-48c1-a55f-2511afbdc573" />
