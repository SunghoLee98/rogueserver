/*
	Copyright (C) 2024 - 2025  Pagefault Games

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package db

import (
	"encoding/gob"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/pagefaultgames/rogueserver/defs"
)

func init() {
	// The server registers these at startup (rogueserver.go) for the untyped save fields.
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})
}

// A co-op session save as the frontend sends it (party and arena trimmed).
const coopSessionJSON = `{
	"seed": "coop-seed",
	"playTime": 120,
	"gameMode": 0,
	"party": [{"id": 101, "species": 1}],
	"enemyParty": [],
	"modifiers": [],
	"enemyModifiers": [],
	"arena": {"biome": 0},
	"pokeballCounts": {"0": 5},
	"money": 1000,
	"score": 0,
	"waveIndex": 2,
	"battleType": 0,
	"trainer": null,
	"gameVersion": "1.12.1",
	"timestamp": 1760000000000,
	"challenges": [],
	"mysteryEncounterType": -1,
	"mysteryEncounterSaveData": {"encounteredEvents": [], "encounterSpawnChance": 3, "queuedEncounters": []},
	"coop": {
		"protocol": 4,
		"participants": {"host": "alice", "guest": "bob"},
		"ownership": {"101": "host", "102": "guest"},
		"accounts": {
			"host": {"trainerId": 11, "secretId": 12, "caughtSpecies": [1, 4], "eggMoves": {"1": 3}, "unlocks": [0, 2]},
			"guest": {"trainerId": 21, "secretId": 22, "caughtSpecies": [7], "eggMoves": {"4": 1}, "unlocks": [1]}
		},
		"timedEvent": "Halloween",
		"decisionSeq": 7,
		"rewardGen": 3,
		"rewardSeq": 9
	}
}`

func roundTrip(t *testing.T, raw string) map[string]interface{} {
	t.Helper()

	var session defs.SessionSaveData
	if err := json.Unmarshal([]byte(raw), &session); err != nil {
		t.Fatalf("decode request: %s", err)
	}

	stored, err := encodeSessionSaveData(session)
	if err != nil {
		t.Fatalf("encode for storage: %s", err)
	}

	loaded, err := decodeSessionSaveData(stored)
	if err != nil {
		t.Fatalf("decode from storage: %s", err)
	}

	out, err := json.Marshal(loaded)
	if err != nil {
		t.Fatalf("encode response: %s", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("parse response: %s", err)
	}

	return result
}

func TestSessionSaveDataKeepsCoopBlock(t *testing.T) {
	result := roundTrip(t, coopSessionJSON)

	var original map[string]interface{}
	if err := json.Unmarshal([]byte(coopSessionJSON), &original); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(result["coop"], original["coop"]) {
		t.Errorf("co-op block changed in the round trip:\n got  %v\n want %v", result["coop"], original["coop"])
	}
}

func TestSoloSessionSaveDataHasNoCoopBlock(t *testing.T) {
	var solo map[string]interface{}
	if err := json.Unmarshal([]byte(coopSessionJSON), &solo); err != nil {
		t.Fatal(err)
	}
	delete(solo, "coop")
	raw, err := json.Marshal(solo)
	if err != nil {
		t.Fatal(err)
	}

	result := roundTrip(t, string(raw))

	if _, ok := result["coop"]; ok {
		t.Errorf("solo save gained a co-op block: %v", result["coop"])
	}
	if result["seed"] != "coop-seed" || result["waveIndex"] != float64(2) {
		t.Errorf("solo fields changed in the round trip: %v", result)
	}
}
