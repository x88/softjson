Simply golang tool for unmarshal json with converting to target structure field data type.

ALPHA RAW version! Does not support pointers

```go
package main

import (
	"github.com/x88/softjson"
	"fmt"
	"log"
	"time"
)

type bioEntry struct {
	Name      string  `json:"name"`
	Age       int     `json:"age"`
	Created   int64   `json:"created"`
	LastOrder float32 `json:"last_order"`
	Total     float64 `json:"total"`
}

func main() {
	jsonData := []byte(`{"name" : "John", "age" : "23", "created" : "1511017709", "last_order" : "52.231", "total" : "8439282.149393"}`)

	entry := bioEntry{}

	err := softjson.UnmarshalSoft(jsonData, &entry)
	if err != nil {
		fmt.Println(`Wrong json`, err)
	}
	log.Printf("Name: %s\nAge: %d\nSigned: %s\nLast order: %f$\nTotal: %f$",
		entry.Name,
		entry.Age,
		time.Unix(entry.Created, 0).Format(time.UnixDate),
		entry.LastOrder,
		entry.Total)
}

// 2017/11/18 18:50:33 Name: John
// Age: 23
// Signed: Sat Nov 18 18:08:29 MSK 2017
// Last order: 52.230999$
// Total: 8439282.149393$

```