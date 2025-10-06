package vo

type Limit int

func (l Limit) Int() int {
	return int(l)
}

const DefaultLimit Limit = 10
