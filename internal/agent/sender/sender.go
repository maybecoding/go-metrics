package sender

import (
	"context"
	"github.com/maybecoding/go-metrics.git/internal/agent/config"
	"net"
)

type Sender struct {
	ctx context.Context
	cfg config.Sender
	ip  net.IP
}

func New(ctx context.Context, cfg config.Sender) *Sender {
	s := &Sender{
		ctx: ctx,
		cfg: cfg,
	}
	s.identifyIP()
	return s
}
