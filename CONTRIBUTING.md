# Contributing to Go Go Gadgets

Thank you for considering working on these gadgets. Please look below to see more information about how to contribute and what is expected when it comes to contributing to this project.

## Getting Started

1. **Fork the repository**: Use the Fork button at the top right of the repository page.
2. **Clone your fork**: Run `git clone https://github.com/your-username/go-go-gadgets.git` to clone your copy to your local machine.
3. **Create a branch**: Use `git checkout -b your-feature-branch` to create a new branch for your work.

## Code Structure

This project is made up of several CLI tools. They each have their own use cases, but they may share some code and packages. The code can be divided into two main areas: CLI tools and common packages.

### CLI Tools

A CLI tool should have its main folder found at the base of repository and have the following structure:
```
cli-folder
├── cmd
│   ├── generate.go
│   └── root.go
├── internal (optional)
├── main.go
├── README.md
└── README.md.tmpl
```

In the above example there are a couple of required elements:
- A `main.go` inside the base folder for the CLI tool.
- A `README.md` which is generated from `README.md.tmpl`
- A `cmd` folder that will house the Cobra CLI commands and any validation for those commands
- A `root.go` inside the `cmd` folder which will be the root command for Cobra
- A `generate.go` inside the `cmd` folder which will be used to programatically generate the README

There is an optional `internal` folder which should house any specific logic that is meant to be a kind of business rule or further logic that is non-trivial, but not shared with other CLI tools in this project. The only real time to use the `internal` folder is for larger pieces of logic that are specific to the CLI tool in question. So logic that will be reused goes elsewhere.

### Common Packages

There are also reusable packages as well. They are located in the `pkg` folder at the base of the repository. These packages contain logic that is meant to be reused across the different CLI tools as needed. For example, the `logger` package allows for a more streamline way to write errors and then immediately exit the program. There is also a special package called `tests` which is meant to house global test helpers. These may be reused across CLI tools or just by one and is meant to be available in the case that other tools need its content.

## Making Changes

### Adding a New CLI Tool

When adding a new CLI tool, go ahead and follow the code structure listed above in [CLI Tools](#cli-tools). Once that is done, you can add the name of the program to the `Makefile`, so it will get properly compiled, installed, and have its README generated. One things you will want to make sure you add the new CLI's base folder name to is `TOOLS`.

### CLI Commands and Flags

The CLI package used for writing the CLIs is [Cobra](https://github.com/spf13/cobra). It allows for easily setting up commands, flags, and subcommands. Make sure that when a command is created, it has a `Use` and `Short` value. The `Use` value is the name of the command. The `Short` is the short description of what it does. Where possible examples should be included in the `Example` property which should be a `heredoc.Doc` string to help cut out the starting whitespace in the multiline string.

Try to keep the nesting of commands to at most 2 layers. So in the case of `epub-lint`, it has a `fix` subcommand which also has its own subcommands. You should not have to type out more than three command names in order to get to the command you want to run.

Names for commands, subcommands, and flags should follow the pattern `NAME NOUN VERB` where possible.

Try to stick with common flag names and abbreviations like `--help` and `-h` or `--file` and `-f`. If a flag is not common or there are a lot of flags for the command, do not use an abbreviation for the flag name.

### Running Tests

When change have been made, tests should always be verified to make sure that they pass. This can be done by running `make test` at the base of the repository. It should let you know if any have failed and if so which ones.

It may help sometimes to add `-v` to the arguments for `go test` to make it show more of the output when a test fails.

## Testing Philosophy

Tests are meant to help make sure that the code is reliable and works as intended for specific scenarios. Tests are not meant to be used to cover all lines of code. You may notice that there is a make rule called `cover` in the Makefile. That rule is not meant to be used for determining the quality of the code. It is just meant to give an idea of where tests exist versus where they do not. Tests are only really needed where logic is prone to be incorrect or may be refactored later and the end result of the logic should not change.

Tests should be table driven tests where possible using a map of a string which is the name of the test to the actual test case struct which will be used in the test function. Test data can be inline in the file or it can come from a `testdata` folder when it is larger and may be better served as its own files.

## Documentation

Documentation is meant to help a user get up to speed or a developer to be more familiar with and able to work with code. In this project, all CLI tools should have a README. Comments in the code are fine, but they should not be required as that just encourages writing poor function and package comments.

Documentation can be found in 4 main places in this project:
- The project [README](README.md)
- The [Contributing Docs](CONTRIBUTING.md)
- Each CLI tool's README/README.md.tmpl
- Each CLI tool's `generate.go` where TODOs and some other specific information can be setup

As of right now there is really no other place for documentation to live in this project.

To generate updated documentation, run `make generate` at the base of the repository.

## Use of AI

The use of AI for contributing to this repo is perfectly fine. However you are expected to make sure that the code works or the documentation is correct. You should be able to answer why you are making a change. Saying "The AI said so" is not a good enough answer for why you are making a change. You should be able to explain in your own words why the change is happening.

AI can be helpful with tasks like adding UTs, bouncing ideas off of it, doing some refactors, as well as other tasks. Use it with discernment if you use it.
