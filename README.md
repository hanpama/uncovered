# uncovered

Show uncovered lines from Go test coverage with context.

## About

`uncovered` analyzes Go coverage profiles and displays lines that weren't executed during testing, along with surrounding context lines for better understanding.

Unlike `go tool cover -html`, which requires opening a browser and checking files one by one, `uncovered` shows all uncovered code directly in your terminal with clear visual highlighting.

## Installation

```bash
go install github.com/hanpama/uncovered@latest
```

## Usage

```
uncovered - Show uncovered lines from Go test coverage with context

Usage:
  uncovered <coverprofile>

Arguments:
  coverprofile    Path to coverage profile file

Examples:
  # Generate coverage and show uncovered lines
  go test -coverprofile=coverage.out ./...
  uncovered coverage.out

  # Use custom coverage file
  uncovered custom.out
```

## Example Output

```
example/calculator.go (29 of 59 lines uncovered)
================================================
     21 }
     22
     23 // Divide returns the quotient of two numbers
>    24 func (c *Calculator) Divide(a, b int) (int, error) {
>    25   if b == 0 {
>    26     return 0, errors.New("division by zero")
>    27   }
>    28   return a / b, nil
     29 }
     30
     31 // Power returns a raised to the power of b
>    32 func (c *Calculator) Power(a, b int) int {
>    33   if b == 0 {
>    34     return 1
>    35   }
>    36   if b < 0 {
>    37     return 0 // simplified: doesn't handle negative powers
     ...
```

Lines marked with `>` are uncovered (displayed in red in the terminal).

## License

MIT
