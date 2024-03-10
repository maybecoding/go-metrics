// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package entity

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

func easyjson9478868cDecodeGithubComMaybecodingGoMetricsGitInternalServerEntity(in *jlexer.Lexer, out *MetricsList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(MetricsList, 0, 1)
			} else {
				*out = MetricsList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Metrics
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9478868cEncodeGithubComMaybecodingGoMetricsGitInternalServerEntity(out *jwriter.Writer, in MetricsList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v MetricsList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9478868cEncodeGithubComMaybecodingGoMetricsGitInternalServerEntity(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MetricsList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9478868cEncodeGithubComMaybecodingGoMetricsGitInternalServerEntity(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MetricsList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9478868cDecodeGithubComMaybecodingGoMetricsGitInternalServerEntity(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MetricsList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9478868cDecodeGithubComMaybecodingGoMetricsGitInternalServerEntity(l, v)
}
func easyjson9478868cDecodeGithubComMaybecodingGoMetricsGitInternalServerEntity1(in *jlexer.Lexer, out *Metrics) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "type":
			out.MType = string(in.String())
		case "delta":
			if in.IsNull() {
				in.Skip()
				out.Delta = nil
			} else {
				if out.Delta == nil {
					out.Delta = new(int64)
				}
				*out.Delta = int64(in.Int64())
			}
		case "value":
			if in.IsNull() {
				in.Skip()
				out.Value = nil
			} else {
				if out.Value == nil {
					out.Value = new(float64)
				}
				*out.Value = float64(in.Float64())
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
func easyjson9478868cEncodeGithubComMaybecodingGoMetricsGitInternalServerEntity1(out *jwriter.Writer, in Metrics) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.MType))
	}
	if in.Delta != nil {
		const prefix string = ",\"delta\":"
		out.RawString(prefix)
		out.Int64(int64(*in.Delta))
	}
	if in.Value != nil {
		const prefix string = ",\"value\":"
		out.RawString(prefix)
		out.Float64(float64(*in.Value))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Metrics) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9478868cEncodeGithubComMaybecodingGoMetricsGitInternalServerEntity1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Metrics) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9478868cEncodeGithubComMaybecodingGoMetricsGitInternalServerEntity1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Metrics) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9478868cDecodeGithubComMaybecodingGoMetricsGitInternalServerEntity1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Metrics) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9478868cDecodeGithubComMaybecodingGoMetricsGitInternalServerEntity1(l, v)
}
