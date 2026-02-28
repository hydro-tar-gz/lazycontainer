package lxd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"lazycontainer/internal/core"
)

type listRecord struct {
	Name   string            `json:"name"`
	Status string            `json:"status"`
	Config map[string]string `json:"config"`
	State  struct {
		Network map[string]struct {
			Addresses []struct {
				Family  string `json:"family"`
				Address string `json:"address"`
			} `json:"addresses"`
		} `json:"network"`
	} `json:"state"`
}

func parseListJSON(in []byte) ([]core.Instance, error) {
	var records []listRecord
	if err := json.Unmarshal(in, &records); err != nil {
		return nil, err
	}

	instances := make([]core.Instance, 0, len(records))
	for _, r := range records {
		ip := "-"
		for ifName, n := range r.State.Network {
			if ifName == "lo" {
				continue
			}
			for _, addr := range n.Addresses {
				if strings.EqualFold(addr.Family, "inet") && addr.Address != "" {
					ip = addr.Address
					break
				}
			}
			if ip != "-" {
				break
			}
		}

		image := "-"
		if r.Config != nil {
			for _, key := range []string{"image.description", "image.os", "volatile.base_image"} {
				if v := strings.TrimSpace(r.Config[key]); v != "" {
					image = v
					break
				}
			}
		}

		instances = append(instances, core.Instance{
			Name:  r.Name,
			State: r.Status,
			IP:    ip,
			Image: image,
		})
	}

	return instances, nil
}

type imageRecord struct {
	Fingerprint string `json:"fingerprint"`
	Aliases     []struct {
		Name string `json:"name"`
	} `json:"aliases"`
	Properties map[string]string `json:"properties"`
}

func parseImageAliasesJSON(in []byte) ([]string, error) {
	var records []imageRecord
	if err := json.Unmarshal(in, &records); err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	aliases := make([]string, 0, len(records))
	for _, r := range records {
		for _, a := range r.Aliases {
			name := strings.TrimSpace(a.Name)
			if name == "" {
				continue
			}
			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
			aliases = append(aliases, name)
		}
	}
	sort.Strings(aliases)
	return aliases, nil
}

func parseImageListJSON(in []byte) (string, error) {
	aliases, err := parseImageAliasesJSON(in)
	if err != nil {
		return "", err
	}
	if len(aliases) == 0 {
		return "No images found", nil
	}

	var records []imageRecord
	if err := json.Unmarshal(in, &records); err != nil {
		return "", err
	}
	lines := make([]string, 0, len(records)+1)
	lines = append(lines, "Available images (remote: images:)")
	for _, r := range records {
		alias := "-"
		for _, a := range r.Aliases {
			if strings.TrimSpace(a.Name) != "" {
				alias = a.Name
				break
			}
		}
		desc := "-"
		if r.Properties != nil {
			if v := strings.TrimSpace(r.Properties["description"]); v != "" {
				desc = v
			} else if v := strings.TrimSpace(r.Properties["os"]); v != "" {
				desc = v
			}
		}
		fp := r.Fingerprint
		if len(fp) > 12 {
			fp = fp[:12]
		}
		lines = append(lines, fmt.Sprintf("- %-24s %-12s %s", alias, fp, desc))
	}
	return strings.Join(lines, "\n"), nil
}
