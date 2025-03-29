# Docker Configuration

This directory contains Docker-related configuration files for the go-toolkit project, designed for quickly deploying development and testing environments.

## File Description

- `docker-compose.yml` - Docker Compose configuration file for starting Redis, Prometheus, and Grafana services
- `prometheus.yml` - Prometheus configuration file defining monitoring targets and scrape intervals
- `grafana/dashboards/` - Directory containing Grafana dashboard configurations

## Quick Start

### Start All Services

```bash
docker-compose up -d
```

### Start Only Redis (for testing)

```bash
docker-compose up -d redis
```

### Stop All Services

```bash
docker-compose down
```

## Accessing Services

- Redis: localhost:6379
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (default username/password: admin/admin)

## Grafana Setup

After logging into Grafana for the first time, you'll need to configure the Prometheus data source:

1. Visit http://localhost:3000 and log in using the default credentials
2. Go to Configuration > Data Sources > Add data source
3. Select Prometheus
4. Set the URL to http://prometheus:9090
5. Click "Save & Test"

Then you can import the preconfigured dashboard:

1. Go to Create > Import
2. Upload the `grafana/dashboards/ratelimit-dashboard.json` file
3. Select the previously created Prometheus data source
4. Click "Import"

## Monitoring Metrics

Key monitoring metrics include:

- `http_requests_total` - Total number of HTTP requests
- `http_request_duration_seconds` - HTTP request duration
- `ratelimit_total` - Number of rate-limited requests