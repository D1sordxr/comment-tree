package errorz

import "errors"

func In(target error, values ...error) bool {
	for _, value := range values {
		if errors.Is(value, target) {
			return true
		}
	}
	return false
}
