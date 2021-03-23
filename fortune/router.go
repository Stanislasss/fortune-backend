package fortune

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/thiagotrennepohl/fortune-backend/models"
)

const (
	RandomFortuneEndpoint         = "/v1/fortune/random"
	AllFortuneEndpoint         = "/v1/fortune/all"
	SaveNewFortuneMessageEndpoint = "/v1/fortune"
	HomeEndpoint                  = "/"
	homepageTemplate              = "homepage.html"
)

var (
	assetsFolder string
)

type FortuneRouter struct {
	fortuneService FortuneService
	templates      *template.Template
}

func (t *FortuneRouter) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func StartFortuneRouter(fortuneService FortuneService, templates *template.Template, router *echo.Echo) {

	fortuneRouter := &FortuneRouter{
		fortuneService: fortuneService,
		templates:      templates,
	}
	router.Renderer = fortuneRouter
	router.GET(RandomFortuneEndpoint, fortuneRouter.GetRandomFortuneMessage)
	router.POST(SaveNewFortuneMessageEndpoint, fortuneRouter.SaveFortuneMessage)
	router.GET(HomeEndpoint, fortuneRouter.home)
	router.GET(AllFortuneEndpoint, fortuneRouter.GetAll)
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

func (router *FortuneRouter) GetAll(ctx echo.Context) error {
	messages, err := router.fortuneService.GetAll()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.Json{"Message": err.Error()})
	}
	return ctx.JSON(http.StatusOK, messages)
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

func (router *FortuneRouter) home(ctx echo.Context) error {
	randomMessage, err := router.fortuneService.FindRandom()
	if err != nil {
		if _, ok := err.(*models.ErrNotFound); ok {
			return ctx.Render(http.StatusOK, homepageTemplate, models.Json{"Message": "No messages found :("})
		}
		return ctx.Render(http.StatusInternalServerError, homepageTemplate, "Ops, something went wrong, we couldn't get any fortune message :(")
	}
	return ctx.Render(http.StatusOK, homepageTemplate, randomMessage)
}
