package application

import (
	"github.com/shopspring/decimal"
	"github.com/tdex-network/tdex-daemon/internal/core/domain"
)

type Market struct {
	BaseAsset  string
	QuoteAsset string
}

type MarketWithFee struct {
	Market
	Fee
}

type Fee struct {
	FeeAsset   string
	BasisPoint int64
}

type MarketWithPrice struct {
	Market
	Price
}

type Price struct {
	BasePrice  decimal.Decimal
	QuotePrice decimal.Decimal
}

type PriceWithFee struct {
	Price
	Fee
	Amount uint64
}

type MarketStrategy struct {
	Market
	Strategy domain.StrategyType
}