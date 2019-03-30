package fortune

import (
	"net/http"

	"github.com/thiagotrennepohl/fortune-backend/models"
	"github.com/labstack/echo"
)

const (
	RandomFortuneEndpoint         = "/v1/fortune/random"
	SaveNewFortuneMessageEndpoint = "/v1/fortune"
)

type FortuneRouter struct {
	fortuneService FortuneService
}

func StartFortuneRouter(fortuneService FortuneService, router *echo.Echo) {

	fortuneRouter := &FortuneRouter{
		fortuneService: fortuneService,
	}

	router.GET(RandomFortuneEndpoint, fortuneRouter.GetRandomFortuneMessage)
	router.POST(SaveNewFortuneMessageEndpoint, fortuneRouter.SaveFortuneMessage)
}

func (router *FortuneRouter) GetRandomFortuneMessage(ctx echo.Context) error {
	message, err := router.fortuneService.FindRandom()
	if err != nil {
		if _, ok := err.(*models.ErrNotFound); ok {
			return ctx.JSON(http.StatusNoContent, models.Json{"Message": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, models.Json{"Message": err.Error()})
	}
	return ctx.JSON(http.StatusOK, message)
}

func (router *FortuneRouter) SaveFortuneMessage(ctx echo.Context) error {
	requestBody := models.FortuneMessage{}
	err := ctx.Bind(&requestBody)
	if err != nil {
		return router.badRequestResponse(err, ctx)
	}

	err = router.fortuneService.Save(requestBody)
	if err != nil {
		if _, ok := err.(*models.ErrMessageAlreadyExists); ok {
			return router.badRequestResponse(err, ctx)
		}
		return router.internalServerErrorResponse(err, ctx)
	}

	return router.statusOKResponse(models.Json{"Message": "ok"}, ctx)
}

func (router *FortuneRouter) badRequestResponse(err error, ctx echo.Context) error {
	return ctx.JSON(http.StatusBadRequest, models.Json{"Message": err.Error()})
}

func (router *FortuneRouter) internalServerErrorResponse(err error, ctx echo.Context) error {
	return ctx.JSON(http.StatusInternalServerError, models.Json{"Message": err.Error()})
}

func (router *FortuneRouter) statusOKResponse(message models.Json, ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, message)
}
