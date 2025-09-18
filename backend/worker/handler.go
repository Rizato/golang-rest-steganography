package worker

import (
	"context"
	"encoding/json"
	"log"
	"steg/shared/models"
	"steg/shared/rabbit"
	"steg/shared/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

type JobMessageHandler struct {
	embedService   *services.EmbedJobService
	extractService *services.ExtractJobService
}

func NewJobMessageHandler(embedService *services.EmbedJobService, extractService *services.ExtractJobService) *JobMessageHandler {
	return &JobMessageHandler{embedService, extractService}
}

func (h *JobMessageHandler) HandleMessage(ctx context.Context, d amqp.Delivery) error {
	var message models.JobMessage
	err := json.Unmarshal(d.Body, &message)
	if err != nil {
		return err
	}
	log.Printf("Received a message: %s", d.Body)
	return h.HandleJob(ctx, message)
}

func (h *JobMessageHandler) HandleJob(ctx context.Context, message models.JobMessage) error {
	switch message.Job {
	case models.Embed:
		return h.HandleEmbed(ctx, message)
	case models.Extract:
		return h.HandleExtract(ctx, message)
	default:
		return rabbit.InvalidMessage
	}
}

func (h *JobMessageHandler) HandleEmbed(ctx context.Context, message models.JobMessage) error {
	job, err := h.embedService.GetEmbedJob(ctx, message.Uuid)
	if err != nil {
		return err
	}
	return h.embedService.Embed(ctx, job)
}

func (h *JobMessageHandler) HandleExtract(ctx context.Context, message models.JobMessage) error {
	job, err := h.extractService.GetExtractJob(ctx, message.Uuid)
	if err != nil {
		return err
	}
	return h.extractService.ExtractToDb(ctx, job)
}
