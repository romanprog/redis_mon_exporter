package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

var redisServers = flag.String("redis.servers",
	getEnv("REDIS_SERVERS", "localhost:6379"),
	"Address list in format 'host:port,host2:port'")

var listenPort = flag.String("listen.port",
	getEnv("LISTEN_PORT", "8080"),
	"Listened port.")

func main() {
	flag.Parse()
	startListen()
}

func startListen() {
	router := mux.NewRouter()
	router.HandleFunc("/metrics", DoChecks).Methods("GET")
	listenUrl := fmt.Sprintf("0.0.0.0:%s", *listenPort)
	log.Printf("Runing listener on %s", listenUrl)
	log.Fatal(http.ListenAndServe(listenUrl, router))
}

func getEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

func checkRedisServer(server string) int {
	c, err := redis.Dial("tcp", server)
	if err != nil {
		log.Printf("Could not connect: endpoint: %s. Error: %v\n", server, err)
		return 1
	}
	defer c.Close()
	tm, _ := time.ParseDuration("1s")
	ret, err := redis.DoWithTimeout(c, tm, "SET", "gomonitor", "1")
	if (err != nil) || (ret != "OK") {
		log.Printf("Could not write: endpoint: %s. Error: %v\n", server, err)
		return 1
	}
	log.Printf("Check OK: endpoint: %s. Result: %s", server, ret)
	return 0
}

func DoChecks(w http.ResponseWriter, r *http.Request) {
	lineTemplate := "redis_mon_write_check{addr=\"%s\"} %b\n"
	log.Printf("Request from host: %s", r.Host)
	for _, serverUrl := range strings.Split(*redisServers, ",") {
		res := checkRedisServer(serverUrl)
		metricLine := fmt.Sprintf(lineTemplate, serverUrl, res)
		w.Write([]byte(metricLine))
	}
}
