package fortune_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
	"github.com/thiagotrennepohl/fortune-backend/fortune"
	"github.com/thiagotrennepohl/fortune-backend/models"
)

var (
	router            = echo.New()
	fortuneRepository fortune.FortuneRepository
	fortuneService    fortune.FortuneService
	dialInfo          = &mgo.DialInfo{
		Addrs:    []string{"localhost:27017"},
		Timeout:  5 * time.Second,
		Database: "test_fortune_app",
	}
)

func init() {
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal(err.Error())
	}
	fortuneRepository = fortune.NewFortuneRepository(session)
	fortuneService = fortune.NewFortuneService(fortuneRepository)
	fortune.StartFortuneRouter(fortuneService, router)

}

func performRequest(t *testing.T, r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Failed to perform request: %s", err.Error())
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestFortuneHandler_FindRandomMessageWithoutAnyMessages(t *testing.T) {
	resp := performRequest(t, router, "GET", fortune.RandomFortuneEndpoint, nil)
	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestFortuneHandler_SaveFortuneMessage(t *testing.T) {
	requestBody, _ := json.Marshal(validFortuneMessageRequestBody)

	resp := performRequest(t, router, "POST", fortune.SaveNewFortuneMessageEndpoint, requestBody)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestFortuneHandler_SaveFortuneMessageWithAnInvalidBody(t *testing.T) {
	requestBody, _ := json.Marshal("Invalid body")

	resp := performRequest(t, router, "POST", fortune.SaveNewFortuneMessageEndpoint, requestBody)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestFortuneHandler_SaveExistentFortuneMessage(t *testing.T) {
	requestBody, _ := json.Marshal(validFortuneMessageRequestBody)
	resp := performRequest(t, router, "POST", fortune.SaveNewFortuneMessageEndpoint, requestBody)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestFortuneHandler_FindRandomMessageWithOnlyOneMessage(t *testing.T) {
	resp := performRequest(t, router, "GET", fortune.RandomFortuneEndpoint, nil)
	assert.Equal(t, http.StatusOK, resp.Code)

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fortuneMessageResponse := models.FortuneMessage{}

	_ = json.Unmarshal(byteArray, &fortuneMessageResponse)

	assert.Equal(t, validFortuneMessage, fortuneMessageResponse)

}
