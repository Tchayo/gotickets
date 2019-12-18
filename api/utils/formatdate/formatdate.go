package formatdate

import (
	"fmt"
	"time"
)

// Layouts
// const (
// 	ANSIC       = "Mon Jan _2 15:04:05 2006"
// 	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
// 	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
// 	RFC822      = "02 Jan 06 15:04 MST"
// 	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
// 	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
// 	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
// 	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
// 	RFC3339     = "2006-01-02T15:04:05Z07:00"
// 	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
// 	Kitchen     = "3:04PM"
// 	Stamp       = "Jan _2 15:04:05"
// 	StampMilli  = "Jan _2 15:04:05.000"
// 	StampMicro  = "Jan _2 15:04:05.000000"
// 	StampNano   = "Jan _2 15:04:05.000000000"
// )

// FormatDate : format date output
func FormatDate(date time.Time, format string) string {

	defaultOutput := date.Format(time.ANSIC)

	switch format {
	case "UNIX":
		output := date.Format(time.UnixDate)
		fmt.Println("Unix Time : " + output)
		return output
	case "RFC":
		output := date.Format(time.RFC822)
		fmt.Println("RFC Time : " + output)
		return output
	case "RFCN":
		output := date.Format(time.RFC822Z)
		fmt.Println("RFCN Time : " + output)
		return output
	default:
		return defaultOutput

	}

}
