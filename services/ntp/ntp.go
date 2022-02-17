package ntp

import (
	"context"
	"github.com/aibotsoft/ntp-service/pkg/config"
	"go.uber.org/zap"
)

type Ntp struct {
	cfg *config.Config
	log *zap.Logger
	ctx context.Context
}

func NewNtp(cfg *config.Config, log *zap.Logger, ctx context.Context) *Ntp {
	return &Ntp{cfg: cfg, log: log, ctx: ctx}

}

func (n *Ntp) Run() error {
	for {
		select {
		case <-n.ctx.Done():
			return n.ctx.Err()
		}
	}
}
func (n *Ntp) Close() error {
	return nil
}
