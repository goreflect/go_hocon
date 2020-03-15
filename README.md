#HOCON (Human-Optimized Config Object Notation)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/goreflect/go_hocon)

[HOCON Docs](https://github.com/typesafehub/config/blob/master/HOCON.md).

> Currently, some features are not implemented, the API might be a little changed in the future.


`example.go`

```go
package main

import (
  "fmt"
  "github.com/goreflect/go_hocon"
)

var configText = `
####################################
# Typesafe HOCON                   #
####################################

config {
  # Comment
  version = "0.0.1"
  one-second = 1s
  one-day = 1day
  array = ["one", "two", "three"] #comment
  bar = "bar"
  foo = foo.${config.bar} 
  number = 1
  object {
    a = "a"
    b = "b"
    c = {
            d = ${config.object.a} //comment
        }
    }
}
// fallback
config.object.a="newA"
config.object.c.f="valueF"

// self reference
self-ref=1
self-ref=[${self-ref}][2]

// byte size
byte-size=10MiB

// system envs
home:${HOME}

plus-equal=foo
plus-equal+=bar

plus-equal-array=[foo]
plus-equal-array+=[bar, ${HOME}]
`

func main() {
	conf, _ := configuration.ParseString(configText)

	duration, _ := conf.GetTimeDuration("config.one-second")
	fmt.Println("config.one-second:", duration)

	timeDuration, _ := conf.GetTimeDuration("config.one-day")
	fmt.Println("config.one-day:", timeDuration)

	list, _ := conf.GetStringList("config.array")
	fmt.Println("config.array:", list)

	getString, _ := conf.GetString("config.bar")
	fmt.Println("config.bar:", getString)

	s, _ := conf.GetString("config.foo")
	fmt.Println("config.foo:", s)

	getInt64, _ := conf.GetInt64("config.number")
	fmt.Println("config.number:", getInt64)

	s2, _ := conf.GetString("config.object.a")
	fmt.Println("config.object.a:", s2)

	s3, _ := conf.GetString("config.object.c.d")
	fmt.Println("config.object.c.d:", s3)

	s4, _ := conf.GetString("config.object.c.f")
	fmt.Println("config.object.c.f:", s4)

	int64List, _ := conf.GetInt64List("self-ref")
	fmt.Println("self-ref:", int64List)

	size, _ := conf.GetByteSize("byte-size")
	fmt.Println("byte-size:", size)

	s5, _ := conf.GetString("home")
	fmt.Println("home:", s5)

	s6, _ := conf.GetString("none", "default-value")
	fmt.Println("default:", s6)

	s7, _ := conf.GetString("plus-equal")
	fmt.Println("plus-equal:", s7)

	stringList, _ := conf.GetStringList("plus-equal-array")
	fmt.Println("plus-equal-array:", stringList)
}

```

```bash
> go run example.go
config.one-second: 1s
config.one-day: 24h0m0s
config.array: [one two three]
config.bar: bar
config.foo: foo.bar
config.number: 1
config.object.a: newA
config.object.c.d: a
config.object.c.f: valueF
self-ref: [1 2]
byte-size: 10485760
home: /home/user
default: default-value
plus-equal: foobar
plus-equal-array: [foo bar /home/user]
```
