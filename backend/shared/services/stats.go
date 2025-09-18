package services

import (
	"context"
	"steg/server/middleware"
	"steg/shared/repositories"
)

type StatsService struct {
	Image   *repositories.ImageRepository
	Embed   *repositories.EmbedJobRepository
	Extract *repositories.ExtractJobRepository
	Stats   *middleware.Stats
}

func NewStatsService(imageRepository *repositories.ImageRepository, embedJobRepo *repositories.EmbedJobRepository, extractJobRepo *repositories.ExtractJobRepository, stats *middleware.Stats) *StatsService {
	return &StatsService{
		imageRepository,
		embedJobRepo,
		extractJobRepo,
		stats,
	}
}

func (s *StatsService) GetStats(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"database": s.GetRepositoryStats(ctx),
		"requests": s.GetMiddlewareStats(),
	}
}

func (service *StatsService) GetRepositoryStats(ctx context.Context) map[string]interface{} {
	images, err := service.Image.List(ctx)
	if err != nil {
		return map[string]interface{}{}
	}
	embeds, err := service.Embed.List(ctx)
	if err != nil {
		return map[string]interface{}{}
	}
	extracts, err := service.Extract.List(ctx)
	if err != nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"images":   len(images),
		"embeds":   len(embeds),
		"extracts": len(extracts),
	}
}

func (service *StatsService) GetMiddlewareStats() map[string]interface{} {
	return service.Stats.GetStats()
}
