package ratelimiter

import (
	"errors"
	"fmt"
    "net"
	"net/http"
    "strings"
	"time"
	"github.com/afonsocraposo/go-rate-limiter/internal/helpers"
)

type Handler func(w http.ResponseWriter, r *http.Request)

func RateLimiter(fn Handler, capacity int, rate int) Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        ip, err := getIp(r)
        if err != nil {
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }
        path := r.URL.EscapedPath()
        limited := isRequestLimited(ip, path, capacity, rate)
        if limited {
            fmt.Printf("%s - %s - blocked IP %s\n", getUTC(), path, ip)
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        } else {
            fn(w,r)
            return
        }
    }
}

func isRequestLimited(ip string, path string, capacity int, rate int) bool {
    key := fmt.Sprintf("%s:%s", path, ip)

    bucket, err := helpers.GetBucket(key)
    if err != nil {
        return true
    }

    if bucket.Timestamp == 0 {
        bucket.Timestamp = getTime();
        bucket.Tokens = capacity - 1
    }else{
        t := getTime()
        elapsedTime := t - bucket.Timestamp
        tokens := int(elapsedTime)/rate

        if tokens > 0 {
            bucket.Tokens += tokens
            bucket.Timestamp = t
            if bucket.Tokens > capacity {
                bucket.Tokens = capacity
            }
        }

        if bucket.Tokens == 0 {
            return true
        }

        bucket.Tokens--
    }

    helpers.SetBucket(key, bucket, capacity * rate)

    return false

}

func getIp(r *http.Request) (string, error) {
    forwarded := r.Header.Get("X-FORWARDED-FOR")
    if forwarded != "" {
        ips := strings.Split(forwarded, ",")
        if len(ips) > 0 {
            return strings.TrimSpace(ips[0]), nil
        }
    }
    remoteAddr := r.RemoteAddr
    if remoteAddr != "" {
        host, _, err := net.SplitHostPort(remoteAddr)
        if err != nil {
            // In case of an error (which means RemoteAddr did not contain a port),
            // simply return the RemoteAddr as the IP.
            return remoteAddr, nil
        }
        return host, nil
    }
    return "", errors.New("Unable to get IP")
}

func getTime() int64 {
    return time.Now().UnixMilli();
}

func getUTC() string {
    t := time.Now().UTC()
    formattedTime := t.Format("2006-01-02 15:04:05.000")
    return formattedTime
}
