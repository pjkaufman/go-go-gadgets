<!-- This file is generated from  https://github.com/pjkaufman/go-go-gadgets/song-converter/README.md.tmpl. Please make any necessary changes there. -->

# Song Converter

This is a program that helps converter some Markdown files with YAML frontmatter into html or csv to help with creating a song book.

## Commands

- [create-csv](#create-csv)
- [create-html](#create-html)

### create-csv

How it works:
- Reads in all of the files in the specified folder.
- Sorts the files alphabetically
- Converts each file into a CSV row
- Writes the content to the specified source


#### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| o | output-file | the file to write the csv to | string |  | false | Should be a file with one of the following extensions: csv |
| d | working-dir | the directory where the Markdown files are located | string |  | true | Should be a directory |

#### Usage

``` bash
# To write the output of converting the files in the specified folder into a csv format to a file:
song-converter create-csv -d working-dir -o churchSongs.csv

# To write the output of converting the files in the specified folder into a csv format to std out:
song-converter create-csv -d working-dir
```

### create-html

How it works:
- Reads in all of the files in the specified folder
- Sorts the files alphabetically
- Adds the cover to the start of the content after converting it to html
- Converts each file into html
- Writes the content to the specified source


#### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| c | cover-file | the markdown cover file to use | string |  | true | Should be a file with one of the following extensions: md |
| o | output | the html file to write the output to | string |  | false | Should be a file with one of the following extensions: html |
| v | version-type | the version descriptor for the type of songs to generate (generally just abridged or unabridged) | string |  | true |  |
| d | working-dir | the directory where the Markdown files are located | string |  | true | Should be a directory |

#### Usage

``` bash
# To write the output of converting the files in the specified folder to html to a file:
song-converter create-html -d working-dir -c cover.md -o songs.html

# To write the output of converting the files in the specified folder to html to std out:
song-converter create-html -d working-dir -s cover.md
```


