package iptable

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
)

var (
	ErrInvalidIpv4Address = errors.New("invalid ip address")
	ErrInvalidIpv4Mask    = errors.New("invalid ip mask")
)

type IPTable struct {
	maskList map[string]*net.IPNet
	ipList   map[string]struct{}
	mu       sync.Mutex
}

func NewIPTable() *IPTable {
	return &IPTable{
		maskList: make(map[string]*net.IPNet),
		ipList:   make(map[string]struct{}),
		mu:       sync.Mutex{},
	}
}

func (a *IPTable) Contains(clientIP string) (bool, error) {
	ip := net.ParseIP(clientIP)
	if ip.To4() == nil {
		return false, fmt.Errorf("%w:%s is not valid ipv4 address", ErrInvalidIpv4Address, clientIP)
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.ipList[clientIP]; ok {
		return true, nil
	}
	for _, mask := range a.maskList {
		if mask.Contains(ip) {
			return true, nil
		}
	}
	return false, nil
}

func (a *IPTable) Add(ipOrMask string) error {
	if a.isMask(ipOrMask) {
		return a.addMask(ipOrMask)
	}
	return a.addIP(ipOrMask)
}

func (a *IPTable) addMask(netmask string) error {
	_, mask, err := net.ParseCIDR(netmask)
	if err != nil {
		return fmt.Errorf("%w: is not valid ipv4 mask %s,  %s", ErrInvalidIpv4Mask, netmask, err)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.maskList[netmask] = mask
	return nil
}

func (a *IPTable) addIP(clientIP string) error {
	ip := net.ParseIP(clientIP)
	if ip.To4() == nil {
		return fmt.Errorf("%w:%s is not valid ipv4 address", ErrInvalidIpv4Address, clientIP)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.ipList[clientIP] = struct{}{}
	return nil
}

func (a *IPTable) GetAll() []string {
	ips := []string{}
	for ip := range a.ipList {
		ips = append(ips, ip)
	}
	for mask := range a.maskList {
		ips = append(ips, mask)
	}
	sort.Strings(ips)
	return ips
}

func (a *IPTable) Remove(ipOrMask string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.isMask(ipOrMask) {
		delete(a.maskList, ipOrMask)
	} else {
		delete(a.ipList, ipOrMask)
	}
	return nil
}

func (a *IPTable) isMask(ipOrMask string) bool {
	return strings.Contains(ipOrMask, "/")
}
