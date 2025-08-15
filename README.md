[![Go Reference](https://pkg.go.dev/badge/github.com/dmikhr/mooc-dsl.svg)](https://pkg.go.dev/github.com/dmikhr/mooc-dsl)
[![License](https://img.shields.io/github/license/dmikhr/mooc-dsl.svg)](LICENSE)

MOOC DSL project aims to simplify the process of creating and managing educational tests across various platforms and systems.

Tutors often face the challenge of adapting their test descriptions to different Learning Management Systems (LMS) and standalone software. Each system typically has its own syntax and format, requiring tutors to invest significant time in understanding and translating their test descriptions.

This project addresses these challenges by introducing a Domain-Specific Language with simple and intuitive syntax for describing tests (see ```assets/sample.txt``` for DSL specifications). 

Tutors describe tests using DSL syntax, which are then processed and transformed into JSON format that can be integrated with third-party educational software, including LMS and custom educational platforms.

Project provides two main features:
* Syntax checker that outputs JSON with errors descriptions
* Parser which transforms DSL into JSON

Check ```assets``` folder for samples of JSON output and sample tests.

**Build instructions**

Install Go

From project directory build binary using `make`

```
make build
```

or build and run simultaneously

```
make build-run
```

Check `Makefile` for other available commands.

Available flags:
* `--fname` to set filename of test
* `--vebose` to print generate JSON in temrinal
* `--showErrors` to print DSL syntax errors in temrinal

In case of errors in the source DSL code, an error report will be generated with the name `[sourcename]_error.json`
