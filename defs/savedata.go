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

package defs

// SessionSlotCount is the number of solo session slots (0 to SessionSlotCount-1).
const SessionSlotCount = 5

// CoopSessionSlot is the one session slot for co-op runs, next to the solo slots.
// It only holds saves with a co-op block, and the solo slots never do.
const CoopSessionSlot = 5

type SystemSaveData struct {
	TrainerId          int                `json:"trainerId"`
	SecretId           int                `json:"secretId"`
	Gender             int                `json:"gender"`
	DexData            DexData            `json:"dexData"`
	StarterData        StarterData        `json:"starterData"`
	StarterMoveData    StarterMoveData    `json:"starterMoveData"`    // Legacy
	StarterEggMoveData StarterEggMoveData `json:"starterEggMoveData"` // Legacy
	GameStats          GameStats          `json:"gameStats"`
	Unlocks            Unlocks            `json:"unlocks"`
	AchvUnlocks        AchvUnlocks        `json:"achvUnlocks"`
	VoucherUnlocks     VoucherUnlocks     `json:"voucherUnlocks"`
	VoucherCounts      VoucherCounts      `json:"voucherCounts"`
	Eggs               []EggData          `json:"eggs"`
	EggPity            []int              `json:"eggPity"`
	UnlockPity         []int              `json:"unlockPity"`
	GameVersion        string             `json:"gameVersion"`
	Timestamp          int                `json:"timestamp"`
	AppliedMigrators   map[string]int     `json:"appliedMigrators"`
}

type DexData map[int]DexEntry

type DexEntry struct {
	SeenAttr     interface{} `json:"seenAttr"`   // integer or string
	CaughtAttr   interface{} `json:"caughtAttr"` // integer or string
	NatureAttr   int         `json:"natureAttr"`
	SeenCount    int         `json:"seenCount"`
	CaughtCount  int         `json:"caughtCount"`
	HatchedCount int         `json:"hatchedCount"`
	Ivs          []int       `json:"ivs"`
	Ribbons      string      `json:"ribbons"`
}

type StarterData map[int]StarterEntry

type StarterEntry struct {
	Moveset         interface{} `json:"moveset"`
	EggMoves        int         `json:"eggMoves"`
	CandyCount      int         `json:"candyCount"`
	Friendship      int         `json:"friendship"`
	AbilityAttr     int         `json:"abilityAttr"`
	PassiveAttr     int         `json:"passiveAttr"`
	ValueReduction  int         `json:"valueReduction"`
	ClassicWinCount int         `json:"classicWinCount"`
}

type StarterMoveData map[int]interface{}

type StarterEggMoveData map[int]int

type GameStats interface{}

type Unlocks map[int]bool

type AchvUnlocks map[string]int

type VoucherUnlocks map[string]int

type VoucherCounts map[string]int

type EggData struct {
	Id                    int       `json:"id"`
	GachaType             GachaType `json:"gachaType"`
	HatchWaves            int       `json:"hatchWaves"`
	Timestamp             int       `json:"timestamp"`
	Tier                  int       `json:"tier"`
	SourceType            int       `json:"sourceType"`
	VariantTier           int       `json:"variantTier"`
	IsShiny               bool      `json:"isShiny"`
	Species               int       `json:"species"`
	EggMoveIndex          int       `json:"eggMoveIndex"`
	OverrideHiddenAbility bool      `json:"overrideHiddenAbility"`
}

type GachaType int

type SessionSaveData struct {
	Seed                     string                   `json:"seed"`
	PlayTime                 int                      `json:"playTime"`
	GameMode                 GameMode                 `json:"gameMode"`
	DailyConfig              DailyConfig              `json:"dailyConfig,omitempty"`
	Party                    []PokemonData            `json:"party"`
	EnemyParty               []PokemonData            `json:"enemyParty"`
	Modifiers                []PersistentModifierData `json:"modifiers"`
	EnemyModifiers           []PersistentModifierData `json:"enemyModifiers"`
	Arena                    ArenaData                `json:"arena"`
	PokeballCounts           PokeballCounts           `json:"pokeballCounts"`
	Money                    int                      `json:"money"`
	Score                    int                      `json:"score"`
	VictoryCount             int                      `json:"victoryCount"`
	FaintCount               int                      `json:"faintCount"`
	ReviveCount              int                      `json:"reviveCount"`
	WaveIndex                int                      `json:"waveIndex"`
	BattleType               BattleType               `json:"battleType"`
	Trainer                  TrainerData              `json:"trainer"`
	GameVersion              string                   `json:"gameVersion"`
	Timestamp                int                      `json:"timestamp"`
	Challenges               []ChallengeData          `json:"challenges"`
	MysteryEncounterType     MysteryEncounterType     `json:"mysteryEncounterType"`
	MysteryEncounterSaveData MysteryEncounterSaveData `json:"mysteryEncounterSaveData"`
	Name                     string                   `json:"name,omitempty"`
	Coop                     *CoopSessionData         `json:"coop,omitempty"`
}

// CoopSessionData is what a co-op run saves besides the solo session fields. Only the host's account keeps the
// save; the guest's client never writes one. The server stores it as is and checks nothing but its presence.
type CoopSessionData struct {
	// Protocol is the co-op protocol version of the client that saved
	Protocol     int              `json:"protocol"`
	Participants CoopParticipants `json:"participants"`
	// Ownership maps party pokemon ids to the role owning them ("host" or "guest")
	Ownership map[int]string `json:"ownership"`
	// Accounts holds the account data each player's simulation used, by role
	Accounts   CoopAccounts `json:"accounts"`
	TimedEvent *string      `json:"timedEvent"`
	// DecisionSeq, RewardGen and RewardSeq are the co-op counters both clients resume from
	DecisionSeq int `json:"decisionSeq"`
	RewardGen   int `json:"rewardGen"`
	RewardSeq   int `json:"rewardSeq"`
}

// CoopParticipants holds the account names of both players.
type CoopParticipants struct {
	Host  string `json:"host"`
	Guest string `json:"guest"`
}

type CoopAccounts struct {
	Host  CoopAccountSnapshot `json:"host"`
	Guest CoopAccountSnapshot `json:"guest"`
}

// CoopAccountSnapshot is the part of an account the co-op battle simulation reads.
type CoopAccountSnapshot struct {
	TrainerId     int         `json:"trainerId"`
	SecretId      int         `json:"secretId"`
	CaughtSpecies []int       `json:"caughtSpecies"`
	EggMoves      map[int]int `json:"eggMoves"`
	Unlocks       []int       `json:"unlocks"`
}

type ChallengeData struct {
	Id       int `json:"id"`
	Value    int `json:"value"`
	Severity int `json:"severity"`
}

type MysteryEncounterType int

type MysteryEncounterTier int

type SeenEncounterData struct {
	Type           MysteryEncounterType `json:"type"`
	Tier           MysteryEncounterTier `json:"tier"`
	WaveIndex      int                  `json:"waveIndex"`
	SelectedOption int                  `json:"selectedOption"`
}

type QueuedEncounter struct {
	Type         MysteryEncounterType `json:"type"`
	SpawnPercent int                  `json:"spawnPercent"`
}

type MysteryEncounterSaveData struct {
	EncounteredEvents    []SeenEncounterData `json:"encounteredEvents"`
	EncounterSpawnChance int                 `json:"encounterSpawnChance"`
	QueuedEncounters     []QueuedEncounter   `json:"queuedEncounters"`
}

type GameMode int

type DailyConfig interface{}

type PokemonData interface{}

type PersistentModifierData interface{}

type ArenaData interface{}

type PokeballCounts map[string]int

type BattleType int

type TrainerData interface{}

type SessionHistoryData struct {
	Seed        string                   `json:"seed"`
	PlayTime    int                      `json:"playTime"`
	Result      SessionHistoryResult     `json:"sessionHistoryResult"`
	GameMode    GameMode                 `json:"gameMode"`
	Party       []PokemonData            `json:"party"`
	Modifiers   []PersistentModifierData `json:"modifiers"`
	Money       int                      `json:"money"`
	Score       int                      `json:"score"`
	WaveIndex   int                      `json:"waveIndex"`
	BattleType  BattleType               `json:"battleType"`
	GameVersion string                   `json:"gameVersion"`
	Timestamp   int                      `json:"timestamp"`
}

type SessionHistoryResult int
