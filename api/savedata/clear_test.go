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

package savedata

import (
	"testing"

	"github.com/pagefaultgames/rogueserver/defs"
)

// A finished classic run: wave 200 won against the final boss (battle type 2).
func completedClassic(save defs.SessionSaveData) defs.SessionSaveData {
	save.GameMode = 0
	save.BattleType = 2
	save.WaveIndex = 200
	return save
}

func TestClearOfACoopRunOnlyDeletesTheSave(t *testing.T) {
	store := &mockSessionStore{}
	response, err := Clear(store, uuid, defs.CoopSessionSlot, "daily", completedClassic(coopSave()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.seedCompletions != 0 || store.dailyRuns != 0 {
		t.Errorf("co-op clear recorded seed completions %d, daily runs %d", store.seedCompletions, store.dailyRuns)
	}
	if response.Success {
		t.Errorf("co-op clear reported a new classic clear")
	}
	if len(store.deleted) != 1 || store.deleted[0] != defs.CoopSessionSlot {
		t.Errorf("deleted slots %v, want [5]", store.deleted)
	}
}

func TestClearOfASoloRunIsUnchanged(t *testing.T) {
	store := &mockSessionStore{}
	response, err := Clear(store, uuid, 0, "daily", completedClassic(soloSave()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.seedCompletions != 1 || !response.Success {
		t.Errorf("solo clear: seed completions %d, success %v", store.seedCompletions, response.Success)
	}
	if len(store.deleted) != 1 || store.deleted[0] != 0 {
		t.Errorf("deleted slots %v, want [0]", store.deleted)
	}
}

func TestNewClearIsNeverTrueForTheCoopSlot(t *testing.T) {
	store := &mockSessionStore{session: coopSave(), seedCompleted: false}
	newClear, err := NewClear(store, uuid, defs.CoopSessionSlot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newClear {
		t.Errorf("a co-op run counted as a new classic clear")
	}
	if len(store.read) != 0 {
		t.Errorf("read slots %v, want none", store.read)
	}

	solo, err := NewClear(&mockSessionStore{session: soloSave()}, uuid, 0)
	if err != nil || !solo {
		t.Errorf("solo new clear: %v, %v", solo, err)
	}
}
