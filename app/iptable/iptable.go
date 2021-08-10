package iptable

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
)

// Validation errors.
var (
	ErrInvalidIpv4Address = errors.New("invalid ip address")
	ErrInvalidIpv4Subnet  = errors.New("invalid ip subnet")
)

// IPTable Container for ip/subnet.
type IPTable struct {
	subnetList map[string]*net.IPNet
	ipList     map[string]struct{}
	mu         sync.Mutex
}

// NewIPTable create new instance of IPTable.
func NewIPTable() *IPTable {
	return &IPTable{
		subnetList: make(map[string]*net.IPNet),
		ipList:     make(map[string]struct{}),
		mu:         sync.Mutex{},
	}
}

// Contains check contain IP in table.
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
	for _, subnet := range a.subnetList {
		if subnet.Contains(ip) {
			return true, nil
		}
	}
	return false, nil
}

// Add add ip or subnet.
func (a *IPTable) Add(ipOrSubnet string) error {
	if a.isSubnet(ipOrSubnet) {
		return a.addSubnet(ipOrSubnet)
	}
	return a.addIP(ipOrSubnet)
}

func (a *IPTable) addSubnet(subnet string) error {
	_, cidr, err := net.ParseCIDR(subnet)
	if err != nil {
		return fmt.Errorf("%w: is not valid ipv4 subnet %s,  %s", ErrInvalidIpv4Subnet, subnet, err)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.subnetList[subnet] = cidr
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

// GetAll get all registered ip or subnet.
func (a *IPTable) GetAll() []string {
	ips := []string{}
	for ip := range a.ipList {
		ips = append(ips, ip)
	}
	for subnet := range a.subnetList {
		ips = append(ips, subnet)
	}
	sort.Strings(ips)
	return ips
}

// Remove remove ip or subnet from table.
func (a *IPTable) Remove(ipOrSubnet string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.isSubnet(ipOrSubnet) {
		delete(a.subnetList, ipOrSubnet)
	} else {
		delete(a.ipList, ipOrSubnet)
	}
}

func (a *IPTable) isSubnet(ipOrSubnet string) bool {
	return strings.Contains(ipOrSubnet, "/")
}
