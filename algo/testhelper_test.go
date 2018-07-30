package algo

import (
	"time"

	gbt "github.com/dirkolbrich/gobacktest"
)

func testHelperTimeMap(dates []string) map[string]time.Time {
	timeMap := make(map[string]time.Time)
	for _, d := range dates {
		time, _ := time.Parse("2006-01-02", d)
		timeMap[d] = time
	}
	return timeMap
}

func testHelperMockData(dates []string) []gbt.DataEvent {
	mockdata := []gbt.DataEvent{}
	for _, d := range dates {
		time, _ := time.Parse("2006-01-02", d)
		symbol := "Test"

		event := &gbt.Event{}
		event.SetSymbol(symbol)
		event.SetTime(time)

		bar := &gbt.Bar{
			Event: *event,
		}
		mockdata = append(mockdata, bar)
	}
	return mockdata
}
