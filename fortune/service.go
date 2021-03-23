package fortune

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/thiagotrennepohl/fortune-backend/models"
)

const invalidMessageIdPattern = `[^A-Za-z0-9\-]`
const ErrMessageAlreadyContainsCheckSum = "Your message already contains a checksum"

type fortuneService struct {
	repository FortuneRepository
}

//NewFortuneService creates a new service responsible for Fortune business rules
func NewFortuneService(forturneRepository FortuneRepository) FortuneService {
	return &fortuneService{
		repository: forturneRepository,
	}
}

func (svc *fortuneService) Save(message models.FortuneMessage) error {
	err := svc.runMessageValidations(message)
	if err != nil {
		return err
	}

	existentMessage, err := svc.repository.FindOne(models.FortuneQuery{"id": message.ID})
	if err != nil {
		if _, ok := err.(*models.ErrNotFound); !ok {
			return err
		}
	}

	message.CheckSum = svc.createFortuneMessageChecksum(message)

	if message == existentMessage {
		return &models.ErrMessageAlreadyExists{Message: "Message already exists"}
	}

	return svc.repository.Save(message)
}

func (svc *fortuneService) GetAll() ([]models.FortuneMessage, error){
	return svc.repository.FindAll()
}

func (svc *fortuneService) FindRandom() (models.FortuneMessage, error) {
	aggregationQuery := []models.FortuneQuery{
		models.FortuneQuery{
			"$sample": models.FortuneQuery{
				"size": 1,
			},
		},
	}
	return svc.repository.FindRandom(aggregationQuery)
}

func (svc *fortuneService) runMessageValidations(message models.FortuneMessage) error {
	if valid := svc.isMessageIDValid(message.ID); !valid {
		return &models.ErrInvalidMessageID{Message: fmt.Sprintf(models.InvalidMessageIDErrMessage, message.ID)}
	}

	if valid := svc.isMessageValid(message.Message); !valid {
		return &models.ErrInvalidMessage{Message: fmt.Sprintf(models.InvalidMessageErrMessage, message.Message)}
	}

	if message.CheckSum != "" {
		return &models.ErrInvalidMessage{Message: ErrMessageAlreadyContainsCheckSum}
	}

	return nil
}

func (svc *fortuneService) isMessageIDValid(messageID string) bool {
	if messageID == "" {
		return false
	}
	if ok, err := regexp.Match(invalidMessageIdPattern, []byte(messageID)); ok || err != nil {
		return false
	}
	return true
}

func (svc *fortuneService) isMessageValid(message string) bool {
	return message != ""
}

func (svc *fortuneService) createFortuneMessageChecksum(message models.FortuneMessage) string {
	md5Hasher := md5.New()
	md5Hasher.Write([]byte(message.ID + message.Message))
	hash := md5Hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
