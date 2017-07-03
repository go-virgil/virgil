// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package virgil

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonBa4d1ddcDecodeGopkgInVirgilV5(in *jlexer.Lexer, out *ValidationInfo) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "token":
			out.Token = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBa4d1ddcEncodeGopkgInVirgilV5(out *jwriter.Writer, in ValidationInfo) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Token != "" {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"token\":")
		out.String(string(in.Token))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ValidationInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBa4d1ddcEncodeGopkgInVirgilV5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ValidationInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBa4d1ddcEncodeGopkgInVirgilV5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ValidationInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBa4d1ddcDecodeGopkgInVirgilV5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ValidationInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBa4d1ddcDecodeGopkgInVirgilV5(l, v)
}
func easyjsonBa4d1ddcDecodeGopkgInVirgilV51(in *jlexer.Lexer, out *SignableRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "content_snapshot":
			if in.IsNull() {
				in.Skip()
				out.Snapshot = nil
			} else {
				out.Snapshot = in.Bytes()
			}
		case "meta":
			(out.Meta).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBa4d1ddcEncodeGopkgInVirgilV51(out *jwriter.Writer, in SignableRequest) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"content_snapshot\":")
	out.Base64Bytes(in.Snapshot)
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"meta\":")
	(in.Meta).MarshalEasyJSON(out)
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SignableRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBa4d1ddcEncodeGopkgInVirgilV51(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SignableRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBa4d1ddcEncodeGopkgInVirgilV51(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SignableRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBa4d1ddcDecodeGopkgInVirgilV51(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SignableRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBa4d1ddcDecodeGopkgInVirgilV51(l, v)
}
func easyjsonBa4d1ddcDecodeGopkgInVirgilV52(in *jlexer.Lexer, out *RequestMeta) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "signs":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Signatures = make(map[string][]uint8)
				} else {
					out.Signatures = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v4 []uint8
					if in.IsNull() {
						in.Skip()
						v4 = nil
					} else {
						v4 = in.Bytes()
					}
					(out.Signatures)[key] = v4
					in.WantComma()
				}
				in.Delim('}')
			}
		case "validation":
			if in.IsNull() {
				in.Skip()
				out.Validation = nil
			} else {
				if out.Validation == nil {
					out.Validation = new(ValidationInfo)
				}
				(*out.Validation).UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBa4d1ddcEncodeGopkgInVirgilV52(out *jwriter.Writer, in RequestMeta) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"signs\":")
	if in.Signatures == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
		out.RawString(`null`)
	} else {
		out.RawByte('{')
		v6First := true
		for v6Name, v6Value := range in.Signatures {
			if !v6First {
				out.RawByte(',')
			}
			v6First = false
			out.String(string(v6Name))
			out.RawByte(':')
			out.Base64Bytes(v6Value)
		}
		out.RawByte('}')
	}
	if in.Validation != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"validation\":")
		if in.Validation == nil {
			out.RawString("null")
		} else {
			(*in.Validation).MarshalEasyJSON(out)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestMeta) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBa4d1ddcEncodeGopkgInVirgilV52(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestMeta) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBa4d1ddcEncodeGopkgInVirgilV52(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestMeta) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBa4d1ddcDecodeGopkgInVirgilV52(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestMeta) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBa4d1ddcDecodeGopkgInVirgilV52(l, v)
}
func easyjsonBa4d1ddcDecodeGopkgInVirgilV53(in *jlexer.Lexer, out *CardModel) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "identity":
			out.Identity = string(in.String())
		case "identity_type":
			out.IdentityType = string(in.String())
		case "public_key":
			if in.IsNull() {
				in.Skip()
				out.PublicKey = nil
			} else {
				out.PublicKey = in.Bytes()
			}
		case "scope":
			out.Scope = Enum(in.String())
		case "data":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Data = make(map[string]string)
				} else {
					out.Data = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v10 string
					v10 = string(in.String())
					(out.Data)[key] = v10
					in.WantComma()
				}
				in.Delim('}')
			}
		case "info":
			if in.IsNull() {
				in.Skip()
				out.DeviceInfo = nil
			} else {
				if out.DeviceInfo == nil {
					out.DeviceInfo = new(DeviceInfo)
				}
				easyjsonBa4d1ddcDecodeGopkgInVirgilV54(in, &*out.DeviceInfo)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBa4d1ddcEncodeGopkgInVirgilV53(out *jwriter.Writer, in CardModel) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"identity\":")
	out.String(string(in.Identity))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"identity_type\":")
	out.String(string(in.IdentityType))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"public_key\":")
	out.Base64Bytes(in.PublicKey)
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"scope\":")
	out.String(string(in.Scope))
	if len(in.Data) != 0 {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"data\":")
		if in.Data == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v13First := true
			for v13Name, v13Value := range in.Data {
				if !v13First {
					out.RawByte(',')
				}
				v13First = false
				out.String(string(v13Name))
				out.RawByte(':')
				out.String(string(v13Value))
			}
			out.RawByte('}')
		}
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"info\":")
	if in.DeviceInfo == nil {
		out.RawString("null")
	} else {
		easyjsonBa4d1ddcEncodeGopkgInVirgilV54(out, *in.DeviceInfo)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CardModel) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonBa4d1ddcEncodeGopkgInVirgilV53(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CardModel) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBa4d1ddcEncodeGopkgInVirgilV53(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CardModel) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonBa4d1ddcDecodeGopkgInVirgilV53(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CardModel) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBa4d1ddcDecodeGopkgInVirgilV53(l, v)
}
func easyjsonBa4d1ddcDecodeGopkgInVirgilV54(in *jlexer.Lexer, out *DeviceInfo) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "device":
			out.Device = string(in.String())
		case "device_name":
			out.DeviceName = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonBa4d1ddcEncodeGopkgInVirgilV54(out *jwriter.Writer, in DeviceInfo) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"device\":")
	out.String(string(in.Device))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"device_name\":")
	out.String(string(in.DeviceName))
	out.RawByte('}')
}