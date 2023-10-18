# go-ini 🚀

`go-ini` is a Go package tailored for manipulation, reading, and writing of INI files. This package is specifically designed to be compliant with [Magic Software's](https://www.magicsoftware.com/) "Magic.ini" format.

<div align="center">
    <img src="./docs/assets/magic.gif" alt="Magic fun img" width="300"/>
</div>

[![](https://img.shields.io/github/actions/workflow/status/devzolo/go-ini/test.yml?branch=main&longCache=true&label=Test&logo=github%20actions&logoColor=fff)](https://github.com/devzolo/go-ini/actions?query=workflow%3ATest)
[![Go Reference](https://pkg.go.dev/badge/github.com/devzolo/go-ini.svg)](https://pkg.go.dev/github.com/devzolo/go-ini)

## 💡 Features

- **Magic Software Compliance**: Adheres to the "Magic.ini" format specifications.
- **Load and Parse**: Efficiently reads INI files.
- **Write and Update**: Capable of writing updates to existing INI files or creating new ones.
- **Section and Key Extraction**: Extract sections and keys from INI content with ease.
- **Value Translation**: Translate values based on translatable sections.

## 🛠 Usage

### 🔧 Installation

To install the package, use:

```bash
go get github.com/devzolo/go-ini
```

### Basic Usage

```go
package main

import (
  "fmt"

  "github.com/devzolo/go-ini"
)

func main() {
  cfg := ini.NewMagicIni()
  err := cfg.LoadIni("path/to/your/magic.ini")
  if err != nil {
    panic(err)
  }
  value := cfg.Get("SomeSection", "SomeKey")
  fmt.Println(value)
}
```

## 📚 Documentation

### 📄 `MagicIni` Structure

Represents a structured format of an INI file compliant with the "Magic.ini" standards. It houses the parsed data, the currently parsed section, and the section with translatable strings.

### 🛠 Primary Methods

- `NewMagicIni()`: Creates and returns a new `MagicIni` instance.
- `LoadIni(path string)`: Reads and parses an INI file, particularly "Magic.ini", from the specified path.
- `LoadAdditionalIni(path string)`: Loads an additional INI file from the specified path, merging its contents with the existing data.
- `Get(section string, key string)`: Fetches a value for a specific section and key.
- `Translate(str string)`: Translates placeholders in the provided string using the corresponding strings from the `TranslatableSection`.
