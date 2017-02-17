package main
import (
    "log"
    "fmt"
    "time"
    "strconv"
    "go_redis_cluster"
    "net/http"
)

const kNumOfRoutine = 50

func redisHandler(w http.ResponseWriter, r *http.Request) {
    cluster, err := redis.NewCluster(
    &redis.Options{
        StartNodes: []string{"127.0.0.1:6380", "127.0.0.1:6381", "127.0.0.1:7382"},
        ConnTimeout: 50 * time.Millisecond,
        ReadTimeout: 50 * time.Millisecond,
        WriteTimeout: 50 * time.Millisecond,
        KeepAlive: 16,
        AliveTime: 60 * time.Second,
    })

    if err != nil {
    log.Fatalf("redis.New error: %s", err.Error())
    }

    chann := make(chan int, kNumOfRoutine)
    for i := 0; i < kNumOfRoutine; i++ {
    go redisTest(cluster, i * 100000, (i+1)*100000, chann)
    }

    for i := 0; i < kNumOfRoutine; i++ {
    _ = <-chann
    }
    fmt.Fprint(w, "redis!")
}

func redisTest(cluster *redis.Cluster, begin, end int, done chan int) {
    prefix := "{mykey}"
    for i := begin; i < end; i++ {
        key := prefix + strconv.Itoa(i)

        _, err := cluster.Do("set", key, i*10)
        if err != nil {
            fmt.Printf("-set %s: %s\n", key, err.Error())
            time.Sleep(100 * time.Millisecond)
            continue
    }
    value, err := redis.Int(cluster.Do("GET", key))
    if err != nil {
        fmt.Printf("-get %s: %s\n", key, err.Error())
        time.Sleep(100 * time.Millisecond)
        continue
    }
    if value != i*10 {
        fmt.Printf("-mismatch %s: %d\n", key, value)
        time.Sleep(100 * time.Millisecond)
        continue
    }
    fmt.Printf("+set %s\n", key)
    time.Sleep(50 * time.Millisecond)
    }

    done <- 1
}


