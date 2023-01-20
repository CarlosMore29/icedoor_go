package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/CarlosMore29/icedoor_go/aes21"
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
	envGLobals()
	logger = logs_cm.NewLogger()
	logger.Info("Inicializacion: Server TCP")

	data := []byte{104, 111, 108, 97}
	logger.Info("Data sin tratar: ", data)
	phraseEncrypt, _ := aes21.Encrypt(data)
	logger.Info("phraseEncrypt: ", phraseEncrypt)
	phraseDencrypt, _ := aes21.Decrypt(phraseEncrypt)
	logger.Info("phraseDencrypt: ", phraseDencrypt)

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

			// Recibir Data del buffer
			buf := make([]byte, 1024)
			len, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("Error reading: %#v\n", err)
				return
			}

			// desencriptar AES 128
			phraseDencrypt, _ := aes21.Decrypt(buf[:len])
			logger.Info("phraseDencrypt: ", phraseDencrypt)

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
			logger.Info("entradas: ", pt)
			conn.Write([]byte(strconv.Itoa(pt)))
			conn.Close()
		}(conn)
	}
}

func Decrypt(b []byte) {
	panic("unimplemented")
}

func envGLobals() {
	env_cm.GetEnvFile()
	cosmosCassandraContactPoint = os.Getenv("CASSANDRA_HOST")
	cosmosCassandraPort = os.Getenv("CASSANDRA_PORT")
	cosmosCassandraUser = os.Getenv("CASSANDRA_USER")
	cosmosCassandraPassword = os.Getenv("CASSANDRA_PASSWORD")
	cosmosCassandraKeySpace = os.Getenv("CASSANDRA_KEYSPACE")
}
