# Contributing to Go Go Gadgets

Thank you for considering working on these gadgets. Please look below to see more information about how to contribute and what is expected when it comes to contributing to this project.

## Getting Started

1. **Fork the repository**: Use the Fork button at the top right of the repository page.
2. **Clone your fork**: Run `git clone https://github.com/your-username/go-go-gadgets.git` to clone your copy to your local machine.
3. **Create a branch**: Use `git checkout -b your-feature-branch` to create a new branch for your work.

## Code Structure

This project is made up of several CLI tools each of which is a gadget. They each have their own use cases, but they may share some code and packages. The code can be divided into two main areas: CLI tools and common packages.

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

There is an optional `internal` folder which should house any specific logic that is meant to be a kind of business rule or further logic that is non-trivial, but not shared with other CLI tools in this project. The only real time to use the `internal` folder is for larger pieces of logic that are specific to the CLI tool in question. So logic that will be reused across CLI tools goes in `pkg` instead.

### Common Packages

There are reusable packages as well. They are located in the `pkg` folder at the base of the repository. These packages contain logic that is meant to be reused across the different CLI tools as needed. For example, the `logger` package allows for a more streamline way to write errors and then immediately exit the program. There is also a special package called `tests` which is meant to house global test helpers. These may be reused across CLI tools or just by one and is meant to be available in the case that other tools need its content.

You may be wondering whether or not something should live in an `internal` folder or inside the `pkg` directory. That will come down to whether or not there is a need for the logic in multiple tools. A lot of times, logic is specific to a single CLI tool and in those cases, the logic can safely go in an `internal` folder. Sometimes however there is logic that gets moved to a `pkg` folder even though it may not be reused. I recommend using your best judgment on this and if it needs moving just going with the flow on that.

## Making Changes

### Adding a New CLI Tool

When adding a new CLI tool, go ahead and follow the code structure listed above in [CLI Tools](#cli-tools). Once that is done, you can add the name of the program to the `Makefile`, so it will get properly compiled, installed, and have its README generated. One things you will want to make sure you add the new CLI's base folder name to is `TOOLS`.

Here is a list of things to make sure you do when you create a new CLI tool:
- [ ] Make sure that the CLI tool is created at the base of the repository
- [ ] Make sure that the CLI tool has a `cmd` folder with a `root.go` and `generate.go` in it
- [ ] Make sure that the CLI tool has a `README.md.tmpl` at the base of CLI tool's directory structure
- [ ] Make sure to add the CLI tool to the `Makefile` under the `TOOLS` variable
- [ ] Make sure that `make install` properly installs the CLI tool
- [ ] Make sure that `make clean` properly removes the CLI tool
- [ ] Make sure that the new CLI tool has been added to the top level README

### CLI Commands and Flags

The CLI package used for writing the CLIs is [Cobra](https://github.com/spf13/cobra). It allows for easily setting up commands, flags, and subcommands. A sample Cobra command that includes what you will need is as follows:

``` Go
var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Compresses and lints all of the epub files in the specified directory even compressing images using imgp if that option is specified.",
  Long: heredoc.Doc(`Gets all of the .epub files in the specified directory.
	Then it lints each epub separately making sure to compress the images if specified.
	Some of the things that the linting includes:
	- Replacing a list of common strings
	- Adds language encoding specified if it is not present already (default is "en")
	- Sets encoding on content files to utf-8 to prevent errors in some readers
	`),
	Example: heredoc.Doc(`To compress images and make general modifications to all epubs in a folder:
	epub-lint optimize -d folder -c
	
	To compress images and make general modifications to all epubs in the current directory:
	epub-lint optimize -c

	To just make general modifications to all epubs in the current directory:
	epub-lint optimize
	`),
	PreRunE: func(cmd *cobra.Command, args []string) error {/* Actual arg and flag validation here */},
	Run: func(cmd *cobra.Command, args []string) {/* Actual logic here*/},
}
```

In this example you will notice that there are several pieces to the Cobra command. The first part is the `Use` property. This is the name of the command. In this case it is called `optimize`.
The next property present `Short` which is the short description of what this command does. `Long` is similar to `Short` in that it describes what the command does, but it should be more descriptive and give the user a better understanding of how to use the command.
The subsequent property is `Example`. This is useful for showing the user sample usages of the command. It will show up under the `help` and `--help` print statements as well. It is not required that you have examples, but it is recommended that you have them especially when there are multiple flags or args.

Now, you may have noticed that there are several instances of `heredoc.Doc` present in the Cobra command. The reason for this is that it strips leading whitespace in the multiline strings. This makes the resulting help much cleaner and more readable while still allowing the source code to be readable since we don't have to visually move the text around to see what it would look like compared to prior or subsequent lines.

The next two things you will notice are `PreRunE` and `Run`. These house the flag and arg validation and general command logic respectively. The former will show the usage menu whenever an error is returned while the latter is just running the core command logic. `PreRunE` should only be used for validation around flags and args to make sure that all preconditions are met that can be calculated without actually getting into the meat of how the command should work. If there is nothing to validate, then do not define the `PreRunE` function.

When dealing with commands, you will find that there are times when it makes sense to nest commands. This is perfectly fine. You may notice in Epub Linter that it has a command called `fix` which has two subcommands: `content` and `validation`. So to run the `content` subcommand you would type `epub-lint fix content`. This is the most nesting that should happen. If a user needs to type more than 3 commands, then there is a problem. One more thing to note, is that creating subcommands like this should be easily spotted by their filenames which should be `PARENT_SUBCOMMAND.go`. So in the case of `content` it would be `fix_content.go`. This allows for spotting these groupings and makes them show up near each other in the file manager.

Try to stick with common flag names and abbreviations like `--help` and `-h` or `--file` and `-f`. If a flag is not common or there are a lot of flags for the command, do not use an abbreviation for the flag name.

### Running Tests

When changes have been made, tests should always be verified to make sure that they pass. This can be done by running `make test` at the base of the repository. It should let you know if any have failed and if so which ones.

It may help sometimes to add `-v` to the arguments for `go test` to make it show more of the output when a test fails.

There are some tests that rely on what are called `golden` files. These files can be regenerated if need by running `make golden`. This will pull any that are not flagged as being kept as is. Then you can run `make test` to verify those tests to make sure nothing changed.

#### Testing Philosophy

Tests are meant to help make sure that the code is reliable and works as intended for specific scenarios. Tests are not meant to be used to cover all lines of code. You may notice that there is a make rule called `cover` in the Makefile. That rule is not meant to be used for determining the quality of the code. It is just meant to give an idea of where tests exist versus where they do not. Tests are only really needed where logic is prone to be incorrect or may be refactored later and the end result of the logic should not change.

Tests should be table driven tests where possible using a map of a string which is the name of the test to the actual test case struct which will be used in the test function. Test data can be inline in the file or it can come from a `testdata` folder when it is larger and may be better served as its own files.

Here is an example of a table driven test format:
``` Go
//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

type updatePlayOrderTestCase struct {
	input          string
	expectedOutput string
}

var updatePlayOrderTestCases = map[string]updatePlayOrderTestCase{
	"Updating the play order works when there are duplicate playOrder values": {
		input: `<ncx>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.html" />
    </navPoint>
    <navPoint id="navPoint-2" playOrder="1">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.html" />
    </navPoint>
  </navMap>
</ncx>`,
		expectedOutput: `<ncx>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.html" />
    </navPoint>
    <navPoint id="navPoint-2" playOrder="2">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.html" />
    </navPoint>
  </navMap>
</ncx>`,
	},
	"Updating the play order works when there is a missing playOrder attribute": {
		input: `<ncx>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.html" />
    </navPoint>
    <navPoint id="navPoint-2">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.html" />
    </navPoint>
  </navMap>
</ncx>`,
		expectedOutput: `<ncx>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.html" />
    </navPoint>
    <navPoint id="navPoint-2" playOrder="2">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.html" />
    </navPoint>
  </navMap>
</ncx>`,
	},
	"Updating the play order does nothing if all playOrders are in order": {
		input: `<ncx>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.html" />
    </navPoint>
    <navPoint id="navPoint-2" playOrder="2">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.html" />
    </navPoint>
  </navMap>
</ncx>`,
		expectedOutput: `<ncx>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel><text>Chapter 1</text></navLabel>
      <content src="chapter1.html" />
    </navPoint>
    <navPoint id="navPoint-2" playOrder="2">
      <navLabel><text>Chapter 2</text></navLabel>
      <content src="chapter2.html" />
    </navPoint>
  </navMap>
</ncx>`,
	},
}

func TestFixPlayOrder(t *testing.T) {
	for name, tc := range updatePlayOrderTestCases {
		t.Run(name, func(t *testing.T) {
			edits := rulefixes.FixPlayOrder(tc.input)

			checkFinalOutputMatches(t, tc.input, tc.expectedOutput, edits...)
		})
	}
}
```

### Linting Files

After changes are made, make sure to run `make lint` to check for any lint errors. If there are any, go ahead and fix them. A lot of the linting errors are meant to help with style and performance. It should allow for focusing more on other things that matter rather than some of the more common gotchas as well as styling things relatively consistently.

### Documentation

Documentation is meant to help a user get up to speed or a developer to be more familiar with and able to work with code. In this project, all CLI tools should have a README. Comments in the code are fine, but they should not be required as that just encourages writing poor function and package comments.

Documentation can be found in 4 main places in this project:
- The project [README](README.md)
- The [Contributing Docs](CONTRIBUTING.md)
- Each CLI tool's README/README.md.tmpl
- Each CLI tool's `generate.go` where TODOs and some other specific information can be setup

As of right now there is really no other place for documentation to live in this project.

To generate updated documentation, run `make generate` at the base of the repository.

Making changes to documentation to improve things or make things clearer is welcome.

#### CLI Tool README

The CLI tool READMEs are generated by running `make generate`. As a part of this, there needs to be a `README.md.tmpl` file at the base of the repository. A bare bones example of a `README.md.tmpl` is as follows:
``` tmpl
<!-- This file is generated from  https://github.com/pjkaufman/go-go-gadgets/cat-acsii/README.md.tmpl. Please make any necessary changes there. -->

# {{ .Title }}

{{ .Description }}

{{- if .Todos }}

## TODOs

{{- range .Todos }}
- {{ . }}
{{- end}}
{{- end}}

## Commands

{{ .CommandStrings }}

```

You don't need to have all of these, but this is a sample that has the name, the description, any todos that might be present, and the commands for the tool itself.

## Use of AI

The use of AI for contributing to this repo is perfectly fine. However you are expected to make sure that the code works or the documentation is correct. You should be able to answer why you are making a change. Saying "The AI said so" is not a good enough answer for why you are making a change. You should be able to explain in your own words why the change is happening.

AI can be helpful with tasks like adding UTs, bouncing ideas off of it, doing some refactors, as well as other tasks. Use it with discernment if you use it.
