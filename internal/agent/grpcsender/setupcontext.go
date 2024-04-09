package grpcsender

import (
	"google.golang.org/grpc/metadata"
	"strings"
)

func (s *Sender) setupContext() {
	mdMap := make(map[string]string)
	if s.ip != nil && s.cfg.IPAddrHeader != "" {
		mdMap[strings.ToLower(s.cfg.IPAddrHeader)] = s.ip.String()
	}
	if len(mdMap) > 0 {
		md := metadata.New(mdMap)
		s.ctx = metadata.NewOutgoingContext(s.ctx, md)
	}

}
