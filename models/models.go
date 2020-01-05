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
	ContextIDType null.String		`json:"context_id_type"`
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