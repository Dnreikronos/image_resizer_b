# Image Resizer

An image resizing application built with Go, Docker, Elasticsearch, Logstash, and Kibana for logging and monitoring. The app processes images by resizing them according to specified dimensions and stores the logs in Elasticsearch.

## Features

- Resize images to custom dimensions
- Logs image processing requests to Elasticsearch
- Configured with Docker for easy deployment
- Uses Logstash for log processing and Kibana for visualization and monitoring

## Project Structure

```
image_resizer_b/
├── cmd/
│   └── main.go                # Entry point of the application
├── configs/
│   └── config.go              # Configuration settings
├── db/
│   ├── connection/
│   │   └── connection.go      # Database connection logic
│   └── migration/
│       └── migration.go       # Database migration logic
├── handlers/
│   └── handlers.go            # HTTP request handlers
├── models/
│   └── image.go               # Image model and business logic
├── tests/
│   └── image_resizer_test.go  # Unit tests
└── utils/
    └── validate.go            # Utility functions
```

## Requirements

- Go 1.18 or higher
- Docker and Docker Compose
- Elasticsearch, Logstash, and Kibana (ELK stack)
- PostgreSQL (for database)

## Setup

### 1. Clone the repository

```bash
git clone https://github.com/Dnreikronos/image_resizer.git
cd image_resizer
```

### 2. Set up environment variables

Create a `.env` file at the root of your project and set the following variables:

```
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
```

### 3. Start the application with Docker Compose

Use Docker Compose to build and start the services:

```bash
docker-compose up --build
```

This will start Elasticsearch, Logstash, Kibana, and your application.

### 4. Access Kibana

Once the services are up, you can access Kibana at [http://localhost:5601](http://localhost:5601) to monitor your logs and use the Discover panel to explore your image processing logs.

### 5. Access the Application

Your application will be available at `http://localhost:9090`.

## Usage

1. Send a POST request to resize an image:
    - Endpoint: `/resize`
    - Method: `POST`
    - Body: `{ "image_url": "<URL_OF_IMAGE>", "width": 200, "height": 200 }`
  
2. The image will be resized, and the process will be logged in Elasticsearch.

## Logs

Logstash processes logs and sends them to Elasticsearch, where you can view and search them using Kibana. You can also view logs in the terminal as they are processed.

## Docker Compose Configuration

Your `docker-compose.yml` includes the following services:

- **elasticsearch**: The Elasticsearch service for storing and indexing logs.
- **logstash**: A Logstash service for processing logs and sending them to Elasticsearch.
- **kibana**: A Kibana service for visualizing and exploring logs in real-time.
- **resizer**: The main application service that resizes images and handles HTTP requests.
- **pgadmin**: A PostgreSQL admin interface for managing the database.
- **postgres**: A PostgreSQL service for database storage.

## Testing

Run the tests for the image resizer application:

```bash
make test
```

## Improvements

Feel free to improve this project by:

- Adding more image processing features (e.g., cropping, rotating).
- Enhancing logging and error handling.
- Optimizing the Docker setup for production use.
- Implementing more robust database migration handling.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Go](https://golang.org/) - The programming language used to build the application.
- [Docker](https://www.docker.com/) - Containerization platform used for deployment.
- [Elasticsearch](https://www.elastic.co/elasticsearch/) - Search engine for storing and searching logs.
- [Logstash](https://www.elastic.co/logstash/) - Log processing pipeline.
- [Kibana](https://www.elastic.co/kibana/) - Visualization tool for analyzing and monitoring logs.
- [PostgreSQL](https://www.postgresql.org/) - Relational database management system for application data.
