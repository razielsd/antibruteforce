package iptable

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAllowFilter(t *testing.T) {
	filter := NewIpTable()
	require.NotNil(t, filter)
}

func TestIpTable_Contains_Found(t *testing.T) {
	tests := []struct {
		name   string
		ips    []string
		search string
	}{
		{
			name:   "single ip",
			ips:    []string{"192.168.1.1"},
			search: "192.168.1.1",
		},
		{
			name:   "single mask",
			ips:    []string{"192.168.1.0/24"},
			search: "192.168.1.1",
		},
		{
			name:   "mask and ip(2)",
			ips:    []string{"192.168.1.0/24", "172.168.10.10", "172.168.10.12"},
			search: "172.168.10.12",
		},
		{
			name:   "mask(2) and ip",
			ips:    []string{"192.168.1.0/24", "172.168.10.0/16", "10.10.10.12"},
			search: "172.168.10.12",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			iptable := NewIpTable()
			for _, ip := range test.ips {
				err := iptable.Add(ip)
				require.NoError(t, err, "Unable add ip/mask for test data")
			}
			found, err := iptable.Contains(test.search)
			require.NoError(t, err, "Error on check ip/mask")
			require.True(t, found)
		})
	}
}

func TestIpTable_Contains_NotFound(t *testing.T) {
	tests := []struct {
		name   string
		ips    []string
		search string
	}{
		{
			name:   "single ip",
			ips:    []string{"192.168.1.2"},
			search: "192.168.1.1",
		},
		{
			name:   "single mask",
			ips:    []string{"192.168.1.0/24"},
			search: "192.168.2.1",
		},
		{
			name:   "mask and ip(2)",
			ips:    []string{"192.168.1.0/24", "172.168.10.10", "172.168.10.12"},
			search: "10.10.10.12",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			iptable := NewIpTable()
			for _, ip := range test.ips {
				err := iptable.Add(ip)
				require.NoError(t, err, "Unable add ip/mask for test data")
			}
			found, err := iptable.Contains(test.search)
			require.NoError(t, err, "Error on check ip/mask")
			require.False(t, found)
		})
	}
}

func TestIpTable_Contains_CheckInvalidIp(t *testing.T) {
	iptable := NewIpTable()
	_, err := iptable.Contains("192.168.1")
	require.ErrorIs(t, err, ErrInvalidIpv4Address)
}

func TestIpTable_Add_InvalidIp(t *testing.T) {
	iptable := NewIpTable()
	err := iptable.Add("192.168.1.")
	require.ErrorIs(t, err, ErrInvalidIpv4Address)
}

func TestIpTable_Add_InvalidMask(t *testing.T) {
	iptable := NewIpTable()
	err := iptable.Add("192.168.1./50")
	require.ErrorIs(t, err, ErrInvalidIpv4Mask)
}
