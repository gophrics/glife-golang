package neo4jd

import (
	"log"
	"time"

	neo4j "github.com/issacnitin/neo4j-go-driver"
)

var (
	driver  neo4j.Driver
	session neo4j.Session
	result  neo4j.Result
	err     error
)

var Instance neo4j.Driver
var Session neo4j.Session

func openDB() {
	Instance, err := neo4j.NewDriver("bolt://neo4j:7687", neo4j.BasicAuth("neo4j", "neo4j", ""))
	if err != nil {
		log.Printf("%s", err)
	}

	if Session, err = Instance.Session(neo4j.AccessModeWrite); err != nil {
		log.Printf("%s", err)
	}
}

func init() {
	openDB()
	go healthChecks()
}

func healthChecks() {
	for true {
		if Instance == nil || Session == nil {
			openDB()
		}
		time.Sleep(10 * time.Second)
	}
}
