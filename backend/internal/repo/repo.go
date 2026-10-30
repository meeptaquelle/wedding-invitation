package repo

import (
	"context"

	"backend/internal/config"
	"backend/internal/sheets"
)

type Repo struct {
	Client *sheets.Client
	Cfg    *config.Config
}

func New(ctx context.Context) (*Repo, error) {
	cfg := config.Load()
	client, err := sheets.NewClient(ctx, cfg.ServiceAccountEmail, cfg.PrivateKey, cfg.SheetID)
	if err != nil {
		return nil, err
	}
	return &Repo{Client: client, Cfg: cfg}, nil
}
