package pac194x5x

import (
	"encoding/binary"
	"fmt"
)

// Codec is an interface for marshalling/unmarshalling data.
type Codec[T any] interface {
	Marshal(value T) ([]byte, error)
	Unmarshal(data []byte) (T, error)
}

var (
	VoidCodec      = &voidCodec{}   // VoidCodec - Codec for Void.
	Uint8Codec     = &uint8Codec{}  // Uint8Codec - Codec for uint8.
	Uint16Codec    = &uint16Codec{} // Uint16Codec - Codec for uint16.
	Uint32Codec    = &uint32Codec{} // Uint32Codec - Codec for uint32.
	Uint64Codec    = &uint64Codec{}
	ProductIDCodec = &productIDCodec{} // ProductIDCodec - Codec for ProductID.
)

type Void any

type voidCodec struct {
}

func (codec *voidCodec) Marshal(_ Void) ([]byte, error) {
	return nil, nil
}

func (codec *voidCodec) Unmarshal(data []byte) (Void, error) {
	if len(data) != 0 {
		return 0, fmt.Errorf("expected 0 bytes, got %d", len(data))
	}
	return nil, nil
}

type uint8Codec struct {
}

func (codec *uint8Codec) Marshal(value uint8) ([]byte, error) {
	return []byte{value}, nil
}

func (codec *uint8Codec) Unmarshal(data []byte) (uint8, error) {
	if len(data) != 1 {
		return 0, fmt.Errorf("expected 1 byte, got %d", len(data))
	}
	return data[0], nil
}

type uint16Codec struct {
}

func (codec *uint16Codec) Marshal(value uint16) ([]byte, error) {
	return binary.BigEndian.AppendUint16(nil, value), nil
}

func (codec *uint16Codec) Unmarshal(data []byte) (uint16, error) {
	if len(data) != 2 {
		return 0, fmt.Errorf("expected 2 bytes, got %d", len(data))
	}
	return binary.BigEndian.Uint16(data), nil
}

type uint32Codec struct {
}

func (codec *uint32Codec) Marshal(value uint32) ([]byte, error) {
	return binary.BigEndian.AppendUint32(nil, value), nil
}

func (codec *uint32Codec) Unmarshal(data []byte) (uint32, error) {
	if len(data) != 4 {
		return 0, fmt.Errorf("expected 4 bytes, got %d", len(data))
	}
	return binary.BigEndian.Uint32(data), nil
}

type uint64Codec struct {
}

func (codec *uint64Codec) Marshal(value uint64) ([]byte, error) {
	return binary.BigEndian.AppendUint64(nil, value)[1:], nil
}

func (codec *uint64Codec) Unmarshal(data []byte) (uint64, error) {
	if len(data) != 7 {
		return 0, fmt.Errorf("expected 7 bytes, got %d", len(data))
	}
	adjustedData := append([]byte{0x00}, data...)
	return binary.BigEndian.Uint64(adjustedData), nil
}

type productIDCodec struct {
}

func (codec *productIDCodec) Marshal(value ProductID) ([]byte, error) {
	return Uint8Codec.Marshal(uint8(value))
}

func (codec *productIDCodec) Unmarshal(data []byte) (ProductID, error) {
	v, err := Uint8Codec.Unmarshal(data)
	if err != nil {
		return 0, err
	}
	return ProductID(v), nil
}
