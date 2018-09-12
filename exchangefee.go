package gobacktest

// ExchangeFeeHandler is a basic interface for managing the exchange fees.
type ExchangeFeeHandler interface {
	Fee() (float64, error)
}

// FixedExchangeFee returns a fixed exchange fee.
type FixedExchangeFee struct {
	ExchangeFee float64
}

// Fee returns the exchange fee of the trade.
func (e *FixedExchangeFee) Fee() (float64, error) {
	return e.ExchangeFee, nil
}
