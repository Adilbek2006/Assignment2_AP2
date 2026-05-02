# AP2 Assignment 3 - Event-Driven Architecture (EDA) with Message Queues

* **Name:** Adilbek Zhetpyspayev
* **Group:** SE-2406

## Contract-First Flow & Repositories
As per the assignment requirements, the Protocol Buffers are managed remotely using a Contract-First approach. The automated generation pipeline is set up via GitHub Actions.

* **Protos Repository (Contracts):** https://github.com/Adilbek2006/protos
* **Generated Code Repository (CI/CD):** https://github.com/Adilbek2006/grpc-generated
* **Main Services Repository (Git History):** https://github.com/Adilbek2006/Assignment2_AP2

## Project Structure
* `order-service/`: Exposes REST API for users and acts as a gRPC Client.
* `payment-service/`: Acts as a gRPC Server processing payments and serves as an **Event Producer** (publishes `payment.completed` messages to RabbitMQ).
* `notification-service/`: New microservice acting as an **Event Consumer**. Listens to the RabbitMQ queue and simulates sending asynchronous email notifications.
* `test-client/`: A console gRPC client built to trigger the Payment process and demonstrate the end-to-end event flow.
* `docker-compose.yml`: Orchestrates the entire infrastructure, including PostgreSQL, RabbitMQ, and all microservices.

## Architecture Highlights
* **Idempotency:** Implemented via an in-memory map (`processedMessages`) in the Notification Service to safely ignore duplicate events based on `OrderID`.
* **Reliability & Manual ACKs:** `auto-ack` is explicitly disabled. The consumer manually acknowledges messages (`d.Ack(false)`) only after successful processing. The queue is configured as `durable=true`.
* **Graceful Shutdown:** Implemented in the services using `os/signal` to ensure safe termination of active connections.
  
* **Evidences:**
Successful EDA Flow (Producer & Consumer logs in Docker):<img width="1788" height="172" alt="image" src="https://github.com/user-attachments/assets/f14447f9-b611-4c95-836b-7629ba16e13b" />


Idempotency handling (Duplicate Ignored):<img width="1793" height="308" alt="image" src="https://github.com/user-attachments/assets/e8f6515a-81b5-47e4-9d97-8c6967c2ca09" />

## How to Run the Project
* **cd test-client:**
* **go run main.go:**

### 1. Prerequisites
* Docker Desktop installed and running locally.

### 2. Start Infrastructure and Services
The entire application environment is containerized. Open a terminal in the root directory and run:
```bash
docker-compose up --build

