# LinTeX: LaTeX document linter
LinTeX is a project that aims to provide easy to work with data structures, functions, and methods to save the hassle of parsing LinTeX yourself. It comes with the option to compile LaTeX locally to LuaLaTeX, using the command line. This project makes use of Parsing Expression Grammar (PEG) scripting via Pigeon.

Steps will be taken to convert this project into a Go package. However, if you would like to use it immediately, installation steps are below.

## Prerequisites
- `lualatex version 1.22.0 (TeX Live 2025)`
    - This can be checked via `lualatex -v` in terminal. The installation used in my context is via TeX Live, but you may be able to install `lualatex` through other methods.
    - It is required to be able to run `lualatex` in terminal as a normal user, as it is used to compile .tex files.
- `go version go1.25.4 windows/amd64`
    - Other similar versions may exist, with this being the one used for the context of this project.
- `pigeon version 1.3.0`
    - `pigeon -h` should showcase a valid help menu in terminal.
    - This is necessary to compile the grammar.peg script to generated its related grammar.go file.

## Installation
1. Install this project either via `git clone` or extracting its zip.
2. Ensure you are in the project root directory.
3. Optionally re-compile the grammar.peg script into a generated grammar.go file via `pigeon -o="grammar.go" grammar.peg"` in the root directory.
4. Run `go build .` within the root directory to build the `lintex` executable.
5. Run the `lintex` executable within the root directory, likely using `./lintex` or `./lintex.exe` depending on your operating system.

## Contributing
If you would like to make a contribution to this project, simply make a pull request with details of the changes being made and the relevant code changes. Some finer details on setting up more nitpicky aspects of the project have been provided below.

### PEG Script
1. Ensure Pigeon has been setup correctly, following the documentation found [here](https://pkg.go.dev/github.com/mna/pigeon).
2. Once changes to the grammar.peg file are complete, simply generate a new grammar.go file using `pigeon -o="grammar.go" grammar.peg` in the root directory.
3. Confirm that the changes made work as expected.