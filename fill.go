package gobacktest

// Fill defines a fill event, when an order is executed and filled.
type Fill interface {
}

// FillDirection indicates the direction of a fill
type FillDirection int

const (
	BOT FillDirection = iota // 0
	SLD
)

// fill is a basic Fill event implementation.
type fill struct {
	Event
	direction   FillDirection // BOT for buy, SLD for sell
	Exchange    string        // exchange symbol
	qty         int64
	price       float64
	commission  float64
	exchangeFee float64
	cost        float64 // the total cost of the filled order incl commission and fees
}

// Direction returns the direction of a Fill.
func (f fill) Direction() FillDirection {
	return f.direction
}

// SetDirection sets the direction of a Fill.
func (f *fill) SetDirection(dir FillDirection) {
	f.direction = dir
}

// Qty returns the quantity of a fill.
func (f fill) Qty() int64 {
	return f.qty
}

// SetQty sets the quantity of a Fill.
func (f *fill) SetQty(i int64) {
	f.qty = i
}

// Price returns the price of a fill.
func (f fill) Price() float64 {
	return f.price
}

// Commission returns the commission of a fill.
func (f fill) Commission() float64 {
	return f.commission
}

// ExchangeFee returns the exchange fee of a fill.
func (f fill) ExchangeFee() float64 {
	return f.exchangeFee
}

// Cost returns the cost of a Fill.
func (f fill) Cost() float64 {
	return f.cost
}

// Value returns the value excluding cost.
func (f fill) Value() float64 {
	value := float64(f.qty) * f.price
	return value
}

// NetValue returns the net value including cost.
func (f fill) NetValue() float64 {
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
