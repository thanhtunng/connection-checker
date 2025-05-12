# Connection Checker

A web application for checking port availability across multiple IP addresses.

## Features

- Check port availability for multiple IP addresses
- Beautiful and responsive UI
- Real-time results with color-coded status
- Support for comma or newline-separated IP addresses

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1. Clone the repository:
```bash
git clone <repository-url>
cd connection-checker
```

2. Build and start the containers:
```bash
docker-compose up --build
```

3. Access the application:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

## Usage

1. Open your web browser and navigate to http://localhost:3000
2. Enter IP addresses in the text area (separated by commas or newlines)
3. Enter the port number you want to check
4. Click "Check Ports" to start the scan
5. Results will be displayed with:
   - Green boxes for open ports
   - Red boxes for closed ports

## API Endpoint

The backend exposes a single endpoint:

- POST `/api/check-ports`
  - Request body:
    ```json
    {
      "ips": ["192.168.1.1", "10.0.0.1"],
      "port": 80
    }
    ```
  - Response:
    ```json
    {
      "results": [
        {
          "ip": "192.168.1.1",
          "status": "open"
        },
        {
          "ip": "10.0.0.1",
          "status": "closed"
        }
      ]
    }
    ```

## Development

The application consists of two main components:

- Frontend: Vue.js with Tailwind CSS
- Backend: Go HTTP server

To modify the application:

1. Make changes to the source code
2. Rebuild the containers:
```bash
docker-compose up --build
``` 