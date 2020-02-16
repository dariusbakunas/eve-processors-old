package models

import (
	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v3"
	"time"
)

type WalletTransaction struct {
	ClientId      int64           `json:"client_id"`
	Quantity      int64           `json:"quantity"`
	UnitPrice     decimal.Decimal `json:"unit_price"`
	Date          time.Time       `json:"date"`
	IsBuy         bool            `json:"is_buy"`
	IsPersonal    bool            `json:"is_personal"`
	JournalRefId  int64           `json:"journal_ref_id"`
	LocationId    int64           `json:"location_id"`
	TransactionId int64           `json:"transaction_id"`
	TypeId        int             `json:"type_id"`
}

type JournalEntry struct {
	Amount        decimal.NullDecimal `json:"amount"`
	Balance       decimal.NullDecimal `json:"balance"`
	ContextID     null.Int            `json:"context_id"`
	ContextIDType null.String         `json:"context_id_type"`
	Date          time.Time           `json:"date"`
	Description   string              `json:"description"`
	FirstPartyID  null.Int            `json:"first_party_id"`
	ID            int64               `json:"id"`
	Reason        null.String         `json:"reason"`
	RefType       string              `json:"ref_type"`
	SecondPartyID null.Int            `json:"second_party_id"`
	Tax           decimal.NullDecimal `json:"tax"`
	TaxReceiverID null.Int            `json:"tax_receiver_id"`
}

type JournalEntriesResponse struct {
	Pages   int
	Entries []JournalEntry
}

type JobLogEntry struct {
	ID            null.Int
	CreatedAt     null.Time
	Category      string
	Status        string
	Message       string
	Error         null.String
	CharacterID   null.Int
	CorporationID null.Int
}

type CharacterSkill struct {
	ActiveSkillLevel int   `json:"active_skill_level"`
	SkillID          int   `json:"skill_id"`
	SP               int64 `json:"skillpoints_in_skill"`
	TrainedLevel     int   `json:"trained_skill_level"`
}

type SkillsResponse struct {
	Skills        []CharacterSkill `json:"skills"`
	TotalSP       int64            `json:"total_sp"`
	UnallocatedSP int              `json:"unallocated_sp"`
}

type SkillQueueItem struct {
	FinishDate      null.Time `json:"finish_date"`
	FinishedLevel   int       `json:"finished_level"`
	LevelEndSP      null.Int  `json:"level_end_sp"`
	LevelStartSP    null.Int  `json:"level_start_sp"`
	QueuePosition   int       `json:"queue_position"`
	SkillID         int       `json:"skill_id"`
	StartDate       null.Time `json:"start_date"`
	TrainingStartSP null.Int  `json:"training_start_sp"`
}

type CharacterMarketOrder struct {
	Duration      int                 `json:"duration"`
	Escrow        decimal.NullDecimal `json:"escrow"`
	IsBuy         bool                `json:"is_buy_order"`
	Issued        time.Time           `json:"issued"`
	IsCorporation bool                `json:"is_corporation"`
	LocationID    int64               `json:"location_id"`
	MinVolume     null.Int            `json:"min_volume"`
	OrderID       int64               `json:"order_id"`
	Price         decimal.Decimal     `json:"price"`
	Range         string              `json:"range"`
	RegionID      int                 `json:"region_id"`
	State         null.String         `json:"state"`
	TypeID        int                 `json:"type_id"`
	VolumeRemain  int                 `json:"volume_remain"`
	VolumeTotal   int                 `json:"volume_total"`
}

type MarketOrderHistoryResponse struct {
	Orders []CharacterMarketOrder
	Pages  int
}

type Blueprint struct {
	ID           int64  `json:"item_id"`
	LocationFlag string `json:"location_flag"`
	LocationID   int64  `json:"location_id"`
	ME           int    `json:"material_efficiency"`
	Quantity     int    `json:"quantity"`
	Runs         int    `json:"runs"`
	TE           int    `json:"time_efficiency"`
	TypeID       int    `json:"type_id"`
}

type BlueprintsResponse struct {
	Blueprints []Blueprint
	Pages      int
}

type MarketOrder struct {
	Duration     int             `json:"duration"`
	IsBuy        bool            `json:"is_buy_order"`
	Issued       time.Time       `json:"issued"`
	LocationID   int64           `json:"location_id"`
	MinVolume    int             `json:"min_volume"`
	OrderID      int64           `json:"order_id"`
	Price        decimal.Decimal `json:"price"`
	SystemID     int             `json:"system_id"`
	Range        string          `json:"range"`
	TypeID       int             `json:"type_id"`
	VolumeRemain int             `json:"volume_remain"`
	VolumeTotal  int             `json:"volume_total"`
}

type IndustryJob struct {
	ActivityID           int                 `json:"activity_id"`
	BlueprintID          int64               `json:"blueprint_id"`
	BlueprintLocationID  int64               `json:"blueprint_location_id"`
	BlueprintTypeID      int                 `json:"blueprint_type_id"`
	CompletedCharacterID null.Int            `json:"completed_character_id"`
	CompletedDate        null.Time           `json:"completed_date"`
	Cost                 decimal.NullDecimal `json:"cost"`
	Duration             int                 `json:"duration"`
	EndDate              time.Time           `json:"end_date"`
	FacilityID           int64               `json:"facility_id"`
	InstalledID          int64               `json:"installer_id"`
	JobID                int64               `json:"job_id"`
	LicensedRuns         null.Int            `json:"licensed_runs"`
	OutputLocationID     int64               `json:"output_location_id"`
	PauseDate            null.Time           `json:"pause_date"`
	Probability          decimal.NullDecimal `json:"probability"`
	ProductTypeID        null.Int            `json:"product_type_id"`
	Runs                 int                 `json:"runs"`
	StartDate            time.Time           `json:"start_date"`
	StationID            int64               `json:"station_id"`
	Status               string              `json:"status"`
	SuccessfulRuns       null.Int            `json:"successful_runs"`
}

type MarketOrderPriceItem struct {
	TypeID    int
	BuyPrice  decimal.Decimal
	SellPrice decimal.Decimal
	SystemID  int
}

type MarketOrderResponse struct {
	Orders []MarketOrder
	Pages  int
}
