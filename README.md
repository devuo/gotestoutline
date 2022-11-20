# gotestoutline

CLI tool that outputs tests and subtests outline of a go test file in JSON format.

## Installing

```sh
go install github.com/devuo/gotestoutline@latest
```

## Using

```sh
> gotestoutline file_test.go
```

This will write an output like the following to stdout:

```json
[
    {
        "name": "TestMain"
        "type": "test",
        "path": [],
        "lbrace": 123,
        "rbrace": 458
    },
    {
        "name": "Succeeds",
        "type": "subtest",
        "path": ["TestMain"],
        "lbrace": 138,
        "rbrace": 235,
    },
    {
        "name": "Succeeds (With Option)",
        "type": "subtest",
        "path": ["TestMain", "Succeeds"],
        "lbrace": 140,
        "rbrace": 210,
    },
    {
        "name": "Fails",
        "type": "subtest",
        "path": ["TestMain"],
        "lbrace": 138,
        "rbrace": 189,
    },
    ...
]
```

## Schema

```go
const (
    TestType           Type = "test"
    SubtestType        Type = "subtest"
    DynamicSubtestType Type = "dynamicsubtest"
)

type Test struct {
    // Name of the test
    Name   string   `json:"name"`
    // Type of test
    Type   Type  `json:"type"`
    // Path to this test, including parent test names
    Path   []string `json:"path"`
    // Line where the test begins
    LBrace int    `json:"lbrace"`
    // Line where the test ends
    RBrace int    `json:"rbrace"`
}
```
