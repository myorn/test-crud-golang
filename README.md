
## How to run the project

```shell
cd compose && sudo docker-compose up -d && cd -
```

## How to test the project

First you need to install `jq` to parse the JSON response from the API.
```shell
sudo apt install jq
```
or for macOS
```shell
brew install jq
```

Then you can run the test script:
```shell
cd ./curl && chmod +x test.sh && ./test.sh && cd -
```
And check if the output is as expected.
```shell
sudo docker-compose logs equipment-registry-service_1
```
And check if messages are being processed by the event monitor service.
```shell
sudo docker exec -it $(sudo docker ps --filter "name=compose_event-monitor-service" --format "{{.ID}}") cat /app/events.log
```

The project uses the following technologies:

- **Docker Compose**: To define and run multi-container Docker applications. The `docker-compose.yml` file specifies services like `postgres`, `rabbitmq`, `equipment-registry-service`, and `event-monitor-service`.
- **PostgreSQL**: A relational database service used to store data.
- **RabbitMQ**: A message broker service used for communication between services.
- **Go**: The programming language used to build the `equipment-registry-service` and `event-monitor-service`.

### Equipment Registry Service
This service is responsible for managing equipment data. It interacts with the PostgreSQL database to perform CRUD (Create, Read, Update, Delete) operations on equipment records.

### Event Monitor Service
This service monitors events related to equipment. It relies on RabbitMQ for message brokering to handle event-driven communication and processing.