# KlipTopia API Documentation

## Authentication

### Login
- Endpoint: `POST /api/auth/login`

- Description: Authenticate and obtain an access token.
- Request:
  - Body:

    ```json
    {
        "username": "user@example.com",
        "password": "securepassword"
    }
    ```

- Response:
  - Success (200 OK):
  
    ```json
    {
        "message": "success",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```

  - Error (401 Unauthorized):

    ```json
    {
        "message": "Invalid credentials"
    }
    ```

### Logout
- Endpoint: `POST /api/auth/logout`
- Description: Invalidate the current access token.
- Request:
  - Headers:
    - `Authorization: Bearer [token]`
- Response:
  - Success (200 OK):
  
    ```json
    {
        "message": "Logout successful"
    }
    ```

### Register
- Endpoint: `POST /api/auth/register`
- Description: Create a new user account.
- Request:
  - Body:

    ```json
    {
        "username": "newuser",
        "email": "newuser@example.com",
        "password": "securepassword"
    }
    ```

- Response:
  - Success (201 Created):

    ```json
    {
        "message": "User registered successfully"
    }
    ```

  - Error (400 Bad Request):

    ```json
    {
        "message": "Validation failed"
    }
    ```

## Clipboard

### Copy
- Endpoint: `POST /api/clipboard/copy`
- Description: Copy content to the clipboard.
- Request:
  - Headers:
    - `Authorization: Bearer [token]`
  - Body:

    ```json
    {
        "deviceIpAddress":"127.0.0.1",
        "content":"content",
        "contentType":"text"
    }
    ```

- Response:
  - Success (200 OK):

    ```json
    {
        "message":"content published"
    }
    ```

  - Error (401 Unauthorized):

    ```json
    {
        "message": "Unauthorized"
    }
    ```

