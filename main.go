package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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

			// Recibir Data del buffer
			buf := make([]byte, 1024)
			len, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("Error reading: %#v\n", err)
				return
			}

			// desencriptar AES 128
			go Decrypt(buf[:len])

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

func Decrypt(data []byte) {
	logger.Info("data: ", data)
	key, _ := hex.DecodeString(os.Getenv("KEY_ENCRYPT"))
	logger.Info("Key:", key)

	// Sacamos el AAD
	aadVector := data[3:9]
	logger.Info("Add Vector:", aadVector)

	// Sacamos el TAG Vector
	tagVector := data[9:25]
	logger.Info("tagVector: ", tagVector)

	// texto Cifrado
	textEncrypt := data[25:]
	logger.Info("textEncrypt", textEncrypt)

	// Obtener el ivVector
	aadData := sha256.New()
	aadData.Write(aadVector)
	aadSha := aadData.Sum(nil)
	logger.Info("aadData: ", aadSha)

	// ivVector
	auxVector0 := aadSha[0:16]
	logger.Info("auxVector0: ", auxVector0)

	auxVector := aadSha[16:]
	logger.Info("auxVector: ", auxVector)

	var arrayVector = make([]byte, 16)
	for i := 0; i < 16; i++ {
		arrayVector[i] = auxVector0[i] * auxVector[i]
	}

	ivVector := arrayVector[:12]
	logger.Info("ivVector: ", ivVector)

	//Dencrypt
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Info("NewCipher: ", err.Error())
		panic("NewCipher: " + err.Error())
	}

	aesgcm, err := cipher.NewGCMWithNonceSize(block, 12)
	if err != nil {
		logger.Info("cipher.NewGCM: ", err.Error())
		panic("cipher.NewGCM: " + err.Error())
	}

	// sz := aesgcm.NonceSize()
	// nonce, cipherText := msg[:sz], msg[sz:]

	textDecrypt, err := aesgcm.Open(nil, ivVector, append(textEncrypt, tagVector...), aadVector)
	if err != nil {
		logger.Info("aesgcm.Open: ", err.Error())
		panic("aesgcm.Open: " + err.Error())
	}

	logger.Info(textDecrypt)

}

func Encrypt(phrase string) (textEncrypt []byte, errorGlobal error) {

	// Generar un aad

	// con un

	return
}

func testAlv() {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	plaintext := []byte("exampleplaintext")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	fmt.Printf("%x\n", ciphertext)
}

func envGLobals() {
	env_cm.GetEnvFile()
	cosmosCassandraContactPoint = os.Getenv("CASSANDRA_HOST")
	cosmosCassandraPort = os.Getenv("CASSANDRA_PORT")
	cosmosCassandraUser = os.Getenv("CASSANDRA_USER")
	cosmosCassandraPassword = os.Getenv("CASSANDRA_PASSWORD")
	cosmosCassandraKeySpace = os.Getenv("CASSANDRA_KEYSPACE")
}
