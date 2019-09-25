kana
====

This is a simple CLI utility to help drill kana (both katakana and hiragana).

## Installation ##

You must have golang installed (see: [https://golang.org/doc/install#install](https://golang.org/doc/install#install))

Then you can simply install with `go get`

```bash
> go get github.com/wcharczuk/kana
```

## Usage

The CLI will ask you to give the romanized versions of a random kana character.

See the source file for the correct answers, but typically it's somewhere from 1-3 latin ascii characters as a correct answer.

## Example Output

```bash
$ go run main.go --katakana=false --hiragana=true
    __ __  ___     _   __ ___
   / //_/ /   |   / | / //   |
  / ,<   / /| |  /  |/ // /| |
 / /| | / ___ | / /|  // ___ |
/_/ |_|/_/  |_|/_/ |_//_/  |_|


へ? he
correct!
ぬ? nu
correct!
ん? n
correct!
は? ^C
Session totals: 3/3 (100.00%)
```
