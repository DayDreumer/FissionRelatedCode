package tag

type Key string

func (k Key) Valid() bool {
	return len(k) != 0
}
