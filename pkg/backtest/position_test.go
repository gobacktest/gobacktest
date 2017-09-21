package backtest

import (
	"reflect"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	// set the example time string in format yyyy-mm-dd
	var exampleTime, _ = time.Parse("2006-01-02", "2017-06-01")

	// testCases is a table for testing creation of a position
	var testCases = []struct {
		msg    string    // error message
		fill   FillEvent // input
		expPos *position // expected Position
	}{
		{"create with buy:",
			&Fill{
				Event:     Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Exchange:  "TEST",
				Direction: "BOT", // BOT for buy or SLD for sell
				Qty:       10, Price: 10,
				Commission: 4, ExchangeFee: 1, Cost: 5,
			},
			&position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: 10, qtyBOT: 10, qtySLD: 0,
				avgPrice: 10, avgPriceNet: 10.5, avgPriceBOT: 10, avgPriceSLD: 0,
				value: -100, valueBOT: 100, valueSLD: 0,
				netValue: -105, netValueBOT: 105, netValueSLD: 0,
				marketPrice: 10, marketValue: 100,
				commission: 4, exchangeFee: 1, cost: 5, costBasis: 105,
				realProfitLoss: 0, unrealProfitLoss: -5, totalProfitLoss: -5,
			},
		},
		{"create with sell:",
			&Fill{
				Event:     Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Exchange:  "TEST",
				Direction: "SLD", // BOT for buy or SLD for sell
				Qty:       10, Price: 10,
				Commission: 4, ExchangeFee: 1, Cost: 5,
			},
			&position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: -10, qtyBOT: 0, qtySLD: 10,
				avgPrice: 10, avgPriceNet: 9.5, avgPriceBOT: 0, avgPriceSLD: 10,
				value: 100, valueBOT: 0, valueSLD: 100,
				netValue: 95, netValueBOT: 0, netValueSLD: 95,
				marketPrice: 10, marketValue: 100,
				commission: 4, exchangeFee: 1, cost: 5, costBasis: -95,
				realProfitLoss: 0, unrealProfitLoss: -5, totalProfitLoss: -5,
			},
		},
	}

	for _, tc := range testCases {
		// initialize new Position ready for use
		var p = new(position)
		p.Create(tc.fill)
		if !reflect.DeepEqual(p, tc.expPos) {
			t.Errorf("%v\nCreate(%v): \nexpected %p %#v, \nactual   %p %#v", tc.msg, tc.fill, tc.expPos, tc.expPos, p, p)
		}
	}
}

func TestUpdate(t *testing.T) {
	// set the example time string in format yyyy-mm-dd
	var exampleTime, _ = time.Parse("2006-01-02", "2017-06-01")

	var posBOT = &position{
		timestamp: exampleTime, symbol: "TEST.DE",
		qty: 10, qtyBOT: 10, qtySLD: 0,
		avgPrice: 10, avgPriceNet: 10.5, avgPriceBOT: 10, avgPriceSLD: 0,
		value: -100, valueBOT: 100, valueSLD: 0,
		netValue: -105, netValueBOT: 105, netValueSLD: 0,
		marketPrice: 10, marketValue: 100,
		commission: 4, exchangeFee: 1, cost: 5, costBasis: 105,
		realProfitLoss: 0, unrealProfitLoss: -5, totalProfitLoss: -5,
	}
	var posSLD = &position{
		timestamp: exampleTime, symbol: "TEST.DE",
		qty: -10, qtyBOT: 0, qtySLD: 10,
		avgPrice: 10, avgPriceNet: 9.5, avgPriceBOT: 0, avgPriceSLD: 10,
		value: 100, valueBOT: 0, valueSLD: 100,
		netValue: 95, netValueBOT: 0, netValueSLD: 95,
		marketPrice: 10, marketValue: 100,
		commission: 4, exchangeFee: 1, cost: 5, costBasis: -95,
		realProfitLoss: 0, unrealProfitLoss: -5, totalProfitLoss: -5,
	}

	// testCases is a table for testing updating a position
	var testCases = []struct {
		msg    string    // error string
		pos    *position // base position
		fill   FillEvent // input
		expPos *position // expected Position
	}{
		{"BOT position, buying stock:",
			posBOT,
			&Fill{
				Event:    Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Exchange: "TEST", Direction: "BOT", // BOT for buy or SLD for sell
				Qty: 15, Price: 15,
				Commission: 6, ExchangeFee: 1, Cost: 7,
			},
			&position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: 25, qtyBOT: 25, qtySLD: 0,
				avgPrice: 13, avgPriceNet: 13.48, avgPriceBOT: 13, avgPriceSLD: 0,
				value: -325, valueBOT: 325, valueSLD: 0,
				netValue: -337, netValueBOT: 337, netValueSLD: 0,
				marketPrice: 15, marketValue: 375,
				commission: 10, exchangeFee: 2, cost: 12, costBasis: 337,
				realProfitLoss: 0, unrealProfitLoss: 38, totalProfitLoss: 38,
			},
		},
		{"BOT position, selling stock:",
			posBOT,
			&Fill{
				Event:    Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Exchange: "TEST", Direction: "SLD", // BOT for buy or SLD for sell
				Qty: 6, Price: 12,
				Commission: 4, ExchangeFee: 1, Cost: 5,
			},
			&position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: 4, qtyBOT: 10, qtySLD: 6,
				avgPrice: 10.75, avgPriceNet: 10.75, avgPriceBOT: 10, avgPriceSLD: 12,
				value: -28, valueBOT: 100, valueSLD: 72,
				netValue: -38, netValueBOT: 105, netValueSLD: 67,
				marketPrice: 12, marketValue: 48,
				commission: 8, exchangeFee: 2, cost: 10, costBasis: 42,
				realProfitLoss: 4, unrealProfitLoss: 6, totalProfitLoss: 10,
			},
		},
		{"BOT position, selling, turning SLD position:",
			posBOT,
			&Fill{
				Event:    Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Exchange: "TEST", Direction: "SLD", // BOT for buy or SLD for sell
				Qty: 15, Price: 5,
				Commission: 4, ExchangeFee: 1, Cost: 5,
			},
			&position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: -5, qtyBOT: 10, qtySLD: 15,
				avgPrice: 7, avgPriceNet: 7, avgPriceBOT: 10, avgPriceSLD: 5,
				value: -25, valueBOT: 100, valueSLD: 75,
				netValue: -35, netValueBOT: 105, netValueSLD: 70,
				marketPrice: 5, marketValue: 25,
				commission: 8, exchangeFee: 2, cost: 10, costBasis: -52.5,
				realProfitLoss: -87.5, unrealProfitLoss: 27.5, totalProfitLoss: -60,
			},
		},
		{"BOT position, exit stock:",
			posBOT,
			&Fill{
				Event:    Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Exchange: "TEST", Direction: "SLD", // BOT for buy or SLD for sell
				Qty: 10, Price: 12,
				Commission: 5, ExchangeFee: 1, Cost: 6,
			},
			&position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: 0, qtyBOT: 10, qtySLD: 10,
				avgPrice: 11, avgPriceNet: 10.95, avgPriceBOT: 10, avgPriceSLD: 12,
				value: 20, valueBOT: 100, valueSLD: 120,
				netValue: 9, netValueBOT: 105, netValueSLD: 114,
				marketPrice: 12, marketValue: 0,
				commission: 9, exchangeFee: 2, cost: 11, costBasis: 0,
				realProfitLoss: 9, unrealProfitLoss: 0, totalProfitLoss: 9,
			},
		},
		{"SLD position, selling stock:",
			posSLD,
			&Fill{
				Event:    Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Exchange: "TEST", Direction: "SLD", // BOT for buy or SLD for sell
				Qty: 15, Price: 15,
				Commission: 6, ExchangeFee: 1, Cost: 7,
			},
			&position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: -25, qtyBOT: 0, qtySLD: 25,
				avgPrice: 13, avgPriceNet: 12.52, avgPriceBOT: 0, avgPriceSLD: 13,
				value: 325, valueBOT: 0, valueSLD: 325,
				netValue: 313, netValueBOT: 0, netValueSLD: 313,
				marketPrice: 15, marketValue: 375,
				commission: 10, exchangeFee: 2, cost: 12, costBasis: -313,
				realProfitLoss: 0, unrealProfitLoss: -62, totalProfitLoss: -62,
			},
		},
		{"SLD position, buying stock:",
			posSLD,
			&Fill{
				Event:    Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Exchange: "TEST", Direction: "BOT", // BOT for buy or SLD for sell
				Qty: 6, Price: 12,
				Commission: 4, ExchangeFee: 1, Cost: 5,
			},
			&position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: -4, qtyBOT: 6, qtySLD: 10,
				avgPrice: 10.75, avgPriceNet: 10.75, avgPriceBOT: 12, avgPriceSLD: 10,
				value: 28, valueBOT: 72, valueSLD: 100,
				netValue: 18, netValueBOT: 77, netValueSLD: 95,
				marketPrice: 12, marketValue: 48,
				commission: 8, exchangeFee: 2, cost: 10, costBasis: -38,
				realProfitLoss: -20, unrealProfitLoss: -10, totalProfitLoss: -30,
			},
		},
		{"SLD position, buying, turning BOT position:",
			posSLD,
			&Fill{
				Event:    Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Exchange: "TEST", Direction: "BOT", // BOT for buy or SLD for sell
				Qty: 15, Price: 5,
				Commission: 4, ExchangeFee: 1, Cost: 5,
			},
			&position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: 5, qtyBOT: 15, qtySLD: 10,
				avgPrice: 7, avgPriceNet: 7, avgPriceBOT: 5, avgPriceSLD: 10,
				value: 25, valueBOT: 75, valueSLD: 100,
				netValue: 15, netValueBOT: 80, netValueSLD: 95,
				marketPrice: 5, marketValue: 25,
				commission: 8, exchangeFee: 2, cost: 10, costBasis: 47.5,
				realProfitLoss: 62.5, unrealProfitLoss: -22.5, totalProfitLoss: 40,
			},
		},
		{"SLD position, exit stock:",
			posSLD,
			&Fill{
				Event:    Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Exchange: "TEST", Direction: "BOT", // BOT for buy or SLD for sell
				Qty: 10, Price: 12,
				Commission: 5, ExchangeFee: 1, Cost: 6,
			},
			&position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: 0, qtyBOT: 10, qtySLD: 10,
				avgPrice: 11, avgPriceNet: 11.05, avgPriceBOT: 12, avgPriceSLD: 10,
				value: -20, valueBOT: 120, valueSLD: 100,
				netValue: -31, netValueBOT: 126, netValueSLD: 95,
				marketPrice: 12, marketValue: 0,
				commission: 9, exchangeFee: 2, cost: 11, costBasis: 0,
				realProfitLoss: -31, unrealProfitLoss: 0, totalProfitLoss: -31,
			},
		},
	}

	for _, tc := range testCases {
		// initialize new Position and copy pointer to struct from testcases
		p := &position{}
		*p = *tc.pos
		p.Update(tc.fill)
		// Check single values of position
		// if p.qty != tc.expPos.qty {
		// 	t.Errorf("%v qty: expected %#v, actual %#v", tc.msg, tc.expPos.qty, p.qty)
		// }
		// if p.qtyBOT != tc.expPos.qtyBOT {
		// 	t.Errorf("%v qtyBOT: expected %#v, actual %#v", tc.msg, tc.expPos.qtyBOT, p.qtyBOT)
		// }
		// if p.qtySLD != tc.expPos.qtySLD {
		// 	t.Errorf("%v qtySLD: expected %#v, actual %#v", tc.msg, tc.expPos.qtySLD, p.qtySLD)
		// }
		// if p.avgPrice != tc.expPos.avgPrice {
		// 	t.Errorf("%v avgPrice: expected %#v, actual %#v", tc.msg, tc.expPos.avgPrice, p.avgPrice)
		// }
		// if p.avgPriceNet != tc.expPos.avgPriceNet {
		// 	t.Errorf("%v avgPriceNet: expected %#v, actual %#v", tc.msg, tc.expPos.avgPriceNet, p.avgPriceNet)
		// }
		// if p.avgPriceBOT != tc.expPos.avgPriceBOT {
		// 	t.Errorf("%v avgPriceBOT: expected %#v, actual %#v", tc.msg, tc.expPos.avgPriceBOT, p.avgPriceBOT)
		// }
		// if p.avgPriceSLD != tc.expPos.avgPriceSLD {
		// 	t.Errorf("%v avgPriceSLD: expected %#v, actual %#v", tc.msg, tc.expPos.avgPriceSLD, p.avgPriceSLD)
		// }
		// if p.value != tc.expPos.value {
		// 	t.Errorf("%v value: expected %#v, actual %#v", tc.msg, tc.expPos.value, p.value)
		// }
		// if p.valueBOT != tc.expPos.valueBOT {
		// 	t.Errorf("%v valueBOT: expected %#v, actual %#v", tc.msg, tc.expPos.valueBOT, p.valueBOT)
		// }
		// if p.valueSLD != tc.expPos.valueSLD {
		// 	t.Errorf("%v valueSLD: expected %#v, actual %#v", tc.msg, tc.expPos.valueSLD, p.valueSLD)
		// }
		// if p.marketPrice != tc.expPos.marketPrice {
		// 	t.Errorf("%v marketPrice: expected %#v, actual %#v", tc.msg, tc.expPos.marketPrice, p.marketPrice)
		// }
		// if p.marketValue != tc.expPos.marketValue {
		// 	t.Errorf("%v marketValue: expected %#v, actual %#v", tc.msg, tc.expPos.marketValue, p.marketValue)
		// }
		// if p.commission != tc.expPos.commission {
		// 	t.Errorf("%v commission: expected %#v, actual %#v", tc.msg, tc.expPos.commission, p.commission)
		// }
		// if p.exchangeFee != tc.expPos.exchangeFee {
		// 	t.Errorf("%v exchangeFee: expected %#v, actual %#v", tc.msg, tc.expPos.exchangeFee, p.exchangeFee)
		// }
		// if p.cost != tc.expPos.cost {
		// 	t.Errorf("%v cost: expected %#v, actual %#v", tc.msg, tc.expPos.cost, p.cost)
		// }
		// if p.costBasis != tc.expPos.costBasis {
		// 	t.Errorf("%v costBasis: expected %#v, actual %#v", tc.msg, tc.expPos.costBasis, p.costBasis)
		// }
		// if p.netValue != tc.expPos.netValue {
		// 	t.Errorf("%v netValue: expected %#v, actual %#v", tc.msg, tc.expPos.netValue, p.netValue)
		// }
		// if p.netValueBOT != tc.expPos.netValueBOT {
		// 	t.Errorf("%v netValueBOT: expected %#v, actual %#v", tc.msg, tc.expPos.netValueBOT, p.netValueBOT)
		// }
		// if p.netValueSLD != tc.expPos.netValueSLD {
		// 	t.Errorf("%v netValueSLD: expected %#v, actual %#v", tc.msg, tc.expPos.netValueSLD, p.netValueSLD)
		// }
		// if p.realProfitLoss != tc.expPos.realProfitLoss {
		// 	t.Errorf("%v realProfitLoss: expected %#v, actual %#v", tc.msg, tc.expPos.realProfitLoss, p.realProfitLoss)
		// }
		// if p.unrealProfitLoss != tc.expPos.unrealProfitLoss {
		// 	t.Errorf("%v unrealProfitLoss: expected %#v, actual %#v", tc.msg, tc.expPos.unrealProfitLoss, p.unrealProfitLoss)
		// }
		// if p.totalProfitLoss != tc.expPos.totalProfitLoss {
		// 	t.Errorf("%v totalProfitLoss: expected %#v, actual %#v", tc.msg, tc.expPos.totalProfitLoss, p.totalProfitLoss)
		// }
		// Check complete position
		if !reflect.DeepEqual(p, tc.expPos) {
			t.Errorf("\n%v Update(%+v): \nexpected %p %#v, \nactual   %p %#v", tc.msg, tc.fill, tc.expPos, tc.expPos, p, p)
		}
	}
}

func TestMultipleUpdate(t *testing.T) {
	// set the example time string in format yyyy-mm-dd
	var exampleTime, _ = time.Parse("2006-01-02", "2017-06-01")

	var p = &position{
		timestamp: exampleTime, symbol: "TEST.DE",
		qty: 10, qtyBOT: 10, qtySLD: 0,
		avgPrice: 10, avgPriceNet: 10.5, avgPriceBOT: 10, avgPriceSLD: 0,
		value: -100, valueBOT: 100, valueSLD: 0,
		netValue: -105, netValueBOT: 105, netValueSLD: 0,
		marketPrice: 10, marketValue: 100,
		commission: 4, exchangeFee: 1, cost: 5, costBasis: 105,
		realProfitLoss: 0, unrealProfitLoss: -5, totalProfitLoss: -5,
	}

	// testCases is a table for testing updating a position
	var testCases = []struct {
		msg     string
		updates []FillEvent
		expPos  *position // expected Position
	}{
		{
			"1. multiple",
			[]FillEvent{
				&Fill{
					Event:    Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
					Exchange: "TEST", Direction: "BOT", // BOT for buy or SLD for sell
					Qty: 15, Price: 15,
					Commission: 6, ExchangeFee: 1, Cost: 7,
				},
				&Fill{
					Event:    Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
					Exchange: "TEST", Direction: "SLD", // BOT for buy or SLD for sell
					Qty: 18, Price: 20,
					Commission: 8, ExchangeFee: 1, Cost: 9,
				},
				&Fill{
					Event:    Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
					Exchange: "TEST", Direction: "BOT", // BOT for buy or SLD for sell
					Qty: 12, Price: 18,
					Commission: 7, ExchangeFee: 1, Cost: 8,
				},
			},
			&position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: 19, qtyBOT: 37, qtySLD: 18,
				avgPrice: 17.2374, avgPriceNet: 17.6842, avgPriceBOT: 14.6216, avgPriceSLD: 20,
				value: -181, valueBOT: 541, valueSLD: 360,
				netValue: -210, netValueBOT: 561, netValueSLD: 351,
				marketPrice: 18, marketValue: 342,
				commission: 25, exchangeFee: 4, cost: 29, costBasis: 318.36,
				realProfitLoss: 108.36, unrealProfitLoss: 23.64, totalProfitLoss: 132,
			},
		},
	}

	for _, tc := range testCases {
		for _, fill := range tc.updates {
			p.Update(fill)
		}
		// checking valus from multiple update
		// if p.qty != tc.expPos.qty {
		// 	t.Errorf("%v qty: expected %#v, actual %#v", tc.msg, tc.expPos.qty, p.qty)
		// }
		// if p.qtyBOT != tc.expPos.qtyBOT {
		// 	t.Errorf("%v qtyBOT: expected %#v, actual %#v", tc.msg, tc.expPos.qtyBOT, p.qtyBOT)
		// }
		// if p.qtySLD != tc.expPos.qtySLD {
		// 	t.Errorf("%v qtySLD: expected %#v, actual %#v", tc.msg, tc.expPos.qtySLD, p.qtySLD)
		// }
		// if p.avgPrice != tc.expPos.avgPrice {
		// 	t.Errorf("%v avgPrice: expected %#v, actual %#v", tc.msg, tc.expPos.avgPrice, p.avgPrice)
		// }
		// if p.avgPriceNet != tc.expPos.avgPriceNet {
		// 	t.Errorf("%v avgPriceNet: expected %#v, actual %#v", tc.msg, tc.expPos.avgPriceNet, p.avgPriceNet)
		// }
		// if p.avgPriceBOT != tc.expPos.avgPriceBOT {
		// 	t.Errorf("%v avgPriceBOT: expected %#v, actual %#v", tc.msg, tc.expPos.avgPriceBOT, p.avgPriceBOT)
		// }
		// if p.avgPriceSLD != tc.expPos.avgPriceSLD {
		// 	t.Errorf("%v avgPriceSLD: expected %#v, actual %#v", tc.msg, tc.expPos.avgPriceSLD, p.avgPriceSLD)
		// }
		// if p.value != tc.expPos.value {
		// 	t.Errorf("%v value: expected %#v, actual %#v", tc.msg, tc.expPos.value, p.value)
		// }
		// if p.valueBOT != tc.expPos.valueBOT {
		// 	t.Errorf("%v valueBOT: expected %#v, actual %#v", tc.msg, tc.expPos.valueBOT, p.valueBOT)
		// }
		// if p.valueSLD != tc.expPos.valueSLD {
		// 	t.Errorf("%v valueSLD: expected %#v, actual %#v", tc.msg, tc.expPos.valueSLD, p.valueSLD)
		// }
		// if p.marketPrice != tc.expPos.marketPrice {
		// 	t.Errorf("%v marketPrice: expected %#v, actual %#v", tc.msg, tc.expPos.marketPrice, p.marketPrice)
		// }
		// if p.marketValue != tc.expPos.marketValue {
		// 	t.Errorf("%v marketValue: expected %#v, actual %#v", tc.msg, tc.expPos.marketValue, p.marketValue)
		// }
		// if p.commission != tc.expPos.commission {
		// 	t.Errorf("%v commission: expected %#v, actual %#v", tc.msg, tc.expPos.commission, p.commission)
		// }
		// if p.exchangeFee != tc.expPos.exchangeFee {
		// 	t.Errorf("%v exchangeFee: expected %#v, actual %#v", tc.msg, tc.expPos.exchangeFee, p.exchangeFee)
		// }
		// if p.cost != tc.expPos.cost {
		// 	t.Errorf("%v cost: expected %#v, actual %#v", tc.msg, tc.expPos.cost, p.cost)
		// }
		// if p.costBasis != tc.expPos.costBasis {
		// 	t.Errorf("%v costBasis: expected %#v, actual %#v", tc.msg, tc.expPos.costBasis, p.costBasis)
		// }
		// if p.netValue != tc.expPos.netValue {
		// 	t.Errorf("%v netValue: expected %#v, actual %#v", tc.msg, tc.expPos.netValue, p.netValue)
		// }
		// if p.netValueBOT != tc.expPos.netValueBOT {
		// 	t.Errorf("%v netValueBOT: expected %#v, actual %#v", tc.msg, tc.expPos.netValueBOT, p.netValueBOT)
		// }
		// if p.netValueSLD != tc.expPos.netValueSLD {
		// 	t.Errorf("%v netValueSLD: expected %#v, actual %#v", tc.msg, tc.expPos.netValueSLD, p.netValueSLD)
		// }
		// if p.realProfitLoss != tc.expPos.realProfitLoss {
		// 	t.Errorf("%v realProfitLoss: expected %#v, actual %#v", tc.msg, tc.expPos.realProfitLoss, p.realProfitLoss)
		// }
		// if p.unrealProfitLoss != tc.expPos.unrealProfitLoss {
		// 	t.Errorf("%v unrealProfitLoss: expected %#v, actual %#v", tc.msg, tc.expPos.unrealProfitLoss, p.unrealProfitLoss)
		// }
		// if p.totalProfitLoss != tc.expPos.totalProfitLoss {
		// 	t.Errorf("%v totalProfitLoss: expected %#v, actual %#v", tc.msg, tc.expPos.totalProfitLoss, p.totalProfitLoss)
		// }
		if !reflect.DeepEqual(p, tc.expPos) {
			t.Errorf("\n%v Update(): \nexpected %p %#v, \nactual   %p %#v", tc.msg, tc.expPos, tc.expPos, p, p)
		}
	}
}

func TestUpdateValue(t *testing.T) {
	// set the example time string in format yyyy-mm-dd
	var exampleTime, _ = time.Parse("2006-01-02", "2017-06-01")
	// initialize new Position ready for use
	var p = &position{
		timestamp: exampleTime, symbol: "TEST.DE",
		qty: 10, qtyBOT: 10, qtySLD: 0,
		avgPrice: 10, avgPriceNet: 10.5, avgPriceBOT: 10, avgPriceSLD: 0,
		value: -100, valueBOT: 100, valueSLD: 0,
		netValue: -105, netValueBOT: 105, netValueSLD: 0,
		marketPrice: 10, marketValue: 100,
		commission: 4, exchangeFee: 1, cost: 5, costBasis: 105,
		realProfitLoss: 0, unrealProfitLoss: -5, totalProfitLoss: -5,
	}
	// testCases is a table for testing updating a position
	var testCases = []struct {
		data   DataEventHandler
		expPos *position // expected Position
	}{
		{
			data: &Bar{
				Event: Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Close: 99,
			},
			expPos: &position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: 10, qtyBOT: 10, qtySLD: 0,
				avgPrice: 10, avgPriceNet: 10.5, avgPriceBOT: 10, avgPriceSLD: 0,
				value: -100, valueBOT: 100, valueSLD: 0,
				netValue: -105, netValueBOT: 105, netValueSLD: 0,
				marketPrice: 99, marketValue: 990,
				commission: 4, exchangeFee: 1, cost: 5, costBasis: 105,
				realProfitLoss: 0, unrealProfitLoss: 885, totalProfitLoss: 885,
			},
		},
		{
			data: &Bar{
				Event: Event{Timestamp: exampleTime, Symbol: "TEST.DE"},
				Close: 45,
			},
			expPos: &position{
				timestamp: exampleTime, symbol: "TEST.DE",
				qty: 10, qtyBOT: 10, qtySLD: 0,
				avgPrice: 10, avgPriceNet: 10.5, avgPriceBOT: 10, avgPriceSLD: 0,
				value: -100, valueBOT: 100, valueSLD: 0,
				netValue: -105, netValueBOT: 105, netValueSLD: 0,
				marketPrice: 45, marketValue: 450,
				commission: 4, exchangeFee: 1, cost: 5, costBasis: 105,
				realProfitLoss: 0, unrealProfitLoss: 345, totalProfitLoss: 345,
			},
		},
	}

	for _, tc := range testCases {
		p.UpdateValue(tc.data)
		if !reflect.DeepEqual(p, tc.expPos) {
			t.Errorf("Create(%v): \nexpected %p %#v, \nactual   %p %#v", tc.data, tc.expPos, tc.expPos, p, p)
		}
	}
}
