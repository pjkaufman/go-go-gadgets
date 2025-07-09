<!-- This file is generated from  https://github.com/pjkaufman/go-go-gadgets/epub-lint/README.md.tmpl. Please make any necessary changes there. -->

# Ebook Linter

This is a program that helps lint and make updates to ebooks.

## Supported File Types
- cbr
- cbz
- epub

## TODOs
- See about removing unused files and images when running epub linting

## Commands

- [cbr](#cbr)
  - [to-cbz](#to-cbz)
- [cbz](#cbz)
  - [compress](#compress)
- [epub](#epub)
  - [compress-and-lint](#compress-and-lint)
  - [fix-validation](#fix-validation)
  - [fixable](#fixable)
  - [replace-strings](#replace-strings)
  - [validate](#validate)

### cbr

Handles operations on cbr files in particular

#### to-cbz

Converts all of the cbr files to cbz files in the specified directory.

##### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| d | directory | the folder where all cbr files should be converted to cbz files | string | . | false |  |

##### Usage

``` bash
# To convert all cbrs to cbzs in a folder:
ebook-lint cbr to-cbz -d folder

# To convert all cbrs to cbzs in the current directory:
ebook-lint cbr to-cbz 
```

### cbz

Handles operations on cbz files in particular

#### compress

Compresses all of the png and jpeg files in the cbz files in the specified directory

##### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| d | directory | the location to run the cbz image compression in | string | . | false |  |
| v | verbose | whether or not to show extra information about the image compression |  | false | false |  |

##### Usage

``` bash
# To compress images in all cbzs in a folder:
ebook-lint cbz compress -d folder

# To compress images in all cbzs in the current directory:
ebook-lint cbz compress
```

### epub

Handles operations on epub files in particular

#### compress-and-lint

Gets all of the .epub files in the specified directory.
Then it lints each epub separately making sure to compress the images if specified.
Some of the things that the linting includes:
- Replacing a list of common strings
- Adds language encoding specified if it is not present already (default is "en")
- Sets encoding on content files to utf-8 to prevent errors in some readers


##### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| i | compress-images | whether or not to also compress images which requires imgp to be installed |  | false | false |  |
| d | directory | the location to run the epub lint logic | string | . | false |  |
| l | lang | the language to add to the xhtml, htm, or html files if the lang is not already specified | string | en | false |  |
|  | removable-file-types | A comma separated list of file extensions of files to remove if they are not in the manifest (i.e. '.jpeg,.jpg') | string | .jpg,.jpeg,.png,.gif,.bmp,.js,.html,.htm,.xhtml,.txt,.css | false |  |
| v | verbose | whether or not to show extra logs like what files were removed from the epub |  | false | false |  |

##### Usage

``` bash
# To compress images and make general modifications to all epubs in a folder:
epub-lint epub compress-and-lint -d folder -i

# To compress images and make general modifications to all epubs in the current directory:
epub-lint epub compress-and-lint -i

# To just make general modifications to all epubs in the current directory:
epub-lint epub compress-and-lint
```

#### fix-validation

Uses the provided epub and EPUBCheck JSON output file to fix auto fixable auto fix issues. Here is a list of all of the error codes that are currently handled:
- OPF-014: add scripted to the list of values in the properties attribute on the manifest item
- OPF-015: remove scripted to the list of values in the properties attribute on the manifest item
- NCX-001: fix discrepancy in identifier between the OPF and NCX files
- OPF-030: add the unique identifier id to the first dc:identifier element that does not have an id already
- RSC-005: seems to be a catch all error id, but the following are handled around it
  - Update ids/attributes to have valid xml ids that conform to the xml and epub spec by removing colons and any other invalid characters with an underscore
    and starting the value with an underscore instead of a number if it currently is started by a number
  - Move attribute properties to their own meta elements that refine the element they were on to fix incorrect scheme declarations or other prefixes
  - Remove empty elements that should not be empty but are empty which is typically an identifier or description that has 0 content in it
- RSC-012: try to fix broken links by removing the id link in the href attribute


##### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
|  | cleanup-jnovels | whether or not to remove JNovels info if it is present |  | false | false |  |
| f | file | the epub file to replace strings in | string |  | true | Should be a file with one of the following extensions: epub |
|  | issue-file | the path to the file with the validation issues | string |  | true | Should be a file with one of the following extensions: json |

##### Usage

``` bash
epub-lint epub fix-validation -f test.epub --issue-file epubCheckOutput.json
will read in the contents of the JSON file and try to fix any of the fixable
validation issues

epub-lint epub fix-validation -f test.epub --issue-file epubCheckOutput.json --cleanup-jnovels
will read in the contents of the JSON file and try to fix any of the fixable
validation issues as well as remove any jnovels specific files
```

#### fixable

Goes through all of the content files and runs the specified fixable actions on them asking
for user input on each value found that matches the potential fix criteria.
Potential things that can be fixed:
- Broken paragraph endings
- Section breaks being hardcoded instead of an hr
- Page breaks being hardcoded instead of an hr
- Oxford commas being missing before or's or and's
- Possible instances of sentences with two subordinate clauses (i.e. have although..., but)
- Possible instances of thoughts that are in parentheses
- Possible instances of conversation encapsulated in square brackets
- Possible instances of words in square brackets that may be necessary for the sentence (i.e. need to have the brackets removed)
- Possible instances of single quotes that should actually be double quotes (i.e. when a word is in single quotes, but is not inside of double quotes)


##### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| a | all | whether to run all of the fixable suggestions |  | false | false |  |
|  | broken-lines | whether to run the logic for getting broken line suggestions |  | false | false |  |
|  | conversation | whether to run the logic for getting conversation suggestions (paragraphs in square brackets may be instances of a conversation) |  | false | false |  |
| f | file | the epub file to find manually fixable issues in | string |  | true | Should be a file with one of the following extensions: epub |
|  | lacking-subordinate-clause | whether to run the logic for getting potentially lacking subordinate clause suggestions |  | false | false |  |
|  | log-file | the place to write debug logs to when using the TUI | string |  | false |  |
|  | necessary-words | whether to run the logic for getting necessary word suggestions (words that are a subset of paragraph content are in square brackets may be instances of necessary words for a sentence) |  | false | false |  |
|  | oxford-commas | whether to run the logic for getting oxford comma suggestions |  | false | false |  |
|  | page-breaks | whether to run the logic for getting page break suggestions (must be used with an epub with a css file) |  | false | false |  |
|  | section-breaks | whether to run the logic for getting section break suggestions (must be used with an epub with a css file) |  | false | false |  |
|  | single-quotes | whether to run the logic for getting incorrect single quote suggestions |  | false | false |  |
|  | thoughts | whether to run the logic for getting thought suggestions (words in parentheses may be instances of a person's thoughts) |  | false | false |  |
| t | use-tui | whether to use the terminal UI for suggesting fixes |  | false | false |  |

##### Usage

``` bash
# To run all of the possible potential fixes:
epub-lint epub fixable -f test.epub -a
Note: this will require a css file to already exist in the epub

# To just fix broken paragraph endings:
epub-lint epub fixable -f test.epub --broken-lines

# To just update section breaks:
epub-lint epub fixable -f test.epub --section-breaks
Note: this will require a css file to already exist in the epub

# To just update page breaks:
epub-lint epub fixable -f test.epub --page-breaks
Note: this will require a css file to already exist in the epub

# To just fix missing oxford commas:
epub-lint epub fixable -f test.epub --oxford-commas

# To just fix potentially lacking subordinate clause instances:
epub-lint epub fixable -f test.epub --lacking-subordinate-clause

# To just fix instances of thoughts in parentheses:
epub-lint epub fixable -f test.epub --thoughts

# To run a combination of options:
epub-lint epub fixable -f test.epub -oxford-commas --thoughts --necessary-words
```

#### replace-strings

Uses the provided epub and extra replace Markdown file to replace a common set of strings and any extra instances specified in the extra file replace. After all replacements are made, the original epub will be moved to a .original file and the new file will take the place of the old file. It will also print out the successful extra replacements with the number of replacements made followed by warnings for any extra strings that it tried to find and replace values for, but did not find any instances to replace.
Note: it only replaces strings in content/xhtml files listed in the opf file.

##### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| f | epub-file | the epub file to replace strings in in | string |  | true | Should be a file with one of the following extensions: epub |
| e | extra-replace-file | the path to the file with extra strings to replace | string |  | true | Should be a file with one of the following extensions: md |

##### Usage

``` bash
epub-lint epub replace-strings -f test.epub -e replacements.md
will replace the common strings and extra strings parsed out of replacements.md in content/xhtml files located in test.epub.
The original test.epub will be moved to test.epub.original and test.epub will have the updated files.

replacements.md is expected to be in the following format:
| Text to replace | Text to replace with |
| --------------- | -------------------- |
| I am typo | I the correct value |
...
| I am another issue to correct | the correction |
```

#### validate

Validates an EPUB file using W3C EPUBCheck tool.
If EPUBCheck is not installed, it will automatically download and install the latest version.

##### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| f | file | the epub file to validate | string |  | true | Should be a file with one of the following extensions: epub |
|  | json-file | specifies that the validation output should be in JSON and in the specified file | string |  | false |  |

##### Usage

``` bash
epub-lint epub validate -f test.epub
will run EPUBCheck against the file specified.
```


