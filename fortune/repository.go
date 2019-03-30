package fortune

import (
	"encoding/json"

	"github.com/thiagotrennepohl/fortune-backend/models"
	"github.com/globalsign/mgo"
)

type fortuneRepository struct {
	dbSession *mgo.Session
}

const fortuneMessagesCollection = "fortune-messages"

func NewFortuneRepository(session *mgo.Session) FortuneRepository {
	return &fortuneRepository{
		dbSession: session,
	}
}

func (db *fortuneRepository) FindOne(query models.FortuneQuery) (models.FortuneMessage, error) {
	var message models.FortuneMessage
	session := db.dbSession.Copy()
	defer session.Close()

	connection := session.DB("").C(fortuneMessagesCollection)

	err := connection.Find(query).One(&message)
	if err != nil {
		if err == mgo.ErrNotFound {
			return message, &models.ErrNotFound{}
		}
		return message, err
	}

	return message, err
}

func (db *fortuneRepository) FindRandom(aggregationQUery []models.FortuneQuery) (models.FortuneMessage, error) {
	var message models.FortuneMessage
	session := db.dbSession.Copy()
	defer session.Close()
	result := make([]models.Json, 1)
	connection := session.DB("").C(fortuneMessagesCollection)

	err := connection.Pipe(aggregationQUery).All(&result)
	if err != nil {
		return message, err
	}

	if len(result) < 1 {
		return message, &models.ErrNotFound{Message: "No messages found"}
	}

	byteArray, err := json.Marshal(result[0])
	if err != nil {
		return message, err
	}
	err = json.Unmarshal(byteArray, &message)
	return message, err
}

func (db *fortuneRepository) Save(message models.FortuneMessage) error {
	session := db.dbSession.Copy()
	defer session.Close()
	connection := session.DB("").C(fortuneMessagesCollection)
	return connection.Insert(message)
}
