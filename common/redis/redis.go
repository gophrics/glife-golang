package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// Instance : Singleton Instance
var Instance *redis.Client

// GeoLocation : redis.GeoLocation
type GeoLocation = redis.GeoLocation

// GeoRadiusQuery : redis.GeoRadiusQuery
type GeoRadiusQuery = redis.GeoRadiusQuery

func openDB() {
	Instance = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func init() {
	openDB()
	go healthChecks()
}

func healthChecks() {
	for true {
		if Instance == nil {
			openDB()
		}
		time.Sleep(100 * time.Millisecond)
	}
}
