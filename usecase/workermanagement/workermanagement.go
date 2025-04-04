package workermanagement

import (
	"context"
	"time"
)

type IWorkerManagementUsecase interface {
	CreateWorker(ctx context.Context, version string) error
	ListWorker(ctx context.Context) ([]ListWorkerOutputItem, error)
	SetWorkerNum(ctx context.Context, num int32) error
	RestartWorker(ctx context.Context) error
}

type ListWorkerOutputItem struct {
	Name      string
	Version   string
	CreatedAt time.Time
	Status    string
}
