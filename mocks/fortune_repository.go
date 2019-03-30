package mocks

import (
	"github.com/thiagotrennepohl/fortune-backend/models"
)

type fortuneRepositoryMock struct {
	SaveFuncErrBehaviour         error
	FindRandomFuncErrBehaviour   error
	FindRandomFuncValueBehaviour models.FortuneMessage
	FindOneFuncValueBehaviour    models.FortuneMessage
	FindOneFuncErrBehaviour      error
}

func NewFortuneRepositoryMock() *fortuneRepositoryMock {
	return &fortuneRepositoryMock{}
}

func (m *fortuneRepositoryMock) SetSaveFuncReturn(err error) {
	m.SaveFuncErrBehaviour = err
}

func (m *fortuneRepositoryMock) SetFindRandomReturn(fortuneMessage models.FortuneMessage, err error) {
	m.FindRandomFuncErrBehaviour = err
	m.FindRandomFuncValueBehaviour = fortuneMessage
}

func (m *fortuneRepositoryMock) SetFindOneReturn(fortuneMessage models.FortuneMessage, err error) {
	m.FindOneFuncErrBehaviour = err
	m.FindOneFuncValueBehaviour = fortuneMessage
}

func (m *fortuneRepositoryMock) Save(message models.FortuneMessage) error {
	return m.SaveFuncErrBehaviour
}

func (m *fortuneRepositoryMock) FindRandom(pipelineQuery []models.FortuneQuery) (models.FortuneMessage, error) {
	return m.FindRandomFuncValueBehaviour, m.FindRandomFuncErrBehaviour
}

func (m *fortuneRepositoryMock) FindOne(query models.FortuneQuery) (models.FortuneMessage, error) {
	return m.FindOneFuncValueBehaviour, m.FindOneFuncErrBehaviour
}
