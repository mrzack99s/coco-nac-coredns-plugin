package coconac

import (
	"errors"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

// init registers this plugin.
func init() { plugin.Register("coconac", setup) }

func setup(c *caddy.Controller) error {
	// c.Next() // Ignore "coconac" and give us the next token.
	if !c.NextArg() {
		return plugin.Error("coconac", c.ArgErr())
	}

	for c.Next() {
		// var t string
		switch c.Val() {
		case "host":
			if !c.NextArg() {
				return c.ArgErr()
			}
			CocoNAC_Host = c.Val()
		case "api_key":
			if !c.NextArg() {
				return c.ArgErr()
			}
			CocoNAC_APIKey = c.Val()
		}
	}

	if CocoNAC_Host == "" {
		return plugin.Error("coconac", errors.New(`need "host" parameter`))
	}

	if CocoNAC_APIKey == "" {
		return plugin.Error("coconac", errors.New(`need "api_key" parameter`))
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return CocoNAC{Next: next}
	})

	return nil
}
