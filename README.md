# goexpert-load-test

## Challenge description:
Objective: Create a CLI system in Go to perform load testing on a web service.

The user should provide the following parameters via CLI:

- `--url`: URL of the service to be tested.
- `--requests`: Total number of requests.
- `--concurrency`: Number of simultaneous calls.

### Test Execution:

- Perform HTTP requests to the specified URL.
- Distribute the requests according to the defined concurrency level.
- Ensure that the total number of requests is met.

### Report Generation:

At the end of the tests, present a report containing:

- Total time spent on execution.
- Total quantity of requests performed.
- Quantity of requests with HTTP status 200.
- Distribution of other HTTP status codes (such as 404, 500, etc.).

### Application Execution:

The application can be executed via docker using the following command:

```
docker run <your docker image> --url=http://example.com --requests=1000 --concurrency=10
```

Replace `<your docker image>` with the name of your Docker image, and adjust the URL, total number of requests, and concurrency level as needed.

### How to run locally with Go
- Clone the repository
- Make sure you have Golang 1.22.3 or higher installed
- Open terminal in project folder
- Run `go run main.go --url https://google.com --requests=100 --concurrency=13`
- You can also specify

### How to run locally with docker
- Clone the repository
- Make sure you have docker installed
- Open terminal in project folder
- Run `docker build -t goexpert-load-test .`
- Run `docker run goexpert-load-test --url https://google.com --requests=100 --concurrency 13`

### Postman mock server
- https://d0a17a76-061a-4e50-b661-74c4bad7b7f6.mock.pstmn.io
- Default return is 200 but you can try /400 and /500 routes for respectively Bad Request and Internal Server Error status codes