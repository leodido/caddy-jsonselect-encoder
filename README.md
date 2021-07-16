# caddy-jsonselect-encoder

> Select what you want to log

By using a selector you can choose what (and how) you want to emit your JSON logs.

## Module

The **module name** is `jsonselect`.

Its syntax is:

```caddyfile
jsonselect <selector> {
    message_key <key>
    level_key   <key>
    time_key    <key>
    name_key    <key>
    caller_key  <key>
    stacktrace_key <key>
    line_ending  <char>
    time_format  <format>
    level_format <format>
}
```

### Selector

The syntax is heavily inspired by [buger/jsonparser](https://github.com/buger/jsonparser).

Thus, you can write a selector like the following one to only output the HTTP status code and the logger name:

```caddyfile
{status} {logger}
```

Or even more complex ones:

```caddyfile
{request>host} {request>method} {request>headers>User-Agent>[0]}
```

The resulting JSON will respect the hierarchy of the selector paths.

## Caddyfile

Log a JSON containing only the level, the timestamp, and the message of the log entry.
Also, use "mex" as a key for the message.

```console
log {
  output stdout
  format jsonselect "{level} {ts} {mex}" {
    message_key mex
  }
}
```

This will output:

```json
{"level":"info","ts":1626453781.3333929,"mex":"handled request"}
```

Log the host of the request and the duration:

```caddyfile
log {
  output stdout
  format jsonselect "{request>host} {duration}"
}
```

Which outputs something like the following JSON respecting the selector's path structure:

```json
{"request":{"host":"localhost:2015"},"duration":0.003321}
```

Maybe you wanna log for Stackdriver ...

```caddyfile
log {
  format if {
      status sw 400
  } jsonselect "{severity} {timestamp} {logName}" {
    level_key "severity"
    level_format "upper"
    time_key "timestamp"
    time_format "rfc3339"
    name_key "logName"
  }
}
```

This outputs:

```json
{"severity":"ERROR","timestamp":"2021-07-16T12:55:10Z","logName":"http.log.access.log0"}
```

## Try it out

From the root directoy of this project, run:

```console
xcaddy run
```

Then open <https://localhost:2015>, go on existing and non-existing pages, and observe the access logs.

To install xcaddy in case you need to, run:

```console
go get -u github.com/caddyserver/xcaddy/cmd/xcaddy
```

## Build

To build [Caddy](https://github.com/caddyserver/caddy) with this module in, execute:

```console
xcaddy build --with github.com/leodido/caddy-jsonselect-encoder
```

---

[![Analytics](https://ga-beacon.appspot.com/UA-49657176-1/caddy-jsonselect-encoder?flat)](https://github.com/igrigorik/ga-beacon)
