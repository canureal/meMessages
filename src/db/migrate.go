package db

import (
	"log"

	"github.com/gocql/gocql"
)

func EnsureKeyspace(hosts string, keyspaces string) {
	cluster := gocql.NewCluster(hosts)
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	err = session.Query(`
		CREATE KEYSPACE IF NOT EXISTS ` + keyspaces + ` WITH REPLICATION = {'class': 'NetworkTopologyStrategy', 'replication_factor':  1} 
	`).Exec()
	if err != nil {
		log.Fatal(err)
	}
	/*
		if no error, we return ensured keyspace
	*/
	log.Printf("keyspace ensured: %v", keyspaces)
}

func Migrate(session *gocql.Session) {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS messages_by_conversation (
			conversation_id uuid,
			message_id timeuuid,
			sender_id uuid,
			content text,
			PRIMARY KEY (conversation_id, message_id)
		) WITH CLUSTERING ORDER BY (message_id DESC)`,
		
		`CREATE TABLE IF NOT EXISTS users_by_id (
			user_id uuid PRIMARY KEY,
			username text,
			created_at timestamp,
			password text
		)`,
		
		`CREATE TABLE IF NOT EXISTS users_by_username (
			username text PRIMARY KEY,
			user_id uuid,
			created_at timestamp,
			password text
		)`,
	}

	for _, statement := range statements {
		if err := session.Query(statement).Exec(); err != nil {
			log.Fatal(err)
		}
	}

	// if no problem we log it
	log.Printf("migrations applied!")
}
