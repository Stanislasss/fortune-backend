package main

import (
	"os"

	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
	"github.com/thiagotrennepohl/fortune-backend/fortune"
)

var mongoDBConnectionString string

func init() {
	if mongodbconnString, ok := os.LookupEnv("MONGO_ADDR"); ok {
		mongoDBConnectionString = mongodbconnString
	} else {
		mongoDBConnectionString = "mongodb://localhost:27017"
	}
}

func main() {

	sess, err := mgo.Dial(mongoDBConnectionString)
	if err != nil {
		panic(err)
	}

	httpRouter := echo.New()

	fortuneRepository := fortune.NewFortuneRepository(sess)
	fortuneService := fortune.NewFortuneService(fortuneRepository)
	fortune.StartFortuneRouter(fortuneService, httpRouter)

	httpRouter.Start(":4000")
}
