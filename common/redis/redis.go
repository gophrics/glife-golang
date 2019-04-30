package redis

import "github.com/go-redis/redis"

// Instance : Singleton Instance
var Instance *redis.Client

// GeoLocation : redis.GeoLocation
type GeoLocation = redis.GeoLocation

func init() {
	Instance = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
