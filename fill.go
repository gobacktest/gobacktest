package gobacktest

// Fill declares a basic fill event
type Fill struct {
	Event
	direction   OrderDirection // BOT for buy, SLD for sell, HLD for hold
	Exchange    string         // exchange symbol
	qty         int64
	price       float64
	commission  float64
	exchangeFee float64
	cost        float64 // the total cost of the filled order incl commission and fees
}

// Direction returns the direction of a Fill
func (f Fill) Direction() OrderDirection {
	return f.direction
}

// SetDirection sets the Directions field of a Fill
func (f *Fill) SetDirection(dir OrderDirection) {
	f.direction = dir
}

// Qty returns the qty field of a fill
func (f Fill) Qty() int64 {
	return f.qty
}

// SetQty sets the Qty field of a Fill
func (f *Fill) SetQty(i int64) {
	f.qty = i
}

// Price returns the Price field of a fill
func (f Fill) Price() float64 {
	return f.price
}

// Commission returns the Commission field of a fill.
func (f Fill) Commission() float64 {
	return f.commission
}

// ExchangeFee returns the ExchangeFee Field of a fill
func (f Fill) ExchangeFee() float64 {
	return f.exchangeFee
}

// Cost returns the Cost field of a Fill
func (f Fill) Cost() float64 {
	return f.cost
}

// Value returns the value without cost.
func (f Fill) Value() float64 {
	value := float64(f.qty) * f.price
	return value
}

// NetValue returns the net value including cost.
func (f Fill) NetValue() float64 {
	if f.direction == BOT {
		// qty * price + cost
		netValue := float64(f.qty)*f.price + f.cost
		return netValue
	}
	// SLD
	//qty * price - cost
	netValue := float64(f.qty)*f.price - f.cost
	return netValue
}
