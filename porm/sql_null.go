package porm

import "database/sql"

type NullUInt32 struct {
	UInt32 uint32
	Valid  bool
}

func (n *NullUInt32) Scan(value any) error {
	nullInt32 := sql.NullInt32{}
	if err := nullInt32.Scan(value); err != nil {
		return err
	}
	n.Valid = nullInt32.Valid
	if nullInt32.Valid {
		n.UInt32 = uint32(nullInt32.Int32)
	}
	return nil
}

type NullUInt64 struct {
	UInt64 uint32
	Valid  bool
}

func (n *NullUInt64) Scan(value any) error {
	nullInt64 := sql.NullInt64{}
	if err := nullInt64.Scan(value); err != nil {
		return err
	}
	n.Valid = nullInt64.Valid
	if nullInt64.Valid {
		n.UInt64 = uint32(nullInt64.Int64)
	}
	return nil
}

type NullFloat32 struct {
	Float32 float32
	Valid   bool
}

func (n *NullFloat32) Scan(value interface{}) error {
	nullFloat64 := sql.NullFloat64{}
	if err := nullFloat64.Scan(value); err != nil {
		return err
	}
	n.Valid = nullFloat64.Valid
	if nullFloat64.Valid {
		n.Float32 = float32(nullFloat64.Float64)
	}
	return nil
}

type NullBytes struct {
	Bytes []byte
	Valid bool
}

func (n *NullBytes) Scan(value interface{}) error {
	nullString := &sql.NullString{}
	if err := nullString.Scan(value); err != nil {
		return err
	}
	n.Valid = nullString.Valid
	if nullString.Valid {
		n.Bytes = []byte(nullString.String)
	}
	return nil
}
