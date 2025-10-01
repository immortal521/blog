package util

import (
	"io"
	"net/http"
	"strings"
)

func FetchCloudflareIPs() ([]string, error) {
	ipv4URL := "https://www.cloudflare.com/ips-v4"
	ipv6URL := "https://www.cloudflare.com/ips-v6"

	var ips []string

	for _, url := range []string{ipv4URL, ipv6URL} {
		resp, err := http.Get(url)
		if err != nil {
			return []string{}, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return []string{}, err
		}

		lines := string(body)
		for line := range strings.SplitSeq(lines, "\n") {
			if line != "" {
				ips = append(ips, strings.TrimSpace(line))
			}
		}
	}
	return ips, nil
}
