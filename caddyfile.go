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
	"strings"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func (e *JSONSelectEncoder) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		args := d.RemainingArgs()
		switch len(args) {
		case 0:
			return d.Errf("%s (%T) requires an argument", moduleID, e)
		default:
			e.Selector = strings.Join(args, " ")
		}

		for n := d.Nesting(); d.NextBlock(n); {
			subdir := d.Val()
			var arg string
			if !d.AllArgs(&arg) {
				return d.ArgErr()
			}
			switch subdir {
			case "message_key":
				e.MessageKey = &arg
			case "level_key":
				e.LevelKey = &arg
			case "time_key":
				e.TimeKey = &arg
			case "name_key":
				e.NameKey = &arg
			case "caller_key":
				e.CallerKey = &arg
			case "stacktrace_key":
				e.StacktraceKey = &arg
			case "line_ending":
				e.LineEnding = &arg
			case "time_format":
				e.TimeFormat = arg
			case "level_format":
				e.LevelFormat = arg
			default:
				return d.Errf("unrecognized subdirective %s", subdir)
			}
		}
	}
	return nil
}
