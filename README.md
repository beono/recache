
[![Build Status](https://travis-ci.org/beono/recache.svg?branch=master)](https://travis-ci.org/beono/recache)
[![GoDoc](https://godoc.org/github.com/beono/recache?status.svg)](https://godoc.org/github.com/beono/recache)
[![Go Report Card](https://goreportcard.com/badge/github.com/beono/recache)](https://goreportcard.com/report/github.com/beono/recache)

# About the project

This is an experimental library.
It has a simple interface that lets you to cache data using redis and invalidate cache by tags.

If you are not familiar with invalidation by tags, then I suggest you reading this article https://symfony.com/blog/new-in-symfony-3-2-tagged-cache.

## How to use

```go
package main

import "github.com/go-redis/redis"
import "github.com/beono/recache"

func main() {

    // initialize go-redis client
    cl = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    recache = recache.NewRedisCache(cl)

    result, err := recache.Get("orders_by_user_1");
    if err == cache.ErrKeyNotFound {
        // get data from the database here
        result := getOrdersByUserID(1)

        // we can tag this cache entry.
        // if user bought iPad then we can use this tag,
        // so later we can invalidate this cache when we update ipad entity
        if err := recache.Set("orders_by_user_1", result, 3600, "orders", "ipad"); err != nil {
            t.Errorf("unexpected error: %q", err)
        }
    }

    if err != nil {
        panic(err)
    }

    fmt.Println(result)
}
```