package main

import "fmt"

//time values in DIGITAL SECONDS. One digital day = 100000 seconds, which is 14% longer than a regular day.
//there are 1000 days in a digital year (a "cycle"), which is 3.17 times longer than a solar year.
//dates are measured relative to the epoch of the DIGITAL AGE, which is chosen semi-arbitrarily as April 5, 2063
const (
	MINUTE int = 100
	HOUR   int = 10000
	DAY    int = 100000
	CYCLE  int = 100000000
)

func GetTimeString(t int) string {
	return fmt.Sprintf("%.1d", (t/HOUR)%10) + "h:" + fmt.Sprintf("%.2d", (t/MINUTE)%100) + "m:" + fmt.Sprintf("%.2d", t%100) + "s"
}

func GetDurationString(t int) string {
	return fmt.Sprintf("%.1d", t/HOUR) + "h:" + fmt.Sprintf("%.2d", (t/MINUTE)%100) + "m:" + fmt.Sprintf("%.2d", t%100) + "s"
}

func GetDateString(t int) string {
	return "Day " + fmt.Sprintf("%.d", (t/DAY)%1000) + ", " + "Cycle " + fmt.Sprintf("%.d", t/CYCLE) + "DE"
}

//TODO: duration functions, etc.
