package toggl

import "time"

func ToDuration(i int) time.Duration {
	return time.Duration(i) * time.Millisecond
}
