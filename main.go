package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/CarlosMore29/icedoor_go/cassandra"
	"github.com/CarlosMore29/icedoor_go/model"
	"github.com/google/uuid"

	"github.com/CarlosMore29/env_cm"
	"github.com/CarlosMore29/logs_cm"
	"github.com/sirupsen/logrus"
)

// Logger
var logger *logrus.Logger

// Cassnadra
var cosmosCassandraContactPoint string
var cosmosCassandraPort string
var cosmosCassandraUser string
var cosmosCassandraPassword string
var cosmosCassandraKeySpace string

func init() {
	logger = logs_cm.NewLogger()
	logger.Info("Inicializacion: Server TCP")

	envGLobals()

}

func main() {

	session, errSession := cassandra.GetSession(cosmosCassandraContactPoint, cosmosCassandraPort, cosmosCassandraUser, cosmosCassandraPassword)
	if errSession != nil {
		logger.Panic(errSession)
	}

	defer session.Close()

	// Listen for incoming connections.
	addr := os.Getenv("SERVER") + ":" + os.Getenv("PORT")
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Panic(err)
	}
	defer l.Close()
	host, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on host: %s, port: %s\n", host, port)

	var pt int = 0

	for {
		// Listen for an incoming connection
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		// Handle connections in a new goroutine
		go func(conn net.Conn) {
			buf := make([]byte, 1024)
			len, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("Error reading: %#v\n", err)
				return
			}
			// fmt.Printf("Message received: %s\n", string(buf[:len]))

			// time.Sleep(8 * time.Second)

			// cassandra.FindAllCassandra(cosmosCassandraKeySpace, "timeline", session)

			var testObj model.Timeline = model.Timeline{
				ID:   uuid.New().String(),
				Data: string(buf[:len]),
				Date: time.Now(),
			}

			// logger.Info(testObj)
			_, error_insert := cassandra.InsertTestCassandra(cosmosCassandraKeySpace, "timeline", session, testObj)

			if error_insert != nil {
				fmt.Println("error_insert: ", error_insert)
			} else {
				// fmt.Println(created)
			}

			pt += 1
			fmt.Println("entradas: ", pt)
			conn.Write([]byte(strconv.Itoa(pt)))
			conn.Close()
		}(conn)
	}
}

func envGLobals() {
	env_cm.GetEnvFile()
	cosmosCassandraContactPoint = os.Getenv("CASSANDRA_HOST")
	cosmosCassandraPort = os.Getenv("CASSANDRA_PORT")
	cosmosCassandraUser = os.Getenv("CASSANDRA_USER")
	cosmosCassandraPassword = os.Getenv("CASSANDRA_PASSWORD")
	cosmosCassandraKeySpace = os.Getenv("CASSANDRA_KEYSPACE")
}
