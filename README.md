# snowid

[![workflow badge](https://github.com/yankeguo/snowid/actions/workflows/go.yml/badge.svg)](https://github.com/yankeguo/snowid/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/yankeguo/snowid.svg)](https://pkg.go.dev/github.com/yankeguo/snowid)
[![codecov](https://codecov.io/gh/yankeguo/snowid/graph/badge.svg?token=MAMV3GWQ8R)](https://codecov.io/gh/yankeguo/snowid)

A concurrent-safe lock-free implementation of snowflake algorithm in Golang

## Install

`go get -u github.com/yankeguo/snowid`

## Usage

```go
// create an unique identifier
id, _ := strconv.ParseUint(os.Getenv("WORKER_ID"), 10, 64)

// create an instance (a sonyflake like instance)
s := snowid.New(snowflake.Options{
    Epoch: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
    ID: id,
    Grain: time.Millisecond*10,
    LeadingBit: true,
})

// get a id
s.NewID()

// stop and release all related resource
s.Stop()

```

## Performance

Less than `1us/op` on **Apple MacBook Air (M1)**

```
goos: darwin
goarch: arm64
pkg: github.com/yankeguo/snowid
BenchmarkGenerator_NewID-8       2465515               469.5 ns/op
PASS
ok      github.com/yankeguo/snowid       1.742s
```

## Credits

GUO YANKE, MIT License
