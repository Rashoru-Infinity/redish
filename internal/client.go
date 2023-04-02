package internal

import (
        "encoding/json"
        "net/http"
        "github.com/go-redis/redis/v8"
)

func GetRedisClient(r *http.Request) *redis.Client {
        var op redis.Options
        js := []byte(r.Header.Get("Redis-Client-Options"))
        if err := json.Unmarshal(js, &op); err != nil {
                return nil
        }
        return redis.NewClient(&op)
}
