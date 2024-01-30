
# KlipTopia API
## Overview

The KlipTopia API provides the backend infrastructure for a clipboard tool named KlipTopia, allowing seamless syncing of clipboard data across multiple devices. This API is designed to be the backbone for a reliable and efficient clipboard synchronization service.

## Features

- **Cross-Device Clipboard Sync:** Sync your clipboard data across various devices in real-time.
- **Secure Data Transmission:** Ensure secure and encrypted communication between devices.
- **User Authentication:** Secure your clipboard data by implementing user authentication.
- **RESTful API:** Follows REST principles for easy integration with various clients.

## Getting Started

These instructions will help you set up the project on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Ensure you have the following software installed on your machine before setting up KlipTopia API:

- [Go](https://golang.org/) (at least version 1.21.5)
- Database (e.g., [PostgreSQL](https://www.postgresql.org/))
- [RabbitMQ](https://www.rabbitmq.com/)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/PreciousNyasulu/KlipTopia-api.git
   ```

2. Install dependencies:

   ```bash
   cd KlipTopia-api
   go mod download
   ```

3. Run Docker containers

    ```bash
    docker compose up -d
    ```

4. Set up the configuration:

   ```bash
   cp .env.example .env
   ```

   Update the `.env` file with your configuration details.

5. Run the application:

   ```bash
   go run cmd/main.go
   ```

5. The API will be available at `http://localhost:9000`.

## API Documentation

Explore the API endpoints and learn how to interact with the KlipTopia clipboard sync service. See [API Documentation](docs/API_DOCUMENTATION.md).

## Contributing

KlipTopia is an open source project. You can contribute to the project by reporting bugs, suggesting features, or submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Hat tip to anyone whose code was used
- Inspiration