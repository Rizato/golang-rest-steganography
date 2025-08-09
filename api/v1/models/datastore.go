package models

import (
	"github.com/google/uuid"
	"sync"
)

type Datastore struct {
	mutex       sync.RWMutex
	embedJobs   map[uuid.UUID]EmbedJob   // Really need row level locks, but that isn't happening
	extractJobs map[uuid.UUID]ExtractJob // Really need row level locks, but that isn't happening
	images      map[uuid.UUID]ServerFile // Really need row level locks, but that isn't happening
}

func NewDatastore() *Datastore {
	return &Datastore{
		embedJobs:   make(map[uuid.UUID]EmbedJob),
		extractJobs: make(map[uuid.UUID]ExtractJob),
		images:      make(map[uuid.UUID]ServerFile),
	}
}

func (ds *Datastore) SaveEmbedJob(job EmbedJob) {
	// TODO No row locks means that this can overwrite someones changes, even though the db level lock prevents corruption
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.embedJobs[job.Uuid] = job
}

func (ds *Datastore) GetEmbedJob(id uuid.UUID) (EmbedJob, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	job, found := ds.embedJobs[id]
	return job, found
}

func (ds *Datastore) GetEmbedJobs() []EmbedJob {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	jobs := make([]EmbedJob, 0, len(ds.embedJobs))
	for _, job := range ds.embedJobs {
		jobs = append(jobs, job)
	}
	return jobs
}

func (ds *Datastore) DeleteEmbedJob(id uuid.UUID) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	delete(ds.embedJobs, id)
}

func (ds *Datastore) SaveExtractJob(job ExtractJob) {
	// TODO No row locks means that this can overwrite someones changes, even though the db level lock prevents corruption
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.extractJobs[job.Uuid] = job
}

func (ds *Datastore) GetExtractJob(id uuid.UUID) (ExtractJob, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	job, found := ds.extractJobs[id]
	return job, found
}

func (ds *Datastore) GetExtractJobs() []ExtractJob {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	jobs := make([]ExtractJob, 0, len(ds.extractJobs))
	for _, job := range ds.extractJobs {
		jobs = append(jobs, job)
	}
	return jobs
}

func (ds *Datastore) DeleteExtractJob(id uuid.UUID) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	delete(ds.extractJobs, id)
}

func (ds *Datastore) SaveImage(image ServerFile) {
	// TODO No row locks means that this can overwrite someones changes, even though the db level lock prevents corruption
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.images[image.Uuid] = image
}

func (ds *Datastore) GetImage(id uuid.UUID) (ServerFile, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	file, found := ds.images[id]
	return file, found
}

func (ds *Datastore) GetImages() []ServerFile {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	images := make([]ServerFile, 0, len(ds.images))
	for _, image := range ds.images {
		images = append(images, image)
	}
	return images
}

func (ds *Datastore) DeleteImage(id uuid.UUID) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	// TODO CASCADE DELETIONS (Do I care? DB will do that when I get there)
	delete(ds.images, id)
}

func (ds *Datastore) GetStats() interface{} {
	// TODO Break down jobs by status
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	return map[string]interface{}{
		"embed-jobs":   len(ds.embedJobs),
		"extract-jobs": len(ds.extractJobs),
		"images":       len(ds.images),
	}
}
