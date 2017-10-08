package backtest

import (
	"reflect"
	"testing"
	"time"
)

func TestResetPortfolio(t *testing.T) {
	var testCases = []struct {
		msg          string
		portfolio    PortfolioHandler
		expPortfolio PortfolioHandler
	}{
		{"testing full portfolio",
			&Portfolio{
				initialCash: 100000,
				cash:        100000,
				holdings: map[string]position{
					"TEST.DE": {qty: 100},
					"BAS.DE":  {qty: 90},
				},
				transactions: []FillEvent{
					&Fill{Direction: "BOT", Qty: 100},
					&Fill{Direction: "BOT", Qty: 90},
				},
				sizeManager: &Size{},
				riskManager: &Risk{},
			},
			&Portfolio{
				initialCash: 100000,
				sizeManager: &Size{},
				riskManager: &Risk{},
			},
		},
		{"testing empty portfolio",
			&Portfolio{
				initialCash:  0,
				cash:         0,
				holdings:     map[string]position{},
				transactions: []FillEvent{},
				sizeManager:  &Size{},
				riskManager:  &Risk{},
			},
			&Portfolio{
				sizeManager: &Size{},
				riskManager: &Risk{},
			},
		},
	}

	for _, tc := range testCases {
		tc.portfolio.Reset()
		if !reflect.DeepEqual(tc.portfolio, tc.expPortfolio) {
			t.Errorf("%v Reset(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expPortfolio, tc.portfolio)
		}
	}
}

func TestOnSignal(t *testing.T) {
	var timestamp, _ = time.Parse("2006-01-02", "2017-09-29")
	var testCases = []struct {
		msg       string
		portfolio PortfolioHandler
		signal    SignalEvent
		data      DataHandler
		expOrder  OrderEvent
		expErr    error
	}{
		{"testing simple signal",
			&Portfolio{
				sizeManager: &Size{},
				riskManager: &Risk{},
			},
			&Signal{
				Event:     Event{Symbol: "TEST.DE", Timestamp: timestamp},
				Direction: "long",
			},
			&Data{
				latest: map[string]DataEventHandler{
					"TEST.DE": &Bar{Close: 100},
				},
			},
			&Order{
				Event:     Event{Symbol: "TEST.DE", Timestamp: timestamp},
				Direction: "long",
				OrderType: "MKT",
			},
			nil},
	}

	for _, tc := range testCases {
		order, err := tc.portfolio.OnSignal(tc.signal, tc.data)
		if !reflect.DeepEqual(order, tc.expOrder) || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("%v OnSignal(): \nexpected %#v %v, \nactual   %#v %v",
				tc.msg, tc.expOrder, tc.expErr, order, err)
		}
	}
}

func TestOnFill(t *testing.T) {
	var timestamp, _ = time.Parse("2006-01-02", "2017-09-29")
	var fillCases = map[string]FillEvent{
		"BOT": &Fill{
			Event:     Event{Symbol: "TEST.DE", Timestamp: timestamp},
			Direction: "BOT",
			Qty:       100,
			Price:     10,
			Cost:      10,
		},
		"SLD": &Fill{
			Event:     Event{Symbol: "TEST.DE", Timestamp: timestamp},
			Direction: "SLD",
			Qty:       100,
			Price:     10,
			Cost:      10,
		},
	}
	var holdingCases = map[string]position{
		"BOT": {
			timestamp:        timestamp,
			symbol:           "TEST.DE",
			qty:              100,
			qtyBOT:           100,
			qtySLD:           0,
			avgPrice:         10,
			avgPriceNet:      10.1,
			avgPriceBOT:      10,
			avgPriceSLD:      0,
			value:            -1000,
			valueBOT:         1000,
			valueSLD:         0,
			netValue:         -1010,
			netValueBOT:      1010,
			netValueSLD:      0,
			marketPrice:      10,
			marketValue:      1000,
			commission:       0,
			exchangeFee:      0,
			cost:             10,
			costBasis:        1010,
			realProfitLoss:   0,
			unrealProfitLoss: -10,
			totalProfitLoss:  -10,
		},
		"SLD": {
			timestamp:        timestamp,
			symbol:           "TEST.DE",
			qty:              -100,
			qtyBOT:           0,
			qtySLD:           100,
			avgPrice:         10,
			avgPriceNet:      9.9,
			avgPriceBOT:      0,
			avgPriceSLD:      10,
			value:            1000,
			valueBOT:         0,
			valueSLD:         1000,
			netValue:         990,
			netValueBOT:      0,
			netValueSLD:      990,
			marketPrice:      10,
			marketValue:      1000,
			commission:       0,
			exchangeFee:      0,
			cost:             10,
			costBasis:        -990,
			realProfitLoss:   0,
			unrealProfitLoss: -10,
			totalProfitLoss:  -10,
		},
	}
	var size = &Size{}
	var risk = &Risk{}

	var testCases = []struct {
		msg          string
		portfolio    PortfolioHandler
		fill         FillEvent
		data         DataHandler
		expPortfolio PortfolioHandler
	}{
		{"testing BOT fill with empty holdings and transactions",
			&Portfolio{
				cash:        10000,
				sizeManager: size,
				riskManager: risk,
			},
			fillCases["BOT"],
			&Data{},
			&Portfolio{
				cash:        8990,
				sizeManager: size,
				riskManager: risk,
				holdings: map[string]position{
					"TEST.DE": {
						timestamp:        timestamp,
						symbol:           "TEST.DE",
						qty:              100,
						qtyBOT:           100,
						qtySLD:           0,
						avgPrice:         10,
						avgPriceNet:      10.1,
						avgPriceBOT:      10,
						avgPriceSLD:      0,
						value:            -1000,
						valueBOT:         1000,
						valueSLD:         0,
						netValue:         -1010,
						netValueBOT:      1010,
						netValueSLD:      0,
						marketPrice:      10,
						marketValue:      1000,
						commission:       0,
						exchangeFee:      0,
						cost:             10,
						costBasis:        1010,
						realProfitLoss:   0,
						unrealProfitLoss: -10,
						totalProfitLoss:  -10,
					},
				},
				transactions: []FillEvent{
					fillCases["BOT"],
				},
			},
		},
		{"testing BOT fill with BOT holdings and transactions",
			&Portfolio{
				cash:        8990,
				sizeManager: size,
				riskManager: risk,
				holdings: map[string]position{
					"TEST.DE": holdingCases["BOT"],
				},
				transactions: []FillEvent{
					fillCases["BOT"],
				},
			},
			fillCases["BOT"],
			&Data{},
			&Portfolio{
				cash:        7980,
				sizeManager: size,
				riskManager: risk,
				holdings: map[string]position{
					"TEST.DE": {
						timestamp:        timestamp,
						symbol:           "TEST.DE",
						qty:              200,
						qtyBOT:           200,
						qtySLD:           0,
						avgPrice:         10,
						avgPriceNet:      10.1,
						avgPriceBOT:      10,
						avgPriceSLD:      0,
						value:            -2000,
						valueBOT:         2000,
						valueSLD:         0,
						netValue:         -2020,
						netValueBOT:      2020,
						netValueSLD:      0,
						marketPrice:      10,
						marketValue:      2000,
						commission:       0,
						exchangeFee:      0,
						cost:             20,
						costBasis:        2020,
						realProfitLoss:   0,
						unrealProfitLoss: -20,
						totalProfitLoss:  -20,
					},
				},
				transactions: []FillEvent{
					fillCases["BOT"],
					fillCases["BOT"],
				},
			},
		},
		{"testing SLD fill with BOT holdings and transactions, should set holding to zero",
			&Portfolio{
				cash:        8990,
				sizeManager: size,
				riskManager: risk,
				holdings: map[string]position{
					"TEST.DE": holdingCases["BOT"],
				},
				transactions: []FillEvent{
					fillCases["BOT"],
				},
			},
			fillCases["SLD"],
			&Data{},
			&Portfolio{
				cash:        9980,
				sizeManager: size,
				riskManager: risk,
				holdings: map[string]position{
					"TEST.DE": {
						timestamp:        timestamp,
						symbol:           "TEST.DE",
						qty:              0,
						qtyBOT:           100,
						qtySLD:           100,
						avgPrice:         10,
						avgPriceNet:      10,
						avgPriceBOT:      10,
						avgPriceSLD:      10,
						value:            0,
						valueBOT:         1000,
						valueSLD:         1000,
						netValue:         -20,
						netValueBOT:      1010,
						netValueSLD:      990,
						marketPrice:      10,
						marketValue:      0,
						commission:       0,
						exchangeFee:      0,
						cost:             20,
						costBasis:        0,
						realProfitLoss:   -20,
						unrealProfitLoss: 0,
						totalProfitLoss:  -20,
					},
				},
				transactions: []FillEvent{
					fillCases["BOT"],
					fillCases["SLD"],
				},
			},
		},
		{"testing SLD fill with empty holdings and transactions",
			&Portfolio{
				cash:        10000,
				sizeManager: size,
				riskManager: risk,
			},
			fillCases["SLD"],
			&Data{},
			&Portfolio{
				cash:        10990,
				sizeManager: size,
				riskManager: risk,
				holdings: map[string]position{
					"TEST.DE": {
						timestamp:        timestamp,
						symbol:           "TEST.DE",
						qty:              -100,
						qtyBOT:           0,
						qtySLD:           100,
						avgPrice:         10,
						avgPriceNet:      9.9,
						avgPriceBOT:      0,
						avgPriceSLD:      10,
						value:            1000,
						valueBOT:         0,
						valueSLD:         1000,
						netValue:         990,
						netValueBOT:      0,
						netValueSLD:      990,
						marketPrice:      10,
						marketValue:      1000,
						commission:       0,
						exchangeFee:      0,
						cost:             10,
						costBasis:        -990,
						realProfitLoss:   0,
						unrealProfitLoss: -10,
						totalProfitLoss:  -10,
					},
				},
				transactions: []FillEvent{
					fillCases["SLD"],
				},
			},
		},
		{"testing SLD fill with SLD holdings and transactions",
			&Portfolio{
				cash:        10990,
				sizeManager: size,
				riskManager: risk,
				holdings: map[string]position{
					"TEST.DE": holdingCases["SLD"],
				},
				transactions: []FillEvent{
					fillCases["SLD"],
				},
			},
			fillCases["SLD"],
			&Data{},
			&Portfolio{
				cash:        11980,
				sizeManager: size,
				riskManager: risk,
				holdings: map[string]position{
					"TEST.DE": {
						timestamp:        timestamp,
						symbol:           "TEST.DE",
						qty:              -200,
						qtyBOT:           0,
						qtySLD:           200,
						avgPrice:         10,
						avgPriceNet:      9.9,
						avgPriceBOT:      0,
						avgPriceSLD:      10,
						value:            2000,
						valueBOT:         0,
						valueSLD:         2000,
						netValue:         1980,
						netValueBOT:      0,
						netValueSLD:      1980,
						marketPrice:      10,
						marketValue:      2000,
						commission:       0,
						exchangeFee:      0,
						cost:             20,
						costBasis:        -1980,
						realProfitLoss:   0,
						unrealProfitLoss: -20,
						totalProfitLoss:  -20,
					},
				},
				transactions: []FillEvent{
					fillCases["SLD"],
					fillCases["SLD"],
				},
			},
		},
		{"testing BOT fill with SLD holdings and transactions, should set holding to zero",
			&Portfolio{
				cash:        10990,
				sizeManager: size,
				riskManager: risk,
				holdings: map[string]position{
					"TEST.DE": holdingCases["SLD"],
				},
				transactions: []FillEvent{
					fillCases["SLD"],
				},
			},
			fillCases["BOT"],
			&Data{},
			&Portfolio{
				cash:        9980,
				sizeManager: size,
				riskManager: risk,
				holdings: map[string]position{
					"TEST.DE": {
						timestamp:        timestamp,
						symbol:           "TEST.DE",
						qty:              0,
						qtyBOT:           100,
						qtySLD:           100,
						avgPrice:         10,
						avgPriceNet:      10,
						avgPriceBOT:      10,
						avgPriceSLD:      10,
						value:            0,
						valueBOT:         1000,
						valueSLD:         1000,
						netValue:         -20,
						netValueBOT:      1010,
						netValueSLD:      990,
						marketPrice:      10,
						marketValue:      0,
						commission:       0,
						exchangeFee:      0,
						cost:             20,
						costBasis:        0,
						realProfitLoss:   -20,
						unrealProfitLoss: 0,
						totalProfitLoss:  -20,
					},
				},
				transactions: []FillEvent{
					fillCases["SLD"],
					fillCases["BOT"],
				},
			},
		},
	}

	for _, tc := range testCases {
		fill, _ := tc.portfolio.OnFill(tc.fill, tc.data)
		if !reflect.DeepEqual(tc.fill, fill) {
		}
		if !reflect.DeepEqual(tc.portfolio, tc.expPortfolio) {
			t.Errorf("%v OnFill(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expPortfolio, tc.portfolio)
		}
	}
}

func TestIsInvested(t *testing.T) {
	// testCases is a table for testing
	var testCases = []struct {
		msg       string    // error message
		symbol    string    // input string
		portfolio Portfolio // input portfolio
		expPos    position  // expected position return
		expOk     bool      // expected bool return
	}{
		{"Portfolio is empty:",
			"TEST.DE",
			Portfolio{},
			position{},
			false,
		},
		{"Portfolio should be invested with long:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 10},
				},
			},
			position{},
			true,
		},
		{"Portfolio should be invested with short:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: -10},
				},
			},
			position{},
			true,
		},
		{"Portfolio should not be invested:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 0},
				},
			},
			position{},
			false,
		},
	}

	for _, tc := range testCases {
		pos, ok := tc.portfolio.IsInvested(tc.symbol)
		if (pos != tc.expPos) && (ok != tc.expOk) {
			t.Errorf("%v\nIsInvested(%v): \nexpected %#v %v, \nactual %#v %v", tc.msg, tc.symbol, tc.expOk, tc.expPos, pos, ok)
		}
	}
}

func TestIsLong(t *testing.T) {
	// testCases is a table for testing
	var testCases = []struct {
		msg       string    // error message
		symbol    string    // input string
		portfolio Portfolio // input portfolio
		expPos    position  // expected position return
		expOk     bool      // expected bool return
	}{
		{"Portfolio is empty:",
			"TEST.DE",
			Portfolio{},
			position{},
			false,
		},
		{"Portfolio should be long:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 10},
				},
			},
			position{},
			true,
		},
		{"Portfolio is short:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: -10},
				},
			},
			position{},
			false,
		},
		{"Portfolio is not invested:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 0},
				},
			},
			position{},
			false,
		},
	}

	for _, tc := range testCases {
		pos, ok := tc.portfolio.IsLong(tc.symbol)
		if (pos != tc.expPos) && (ok != tc.expOk) {
			t.Errorf("%v\nIsLong(%v): \nexpected %#v %v, \nactual %#v %v", tc.msg, tc.symbol, tc.expOk, tc.expPos, pos, ok)
		}
	}
}

func TestIsShort(t *testing.T) {
	// testCases is a table for testing
	var testCases = []struct {
		msg       string    // error message
		symbol    string    // input string
		portfolio Portfolio // input portfolio
		expPos    position  // expected position return
		expOk     bool      // expected bool return
	}{
		{"Portfolio is empty:",
			"TEST.DE",
			Portfolio{},
			position{},
			false,
		},
		{"Portfolio is long:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 10},
				},
			},
			position{},
			false,
		},
		{"Portfolio should be short:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: -10},
				},
			},
			position{},
			true,
		},
		{"Portfolio is not invested:",
			"TEST.DE",
			Portfolio{
				holdings: map[string]position{
					"TEST.DE": {qty: 0},
				},
			},
			position{},
			false,
		},
	}

	for _, tc := range testCases {
		pos, ok := tc.portfolio.IsShort(tc.symbol)
		if (pos != tc.expPos) && (ok != tc.expOk) {
			t.Errorf("%v\nIsShort(%v): \nexpected %#v %v, \nactual %#v %v", tc.msg, tc.symbol, tc.expOk, tc.expPos, pos, ok)
		}
	}
}

func TestPortfolioValue(t *testing.T) {
	var testCases = []struct {
		msg       string
		portfolio PortfolioHandler
		expValue  float64
	}{
		{"testing value of positiv holdings",
			&Portfolio{
				cash: 10000,
				holdings: map[string]position{
					"TEST.DE": {qty: 100, marketValue: 200},
					"BAS.DE":  {qty: 100, marketValue: 300},
					"APPL":    {qty: 100, marketValue: 500},
				},
			},
			11000,
		},
		{"testing value of negativ holdings",
			&Portfolio{
				cash: 10000,
				holdings: map[string]position{
					"TEST.DE": {qty: 100, marketValue: -200},
					"BAS.DE":  {qty: 100, marketValue: -300},
					"APPL":    {qty: 100, marketValue: -500},
				},
			},
			9000,
		},
		{"testing value of mixed holdings",
			&Portfolio{
				cash: 10000,
				holdings: map[string]position{
					"TEST.DE": {qty: 100, marketValue: 200},
					"BAS.DE":  {qty: 100, marketValue: -300},
					"APPL":    {qty: 100, marketValue: 500},
				},
			},
			10400,
		},
		{"testing value of empty holdings",
			&Portfolio{},
			0,
		},
	}

	for _, tc := range testCases {
		value := tc.portfolio.Value()
		if value != tc.expValue {
			t.Errorf("%v Value(): \nexpected %#v, \nactual %#v",
				tc.msg, tc.expValue, value)
		}
	}
}
