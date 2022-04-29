package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	_ "encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
	"topsc/mongo/maindbv2"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	// cors "github.com/itsjamie/gin-cors"
)

type Gqlquery struct {
	Query string `json:"query"`
}

type AddScore struct {
	Address string `json:"address"`
	Score   string `json:"score"`
	Egg     string `json:"egg"`
}
type ReturnGetTop100 struct {
	Address string `json:"address"`
	Score   string `json:"score"`
	Egg     string `json:"egg"`
}
type GetRank struct {
	Address string `json:"address"`
}

type ReturnGetRank struct {
	Address string `json:"address"`
	Score   string `json:"score"`
	Egg     string `json:"egg"`
}

var (
	dbmain     = `TOPSCORE`
	Collection = `MAIN`
	key        = "ffb264d250fe591f3ea32170b82ed7daa0950c8fb4e274cb6b8ad6cf5c2af307"
)

func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, er1 := hex.DecodeString(keyString)
	if er1 != nil {
		return `err`
	}
	enc, er2 := hex.DecodeString(encryptedString)
	if er2 != nil {
		return `err`
	}
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

func main() {
	r := gin.Default()
	// r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {

		c.JSON(200, "topsc V0-00")
	})

	// r.POST("/test", func(c *gin.Context) {
	// 	var input AddScore
	// 	c.BindJSON(&input)
	// 	// c.ShouldBind(&input)
	// 	fmt.Println(input.Address)
	// 	fmt.Println(input.Score)
	// 	fmt.Println(input.Egg)
	// 	//=========================================
	// 	dbtest := maindbv2.Finddb(c, dbmain, Collection, bson.M{}, "score", -1, 100, 0)
	// 	fmt.Println(dbtest)
	// 	//=========================================
	// 	// var output []AddScore

	// 	// output = append(output, input)

	// 	c.JSON(200, dbtest)
	// })

	r.POST("/AddScore-api", func(c *gin.Context) {
		var input AddScore
		c.BindJSON(&input)
		// c.ShouldBind(&input)
		// fmt.Println(input.Address)
		// fmt.Println(input.Score)
		// fmt.Println(input.Egg)
		//================================================================
		var output string
		output = ""

		AddressDecryptAES := decrypt(input.Address, key)
		if AddressDecryptAES == `err` {
			output = `Decrypt err`
			c.JSON(200, output)
		}
		ScoreDecryptAES := decrypt(input.Score, key)
		if ScoreDecryptAES == `err` {
			output = `Decrypt err`
			c.JSON(200, output)
		}

		fmt.Println(ScoreDecryptAES)

		DBfindADD := maindbv2.Finddb(c, dbmain, Collection, bson.M{"address": AddressDecryptAES}, "_id", 1, 0, 0)

		// fmt.Println(len(DBfindADD))
		if len(DBfindADD) > 0 {

			// fmt.Println(reflect.TypeOf(DBfindADD[0][`address`]))

			NEWsc, err := strconv.ParseFloat(ScoreDecryptAES, 64)
			if err != nil {
				NEWsc = 0
			}

			CURscST := fmt.Sprintf("%f", DBfindADD[0][`score`])

			CURsc, err := strconv.ParseFloat(CURscST, 64)
			if err != nil {
				CURsc = 0
			}

			fmt.Println(CURsc)

			if NEWsc > CURsc {

				NEWscST := fmt.Sprintf("%f", NEWsc)
				sc, err := strconv.ParseFloat(NEWscST, 64)
				if err != nil {
					sc = 0
				}

				updatesocre := bson.M{
					"score": sc,
				}
				// fmt.Println(DBfindADD[0][`address`])
				updatedatadb := maindbv2.UpdateArchive(c, dbmain, Collection, bson.M{"address": AddressDecryptAES}, updatesocre)
				if updatedatadb == `nok` {
					output = `database have some problem`
					c.JSON(200, output)
				} else {
					output = `new high`
				}

			} else {
				output = `not new high`
			}

		} else {
			sc, err := strconv.ParseFloat(ScoreDecryptAES, 64)
			if err != nil {
				sc = 0
			}

			insertsocre := bson.M{
				"address": AddressDecryptAES,
				"score":   sc,
				"egg":     input.Egg,
			}

			DBinsertsocre := maindbv2.Insertdb(c, dbmain, Collection, insertsocre)
			if DBinsertsocre == `nok` {
				output = `database have some problem`
				c.JSON(200, output)
			} else {
				output = `new address`
			}
		}

		c.JSON(200, output)
	})

	r.POST("/GetTop100-api", func(c *gin.Context) {
		var input AddScore
		c.BindJSON(&input)

		//=========================================
		dbtest := maindbv2.Finddb(c, dbmain, Collection, bson.M{}, "score", -1, 100, 0)
		fmt.Println(dbtest)
		//=========================================

		c.JSON(200, dbtest)
	})

	r.POST("/GetRank-api", func(c *gin.Context) {
		var input GetRank
		c.BindJSON(&input)

		//=========================================
		dbtest := maindbv2.Finddb(c, dbmain, Collection, bson.M{"address": input.Address}, "score", -1, 100, 0)
		fmt.Println(dbtest)
		//=========================================

		c.JSON(200, dbtest)
	})

	r.Run(":9105")
	// r.RunTLS(":9105", "./testdata/server.pem", "./testdata/server.key")
}
