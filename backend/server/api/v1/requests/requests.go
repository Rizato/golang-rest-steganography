package requests

import (
	"errors"

	"github.com/google/uuid"
)

// These are request definitions, which should mostly match the model, just without the user specified fields

var (
	MissingUUIDError     = errors.New("missing uuid")
	InvalidUUIDError     = errors.New("invalid uuid")
	MissingMessageError  = errors.New("missing message")
	MessageTooLargeError = errors.New("message too large")
)

func ValidateMessage(message string) error {
	if message == "" {
		return MissingMessageError
	}

	// todo make a little config for this
	if len(message) > 1000 {
		return MessageTooLargeError
	}
	return nil
}

func ValidateUUID(value string) (uuid.UUID, error) {
	var parsed uuid.UUID
	if value == "" {
		return parsed, MissingUUIDError
	}

	if len(value) != 36 {
		return parsed, InvalidUUIDError
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return parsed, InvalidUUIDError
	}
	return parsed, nil
}

type RequestValidator interface {
	Validate() error
}

// TODO Algorithm choices
type CreateEmbedRequest struct {
	ImageUUID     string `json:"image-uuid"`
	Message       string `json:"message"`
	ValidatedUUID uuid.UUID
}

func (request *CreateEmbedRequest) Validate() error {
	err := ValidateMessage(request.Message)
	if err != nil {
		return err
	}
	parsed, err := ValidateUUID(request.ImageUUID)
	if err != nil {
		return err
	}
	request.ValidatedUUID = parsed
	return nil
}

type CreateExtractRequest struct {
	ImageUUID     string `json:"image-uuid"`
	ValidatedUUID uuid.UUID
}

func (request *CreateExtractRequest) Validate() error {
	parsed, err := ValidateUUID(request.ImageUUID)
	if err != nil {
		return err
	}
	request.ValidatedUUID = parsed
	return nil
}
