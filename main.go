package main

import (
	"github.com/fortune-backend/fortune"
	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
)

func main() {

	sess, err := mgo.Dial("mongodb://localhost:27017/fortune")
	if err != nil {
		panic(err)
	}

	httpRouter := echo.New()

	fortuneRepository := fortune.NewFortuneRepository(sess)
	fortuneService := fortune.NewFortuneService(fortuneRepository)
	fortune.StartFortuneRouter(fortuneService, httpRouter)

	httpRouter.Start(":4000")
}
