# Receipt Processor

Build a webservice that fulfils the documented API.[api.yml](./api.yml)

## Setup and Usage

Follow these steps to set up and run the application:

1. **Clone Repository**: Download this repository to your local machine.
   With SSH URL

   ```bash
   git clone git@github.com:samliu98/receipt-processor-challenge.git
   ```

   With HTTP URL

   ```bash
   git clone https://github.com/samliu98/receipt-processor-challenge.git
   ```

2. **Navigate to Repository**: Navigate to the repository folder.

   ```bash
   cd receipt-processor-challeng
   ```

3. **Install Dependencies**: Use the following command to download and install the project's dependencies:

   ```bash
   go mod tidy
   ```

4. **Run Application**: Execute the following command to run the application:

   ```bash
   go run .
   ```

   After running the application, you can test if it's working by calling the `/ping` endpoint. Open a new terminal window and use the following `curl` command:

   ```bash
   curl http://localhost:8080/ping
   ```

   If the application is running, you should see a response like this:

   ```bash
   {"message": "pong"}
   ```

## Testing

### Unit Tests

To verify the application’s performance and reliability, follow these steps to run the unit tests:

1. **Open Terminal**: Open a new terminal window.
2. **Run Unit Tests**: Run the unit tests by typing:
   ```bash
   go test ./...
   ```
3. **Get Test Coverage Rate**: To view test coverage, use:
   ```bash
   go test -cover ./...
   ```

### Manual Testing

For a practical examination of the application's APIs, manual testing can be done using tools like `curl`. Here's a guide to testing the APIs with sample data:

1. Process Receipts:

```bash
curl -X POST -H "Content-Type: application/json" -d '{
   "retailer": "M&M Corner Market",
   ...
}' http://localhost:8080/receipts/process
```

2. Retrieve Points:

```bash
curl http://localhost:8080/receipts/{id}/points
```

Make sure to substitute `{id}` with the actual ID you receive from the "Process Receipts" API.

Enjoy the work!
