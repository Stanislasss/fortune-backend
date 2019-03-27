package models

const NotFound = "Not Found"
const InvalidMessageIDErrMessage = "Sorry %s is not valid as Message id field"

type ErrInvalidMessageID struct {
	Message string
}

func (e *ErrInvalidMessageID) Error() string {
	return e.Message
}

const InvalidMessageErrMessage = "Sorry %s is not valid as Message field"

type ErrInvalidMessage struct {
	Message string
}

func (e *ErrInvalidMessage) Error() string {
	return e.Message
}

type ErrNotFound struct {
	Message string
}

func (e *ErrNotFound) Error() string {
	return e.Message
}

type ErrMessageAlreadyExists struct {
	Message string
}

func (e *ErrMessageAlreadyExists) Error() string {
	return e.Message
}
