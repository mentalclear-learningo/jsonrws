package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	exampleOne()
	exampleTwo()
	customParsing()
}

func exampleOne() {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	toFile := Person{
		Name: "Fred",
		Age:  40,
	}

	// Write it out!
	tmpFile, err := ioutil.TempFile(os.TempDir(), "sample-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())
	err = json.NewEncoder(tmpFile).Encode(toFile)
	if err != nil {
		panic(err)
	}
	err = tmpFile.Close()
	if err != nil {
		panic(err)
	}

	// Let's read it in!
	tmpFile2, err := os.Open(tmpFile.Name())
	if err != nil {
		panic(err)
	}
	var fromFile Person
	err = json.NewDecoder(tmpFile2).Decode(&fromFile)
	if err != nil {
		panic(err)
	}
	err = tmpFile2.Close()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", fromFile)
}

func exampleTwo() {
	const data = `
		{"name": "Fred", "age": 40}
		{"name": "Mary", "age": 21}
		{"name": "Pat", "age": 30}
	`
	var t struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	dec := json.NewDecoder(strings.NewReader(data))
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	for dec.More() {
		err := dec.Decode(&t)
		if err != nil {
			panic(err)
		}
		fmt.Println(t)
		err = enc.Encode(t)
		if err != nil {
			panic(err)
		}
	}
	out := b.String()
	fmt.Println(out)
}

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Order struct {
	ID          string      `json:"id"`
	Items       []Item      `json:"items"`
	DateOrdered RFC822ZTime `json:"date_ordered"`
	CustomerID  string      `json:"customer_id"`
}

type RFC822ZTime struct {
	time.Time
}

func (rt RFC822ZTime) MarshalJSON() ([]byte, error) {
	out := rt.Time.Format(time.RFC822Z)
	return []byte(`"` + out + `"`), nil
}

func (rt *RFC822ZTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	t, err := time.Parse(`"`+time.RFC822Z+`"`, string(b))
	if err != nil {
		return err
	}
	*rt = RFC822ZTime{t}
	return nil
}

func customParsing() {
	data := `
	{
		"id": "12345",
		"items": [
			{
				"id": "xyz123",
				"name": "Thing 1"
			},
			{
				"id": "abc789",
				"name": "Thing 2"
			}
		],
		"date_ordered": "01 May 20 13:01 +0000",
		"customer_id": "3"
	}`

	var o Order
	err := json.Unmarshal([]byte(data), &o)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", o)
	fmt.Println(o.DateOrdered.Month())
	out, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
