# CRUD Go SQLite

A simple CRUD (Create, Read, Update, Delete) web application built with Go and SQLite for managing users.

## Features

- List all users
- Create a new user
- Edit an existing user
- Delete a user
- Web interface with HTML templates and static CSS

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/mat1520/crud-go-sqlite.git
   ```

2. Navigate to the project directory:
   ```
   cd crud-go-sqlite
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

4. Run the application:
   ```
   go run main.go
   ```

The server will start on `http://localhost:8080`.

## Usage

- Open your browser and navigate to `http://localhost:8080`.
- Use the web interface to perform CRUD operations on users.

## Technologies Used

- **Go**: Backend language
- **SQLite**: Database
- **HTML Templates**: For rendering views
- **CSS**: For styling

## Project Structure

- `main.go`: Main application file
- `templates/`: HTML templates for the web interface
- `static/`: Static files (CSS, etc.)
- `users.db`: SQLite database file

## License

This project is licensed under the MIT License.
