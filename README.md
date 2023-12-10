# Сonstellation Demo Backend

## Installation

To install your project, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/cexles/api-service.git
   cd api-service

2. Install dependencies

    ```bash
   go mod download
   
3. Create config file

    ```bash
   cp config.example.json config.json
   
Fill rpc endpoint and database

4. Use migrations from migration folder

    ```bash
   migrate -path migration/ -database "postgresql://db_user:db_pass@localhost:db_port/db_name?sslmode=disable" -verbose up

5. Run backend

    ```bash
   go run main.go