package backtest

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// StatisticHandler is a basic statistic interface
type StatisticHandler interface {
	EventTracker
	TransactionTracker
	StatisticPrinter
	Reseter
	StatisticUpdater
	Resulter
}

// EventTracker is responsible for all event tracking during a backtest
type EventTracker interface {
	TrackEvent(EventHandler)
	Events() []EventHandler
}

// TransactionTracker is responsible for all transaction tracking during a backtest
type TransactionTracker interface {
	TrackTransaction(FillEvent)
	Transactions() []FillEvent
}

// StatisticPrinter handles printing of the statistics to screen
type StatisticPrinter interface {
	PrintResult()
}

// StatisticUpdater handles the updateing of the statistics
type StatisticUpdater interface {
	Update(DataEventHandler, PortfolioHandler)
}

// Resulter bundles all methods which return the results of the backtest
type Resulter interface {
	TotalEquityReturn() (float64, error)
}

// Statistic is a basic test statistic, which holds simple lists of historic events
type Statistic struct {
	eventHistory       []EventHandler
	transactionHistory []FillEvent
	equity             []equityPoint
}

type equityPoint struct {
	timestamp    time.Time
	equity       float64
	equityReturn float64
}

// Update updates the complete statistics to a given data event
func (s *Statistic) Update(d DataEventHandler, p PortfolioHandler) {
	e := equityPoint{}
	e.timestamp = d.GetTime()
	e.equity = p.Value()

	if len(s.equity) > 0 {
		lastEquity := decimal.NewFromFloat(s.equity[len(s.equity)-1].equity)
		equity := decimal.NewFromFloat(e.equity)
		equityReturn := equity.Sub(lastEquity).Div(lastEquity)
		e.equityReturn, _ = equityReturn.Round(4).Float64()
	}
	// append new quity point
	s.equity = append(s.equity, e)
}

// TrackEvent tracks an event
func (s *Statistic) TrackEvent(e EventHandler) {
	s.eventHistory = append(s.eventHistory, e)
}

// Events returns the complete events history
func (s Statistic) Events() []EventHandler {
	return s.eventHistory
}

// TrackTransaction tracks a transaction aka a fill event
func (s *Statistic) TrackTransaction(f FillEvent) {
	s.transactionHistory = append(s.transactionHistory, f)
}

// Transactions returns the complete events history
func (s Statistic) Transactions() []FillEvent {
	return s.transactionHistory
}

// Reset the statistic to a clean state
func (s *Statistic) Reset() {
	s.eventHistory = nil
	s.transactionHistory = nil
	s.equity = nil
}

// PrintResult prints the backtest statistics to the screen
func (s Statistic) PrintResult() {
	fmt.Println("Printing backtest results:")
	fmt.Printf("Counted %d total events.\n", len(s.Events()))

	fmt.Printf("Counted %d total transactions:\n", len(s.Transactions()))
	for k, v := range s.Transactions() {
		fmt.Printf("%d. Transaction: %v Action: %s Price: %f Qty: %d\n", k+1, v.GetTime().Format("2006-01-02"), v.GetDirection(), v.GetPrice(), v.GetQty())
	}
}

// TotalEquityReturn calculates the the total return on the first and last equity point
func (s Statistic) TotalEquityReturn() (r float64, err error) {
	firstEquityPoint, ok := s.firstEquityPoint()
	if !ok {
		return r, errors.New("could not calculate totalEquityReturn, no equity points found")
	}
	firstEquity := decimal.NewFromFloat(firstEquityPoint.equity)

	lastEquityPoint, _ := s.lastEquityPoint()
	// if !ok {
	// 	return r, errors.New("could not calculate totalEquityReturn, no last equity point")
	// }
	lastEquity := decimal.NewFromFloat(lastEquityPoint.equity)

	totalEquityReturn := lastEquity.Sub(firstEquity).Div(firstEquity)
	total, _ := totalEquityReturn.Round(4).Float64()
	return total, nil
}

// returns the first equityPoint
func (s Statistic) firstEquityPoint() (ep equityPoint, ok bool) {
	if len(s.equity) <= 0 {
		return ep, false
	}
	ep = s.equity[0]

	return ep, true
}

// returns the last equityPoint
func (s Statistic) lastEquityPoint() (ep equityPoint, ok bool) {
	if len(s.equity) <= 0 {
		return ep, false
	}
	ep = s.equity[len(s.equity)-1]

	return ep, true
}
