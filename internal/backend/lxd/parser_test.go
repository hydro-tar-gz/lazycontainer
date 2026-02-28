package lxd

import "testing"

func TestParseListJSON(t *testing.T) {
	in := []byte(`[
		{
			"name": "web-1",
			"status": "Running",
			"config": {"image.description": "ubuntu/24.04"},
			"state": {
				"network": {
					"lo": {"addresses": [{"family": "inet", "address": "127.0.0.1"}]},
					"eth0": {"addresses": [{"family": "inet", "address": "10.10.10.22"}]}
				}
			}
		},
		{
			"name": "db-1",
			"status": "Stopped",
			"config": {"volatile.base_image": "f0f1f2"},
			"state": {
				"network": {
					"lo": {"addresses": [{"family": "inet", "address": "127.0.0.1"}]}
				}
			}
		}
	]`)

	instances, err := parseListJSON(in)
	if err != nil {
		t.Fatalf("parseListJSON returned error: %v", err)
	}
	if len(instances) != 2 {
		t.Fatalf("expected 2 instances, got %d", len(instances))
	}
	if instances[0].Name != "web-1" || instances[0].IP != "10.10.10.22" || instances[0].Image != "ubuntu/24.04" {
		t.Fatalf("unexpected first instance: %+v", instances[0])
	}
	if instances[1].Name != "db-1" || instances[1].IP != "-" || instances[1].Image != "f0f1f2" {
		t.Fatalf("unexpected second instance: %+v", instances[1])
	}
}

func TestParseListJSONInvalid(t *testing.T) {
	if _, err := parseListJSON([]byte(`not-json`)); err == nil {
		t.Fatal("expected parse error, got nil")
	}
}
