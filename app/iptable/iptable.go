package iptable

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

var (
	ErrInvalidIpv4Address = errors.New("invalid ip address")
	ErrInvalidIpv4Mask    = errors.New("invalid ip mask")
)

type IPTable struct {
	maskList []*net.IPNet
	ipList   []net.IP
}

func NewIPTable() *IPTable {
	return &IPTable{}
}

func (a *IPTable) Contains(clientIP string) (bool, error) {
	ip := net.ParseIP(clientIP)
	if ip.To4() == nil {
		return false, fmt.Errorf("%w:%s is not valid ipv4 address", ErrInvalidIpv4Address, clientIP)
	}
	for _, ipdb := range a.ipList {
		if ipdb.Equal(ip) {
			return true, nil
		}
	}
	for _, mask := range a.maskList {
		if mask.Contains(ip) {
			return true, nil
		}
	}
	return false, nil
}

func (a *IPTable) Add(ipOrMask string) error {
	if strings.Contains(ipOrMask, "/") {
		return a.addMask(ipOrMask)
	}
	return a.addIP(ipOrMask)
}

func (a *IPTable) addMask(netmask string) error {
	_, mask, err := net.ParseCIDR(netmask)
	if err != nil {
		return fmt.Errorf("%w: is not valid ipv4 mask %s,  %s", ErrInvalidIpv4Mask, netmask, err)
	}
	a.maskList = append(a.maskList, mask)
	return nil
}

func (a *IPTable) addIP(clientIP string) error {
	ip := net.ParseIP(clientIP)
	if ip.To4() == nil {
		return fmt.Errorf("%w:%s is not valid ipv4 address", ErrInvalidIpv4Address, clientIP)
	}
	a.ipList = append(a.ipList, ip)

	return nil
}
