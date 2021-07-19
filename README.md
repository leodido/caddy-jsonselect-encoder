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

A selector represents the JSON path of the log entry you want to select.

The syntax is heavily inspired by [buger/jsonparser](https://github.com/buger/jsonparser). With some additions...

So, you can write a selector like the following one to only output the HTTP status code and the logger name:

```caddyfile
{status} {logger}
```

Or you can even select deeper in the JSON log entry:

```caddyfile
{request>host} {request>method}
```

The resulting JSON will respect the hierarchy of the selector paths.

Thus, for a selector like `{request>method}` the resulting JSON log entry will look like this:

```json
{"request":{"method":"GET"}}
```

Notice that the parsing of selectors happens at provisioning time to do not impact encoding performances.

Finally, I extended the syntax to support the reshaping of the resulting JSON keys.

To define a key for a given selector, you can use the following syntax:

```caddyfile
{key:selector}
```

For example, to store the status of a log entry in a a `httpRequest.responseStatus` JSON path you can write:

```caddyfile
{httpResponse>responseStatus:status}
```

Which will output the following JSON:

```json
{"httpRequest":{"responseSize":17064}}
```

This is particularly useful to adapt your log entries to different JSON structures like the Stackdriver one.

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

Even more, you can define keys for the resulting output to better match the [Stackdriver log entry format](https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry).

Like this:

```caddyfile
log {
  format jsonselect "{level} {timestamp:ts} {httpRequest>requestMethod:request>method} {httpRequest>protocol:request>proto} {httpRequest>status:status} {httpRequest>responseSize:size}" {
    time_format "rfc3339_nano"
  }
}
```

Which outputs:

```json
{"level":"info","timestamp":"2021-07-19T14:48:56.262966Z","httpRequest":{"protocol":"HTTP/2.0","requestMethod":"GET","responseSize":17604,"status":200}}
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
