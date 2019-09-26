package ep6v2

import (
	"encoding/binary"
	"math"
	"unicode/utf16"
	"unsafe"
)

func DecodeUtf16String(data []byte) string {
	var buf []uint16
	for i := 0; i < len(data)/2; i++ {
		word := binary.BigEndian.Uint16(data[i*2:])
		if word == 0 {
			break
		}
		buf = append(buf, word)
	}
	return string(utf16.Decode(buf))
}

func DecodeAsciiString(data []byte) string {
	var buf []byte
	for i := 0; i < len(data)-1; i += 2 {
		if data[i+1] == 0 {
			break
		}
		buf = append(buf, data[i+1])
		if data[i] == 0 {
			break
		}
		buf = append(buf, data[i])
	}
	return string(buf)
}

func ToSingle(data []byte) float32 {
	//参见C# ToSingle源码 https://referencesource.microsoft.com/#mscorlib/system/bitconverter.cs,7d2958fc09cde954,references
	var val = [4]byte{data[1], data[0], data[3], data[2]}
	var pFloat = (*float32)(unsafe.Pointer(&val))
	return *pFloat
}

func ToFloat32(v float32, point int) float32 {
	return float32(math.Floor(float64(v)*math.Pow10(point)) / math.Pow10(point))
}
