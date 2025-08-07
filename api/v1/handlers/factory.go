package handlers

import "steg/api/v1/models"

type EncodeFactory struct {
}

func (h *EncodeFactory) Create() (*models.EncodeJob, error) {
	return models.NewEncodeJob(), nil
}

func NewEncodeFactory() *EncodeFactory {
	return &EncodeFactory{}
}

type DecodeFactory struct {
}

func (h *DecodeFactory) Create() (*models.DecodeJob, error) {
	return models.NewDecodeJob(), nil
}

func NewDecodeFactory() *DecodeFactory {
	return &DecodeFactory{}
}
