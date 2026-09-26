package db

import (
	"log"

	"github.com/gocql/gocql"
)

func Connect(hosts string, keyspaces string) *gocql.Session {
	cluster := gocql.NewCluster(hosts)
	cluster.Keyspace = keyspaces
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to scylla")
	return session
}
