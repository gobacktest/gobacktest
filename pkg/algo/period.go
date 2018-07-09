package algo

import (
	"time"

	bt "github.com/dirkolbrich/gobacktest/pkg/backtest"
)

// runOnce which returns true once and then returns false.
type runOnce struct {
	bt.Algo
	hasRun bool
}

// RunOnce returns a runOnce algo ready to use.
func RunOnce() bt.AlgoHandler {
	return &runOnce{}
}

// Run runs the RunOnce() algo.
func (ro *runOnce) Run(s bt.StrategyHandler) (bool, error) {
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
type runPeriod struct {
	PeriodRunner
	bt.Algo
	runOnFirstDate bool
	runOnLastDate  bool
	runEndOfPeriod bool
}

// // OnFirstDate sets if RunPeriod runs on the first date it encounters.OnFirstDate
// // Default is true. If set to false, it will ignore intermediate dates of a period and runs,
// // if a period change happend, e.g. change of month
// func (rp *runPeriod) OnFirstDate(b bool) bt.AlgoHandler {
// 	rp.runOnFirstDate = b
// 	return rp
// }

// func (rp *runPeriod) OnLastDate(b bool) bt.AlgoHandler {
// 	rp.runOnLastDate = b
// 	return rp
// }

// // OnEndOfPeriod sets if RunPeriod runs on the end of e.g. a week or a month
// func (rp *runPeriod) OnEndOfPeriod(b bool) bt.AlgoHandler {
// 	rp.runEndOfPeriod = b
// 	return rp
// }

// Run runs the algo. It retrieves the current date and the date before.
// Then calls the specific implementations to compare these two dates.
func (rp *runPeriod) Run(s bt.StrategyHandler) (bool, error) {
	now := rp.getNow(s)

	toCompare, ok := rp.getDateToCompare(s)
	if !ok {
		return false, nil
	}

	// call the specific implementation of each algo
	return rp.CompareDates(now, toCompare)
}

func (rp *runPeriod) getNow(s bt.StrategyHandler) time.Time {
	// get current data event date
	event, _ := s.Event()
	now := event.GetTime()

	return now
}

func (rp *runPeriod) getDateToCompare(s bt.StrategyHandler) (time.Time, bool) {
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

func runPeriodWithOptions(opt ...string) runPeriod {
	rp := runPeriod{}

	if len(opt) == 0 {
		rp.runOnFirstDate = true
		return rp
	}

	for _, option := range opt {
		switch option {
		case "onFirstDate":
			rp.runOnFirstDate = true
		case "onLastDate":
			rp.runOnLastDate = true
		case "endOfPeriod":
			rp.runEndOfPeriod = true
		}
	}

	return rp
}

// runDaily returns true on month change.
type runDaily struct {
	runPeriod
}

// RunDaily returns a RunDaily Algo ready to use.
func RunDaily(opt ...string) bt.AlgoHandler {
	rp := runPeriodWithOptions(opt...)

	runDaily := runDaily{runPeriod: rp}
	runDaily.runPeriod.PeriodRunner = runDaily
	return &runDaily
}

// CompareDates compares two dates if the day is different.
func (rd runDaily) CompareDates(now, toCompare time.Time) (bool, error) {
	if now.Day() == toCompare.Day() {
		return false, nil
	}

	return true, nil
}

// RunWeekly return true on month change.
type runWeekly struct {
	runPeriod
}

// RunWeekly return a RunDaily Algo ready to use.
func RunWeekly(opt ...string) bt.AlgoHandler {
	rp := runPeriodWithOptions(opt...)

	runWeekly := runWeekly{runPeriod: rp}
	runWeekly.runPeriod.PeriodRunner = runWeekly
	return &runWeekly
}

// CompareDates compares two dates if in same week.
func (rw runWeekly) CompareDates(now, toCompare time.Time) (bool, error) {
	nowYear, nowWeek := now.ISOWeek()
	compareYear, compareWeek := toCompare.ISOWeek()
	if (nowWeek == compareWeek) && (nowYear == compareYear) {
		return false, nil
	}

	return true, nil
}

// RunMonthly return true on month change.
type runMonthly struct {
	runPeriod
}

// RunMonthly return a RunMonthly Algo ready to use.
func RunMonthly(opt ...string) bt.AlgoHandler {
	rp := runPeriodWithOptions(opt...)

	runMonthly := runMonthly{runPeriod: rp}
	runMonthly.runPeriod.PeriodRunner = runMonthly
	return &runMonthly
}

// CompareDates compares two dates if in same month.
func (rm runMonthly) CompareDates(now, toCompare time.Time) (bool, error) {
	if (now.Month() == toCompare.Month()) && (now.Year() == toCompare.Year()) {
		return false, nil
	}

	return true, nil
}

// RunQuarterly return true on quarter change.
// March-31 to April-01, June-30 to July-01,
// September-31 to October-01, December-31 to January-01
type runQuarterly struct {
	runPeriod
}

// RunQuarterly return a RunMonthly Algo ready to use.
func RunQuarterly(opt ...string) bt.AlgoHandler {
	rp := runPeriodWithOptions(opt...)

	runQuarterly := runQuarterly{runPeriod: rp}
	runQuarterly.runPeriod.PeriodRunner = runQuarterly
	return &runQuarterly
}

// CompareDates compares two dates if in same quarter.
func (rq runQuarterly) CompareDates(now, toCompare time.Time) (bool, error) {
	if (now.Month() == toCompare.Month()) && (now.Year() == toCompare.Year()) {
		return false, nil
	}

	return true, nil
}

// RunYearly return true on year change.
type runYearly struct {
	runPeriod
}

// RunYearly return a RunMonthly Algo ready to use.
func RunYearly(opt ...string) bt.AlgoHandler {
	rp := runPeriodWithOptions(opt...)

	runYearly := runMonthly{runPeriod: rp}
	runYearly.runPeriod.PeriodRunner = runYearly
	return &runYearly
}

// CompareDates compares two dates if in same year.
func (ry runYearly) CompareDates(now, toCompare time.Time) (bool, error) {
	if now.Year() == toCompare.Year() {
		return false, nil
	}

	return true, nil
}
