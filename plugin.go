// Copyright 2021 Leonardo Di Donato
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonselect

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
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

	getters [][]string
	setters [][]string
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

	e.setters = [][]string{}
	e.getters = [][]string{}
	r := caddy.NewReplacer()
	r.Map(func(sel string) (interface{}, bool) {
		var set, get string

		parts := strings.Split(sel, ":")
		if len(parts) == 1 {
			set = parts[0]
			get = set
		} else if len(parts) == 2 {
			set = parts[0]
			get = parts[1]
		} else {
			// todo > error out - how?
			return nil, false
		}

		e.setters = append(e.setters, strings.Split(set, ">"))
		e.getters = append(e.getters, strings.Split(get, ">"))
		return nil, false
	})
	r.ReplaceAll(e.Selector, "")

	if len(e.setters) != len(e.getters) {
		return fmt.Errorf("selector must have the same number of setters and getters")
	}

	e.Encoder = zapcore.NewJSONEncoder(e.ZapcoreEncoderConfig())
	return nil
}

func (e JSONSelectEncoder) Clone() zapcore.Encoder {
	return JSONSelectEncoder{
		LogEncoderConfig: e.LogEncoderConfig,
		Encoder:          e.Encoder.Clone(),
		Selector:         e.Selector,
		getters:          e.getters,
		setters:          e.setters,
	}
}

func (e JSONSelectEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return buf, err
	}

	res := []byte{'{', '}'}
	// Temporary workaround the bug https://github.com/buger/jsonparser/issues/232
	// todo > switch back to EachKey (see git history) for perf reasons when fixed
	for idx, paths := range e.getters {
		val, typ, _, err := jsonparser.Get(buf.Bytes(), paths...)
		if err == jsonparser.KeyPathNotFoundError {
			// Path not found, skip
			continue
		}
		if err != nil {
			return nil, err
		}
		switch typ {
		case jsonparser.NotExist:
			// path not found, skip
		case jsonparser.String:
			res, _ = jsonparser.Set(res, append(append([]byte{'"'}, val...), '"'), e.setters[idx]...)
		default:
			res, _ = jsonparser.Set(res, val, e.setters[idx]...)
		}
	}

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

// Interface guards
var (
	_ zapcore.Encoder       = (*JSONSelectEncoder)(nil)
	_ caddy.Provisioner     = (*JSONSelectEncoder)(nil)
	_ caddyfile.Unmarshaler = (*JSONSelectEncoder)(nil)
)
