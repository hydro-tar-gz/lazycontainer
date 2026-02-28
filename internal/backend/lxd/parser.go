package lxd

import (
	"encoding/json"
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
