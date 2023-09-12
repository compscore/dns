package main

import (
	"context"
	"fmt"
	"net"
	"strings"
)

func Run(ctx context.Context, target string, command string, expectedOutput string, username string, password string) (bool, string) {
	// Set up the custom resolver with the provided DNS server
	deadline, ok := ctx.Deadline()
	if !ok {
		return false, "failed to get deadline from context"
	}

	if !strings.Contains(target, ":") {
		target = target + ":53"
	}

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Deadline: deadline,
			}
			return d.DialContext(ctx, "udp", target)
		},
	}

	// Resolve based on the record type
	var addresses []string

	commandSplit := strings.Split(command, " ")
	if len(commandSplit) != 2 {
		return false, fmt.Sprintf("invalid command; format should be \"[record type] [domain]\"; provided: \"%s\"", command)
	}
	command = commandSplit[0]
	domain := commandSplit[1]

	switch command {
	case "A":
		ips, err := r.LookupIP(ctx, "ip4", domain)
		if err != nil {
			return false, err.Error()
		}
		for _, ip := range ips {
			addresses = append(addresses, ip.String())
		}
	case "AAAA":
		ips, err := r.LookupIP(ctx, "ip6", domain)
		if err != nil {
			return false, err.Error()
		}
		for _, ip := range ips {
			addresses = append(addresses, ip.String())
		}
	case "MX":
		mxs, err := r.LookupMX(ctx, domain)
		if err != nil {
			return false, err.Error()
		}
		for _, mx := range mxs {
			addresses = append(addresses, mx.Host)
		}
	case "TXT":
		txts, err := r.LookupTXT(ctx, domain)
		if err != nil {
			return false, err.Error()
		}
		addresses = append(addresses, txts...)
	case "CNAME":
		cname, err := r.LookupCNAME(ctx, domain)
		if err != nil {
			return false, err.Error()
		}
		addresses = append(addresses, cname)
	case "NS":
		nss, err := r.LookupNS(ctx, domain)
		if err != nil {
			return false, err.Error()
		}
		for _, ns := range nss {
			addresses = append(addresses, ns.Host)
		}
	default:
		return false, "unsupported record type: " + command
	}

	// Check if the expected result is in the list of resolved addresses
	for _, addr := range addresses {
		if addr == expectedOutput {
			return true, ""
		}
	}

	return false, fmt.Sprintf(
		"expected output \"%s\" not found in resolved addresses: [%s]",
		expectedOutput,
		func() string {
			var s []string

			for _, addr := range addresses {
				s = append(s, "\""+addr+"\"")
			}

			return strings.Join(s, ", ")
		}(),
	)
}
