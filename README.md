# redis_mon_exporter
Prometheus exporter for monitoring redis servers (connection, write, etc)

Build as docker image:

```
git clone https://github.com/romanprog/redis_mon_exporter.git

cd redis_mon_exporter/

docker build -t redis_mon_exporter .
```

Run:

```
docker run -d -p 8080:8080 redis_mon_exporter:latest -redis.servers="redis-host:6379,redis-host2:6378"
```

Test:
```
curl localhost:8080/metrics
```

Run in swarm (compose): 
```
version: '3.6'

services:
  redis-mon-exporter:
    image: artiloop/redis_mon_exporter
    environment:
        REDIS_SERVERS: "redis-node1:6379,redis-node2:6379"
        LISTEN_PORT: "8081"

```

Prometheus config:
```
  - job_name: redis_exporter
    static_configs:
      - targets: ['redis-mon-exporter:8081']
```

