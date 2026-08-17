# Greenlight

A REST API for managing movie data, built with Go and MySQL while working through *Let's Go Further* by Alex Edwards.

The original book uses PostgreSQL. This implementation adapts the database layer and migrations to MySQL.

## Features

Currently implemented:

* Create movies
* Get a movie by ID
* Update movies
* Delete movies
* JSON request and response handling
* Input validation
* MySQL persistence
* Database migrations
* Structured logging
* API versioning with `/v1`

## Tech Stack

* Go
* MySQL
* `net/http`
* `database/sql`
* golang-migrate

## Project Structure

```text
greenlight/
├── cmd/
│   └── api/
│       ├── main.go
│       ├── routes.go
│       └── ...
│
├── internal/
│   ├── data/
│   └── validator/
│
├── migrations/
├── .env.example
├── Makefile
├── go.mod
└── go.sum
```

## Configuration

Create a `.env` file based on `.env.example`:

```env
DB_DSN=username:password@tcp(localhost:3306)/greenlight?parseTime=true
```

Do not commit `.env` to version control.

## Database Migrations

Run all pending migrations:

```bash
migrate \
  -path=./migrations \
  -database='mysql://username:password@tcp(localhost:3306)/greenlight' \
  up
```

Rollback the most recent migration:

```bash
migrate \
  -path=./migrations \
  -database='mysql://username:password@tcp(localhost:3306)/greenlight' \
  down 1
```

## Running the API

```bash
go run ./cmd/api
```

The server runs on:

```text
http://localhost:4000
```

## Example

Create a movie:

```bash
curl -X POST localhost:4000/v1/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Black Panther",
    "year": 2019,
    "runtime": 135,
    "genres": ["sci-fi", "action", "adventure"]
  }'
```

Get a movie:

```bash
curl localhost:4000/v1/movies/1
```

## Status

This project is currently under development as I continue implementing the Greenlight API and exploring production-oriented Go backend patterns.
