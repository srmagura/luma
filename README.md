# The Luma Programming Language

A programming language with compile-time type safety that compiles into bytecode. The compiler is written in Go. The runtime is a stack-based virtual machine written in Rust.

## Example Programs

- [Fibonacci numbers](./test_runner/fibonacci.luma)

## Development

Development is currently being done on Linux and macOS.

- Build compiler: `go build -C compiler`
- Test compiler: `go test -C compiler`
- Run end-to-end tests: `./test.sh`
