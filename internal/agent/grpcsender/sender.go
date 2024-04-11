package grpcsender

import (
	"context"
	"github.com/maybecoding/go-metrics.git/api/metric/v1/pb"
	"github.com/maybecoding/go-metrics.git/internal/agent/config"
	"google.golang.org/grpc"
	"net"
)

type Sender struct {
	cfg        config.Sender
	ctx        context.Context
	ip         net.IP
	clientConn *grpc.ClientConn
	clintGRPC  pb.MetricsV1Client
}

func New(ctx context.Context, cfg config.Sender) *Sender {
	s := &Sender{ctx: ctx, cfg: cfg}
	s.identifyIP()
	return s
}
