package neo4jd

import (
	"log"

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

func init() {
	Instance, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "abc", ""))
	if err != nil {
		log.Fatal(err)
	}

	if Session, err = Instance.Session(neo4j.AccessModeWrite); err != nil {
		log.Fatal(err)
	}
}
