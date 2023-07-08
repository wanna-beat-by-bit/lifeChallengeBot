package storage

import (
	"context"
	"time"
)

type Storage interface {
	Init(ctx context.Context) error
	CreateMission(ctx context.Context, mission *Mission) error
	ReadLatestMissions(ctx context.Context) ([]Mission, error)
	RemoveMission(ctx context.Context, id int) error
}

type Mission struct {
	Id       int
	Text     string
	Deadline time.Time
}
