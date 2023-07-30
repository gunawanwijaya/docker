package z

import "encoding/base64"

// Base64
//
//	mode = 0 -> Std().Strict()
//	mode = 1 -> URL().Strict()
func (encoding) Base64(mode int) EncoderDecoder {
	var e *base64.Encoding
	switch mode {
	default:
		e = base64.RawStdEncoding.Strict()
	case 1:
		e = base64.RawURLEncoding.Strict()
	}

	return encoding_base64{e}
}

type encoding_base64 struct {
	e *base64.Encoding
}

func (x encoding_base64) Encode(msg []byte) (e EncodedValue, err error) {
	return encodedValue{Text(x.e.EncodeToString(msg))}, nil

}
func (x encoding_base64) Decode(e EncodedValue) (msg []byte, err error) {
	return x.e.DecodeString(e.String())
}
