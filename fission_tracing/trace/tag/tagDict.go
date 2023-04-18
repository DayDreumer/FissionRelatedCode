package tag

type TagPair struct {
	key   Key
	value Value
}

func (tp TagPair) Valid() bool {
	return tp.key.Valid() && tp.value.Type() != NONE
}

type TagDict struct {
	tagMap map[Key]Value
}

type TagDictInfo struct {
	Key   string
	Value string
}

func NewTagDict() *TagDict {
	td := &TagDict{
		tagMap: map[Key]Value{},
	}
	return td
}

func (tg *TagDict) Insert(k Key, v Value) bool {
	tg.tagMap[k] = v
	return true
}

func (tg *TagDict) Search(k Key) (Value, bool) {
	if v, ok := tg.tagMap[k]; ok {
		return v, true
	}
	return Value{}, false
}

// func (tg TagDict) MarshalJSON() ([]byte, bool) {
// 	// keyArray := make([]string, 0)
// 	// valueArray := make([]string, 0)
// 	t := TagDictInfo{
// 		Key:   "1",
// 		Value: "2",
// 	}
// 	r := json.Marshal{t}
// 	return r
// }
