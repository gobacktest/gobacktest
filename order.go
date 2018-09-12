package gobacktest

// OrderStatus defines an order status.
type OrderStatus int

// different types of order status
const (
	OrderNone OrderStatus = iota // 0
	OrderNew
	OrderSubmitted
	OrderPartiallyFilled
	OrderFilled
	OrderCanceled
	OrderCancelPending
	OrderInvalid
)

// OrderType defines the type of an order.
type OrderType int

// different types of orders
const (
	MarketOrder OrderType = iota // 0
	MarketOnOpenOrder
	MarketOnCloseOrder
	StopMarketOrder
	LimitOrder
	StopLimitOrder
)

// Order declares a basic order event.
type Order struct {
	Event
	id           int
	orderType    OrderType // market or limit
	status       OrderStatus
	direction    Direction // buy or sell
	assetType    string
	qty          int64 // quantity of the order
	qtyFilled    int64
	avgFillPrice float64
	limitPrice   float64 // limit for the order
	stopPrice    float64
}

// ID returns the id of the Order.
func (o Order) ID() int {
	return o.id
}

// SetID of the Order.
func (o *Order) SetID(id int) {
	o.id = id
}

// Direction returns the direction of an Order.
func (o Order) Direction() Direction {
	return o.direction
}

// SetDirection sets the direction of an Order.
func (o *Order) SetDirection(dir Direction) {
	o.direction = dir
}

// Qty returns the quantity of an Order.
func (o Order) Qty() int64 {
	return o.qty
}

// SetQty sets the quantity of an Order.
func (o *Order) SetQty(i int64) {
	o.qty = i
}

// Status returns the status of an Order.
func (o Order) Status() OrderStatus {
	return o.status
}

// Limit returns the limit price of an Order.
func (o Order) Limit() float64 {
	return o.limitPrice
}

// Stop returns the stop price of an Order.
func (o Order) Stop() float64 {
	return o.stopPrice
}

// Cancel an order.
func (o *Order) Cancel() {
	o.status = OrderCancelPending
}

// Update an order on a fill event.
func (o *Order) Update(fill FillEvent) {
	// not implemented
}
