# Go REST API with Gin and PostgreSQL

A simple yet robust RESTful API built with Go, the Gin framework, and PostgreSQL. This project is fully containerized using Docker and Docker Compose, providing a clean and reproducible development environment.

## ✨ Features

- **Health Check:** A `/health` endpoint to monitor the API's status.
- **User Management:** Full CRUD (Create, Read, Update, Delete) functionality for users.
- **Containerized:** Runs entirely within Docker containers for consistency and ease of deployment.
- **Structured Layout:** Follows the standard Go project layout for better organization.

## 🛠️ Technologies Used

- **Go:** The core programming language.
- **Gin:** A high-performance HTTP web framework for Go.
- **PostgreSQL:** A powerful, open-source object-relational database system.
- **sqlx:** A library providing a set of extensions on top of `database/sql`.
- **Docker & Docker Compose:** For containerizing and orchestrating the application and database services.

## 🚀 Getting Started

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- [Go](https://go.dev/doc/install) (v1.24 or newer)
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)

### Installation & Running

1.  **Clone the repository:**
    ```sh
    git clone git@github.com:raingrave/go-rest-api.git
    cd go-rest-api
    ```

2.  **Run the application with Docker Compose:**
    This single command will build the API image, start the API and database containers, and connect them.
    ```sh
    docker compose up --build -d
    ```
    The API will be available at `http://localhost:3000`.

3.  **Set up the database:**
    Connect to the PostgreSQL database (running on `localhost:5432`) and execute the following SQL command to create the `users` table.
    ```sql
    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE
    );
    ```

## Endpoints API

The base URL is `http://localhost:3000`.

| Method   | Endpoint      | Description                  | Request Body Example                             |
| :------- | :------------ | :--------------------------- | :----------------------------------------------- |
| `GET`    | `/health`     | Checks the API status.       | `N/A`                                            |
| `GET`    | `/users`      | Retrieves a list of all users. | `N/A`                                            |
| `POST`   | `/users`      | Creates a new user.          | `{"name": "John Doe", "email": "john@doe.com"}`   |
| `GET`    | `/users/{id}` | Retrieves a single user by ID. | `N/A`                                            |
| `PUT`    | `/users/{id}` | Updates an existing user.    | `{"name": "Jane Doe", "email": "jane@doe.com"}`   |
| `DELETE` | `/users/{id}` | Deletes a user by ID.        | `N/A`                                            |
