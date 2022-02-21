package tvm

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"tongo/boc"
)

type EntryType int

const (
	Int EntryType = iota
	Null
	Cell
	Tuple
)

type TvmStackEntry struct {
	Type     EntryType
	intVal   big.Int
	cellVal  *boc.Cell
	tupleVal []TvmStackEntry
}

func NewIntStackEntry(val big.Int) TvmStackEntry {
	return TvmStackEntry{
		Type:   Int,
		intVal: val,
	}
}

func NewNullStackEntry() TvmStackEntry {
	return TvmStackEntry{
		Type: Null,
	}
}

func (e *TvmStackEntry) Int() big.Int {
	return e.intVal
}

func (e *TvmStackEntry) Int64() int64 {
	return e.intVal.Int64()
}

func (e *TvmStackEntry) Uint64() uint64 {
	return e.intVal.Uint64()
}

func (e *TvmStackEntry) Cell() *boc.Cell {
	return e.cellVal
}

func (e *TvmStackEntry) Tuple() []TvmStackEntry {
	return e.tupleVal
}

func (e *TvmStackEntry) IsNull() bool {
	return e.Type == Null
}

func (e *TvmStackEntry) IsInt() bool {
	return e.Type == Int
}

func (e *TvmStackEntry) IsCell() bool {
	return e.Type == Cell
}

func (e *TvmStackEntry) IsTuple() bool {
	return e.Type == Tuple
}

func (e *TvmStackEntry) UnmarshalJSON(data []byte) error {
	var m map[string]json.RawMessage

	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	var entryType string
	err = json.Unmarshal(m["type"], &entryType)
	if err != nil {
		return err
	}

	if entryType == "int" {
		var intEntry struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}
		err = json.Unmarshal(data, &intEntry)
		if err != nil {
			return err
		}

		e.Type = Int
		e.intVal.SetString(intEntry.Value, 10)
	} else if entryType == "null" {
		e.Type = Null
	} else if entryType == "tuple" {
		var tupleEntry struct {
			Type  string          `json:"type"`
			Value []TvmStackEntry `json:"value"`
		}
		err = json.Unmarshal(data, &tupleEntry)
		if err != nil {
			return err
		}
		e.Type = Tuple
		e.tupleVal = tupleEntry.Value
	} else if entryType == "cell" {
		var cellEntry struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}
		err = json.Unmarshal(data, &cellEntry)
		if err != nil {
			return err
		}

		e.Type = Cell
		cellData, err := base64.StdEncoding.DecodeString(cellEntry.Value)
		if err != nil {
			return err
		}
		parsedBoc, err := boc.DeserializeBoc(cellData)
		if err != nil {
			return err
		}
		e.cellVal = parsedBoc[0]
	} else {
		return errors.New("unknown stack entry type")
	}

	return nil
}
