package cli

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/razielsd/antibruteforce/app/config"
)

type Cli struct {
	client *clientAPI
}

func NewCli(cfg config.AppConfig) *Cli {
	return &Cli{
		client: newClientAPI(cfg.Addr),
	}
}

func (c *Cli) ShowWhitelist() (string, error) {
	l, err := c.client.GetWhitelist()
	if err != nil {
		return "", fmt.Errorf("%w: unable get blacklist", err)
	}
	return c.drawList("Whitelist", l), nil
}

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

func (c *Cli) AppendBlacklist(clientIP string) (string, error) {
	err := c.client.appendBlacklist(clientIP)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func (c *Cli) AppendWhitelist(clientIP string) (string, error) {
	err := c.client.appendWhitelist(clientIP)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func (c *Cli) RemoveBlacklist(clientIP string) (string, error) {
	err := c.client.removeBlacklist(clientIP)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func (c *Cli) RemoveWhitelist(clientIP string) (string, error) {
	err := c.client.removeWhitelist(clientIP)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func (c *Cli) DropBucketByLogin(key string) (string, error) {
	err := c.client.dropBucketByLogin(key)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func (c *Cli) DropBucketByPwd(key string) (string, error) {
	err := c.client.dropBucketByPasswd(key)
	if err != nil {
		return "", err
	}
	return "OK", nil
}

func (c *Cli) DropBucketByIP(key string) (string, error) {
	err := c.client.dropBucketByIP(key)
	if err != nil {
		return "", err
	}
	return "OK", nil
}
