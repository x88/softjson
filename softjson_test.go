package softjson

import (
	"testing"
)

var (
	jsonData = []byte(`{"name" : "John", "age" : "23", "created" : "1511017709", "last_order" : "52.231", "total" : "8439282.149393"}`)
)

type bioEntry struct {
	Name      string  `json:"name"`
	Age       int     `json:"age"`
	Created   int64   `json:"created"`
	LastOrder float32 `json:"last_order"`
	Total     float64 `json:"total"`
}

func TestUnmarshalSoft(t *testing.T) {
	t1 := bioEntry{
		Name:      "John",
		Age:       23,
		Created:   1511017709,
		LastOrder: 52.231,
		Total:     8439282.149393,
	}

	t2 := bioEntry{}

	err := UnmarshalSoft(jsonData, &t2)

	if err != nil {
		t.Fatal(err.Error())
	}

	if t1.Name != t2.Name {
		t.Fatalf(`Value isn't unmarshalled correctly: string`)
	}

	if t1.Age != t2.Age {
		t.Fatalf(`Value isn't unmarshalled correctly: int`)
	}

	if t1.Created != t2.Created {
		t.Fatalf(`Value isn't unmarshalled correctly: int64`)
	}

	if t1.LastOrder != t2.LastOrder {
		t.Fatalf(`Value isn't unmarshalled correctly: float32`)
	}

	if t1.Total != t2.Total {
		t.Fatalf(`Value isn't unmarshalled correctly: float64`)
	}

}
