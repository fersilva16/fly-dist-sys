package main

import "encoding/json"

var (
	READ  string = "r"
	WRITE string = "w"
)

type Op struct {
	fn    string
	key   int
	value interface{}
}

type Txn []Op

func (o *Op) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{o.fn, o.key, o.value})
}

func (o *Op) UnmarshalJSON(b []byte) error {
	var op []interface{}

	if err := json.Unmarshal(b, &op); err != nil {
		return err
	}

	o.fn = op[0].(string)
	o.key = int(op[1].(float64))
	o.value = op[2]

	return nil
}
