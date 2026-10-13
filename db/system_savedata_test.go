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
	"encoding/json"
	"reflect"
	"slices"
	"testing"

	"github.com/pagefaultgames/rogueserver/defs"
)

// A system save as a co-op player's client sends it (dex and starters trimmed): co-op run counts in the game stats
// and the co-op clear achievement.
const coopSystemJSON = `{
	"trainerId": 1234,
	"secretId": 5678,
	"gender": 0,
	"dexData": {},
	"starterData": {},
	"gameStats": {"playTime": 3600, "battles": 12, "classicSessionsPlayed": 2, "sessionsWon": 0, "coopSessionsPlayed": 3, "coopSessionsWon": 1},
	"unlocks": {"0": true},
	"achvUnlocks": {"CLASSIC_VICTORY": 0, "COOP_VICTORY": 1760000000000},
	"voucherUnlocks": {},
	"voucherCounts": {"0": 1, "1": 0, "2": 0, "3": 0},
	"eggs": [],
	"eggPity": [0, 0, 0, 0],
	"unlockPity": [0, 0, 0, 0],
	"gameVersion": "1.12.1",
	"timestamp": 1760000000000
}`

func TestSystemSaveDataKeepsCoopStatsAndAchievement(t *testing.T) {
	var system defs.SystemSaveData
	if err := json.Unmarshal([]byte(coopSystemJSON), &system); err != nil {
		t.Fatalf("decode request: %s", err)
	}

	stored, err := encodeSystemSaveData(system)
	if err != nil {
		t.Fatalf("encode for storage: %s", err)
	}
	loaded, err := decodeSystemSaveData(stored)
	if err != nil {
		t.Fatalf("decode from storage: %s", err)
	}

	stats, ok := loaded.GameStats.(map[string]interface{})
	if !ok {
		t.Fatalf("game stats came back as %T", loaded.GameStats)
	}
	if stats["coopSessionsPlayed"] != 3.0 || stats["coopSessionsWon"] != 1.0 {
		t.Fatalf("co-op stats were not kept: %v", stats)
	}
	if loaded.AchvUnlocks["COOP_VICTORY"] != 1760000000000 {
		t.Fatalf("co-op achievement was not kept: %v", loaded.AchvUnlocks)
	}
	if !reflect.DeepEqual(loaded.GameStats, system.GameStats) {
		t.Fatalf("game stats changed: %v != %v", loaded.GameStats, system.GameStats)
	}
}

func TestAccountStatColumnsSkipCoopStats(t *testing.T) {
	stats := map[string]interface{}{
		"playTime":           3600.0,
		"sessionsWon":        0.0,
		"coopSessionsPlayed": 3.0,
		"coopSessionsWon":    1.0,
	}

	columns, values, err := accountStatColumns(stats, map[string]int{"0": 1, "9": 4})
	if err != nil {
		t.Fatal(err)
	}

	got := map[string]interface{}{}
	for i, column := range columns {
		got[column] = values[i]
	}
	want := map[string]interface{}{"playTime": 3600.0, "sessionsWon": 0.0, "regularVouchers": 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("columns %v, want %v", got, want)
	}
	if slices.Contains(columns, "coopSessionsPlayed") {
		t.Fatal("co-op stats have no account stats column")
	}
}

func TestAccountStatColumnsRejectNonNumbers(t *testing.T) {
	if _, _, err := accountStatColumns(map[string]interface{}{"coopSessionsPlayed": "3"}, nil); err == nil {
		t.Fatal("expected an error for a stat that is not a number")
	}
	if _, _, err := accountStatColumns([]interface{}{}, nil); err == nil {
		t.Fatal("expected an error for stats that are not an object")
	}
}
