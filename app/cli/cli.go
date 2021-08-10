package cli

import (
	"bytes"
	"fmt"

	"github.com/razielsd/antibruteforce/app/config"
)

// Cli - cli api.
type Cli struct {
	client *clientAPI
}

// NewCli - create new instance of Cli.
func NewCli(cfg config.AppConfig) *Cli {
	return &Cli{
		client: newClientAPI(cfg.Addr),
	}
}

// ShowWhitelist show whitelist.
func (c *Cli) ShowWhitelist() (string, error) {
	l, err := c.client.getWhitelist()
	if err != nil {
		return "", fmt.Errorf("%w: unable get whitelist", err)
	}
	return c.drawList("Whitelist", l), nil
}

// ShowBlacklist show blacklist.
func (c *Cli) ShowBlacklist() (string, error) {
	l, err := c.client.getBlacklist()
	if err != nil {
		return "", fmt.Errorf("%w: unable get blacklist", err)
	}
	return c.drawList("Blacklist", l), nil
}

func (c *Cli) drawList(title string, l []string) string {
	buf := bytes.NewBuffer([]byte{})

	fmt.Fprintf(buf, "--=== %s ===--\n", title)
	for _, ip := range l {
		_, _ = fmt.Fprintf(buf, "%s\n", ip)
	}
	if len(l) == 0 {
		_, _ = fmt.Fprintln(buf, "Empty")
	}

	return buf.String()
}

// AppendBlacklist add ip/subnet to blacklist.
func (c *Cli) AppendBlacklist(clientIP string) (string, error) {
	err := c.client.appendBlacklist(clientIP)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

// AppendWhitelist add ip/subnet to whitelist.
func (c *Cli) AppendWhitelist(clientIP string) (string, error) {
	err := c.client.appendWhitelist(clientIP)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

// RemoveBlacklist remove ip/subnet from blacklist.
func (c *Cli) RemoveBlacklist(clientIP string) (string, error) {
	err := c.client.removeBlacklist(clientIP)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

// RemoveWhitelist - remove ip/subnet from whitelist.
func (c *Cli) RemoveWhitelist(clientIP string) (string, error) {
	err := c.client.removeWhitelist(clientIP)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

// DropBucketByLogin remove bucket by login.
func (c *Cli) DropBucketByLogin(key string) (string, error) {
	err := c.client.dropBucketByLogin(key)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

// DropBucketByPwd remove bucket by password.
func (c *Cli) DropBucketByPwd(key string) (string, error) {
	err := c.client.dropBucketByPasswd(key)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

// DropBucketByIP remove bucket by IP.
func (c *Cli) DropBucketByIP(key string) (string, error) {
	err := c.client.dropBucketByIP(key)
	if err != nil {
		return "", err
	}
	return "OK", nil
}
