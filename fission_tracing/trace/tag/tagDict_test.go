package tag_test

import (
	"encoding/json"
	"testing"

	"fission.tracing/tag"
)

func TestDict(t *testing.T) {
	dict := tag.NewTagDict()
	keyArray := []tag.Key{
		"isEnd",
		"childs",
	}
	ValueArray := []tag.Value{
		tag.GetBoolValue(true),
		tag.GetIntValue(3),
	}
	for i := 0; i < 2; i++ {
		dict.Insert(keyArray[i], ValueArray[i])
	}
	if v, ok := dict.Search("isEnd"); ok {
		t.Logf("isEnd tag found, value is %s", v.String())
	} else {
		t.Fatalf("isEnd not found.")
	}
	if v, ok := dict.Search("childs"); ok {
		t.Logf("childs tag found, value is %s", v.String())
	} else {
		t.Fatalf("childs not found.")
	}
}

func TestDictSerialization(t *testing.T) {
	dict := tag.NewTagDict()
	keyArray := []tag.Key{
		"isEnd",
		"childs",
	}
	ValueArray := []tag.Value{
		tag.GetBoolValue(true),
		tag.GetIntValue(3),
	}
	for i := 0; i < 2; i++ {
		dict.Insert(keyArray[i], ValueArray[i])
	}
	info, err := json.Marshal(dict)
	if err != nil {
		t.Fatalf("Dict Marshal failed.")
	} else {
		t.Logf("dict info: %s", string(info))
	}

	recvDict := tag.NewTagDict()
	json.Unmarshal(info, &recvDict)
	if v, ok := dict.Search("isEnd"); ok {
		t.Logf("isEnd tag exists, value is %s", v.String())
	} else {
		t.Fatalf("isEnd not found, deserialization failed.")
	}
	if v, ok := dict.Search("childs"); ok {
		t.Logf("childs tag exists, value is %s", v.String())
	} else {
		t.Fatalf("childs not found, deserialization failed.")
	}
}
