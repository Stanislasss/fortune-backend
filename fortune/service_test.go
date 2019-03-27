package fortune_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fortune-backend/fortune"
	"github.com/fortune-backend/models"

	"github.com/fortune-backend/mocks"
)

var (
	validFortuneMessageRequestBody = models.FortuneMessage{
		ID:      "jdd8192jnf1924q",
		Message: "CI/CD without testing it's just bug delivery.",
	}
	validFortuneMessage = models.FortuneMessage{
		ID:       "jdd8192jnf1924q",
		Message:  "CI/CD without testing it's just bug delivery.",
		CheckSum: "0ca243cf098339bbae2e10552ff12ae0",
	}
	invalidIDFortuneMessage = models.FortuneMessage{
		ID:      "#44551--!",
		Message: "CI/CD without testing it's just bug delivery.",
	}
	invalidMessageFortuneMessage = models.FortuneMessage{
		ID:      "jdd8192jnf1924q",
		Message: "",
	}
	invalidCheckSumFortuneMessage = models.FortuneMessage{
		ID:       "jdd8192jnf1924q",
		Message:  "CI/CD without testing it's just bug delivery.",
		CheckSum: "123k3jdnf129947716t5712678wdq9",
	}
)

func TestSaveFortuneMessage(t *testing.T) {

	fortuneRepo := mocks.NewFortuneRepositoryMock()
	fortuneService := fortune.NewFortuneService(fortuneRepo)

	//HappyPath
	fortuneRepo.SetSaveFuncReturn(nil)
	fortuneRepo.SetFindOneReturn(models.FortuneMessage{}, nil)
	err := fortuneService.Save(validFortuneMessageRequestBody)
	assert.NoError(t, err)

	//Existent message in the database
	fortuneRepo.SetFindOneReturn(validFortuneMessage, nil)
	err = fortuneService.Save(validFortuneMessageRequestBody)
	assert.Error(t, err)
	assert.IsType(t, err, &models.ErrMessageAlreadyExists{})

	//Generic Mongodb Error
	fortuneRepo.SetFindOneReturn(validFortuneMessage, fmt.Errorf("Mongodb Error"))
	err = fortuneService.Save(validFortuneMessageRequestBody)
	assert.Error(t, err)

	//Invalid ID
	err = fortuneService.Save(invalidIDFortuneMessage)
	assert.Error(t, err)
	assert.Errorf(t, err, "Sorry "+invalidIDFortuneMessage.ID+" is not valid as Message id field")

	//Empty ID
	invalidIDFortuneMessage.ID = ""
	err = fortuneService.Save(invalidIDFortuneMessage)
	assert.Error(t, err)
	assert.Errorf(t, err, "Sorry "+invalidIDFortuneMessage.ID+" is not valid as Message id field")

	//Invalid Message
	err = fortuneService.Save(invalidMessageFortuneMessage)
	assert.Error(t, err)
	assert.Errorf(t, err, "Sorry "+invalidMessageFortuneMessage.Message+" is not valid as Message id field")

	//Invalid CheckSUm
	err = fortuneService.Save(invalidCheckSumFortuneMessage)
	assert.Error(t, err)
	assert.Errorf(t, err, "Your message already contains a checksum")
}

func TestFindRandomMessage(t *testing.T) {
	fortuneRepo := mocks.NewFortuneRepositoryMock()
	fortuneService := fortune.NewFortuneService(fortuneRepo)

	//HappyPath
	fortuneRepo.SetFindRandomReturn(validFortuneMessage, nil)
	fortuneMessage, err := fortuneService.FindRandom()
	assert.NoError(t, err)
	assert.Equal(t, fortuneMessage, validFortuneMessage)

	//No content found
	fortuneRepo.SetFindRandomReturn(models.FortuneMessage{}, &models.ErrNotFound{Message: "No messages found"})
	fortuneMessage, err = fortuneService.FindRandom()
	assert.IsType(t, err, &models.ErrNotFound{})
	assert.Equal(t, fortuneMessage, models.FortuneMessage{})

	//Generic Error
	fortuneRepo.SetFindRandomReturn(models.FortuneMessage{}, fmt.Errorf("This could be any error like json binding, pipeline failing, connection failed, etc."))
	fortuneMessage, err = fortuneService.FindRandom()
	assert.Error(t, err)
	assert.Equal(t, fortuneMessage, models.FortuneMessage{})
}
