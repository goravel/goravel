package mathutil

import "fmt"

// DataSize format bytes number friendly. eg: 1024 => 1KB, 1024*1024 => 1MB
//
// Usage:
//
//	file, err := os.Open(path)
//	fl, err := file.Stat()
//	fmtSize := DataSize(fl.Size())
func DataSize(size uint64) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%dB", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.2fK", float64(size)/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.2fM", float64(size)/1024/1024)
	default:
		return fmt.Sprintf("%.2fG", float64(size)/1024/1024/1024)
	}
}

var timeFormats = [][]int{
	{0},
	{1},
	{2, 1},
	{60},
	{120, 60},
	{3600},
	{7200, 3600},
	{86400},
	{172800, 86400}, // second elem is unit.
	{2592000},
	{2592000 * 2, 2592000},
}

var timeMessages = []string{
	"< 1 sec", "1 sec", "secs", "1 min", "mins", "1 hr", "hrs", "1 day", "days", "1 month", "months",
}

// HowLongAgo format a seconds, get how lang ago. eg: 1 day, 1 week
func HowLongAgo(sec int64) string {
	intVal := int(sec)
	length := len(timeFormats)

	for i, item := range timeFormats {
		if intVal >= item[0] {
			ni := i + 1
			match := false

			if ni < length { // next exists
				next := timeFormats[ni]
				if intVal < next[0] { // current <= intVal < next
					match = true
				}
			} else if ni == length { // current is last
				match = true
			}

			if match { // match success
				if len(item) == 1 {
					return timeMessages[i]
				}
				return fmt.Sprintf("%d %s", intVal/item[1], timeMessages[i])
			}
		}
	}

	return "unknown" // He should never happen
}
