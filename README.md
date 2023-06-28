# snowid

![workflow badge](https://github.com/guoyk93/snowid/actions/workflows/go.yml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/guoyk93/snowid.svg)](https://pkg.go.dev/github.com/guoyk93/snowid)

A concurrent-safe lock-free implementation of snowflake algorithm in Golang

## Install

`go get -u github.com/guoyk93/snowflake`

## Usage

```go
// create an unique identifier
id, _ := strconv.ParseUint(os.Getenv("WORKER_ID"), 10, 64)

// create an instance (a sonyflake like instance)
s := snowflake.New(snowflake.Options{
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

Less than `1us/op` on Apple MacBook Air (M1)

```
goos: darwin
goarch: arm64
pkg: github.com/guoyk93/snowid
BenchmarkGenerator_NewID-8       2465515               469.5 ns/op
PASS
ok      github.com/guoyk93/snowid       1.742s
```

## Donation

View <https://guoyk.xyz/donation>

## Credits

Guo Y.K., MIT License
