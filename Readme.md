# adLaCarte History
Queries data from https://adlacarte.adesso.de/ periodically for Prometheus.

## Setup
### Go application
Encode `<username>:<password>` in Base64 and add it to the `credentials` file.
Then you can run the following command to see current values as well as expose them on http://localhost:8080 for Prometheus:
```
$ go run adLaCarteHistory.go
2020/01/08 09:50:46 | Entenhaus | ChiliPeppers | PiDoe | ChinaImbissBUI | PizzeriaMammaMia | Pinnochio | PizzariabeiMarco |
2020/01/08 09:50:47 |      0.00 |         0.00 |  0.00 |           0.00 |             0.00 |      0.00 |             0.00 |
2020/01/08 09:55:47 |      0.00 |         0.00 |  0.00 |           0.00 |             0.00 |      0.00 |             0.00 |
2020/01/08 10:00:47 |      0.00 |         0.00 |  0.00 |           0.00 |             0.00 |      0.00 |             0.00 |
```

### Using Docker
The same can be achieved using the Docker image of this repo:
```
docker run -v $PWD/credentials:/credentials docker.pkg.github.com/juliansauer/adlacartehistory/adlacarte-history:latest
2020/01/08 09:50:46 | Entenhaus | ChiliPeppers | PiDoe | ChinaImbissBUI | PizzeriaMammaMia | Pinnochio | PizzariabeiMarco |
2020/01/08 09:50:47 |      0.00 |         0.00 |  0.00 |           0.00 |             0.00 |      0.00 |             0.00 |
2020/01/08 09:55:47 |      0.00 |         0.00 |  0.00 |           0.00 |             0.00 |      0.00 |             0.00 |
2020/01/08 10:00:47 |      0.00 |         0.00 |  0.00 |           0.00 |             0.00 |      0.00 |             0.00 |
```

### Docker Compose
The provided docker-compose file will start the crawler as well as Prometheus and Grafana instances.
Edit the `credentials` file as described above and add a password for Grafana in the `.env` file.
Then run `docker-compose up -d` to access them:

| Service           | Port                          |
| ----------------- | ----------------------------- |
| adLaCarte History | [8080](http://localhost:8080) |
| Prometheus        | [9090](http://localhost:9090) |
| Grafana           | [3000](http://localhost:3000) |
