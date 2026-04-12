# Assignment 1
**Student:** Adilbek Zhetpyspayev  
**Group:** SE-2406

## Architecture Decisions
**Clean Architecture**: Each service is divided into Domain, UseCase, Repository, and Transport layers to ensure Separation of Concerns
- Microservices Principles:
    - Database per Service: Order Service and Payment Service use separate PostgreSQL databases 
    - No Shared Code: Services are fully decoupled. Domain models are duplicated to avoid "Distributed Monolith"
- Manual Dependency Injection: All dependencies are initialized in the Composition Root

## Failure Handling
- Resilient Communication: The Order Service uses a custom http.Client with a 2-second timeout when calling the Payment Service
- Graceful Degradation: If the Payment Service is unavailable, the Order Service returns a 503 Service Unavailable error and marks the order status as "Failed" to maintain consistency

## How to Run
1. Create orders_db and payments_db in PostgreSQL.
2. Run SQL scripts from /migrations in respective databases.
3. Start Payment Service: cd payment-service && go run cmd/payment/main.go
4. Start Order Service: cd order-service && go run cmd/order/main.go