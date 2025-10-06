package persistence

import "github.com/aarondl/opt/omit"

func ID(id int64) omit.Val[int64] {
	return omit.From(id)
}
