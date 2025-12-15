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

    http://localhost:8080

## Project Status

This is a guided learning project from Boot.dev, built to practice web development concepts.
Not intended for production use.
