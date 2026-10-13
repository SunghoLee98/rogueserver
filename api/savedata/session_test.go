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

func TestUpdateSessionKeepsSoloAndCoopApart(t *testing.T) {
	tests := []struct {
		name    string
		slot    int
		save    defs.SessionSaveData
		wantErr bool
	}{
		{name: "solo save in a solo slot", slot: 0, save: soloSave()},
		{name: "solo save in the last solo slot", slot: 4, save: soloSave()},
		{name: "co-op save in the co-op slot", slot: defs.CoopSessionSlot, save: coopSave()},
		{name: "co-op save in a solo slot", slot: 2, save: coopSave(), wantErr: true},
		{name: "solo save in the co-op slot", slot: defs.CoopSessionSlot, save: soloSave(), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockSessionStore{}
			err := UpdateSession(store, uuid, tt.slot, tt.save)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error")
				}
				if len(store.stored) != 0 {
					t.Errorf("rejected save was stored in slots %v", store.stored)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(store.stored) != 1 || store.stored[0] != tt.slot {
				t.Errorf("stored slots %v, want [%d]", store.stored, tt.slot)
			}
		})
	}
}

// /savedata/updateall goes through Update; it must apply the same rules.
func TestUpdateKeepsSoloAndCoopApart(t *testing.T) {
	store := &mockSessionStore{}
	if err := Update(store, uuid, 1, coopSave()); err == nil {
		t.Errorf("a co-op save in solo slot 1 should be rejected")
	}
	if err := Update(store, uuid, defs.CoopSessionSlot, soloSave()); err == nil {
		t.Errorf("a solo save in the co-op slot should be rejected")
	}
	if err := Update(store, uuid, defs.CoopSessionSlot, coopSave()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := Update(store, uuid, 6, coopSave()); err == nil {
		t.Errorf("slot 6 should be out of range")
	}
	if len(store.stored) != 1 || store.stored[0] != defs.CoopSessionSlot {
		t.Errorf("stored slots %v, want [5]", store.stored)
	}
}
