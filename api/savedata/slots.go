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
	"fmt"

	"github.com/pagefaultgames/rogueserver/defs"
)

// ValidateSessionSlot checks a session slot id: one of the solo slots or the co-op slot.
func ValidateSessionSlot(slot int) error {
	if (slot < 0 || slot >= defs.SessionSlotCount) && !IsCoopSessionSlot(slot) {
		return fmt.Errorf("slot id %d out of range", slot)
	}

	return nil
}

// IsCoopSessionSlot reports whether the slot is the one for co-op runs.
func IsCoopSessionSlot(slot int) bool {
	return slot == defs.CoopSessionSlot
}

// ValidateSessionSlotData checks that a session save fits its slot:
// the co-op slot only takes saves with a co-op block, and the solo slots never do.
func ValidateSessionSlotData(slot int, save defs.SessionSaveData) error {
	err := ValidateSessionSlot(slot)
	if err != nil {
		return err
	}

	if IsCoopSessionSlot(slot) && save.Coop == nil {
		return fmt.Errorf("slot id %d only accepts co-op saves", slot)
	}

	if !IsCoopSessionSlot(slot) && save.Coop != nil {
		return fmt.Errorf("slot id %d does not accept co-op saves", slot)
	}

	return nil
}
