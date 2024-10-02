# go-lox
An interperter for the [Lox](https://craftinginterpreters.com/the-lox-language.html) programming language, written in Go.

Currently suppports:
- Arithmetic operators `+`, `-`, `*`, `\`
- Comparision operators `>`, `>=`, `<`, `<=`
- Equality operators `==`, `!=`
- Grouping with `()`
- Variable declarations with `var`
- Print statements with `print`
- Code blocks with localized variable scopes `{}`

Check [run_tests.yaml](cmd/myinterpreter/run_tests.yaml) for example programs

## Usage

Use the `go-lox.sh` script to build and run the interpreter.

Run programs with:
```sh
$ ./go-lox.sh run [file]
```

## Plumbing commands
### Tokenization

```sh
# input-file.lox
({*.,+*});

# tokenize the input file
$ ./go-lox.sh tokenize input-file.lox
LEFT_PAREN ( null
LEFT_BRACE { null
STAR * null
DOT . null
COMMA , null
PLUS + null
STAR * null
RIGHT_BRACE } null
RIGHT_PAREN ) null
SEMICOLON ; null
EOF  null
```

### Parsing and expression generation
```sh
# input-file.lox
var x = 5 + 1;
print x;

# parse the input file
$ ./go-lox.sh parse input-file.lox
(var (= x (+ 5 1)))
(print x)
```