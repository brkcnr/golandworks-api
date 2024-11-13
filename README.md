# Golandworks-API

Simple api server for making your own to-do list.

Docker and PostgreSQL needs to be installed.

Set up your own PostgreSQL server through docker-compose.yaml file's settings.

- Hostname: localhost
- User: postgres
- Password: Get it from docker-compose.yaml file.
- Port: Get it from docker-compose.yaml file.

---

## Usage

To run docker;

```bash
docker-compose up
```

or you can right click docker-compose.yaml file and select compose up.

---

To run server locally;

```go
go run main.go
```