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

type IpTable struct {
	maskList []*net.IPNet
	ipList   []net.IP
}

func NewIpTable() *IpTable {
	return &IpTable{}
}

func (a *IpTable) Contains(clientIp string) (bool, error) {
	ip := net.ParseIP(clientIp)
	if ip.To4() == nil {
		return false, fmt.Errorf("%w:%s is not valid ipv4 address", ErrInvalidIpv4Address, clientIp)
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

func (a *IpTable) Add(ipOrMask string) error {
	if strings.Contains(ipOrMask, "/") {
		return a.addMask(ipOrMask)
	}
	return a.addIp(ipOrMask)
}

func (a *IpTable) addMask(netmask string) error {
	_, mask, err := net.ParseCIDR(netmask)
	if err != nil {
		return fmt.Errorf("%w: is not valid ipv4 mask %s,  %s", ErrInvalidIpv4Mask, netmask, err)
	}
	a.maskList = append(a.maskList, mask)
	return nil
}

func (a *IpTable) addIp(clientIp string) error {
	ip := net.ParseIP(clientIp)
	if ip.To4() == nil {
		return fmt.Errorf("%w:%s is not valid ipv4 address", ErrInvalidIpv4Address, clientIp)
	}
	a.ipList = append(a.ipList, ip)

	return nil
}
