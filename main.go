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

var redisServers = flag.String("redis_servers", getEnv("REDIS_SERVERS", "localhost:6379"), "Address list in format 'host:port,host2:port'")

// main function to boot up everything
func main() {
	flag.Parse()
	fmt.Println(*redisServers)

	router := mux.NewRouter()
	router.HandleFunc("/metrics", DoTests).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

func testRedis(server string) int {
	//srvArr := strings.Split(server, ":")
	fmt.Println("Test")
	c, err := redis.Dial("tcp", server)
	if err != nil {
		log.Printf("Could not connect: %v\n", err)
		return 1
	}
	defer c.Close()
	tm, _ := time.ParseDuration("1s")
	ret, err := redis.DoWithTimeout(c, tm, "SET", "gomonitor", "1")
	if (err != nil) || (ret != "OK") {
		return 1
	}
	return 0
}

func DoTests(w http.ResponseWriter, r *http.Request) {
	lineTemplate := "redis_mon_write_check{instance=\"%s\"} %b"
	for _, serverUrl := range strings.Split(*redisServers, ",") {
		res := testRedis(serverUrl)
		metricLine := fmt.Sprintf(lineTemplate, serverUrl, res)
		w.Write([]byte(metricLine))
	}
	return
}
