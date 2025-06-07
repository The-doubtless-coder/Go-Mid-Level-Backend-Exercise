# Go-Mid-Level-Backend-Exercise
This is a backend microservice built with **Go (Gin)**, **PostgreSQL**, and **Keycloak** for user authentication via **OpenID Connect**. It manages customers, hierarchical product categories, and orders, complete with **email notifications**, **SMS alerts**, and **token-based authentication**.

---

## ✨ Features

- 🔐 **Authentication & Authorization** (OpenID Connect with Keycloak)
- 🧍 **Customers CRUD**
- 🗂️ **Nested Product Categories** (unlimited depth)
- 🛒 **Order Management**
- 📊 **Average Product Price** per category (recursive)
- ✉️ **Email Notifications** to admin on order placement
- 📱 **SMS Alerts** to customer using [Africa’s Talking API](https://africastalking.com/)
- 🧪 **Unit Tests with Coverage Checking**
- 🤖 **CI/CD Pipeline** (GitHub Actions)
- 📄 **Swagger Documentation**

---

## 📐 System Architecture

+------------------+ +---------------+
| Customer | <---> | Keycloak |
+------------------+ +---------------+
| |
| Login/Register |
v |
+-------------------------------+
| Go Backend |
| Gin + GORM + PostgreSQL |
+-------------------------------+
| | |
v v v
+------+ +------+ +------+
| Order| |Product| |Category|
+------+ +------+ +------+
|
| On order placed
|
|--> Send SMS (Africa’s Talking)
|--> Send Email (Admin)


---

## 🗃️ Data Models

### 🧍 Customers
- UUID `id`
- `name`, `email`, `phone`
- Authenticated via Keycloak

### 🛒 Products
- UUID `id`
- `name`, `price`, `category_id`
- Belongs to a nested category

### 🗂️ Categories
- UUID `id`
- `name`, `parent_id` (self-referencing for hierarchy)

### 📦 Orders
- UUID `id`
- `customer_id`, `product_id`, `quantity`, `total_price`, `created_at`

---

## 🔐 Authentication

- Keycloak realm: `savannah`
- Clients:
    - `go-customer`: For login (standard + direct access)
    - `go-backend`: For service accounts (used by backend)
- Roles:
    - `customer`
    - `backend_admin`

Token lifespan: **15 minutes**

---

## 📲 SMS via Africa’s Talking

- Sandbox number: `+254 790 360 360`
- Alerts customers immediately after placing orders.
- Configure credentials in `.env`:
  ```env
  AT_API_KEY=your_sandbox_api_key
  AT_USERNAME=sandbox
  AT_SENDER_ID=your_sender_id

## 📲 EMAIL notifs via GMAIL or any MAIL provider
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your@email.com
SMTP_PASS=yourpassword
ADMIN_EMAIL=admin@example.com

## 📲 Unit Testing (service layer) with Coverage
go test ./... -cover
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

## CI/CD
## GitHub Actions workflow:
Run tests
Check coverage
Build Docker image -(not included since I don't have a registry and a VM to deploy to)
Push to registry
Deploy (optional)
CI config in .github/workflows/ci.yml

## Deployment
## With Docker Compose (Dev)
Find file in root directory
docker-compose up --build

## 📄 API Documentation
Swagger docs auto-generated via annotations.
Accessible at: http://localhost:8083/swagger/index.html or the port that you are running the project at.





