package dns

import (
	"log"
	"net"
	"strconv"
	"time"
)

type Service struct {
	IP   string
	Host string
	Port uint16
}

type ChechService = func(string) (bool, error)

func TCPCheck(addr string) (bool, error) {
	sock, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		log.Printf("%v", err)
		return false, err
	}
	defer sock.Close()
	return true, nil
}

func ServiceDiscovery(service string, checkSrvFnc ChechService) *Service {
	_, addrs, _ := net.LookupSRV(service, "tcp", "fr123k.uk")

	for _, a := range addrs {
		addr := net.JoinHostPort(a.Target, strconv.Itoa(int(a.Port)))
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			log.Printf("ResolveTCPAddr(%s): %s", tcpAddr, err.Error())
			continue
		}
		valid, err := checkSrvFnc(addr)
		if !valid {
			continue
		}
		return &Service{
			IP:   tcpAddr.IP.String(),
			Host: a.Target,
			Port: a.Port,
		}
	}
	return nil
}
