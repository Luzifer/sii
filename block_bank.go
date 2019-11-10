package sii

func init() {
	RegisterBlock(&Bank{})
}

type Bank struct {
	MoneyAccount            int64   `sii:"money_account"`
	CoinsuranceFixed        int64   `sii:"coinsurance_fixed"`
	CoinsuranceRatio        float32 `sii:"coinsurance_ratio"`
	AccidentSeverity        float32 `sii:"accident_severity"`
	Loans                   int64   `sii:"loans"`
	AppEnabled              bool    `sii:"app_enabled"`
	LoanLimit               int64   `sii:"loan_limit"`
	PaymentTimer            float32 `sii:"payment_timer"`
	Overdraft               bool    `sii:"overdraft"`
	OverdraftTimer          float32 `sii:"overdraft_timer"`
	OverdraftWarnCount      int64   `sii:"overdraft_warn_count"`
	SellPlayersTruckLater   bool    `sii:"sell_players_truck_later"`
	SellPlayersTrailerLater bool    `sii:"sell_players_trailer_later"`

	blockName string
}

func (Bank) Class() string { return "bank" }

func (b *Bank) Init(class, name string) {
	b.blockName = name
}

func (b Bank) Name() string { return b.blockName }
