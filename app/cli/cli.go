package cli

import (
	"bytes"
	"fmt"
	"text/tabwriter"

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
	l, err := c.client.GetWhitelist()
	if err != nil {
		return "", fmt.Errorf("%w: unable get blacklist", err)
	}
	return c.drawList("Whitelist", l), nil
}

// ShowBlacklist show blacklist.
func (c *Cli) ShowBlacklist() (string, error) {
	l, err := c.client.GetBlacklist()
	if err != nil {
		return "", fmt.Errorf("%w: unable get blacklist", err)
	}
	return c.drawList("Blacklist", l), nil
}

func (c *Cli) drawList(title string, l map[string]bwItem) string {
	buf := bytes.NewBuffer([]byte{})

	fmt.Fprintf(buf, "--=== %s ===--\n", title)
	w := tabwriter.NewWriter(buf, 1, 0, 2, ' ', tabwriter.Debug)
	_, _ = fmt.Fprintf(w, "%s\t %s\t %s\n", "IP/Mask", "Count", "LastAccess.")
	for ip, stat := range l {
		_, _ = fmt.Fprintf(w, "%s\t %d\t %s\n", ip, stat.Counter, "stat.LastAccess.")
	}
	if len(l) == 0 {
		_, _ = fmt.Fprintln(w, "Empty")
	}

	w.Flush()
	return buf.String()
}

// AppendBlacklist add ip/mask to blacklist.
func (c *Cli) AppendBlacklist(clientIP string) (string, error) {
	err := c.client.appendBlacklist(clientIP)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

// AppendWhitelist add ip/mask to whitelist.
func (c *Cli) AppendWhitelist(clientIP string) (string, error) {
	err := c.client.appendWhitelist(clientIP)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

// RemoveBlacklist remove ip/mask from blacklist.
func (c *Cli) RemoveBlacklist(clientIP string) (string, error) {
	err := c.client.removeBlacklist(clientIP)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

// RemoveWhitelist - remove ip/mask from whitelist.
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
