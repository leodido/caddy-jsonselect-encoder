package jsonselect

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/logging"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const (
	moduleName = "jsonselect"
	moduleID   = "caddy.logging.encoders." + moduleName
)

func init() {
	caddy.RegisterModule(JSONSelectEncoder{})
}

type JSONSelectEncoder struct {
	logging.LogEncoderConfig
	zapcore.Encoder `json:"-"`
	Selector        string `json:"selector,omitempty"`

	keys [][]string
}

func (JSONSelectEncoder) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: moduleID,
		New: func() caddy.Module {
			return &JSONSelectEncoder{
				Encoder: new(logging.JSONEncoder),
			}
		},
	}
}

func (e *JSONSelectEncoder) Provision(ctx caddy.Context) error {
	if e.Selector == "" {
		return fmt.Errorf("selector is mandatory")
	}

	e.keys = [][]string{}
	r := caddy.NewReplacer()
	r.Map(func(key string) (interface{}, bool) {
		e.keys = append(e.keys, strings.Split(key, ">"))
		return nil, false
	})
	r.ReplaceAll(e.Selector, "")

	e.Encoder = zapcore.NewJSONEncoder(e.ZapcoreEncoderConfig())
	return nil
}

func (e JSONSelectEncoder) Clone() zapcore.Encoder {
	return JSONSelectEncoder{
		LogEncoderConfig: e.LogEncoderConfig,
		Encoder:          e.Encoder.Clone(),
		Selector:         e.Selector,
		keys:             e.keys,
	}
}

func (e JSONSelectEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return buf, err
	}

	res := []byte{'{', '}'}
	jsonparser.EachKey(
		buf.Bytes(),
		func(idx int, val []byte, typ jsonparser.ValueType, err error) {
			// todo > handle error
			switch typ {
			case jsonparser.String:
				res, _ = jsonparser.Set(res, append(append([]byte{'"'}, val...), '"'), e.keys[idx]...)
			default:
				res, _ = jsonparser.Set(res, val, e.keys[idx]...)
			}
		},
		e.keys...,
	)

	// Reset the buffer to output our own content
	buf.Reset()
	// Insert the new content
	nl := []byte("\n")
	if !bytes.HasSuffix(res, nl) {
		res = append(res, nl...)
	}
	buf.Write(res)

	return buf, err
}
