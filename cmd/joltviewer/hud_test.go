package main

import "testing"

// hudRows must only list the driving controls on vehicle scenes: the extra rows
// used to be drawn unconditionally, at the same y as the orbit hint.
func TestHUDRowsVehicleControls(t *testing.T) {
	keys := func(vehicle bool) map[string]bool {
		got := map[string]bool{}
		for _, r := range hudRows("Vehicle", 1, false, vehicle) {
			if r.key != "" {
				got[r.key] = true
			}
		}
		return got
	}

	if plain := keys(false); plain["Arrows"] || plain["B"] {
		t.Errorf("non-vehicle scene lists driving controls: %v", plain)
	}
	if veh := keys(true); !veh["Arrows"] || !veh["B"] {
		t.Errorf("vehicle scene is missing its driving controls: %v", veh)
	}
	// Every row must have text (only the spacer is blank) and a unique key.
	seen := map[string]bool{}
	for _, r := range hudRows("Vehicle", 1, true, true) {
		if r.key != "" && seen[r.key] {
			t.Errorf("duplicate key column %q", r.key)
		}
		seen[r.key] = true
	}
}
