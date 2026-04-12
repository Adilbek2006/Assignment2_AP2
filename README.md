 AP2 Assignment 2 - gRPC Migration & Contract-First Development

* **Name:** Adilbek Zhetpyspayev
* **Group:** SE-2406

## Contract-First Flow & Repositories
As per the assignment requirements, the Protocol Buffers are managed remotely using a Contract-First approach.The automated generation pipeline is set up via GitHub Actions

* **Protos Repository (Contracts):** https://github.com/Adilbek2006/protos
* **Generated Code Repository (CI/CD):** https://github.com/Adilbek2006/grpc-generated
* **Main Services Repository (Git History):** https://github.com/Adilbek2006/Assignment2_AP2

## Project Structure
* `order-service/`: Exposes REST API for users, acts as a gRPC Client to communicate with Payment Service, and acts as a gRPC Server for Server-side Streaming
* `payment-service/`: Acts as a gRPC Server handling payment processing
* `test-client/`: A console gRPC client built to demonstrate real-time Server-side Streaming tied to database updates

## How to Run the Project

### 1. Prerequisites
* Go 1.25.5
* PostgreSQL running locally

### 2. Environment Configuration
Create `.env` files in both `order-service` and `payment-service` to manage ports and database DSNs securely 

### 3. Start Payment Service (gRPC Server)
cd payment-service
go mod tidy
go run main.go

### 4.Start Order Service (REST + gRPC Streaming Server)
cd order-service
go mod tidy
go run main.go


### 5. Run Streaming Client
Bash
cd test-client
go run main.go

Evidences
<img width="1218" height="127" alt="image" src="https://github.com/user-attachments/assets/3c9d4249-7601-4643-8959-52183117c197" />
Before: 
<img width="1140" height="84" alt="image" src="https://github.com/user-attachments/assets/f5640c2b-c382-4409-91a4-e66f352780b4" />
<img width="1099" height="519" alt="image" src="https://github.com/user-attachments/assets/d02070a6-7924-4884-ba6c-008bdc80d8d6" />

After: 
<img width="1196" height="111" alt="image" src="https://github.com/user-attachments/assets/10450b8f-1049-4cb0-ad07-3862821dcb04" />

