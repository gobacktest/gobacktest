package algo

import (
	"time"

	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

// RunOnce which returns true once and then returns false.
type RunOnce struct {
	bt.Algo
	hasRun bool
}

// Run runs the RunOnce() algo.
func (ro *RunOnce) Run(s bt.StrategyHandler) (bool, error) {
	if ro.hasRun {
		return false, nil
	}

	ro.hasRun = true
	return true, nil
}

// PeriodRunner defines how the function to compare two dates.
type PeriodRunner interface {
	CompareDates(time.Time, time.Time) (bool, error)
}

// RunPeriod is a base algo for time structures.
type RunPeriod struct {
	PeriodRunner
	bt.Algo
	runOnFirstDate bool
	runOnLastDate  bool
	runEndOfPeriod bool
}

// OnFirstDate sets if RunPeriod runs on the th first date it encounters.OnFirstDate
// Default is true. If set to false, it will ignore intermediate dates of a period and runs,
// if a period change happend, e.g. change of month
func (rp *RunPeriod) OnFirstDate(b bool) *RunPeriod {
	rp.runOnFirstDate = b
	return rp
}

// OnEndOfPeriod sets if RunPeriod runs on the end of e.g. a week or a month
func (rp *RunPeriod) OnEndOfPeriod(b bool) *RunPeriod {
	rp.runEndOfPeriod = b
	return rp
}

// Run runs the algo. It retrieves the current date and the date before.
// Then calls the specific implementations to compare these two dates.
func (rp *RunPeriod) Run(s bt.StrategyHandler) (bool, error) {
	now := rp.getNow(s)

	toCompare, ok := rp.getDateToCompare(s)
	if !ok {
		return false, nil
	}

	// call the specific implementation of each algo
	return rp.CompareDates(now, toCompare)
}

func (rp *RunPeriod) getNow(s bt.StrategyHandler) time.Time {
	// get current data event date
	event, _ := s.Event()
	now := event.GetTime()

	return now
}

func (rp *RunPeriod) getDateToCompare(s bt.StrategyHandler) (time.Time, bool) {
	data, ok := s.Data()
	// no history yet, so nothing to compare
	if !ok {
		return time.Time{}, false
	}

	history := data.History()
	// no more elemtns in history, so nothing to compare
	if len(history) <= 1 {
		return time.Time{}, false
	}
	dateToCompare := history[len(history)-2].GetTime()

	return dateToCompare, true
}

// RunDaily returns true on month change.
type RunDaily struct {
	RunPeriod
}

// NewRunDaily return a RunDaily Algo ready to use.
func NewRunDaily() *RunDaily {
	runDaily := RunDaily{RunPeriod{
		runOnFirstDate: true,
	}}
	runDaily.RunPeriod.PeriodRunner = runDaily
	return &runDaily
}

// CompareDates compares two t+dates if the day is different.
func (rd RunDaily) CompareDates(now, toCompare time.Time) (bool, error) {
	if now.Day() == toCompare.Day() {
		return false, nil
	}

	return true, nil
}

// RunWeekly return true on month change.
type RunWeekly struct {
	RunPeriod
}

// CompareDates compares two dates if in same week.
func (rw RunWeekly) CompareDates(now, toCompare time.Time) (bool, error) {
	nowYear, nowWeek := now.ISOWeek()
	compareYear, compareWeek := toCompare.ISOWeek()
	if (nowWeek == compareWeek) && (nowYear == compareYear) {
		return false, nil
	}

	return true, nil
}

// RunMonthly return true on month change.
type RunMonthly struct {
	RunPeriod
}

// NewRunMonthly return a RunMonthly Algo ready to use.
func NewRunMonthly() *RunMonthly {
	runMonthly := RunMonthly{RunPeriod{
		runOnFirstDate: true,
	}}
	runMonthly.RunPeriod.PeriodRunner = runMonthly
	return &runMonthly
}

// CompareDates compares two dates if in same month.
func (rm RunMonthly) CompareDates(now, toCompare time.Time) (bool, error) {
	if (now.Month() == toCompare.Month()) && (now.Year() == toCompare.Year()) {
		return false, nil
	}

	return true, nil
}

// RunQuarterly return true on quarter change.
// March-31 to April-01, June-30 to July-01,
// September-31 to October-01, December-31 to January-01
type RunQuarterly struct {
	RunPeriod
}

// NewRunQuarterly return a RunMonthly Algo ready to use.
func NewRunQuarterly() *RunQuarterly {
	runQuarterly := RunQuarterly{RunPeriod{
		runOnFirstDate: true,
	}}
	runQuarterly.RunPeriod.PeriodRunner = runQuarterly
	return &runQuarterly
}

// CompareDates compares two dates if in same quarter.
func (rq RunQuarterly) CompareDates(now, toCompare time.Time) (bool, error) {
	if (now.Month() == toCompare.Month()) && (now.Year() == toCompare.Year()) {
		return false, nil
	}

	return true, nil
}

// RunYearly return true on year change.
type RunYearly struct {
	RunPeriod
}

// NewRunYearly return a RunMonthly Algo ready to use.
func NewRunYearly() *RunYearly {
	runYearly := RunYearly{RunPeriod{
		runOnFirstDate: true,
	}}
	runYearly.RunPeriod.PeriodRunner = runYearly
	return &runYearly
}

// CompareDates compares two dates if in same year.
func (ry RunYearly) CompareDates(now, toCompare time.Time) (bool, error) {
	if now.Year() == toCompare.Year() {
		return false, nil
	}

	return true, nil
}
