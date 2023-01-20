package aes21

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"os"

	"github.com/CarlosMore29/env_cm"
	"github.com/CarlosMore29/logs_cm"
	"github.com/sirupsen/logrus"
)

// Logger
var logger *logrus.Logger

func init() {
	env_cm.GetEnvFile()
	logger = logs_cm.NewLogger()
}

func Decrypt(data []byte) (textDecrypt []byte, errorGlobal error) {

	key, _ := hex.DecodeString(os.Getenv("KEY_ENCRYPT"))
	// logger.Info("Key:", key)

	// Sacamos el AAD
	aadVector := data[3:9]
	//logger.Info("Add Vector:", aadVector)

	// Sacamos el TAG Vector
	tagVector := data[9:25]
	// logger.Info("tagVector: ", tagVector)

	// texto Cifrado
	textEncrypt := data[25:]
	// logger.Info("textEncrypt", textEncrypt)

	// Obtener el ivVector
	aadData := sha256.New()
	aadData.Write(aadVector)
	aadSha := aadData.Sum(nil)
	// logger.Info("aadData: ", aadSha)

	// ivVector
	auxVector0 := aadSha[0:16]
	auxVector := aadSha[16:]

	var arrayVector = make([]byte, 16)
	for i := 0; i < 16; i++ {
		arrayVector[i] = auxVector0[i] * auxVector[i]
	}

	ivVector := arrayVector[:12]
	// logger.Info("ivVector: ", ivVector)

	// Dencrypt
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error("NewCipher: ", err.Error())
		panic("NewCipher: " + err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error("cipher.NewGCM: ", err.Error())
		panic("cipher.NewGCM: " + err.Error())
	}

	textDecrypt, errorGlobal = aesgcm.Open(nil, ivVector, append(textEncrypt, tagVector...), aadVector)
	if err != nil {
		logger.Error("aesgcm.Open: ", err.Error())
		panic("aesgcm.Open: " + err.Error())
	}

	// logger.Info(textDecrypt)
	return
}

func Encrypt(phrase []byte) (cipherTrama []byte, errorGlobal error) {
	// Obtener la Key
	key, _ := hex.DecodeString(os.Getenv("KEY_ENCRYPT"))
	// logger.Info("key: ", key)

	// Generar un aad
	aadVector := make([]byte, 6)
	if _, err := rand.Read(aadVector); err != nil {
		panic("No se pudo generar un numero random")
	}
	// logger.Info("aadVector: ", aadVector)

	// Generar el (IV) Vector de inicializacion
	aadData := sha256.New()
	aadData.Write(aadVector)
	aadSha := aadData.Sum(nil)
	// logger.Info("aadData: ", aadSha)

	auxVector0 := aadSha[0:16]
	auxVector := aadSha[16:]

	var arrayVector = make([]byte, 16)
	for i := 0; i < 16; i++ {
		arrayVector[i] = auxVector0[i] * auxVector[i]
	}

	ivVector := arrayVector[:12]
	// logger.Info("ivVector: ", ivVector)

	// Encrypt
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error("Error NewCipher: ", err.Error())
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error("Error NewGCM: ", err.Error())
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, ivVector, phrase, aadVector)
	// logger.Info("%x\n", ciphertext)

	// La encrypt da C - Tag y lo queremos Tag - C EL TAMAÑO de TAG ES DE 16
	textEncrypt := ciphertext[:len(ciphertext)-16]
	tagEncrypt := ciphertext[len(ciphertext)-16:]

	cipherTrama = createTrama(append(tagEncrypt, textEncrypt...), aadVector)

	return
}

func createTrama(ciphertext []byte, aadVector []byte) (trama []byte) {
	// Agrego el header
	headEncrypt := []byte{0, 0, 22}
	// Agrego el AAD
	trama = append(headEncrypt, aadVector...)
	//Agrego el texto encriptado
	trama = append(trama, ciphertext...)
	return
}
