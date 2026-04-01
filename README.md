# LinTeX: Enhancing LaTeX document accessibility for screen reader users
LinTeX is a tool under development that aims to identify accessibility concerns for screen reader users, and suggest solutions for them. It is outwith the scope of this project to understand and acount for every LaTeX command and its usage, so the majority of commands are, in essence, ignored by LinTeX. LinTeX focuses on ensuring the document's general structure is valid, and that tagging and image descriptors are appropriate for a screen reader to make use of.

## Prerequisites
- `lualatex version 1.22.0 (TeX Live 2025)`
    - This can be checked via `lualatex -v` in terminal. The installation used in my context is via TeX Live, but you may be able to install `lualatex` through other methods.
    - It is required to be able to run `lualatex` in terminal as a normal user, as it is used to compile .tex files.
- `npm (Node package manager) version 10.8.2`
    - This exists so related CSS frameworks can be imported for usage.
    - `npx` should also be supported, with the same version. This may be downloaded together with `npm` by default, or require separate installation, depending on your context.
- `tailwindcss/cli version 4.2.2`
    - This is utilised to use Tailwind CSS without using JavaScript.
    - `npm install tailwindcss @tailwindcss/cli` can be used to install the relevant package, with `npm install` within the root directory being another option.
- `flowbite version 4.0.1`
    - This is the components library used together with Tailwind CSS to utilise pre-built components.
    - `npm install` within the root directory will install this for you.
- `go version go1.25.4 windows/amd64`
    - Other similar versions may exist, with this being the one used for the context of this project.
- `pigeon version 1.3.0`
    - `pigeon -h` should showcase a valid help menu in terminal.
    - This is necessary to compile the grammar.peg script to generated its related grammar.go file.

## Installation
1. Install this project either via `git clone` or extracting its zip.
2. Ensure you are in the project root directory.
3. Run `npm install` to install necessary packages.
4. Optionally re-compile the grammar.peg script into a generated grammar.go file via `pigeon -o="grammar.go" grammar.peg"` in the root directory.
5. Optionally re-compile the css files via `npx @tailwindcss/cli -i ./styles/main.css -o ./styles/output.css` in the root directory.
6. Run `go build .` within the root directory to build the `lintex` executable.
7. Run the `lintex` executable within the root directory, likely using `./lintex` or `./lintex.exe` depending on your operating system.
8. Visit `http://localhost:8080/` in your browser of choice. You are now loaded in the LinTeX website and can start exploring it.

## Contributing
If you would like to make a contribution to this project, simply make a pull request with details of the changes being made and the relevant code changes. Some finer details on setting up more nitpicky aspects of the project have been provided below.

### PEG Script
1. Ensure Pigeon has been setup correctly, following the documentation found [here](https://pkg.go.dev/github.com/mna/pigeon).
2. Once changes to the grammar.peg file are complete, simply generate a new grammar.go file using `pigeon -o="grammar.go" grammar.peg` in the root directory.
3. Confirm that the changes made work as expected.

### CSS Setup
1. Ensure Tailwind CSS's CLI has been setup. You may use `npm install` to speed up this process, or visit documentation [here](https://tailwindcss.com/docs/installation/tailwind-cli). Ensure you install via `npm` rather than using a `cdn`.
2. Ensure Flowbite has been setup. You may use `npm install` to speed up this process, or visit documentation [here](https://flowbite.com/docs/getting-started/quickstart/)
3. You may choose to either watch for changes in CSS (using `npx @tailwindcss/cli -i ./styles/main.css -o ./styles/output.css --watch`) or compile it once (using `npx @tailwindcss/cli -i ./styles/main.css -o ./styles/output.css`)
    - Any newly created CSS files should be created within the `styles` folder.
    - This is required for any new styling changes, to re-compile the `main.css` file within the `styles` folder, or to create a new CSS file that uses Tailwind CSS.
4. Ensure any newly created files follow the Tailwind CSS CLI documentation steps, with the folder being changed to `styles`.