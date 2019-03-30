package fortune

import "github.com/thiagotrennepohl/fortune-backend/models"

//MessageRepository stands for "database" operations
type FortuneRepository interface {
	Save(message models.FortuneMessage) error
	FindOne(query models.FortuneQuery) (models.FortuneMessage, error)
	FindRandom(aggregationQuery []models.FortuneQuery) (models.FortuneMessage, error)
}

type FortuneService interface {
	Save(message models.FortuneMessage) error
	FindRandom() (models.FortuneMessage, error)
}
