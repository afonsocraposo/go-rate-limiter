package helpers

import (
	"context"
	"fmt"
	"time"
	"github.com/redis/go-redis/v9"
    "encoding/json"
)

var client = redis.NewClient(&redis.Options{
    Addr:	  "127.0.0.1:6379",
    Password: "", // no password set
    DB:		  0,  // use default DB
})
var ctx = context.Background()

type Bucket struct {
    Tokens      int `json:"tokens"`
    Timestamp   int64 `json:"timestamp"`
}

func SetVar(key string, value string, ttl int) error {
    err := client.Set(ctx, key, value, time.Duration(ttl) * time.Millisecond).Err()
    return err
}

func GetVar(key string) ([]byte, error) {
    data, err := client.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return []byte{}, nil
    } else if err != nil {
        fmt.Println(err.Error())
        return []byte{}, err
    }
    return data, nil
}

func GetBucket(ip string) (Bucket, error) {
    data, err := GetVar(ip)
    if err != nil {
        return Bucket{}, err
    } else if (len(data) == 0) {
        return Bucket{}, nil
    } else {
        var bucket Bucket
        err2 := json.Unmarshal(data, &bucket)
        if err2 != nil {
            return Bucket{}, err2
        }
        return bucket, nil
    }
}

func SetBucket(ip string, bucket Bucket, ttl int) error {
    data, err := json.Marshal(bucket)
    if err != nil {
        return err
    }

    value := string(data)

    err2 := SetVar(ip, value, ttl)
    if err2 != nil {
        return err2
    }

    return nil
}
