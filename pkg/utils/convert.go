package utils

func I2I(is []int64) []any {
	ai := make([]any, 0, len(is))
	for _, e := range is {
		ai = append(ai, e)
	}
	return ai
}
