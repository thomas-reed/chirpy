# Chirpy
Boot.dev project - Chirpy is a simple Twitter-style microblogging app built as a learning project.  
Users can sign up and post short text updates (“chirps”).

## Features

- Create an account
- Log in and log out
- Create, view, update, and delete chirps
- View a feed of recent chirps

## Tech Stack

- Language: Go
- HTTP Server: Go standard library (`net/http`)
- JSON Handling: Go standard library (`encoding/json`)
- Database: PostgreSQL

## Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/thomas-reed/chirpy.git
   cd chirpy
   ```

2. Configure environment:

    ```bash
    cp .env.example .env
    # edit .env with your own settings
    ```

3. Build the server:

    ```bash
    go build -o server
    ```

4. Run the server:

    ```bash
    ./server
    ```

5. Open:

    http://localhost:8080/app

    Be advised, this was a backend project - the frontend is one page with a title on it.  Do this just to confirm it's working.
    
## API Overview

Chirpy is primarily a backend JSON API built with Go’s `net/http` and `http.ServeMux`.  All endpoints return JSON responses and use standard HTTP status codes.

### UI

- `GET /app/*` – serves the frontend assets (from `FILEPATHROOT`) via a static file server

### Health

- `GET /api/healthz` – simple health check endpoint

### Auth & Users

Authenticated requests use a Bearer token in the header:

  ```http
  Authorization: Bearer <access_token>
  ```

- `POST /api/users` – create a new user
- `PUT /api/users` – update user credentials (e.g. email/password), requires auth
- `POST /api/login` – log in with email/password  
  - Returns an access token (JWT) and a refresh token
- `POST /api/refresh` – exchange a valid refresh token for a new access token
- `POST /api/revoke` – revoke a refresh token

### Chirps

- GET /api/chirps – list chirps (optionally with filter by `author_id`, or sorting using `sort` 'asc' or 'desc')
- GET /api/chirps/{id} – get a single chirp by ID
- POST /api/chirps – create a new chirp (requires auth)
- DELETE /api/chirps/{id} – delete a chirp (requires auth and ownership)

### Webhooks

- POST /api/polka/webhooks – webhook endpoint for Polka (e.g. to handle events and unlock “red” chirps)

### Admin

- GET /admin/metrics – view basic request metrics (e.g. requests count)
- POST /admin/reset – reset server state (for local development/testing)


## Project Status

This is a guided learning project from Boot.dev, built to practice web development concepts.
Not intended for production use.
