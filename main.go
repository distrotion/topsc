package main

import (
	_ "encoding/json"
	"fmt"
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
)

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

	r.POST("/test", func(c *gin.Context) {
		var input AddScore
		c.BindJSON(&input)
		// c.ShouldBind(&input)
		fmt.Println(input.Address)
		fmt.Println(input.Score)
		fmt.Println(input.Egg)
		//=========================================
		dbtest := maindbv2.Finddb(c, dbmain, Collection, bson.M{}, "score", -1, 100, 0)
		fmt.Println(dbtest)
		//=========================================
		// var output []AddScore

		// output = append(output, input)

		c.JSON(200, dbtest)
	})

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

		DBfindADD := maindbv2.Finddb(c, dbmain, Collection, bson.M{"address": input.Address}, "_id", 1, 0, 0)

		// fmt.Println(len(DBfindADD))
		if len(DBfindADD) > 0 {

			// fmt.Println(reflect.TypeOf(DBfindADD[0][`address`]))

			NEWsc, err := strconv.ParseFloat(input.Score, 64)
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
				updatedatadb := maindbv2.UpdateArchive(c, dbmain, Collection, bson.M{"address": input.Address}, updatesocre)
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
			sc, err := strconv.ParseFloat(input.Score, 64)
			if err != nil {
				sc = 0
			}

			insertsocre := bson.M{
				"address": input.Address,
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
