# Guntar

Guntar is a CLI tool for working with tar archives. It allows you to read tar archives and browse, read, and extract files directly in memory.

## Features

- Browse tar archives in memory
- Extract files from tar archives
- List files within a tar archive


## Installation

__Work in Progress__ (I will add a docker image, packages, etc...)

For now you can clone the repository and run

```
make install-binary
```

## Usage

### Basic Command Structure

```sh
guntar [command]
```

### Available Commands

#### `explore`
Explore your tar archive in memory directly in your CLI. You can browse, look into files, and extract selected files/folders.
This interactive cli is based on [bubbletea](https://github.com/charmbracelet/bubbletea) project

Usage:
```sh
guntar explore <archive file> [flags]
```

Flags:
- `-h`, `--help`: Help for explore
- `-o`, `--output string`: Output directory to extract archive

Example:
```sh
guntar explore archive.tar -o output_directory
```

- Navigate through directories and files with arrows
- Select files or directory to extract with ctrl+a
    - no checkmark -> file or directory not selected
    - $\color{Green}{\textsf{✓}}$ -> file selected / all child in directory selected
    - $\color{Orange}{\textsf{✓}}$ -> some files are selected in the directory
- Extract files with ctrl+s

_Known Issues:_
- big files can break the textbox view -> will set a max size preview


#### `extract`
Extract files from a tar archive.

Usage:
```sh
guntar extract <archive file> [flags]
```

Flags:
- `-e`, `--ext []string`: List of files to extract
- `-h`, `--help`: Help for extract

Example:
```sh
guntar extract archive.tar -e file1.txt -e file2.txt
```

#### `help`
Display help information about any command.

Usage:
```sh
guntar help [command]
```

#### `list`
List all files in the current archive.

Usage:
```sh
guntar list <archive file> [flags]
```

Flags:
- `-h`, `--help`: Help for list

Example:
```sh
guntar list archive.tar
```

### Global Flags

- `-h`, `--help`: Display help information for Guntar.

## Examples

### List All Files in a Tar Archive

```sh
guntar list archive.tar
```

### Extract Files from a Tar Archive

```sh
guntar extract archive.tar -e file1.txt -e file2.txt
```

### Explore a Tar Archive in Memory

```sh
guntar explore archive.tar -o output_directory
```

## Getting Help

For more information about a specific command, use:

```sh
guntar [command] --help
```

## Contributing

TODO

## License

this project use MIT license