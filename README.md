# Project Sprint ☺

Lets get code together in this project sprint.
This repo should be tested 7 day from now.

## Installation

Clone project

```bash
  git clone https://github.com/anti-social-cat-social/halosuster.git my-project
  cd my-project
```

After cloning the project and inside the project folder, run these command

1. Install dependecies & library needed

```go
    go mod install
```

2. Copy / create env file
   Create .env file or run this command

```bash
    cp .env.example .env
```

This service depends on `golang-migrate cli`. Install it first globally before do migrations (step 3 and 4).

3. Creating migration
    Create migration for the project using [Golang Migrate](https://github.com/golang-migrate/migrate)
    
    ```bash
     migrate create -ext sql -dir database/migrations {table_name}
    ```

4. Running migration
   Run the project migration to get updated with the table

```
   migrate -path database/migrations -database "postgres://testing:testing@localhost:5433/testing?sslmode=disable" -verbose up
```

5. Run the project
   After all the step above is fulfilled, you can run this project using this command.
   (Make sure you are in root folder of the project)

```bash
    go run .
```

5. Optional (Running database using docker-compose for development)
    If you want to run the database using docker compose, you can run this command to start the database
    
    ```bash
    ./scripts/setup.sh
    ```