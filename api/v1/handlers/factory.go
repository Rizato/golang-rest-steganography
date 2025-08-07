package handlers

import "steg/api/v1/models"

type EmbedJobFactory struct {
}

func (h *EmbedJobFactory) Create() (*models.EmbedJob, error) {
	return models.NewEmbedJob(), nil
}

func NewEmbedJobFactory() *EmbedJobFactory {
	return &EmbedJobFactory{}
}

type ExtractJobFactory struct {
}

func (h *ExtractJobFactory) Create() (*models.ExtractJob, error) {
	return models.NewExtractJob(), nil
}

func NewExtractJobFactory() *ExtractJobFactory {
	return &ExtractJobFactory{}
}
