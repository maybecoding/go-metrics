package sender

import (
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"net"
)

func (j *Sender) identifyIP() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Error().Err(err).Msg("failed to identify ip address of host InterfaceAddrs error")
		return
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			j.ip = ipNet.IP
			return
		}
	}
	logger.Error().Msg("ip address is not found")
}
