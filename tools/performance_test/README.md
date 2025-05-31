# Performance Testing with k6

This directory contains a k6 load testing setup to test the performance of various URLs in your application using k6's built-in web dashboard.

## What it tests

The script randomly visits one of the following URL patterns:
- Base URL (homepage)
- Tag pages: `/tag/tag[random number]`
- Post pages: `/test-post-[random number]`
- Pagination: `/?page=[random number]`
- Search results: `/search?t=random_words_[random number]`

## Running the tests

### Prerequisites

- Docker and Docker Compose installed

### Steps to run

1. Start the test with the following command from this directory:

```bash
docker-compose up
```

This will:
- Run the k6 load test
- Start the k6 web dashboard on port 5665
- Automatically generate an HTML report in the `reports` directory when the test completes

2. View the results in the web dashboard at http://localhost:5665

3. When the test completes, find the HTML report at `./reports/k6-report.html`

### Configuration

You can customize the test by setting environment variables:

```bash
docker-compose up
```

### Modifying test parameters

To change the test parameters, edit the `script.js` or `docker-compose.yml` file.

## Web Dashboard Options

The following environment variables can be configured in the docker-compose.yml file:

- `K6_WEB_DASHBOARD`: Enable the web dashboard (default: true)
- `K6_WEB_DASHBOARD_HOST`: Host to bind the web dashboard to (default: 0.0.0.0)
- `K6_WEB_DASHBOARD_PORT`: Port to bind the web dashboard to (default: 5665)
- `K6_WEB_DASHBOARD_PERIOD`: Period in seconds to update the web dashboard (default: 10s)
- `K6_WEB_DASHBOARD_OPEN`: Open the web dashboard in the default browser (default: false)
- `K6_WEB_DASHBOARD_EXPORT`: Filename to export the HTML test report to (default: `/reports/k6-report.html`)

## Generating HTML Reports

The setup is configured to automatically generate an HTML report in the `reports` directory when the test completes. This report is a complete, self-contained file that you can open in any browser and share with your team.

You can also manually generate a report by clicking the "Report" button in the web dashboard UI while the test is running. 