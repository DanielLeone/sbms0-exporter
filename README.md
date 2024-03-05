A prometheus exporter for the https://electrodacus.com/ SBMS0

### how to run it (in docker)

```shell
URL='IP_OF_SBMS0' docker compose up --build
```

### scrape using Prometheus

```shell
curl localhost:9000/metrics
```

### how does it work?

Calls the `rawData` endpoint of the SBMS0 device (requires Wifi Module)
and parses all the data; similar to that of the `legacy` html page
