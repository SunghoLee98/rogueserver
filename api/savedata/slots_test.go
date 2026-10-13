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

// mockSessionStore records the session calls the savedata functions make.
type mockSessionStore struct {
	stored          []int
	deleted         []int
	read            []int
	seedCompletions int
	dailyRuns       int
	session         defs.SessionSaveData
	seedCompleted   bool
}

func (m *mockSessionStore) UpdateAccountLastActivity(uuid []byte) error { return nil }

func (m *mockSessionStore) StoreSessionSaveData(uuid []byte, data defs.SessionSaveData, slot int) error {
	m.stored = append(m.stored, slot)
	return nil
}

func (m *mockSessionStore) DeleteSessionSaveData(uuid []byte, slot int) error {
	m.deleted = append(m.deleted, slot)
	return nil
}

func (m *mockSessionStore) ReadSessionSaveData(uuid []byte, slot int) (defs.SessionSaveData, error) {
	m.read = append(m.read, slot)
	return m.session, nil
}

func (m *mockSessionStore) ReadSeedCompleted(uuid []byte, seed string) (bool, error) {
	return m.seedCompleted, nil
}

func (m *mockSessionStore) TryAddSeedCompletion(uuid []byte, seed string, mode int) (bool, error) {
	m.seedCompletions++
	return true, nil
}

func (m *mockSessionStore) AddOrUpdateAccountDailyRun(uuid []byte, score int, waveCompleted int) error {
	m.dailyRuns++
	return nil
}

func (m *mockSessionStore) SetAccountBanned(uuid []byte, banned bool) error { return nil }

func (m *mockSessionStore) UpdateSystemSaveData(uuid []byte, data defs.SystemSaveData) error {
	return nil
}

func (m *mockSessionStore) UpdateAccountStats(uuid []byte, stats defs.GameStats, voucherCounts map[string]int) error {
	return nil
}

func (m *mockSessionStore) StoreSystemSaveDataS3(uuid []byte, data defs.SystemSaveData) error {
	return nil
}

func (m *mockSessionStore) StoreSystemSaveData(uuid []byte, data defs.SystemSaveData) error {
	return nil
}

var uuid = []byte("0123456789abcdef")

func TestValidateSessionSlot(t *testing.T) {
	for slot := 0; slot <= 5; slot++ {
		if err := ValidateSessionSlot(slot); err != nil {
			t.Errorf("slot %d: unexpected error %v", slot, err)
		}
	}
	for _, slot := range []int{-1, 6, 100} {
		if err := ValidateSessionSlot(slot); err == nil {
			t.Errorf("slot %d: expected an out of range error", slot)
		}
	}
}

func TestIsCoopSessionSlot(t *testing.T) {
	if !IsCoopSessionSlot(defs.CoopSessionSlot) {
		t.Errorf("slot %d should be the co-op slot", defs.CoopSessionSlot)
	}
	for slot := 0; slot < defs.SessionSlotCount; slot++ {
		if IsCoopSessionSlot(slot) {
			t.Errorf("solo slot %d is not the co-op slot", slot)
		}
	}
}

func TestDeleteAcceptsTheCoopSlot(t *testing.T) {
	store := &mockSessionStore{}
	if err := Delete(store, uuid, 1, defs.CoopSessionSlot); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := Delete(store, uuid, 1, 6); err == nil {
		t.Errorf("slot 6 should be out of range")
	}
	if len(store.deleted) != 1 || store.deleted[0] != defs.CoopSessionSlot {
		t.Errorf("deleted slots %v, want [5]", store.deleted)
	}
}

func TestClearAcceptsTheCoopSlot(t *testing.T) {
	store := &mockSessionStore{}
	if _, err := Clear(store, uuid, defs.CoopSessionSlot, "daily", coopSave()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := Clear(store, uuid, 6, "daily", coopSave()); err == nil {
		t.Errorf("slot 6 should be out of range")
	}
}

func TestNewClearAcceptsTheCoopSlot(t *testing.T) {
	store := &mockSessionStore{}
	if _, err := NewClear(store, uuid, defs.CoopSessionSlot); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := NewClear(store, uuid, 6); err == nil {
		t.Errorf("slot 6 should be out of range")
	}
}

func coopSave() defs.SessionSaveData {
	return defs.SessionSaveData{
		Seed: "coop-seed",
		Coop: &defs.CoopSessionData{
			Protocol:     4,
			Participants: defs.CoopParticipants{Host: "alice", Guest: "bob"},
		},
	}
}

func soloSave() defs.SessionSaveData {
	return defs.SessionSaveData{Seed: "solo-seed"}
}
