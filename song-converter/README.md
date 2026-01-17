<!-- This file is generated from  https://github.com/pjkaufman/go-go-gadgets/song-converter/README.md.tmpl. Please make any necessary changes there. -->

# Song Converter

This is a program that helps converter some Markdown files with YAML frontmatter into html or csv to help with creating a song book.

## Commands

- [compare](#compare)
- [create](#create)
  - [book](#book)
  - [csv](#csv)
  - [html](#html)

### compare

Compares the provided html and pdf file to see if there are any potentially meaningful difference like linebreaks and whitespace differences

#### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| f | file | the pdf file to compare with the html file | string |  | true | Should be a file with one of the following extensions: pdf |
|  | join-lines | the number of lines at the start of the pdf to join together to help make the html and pdf content as similar as possible | int | 0 | false |  |
| s | source | the html file that was used to generate the pdf file | string |  | true | Should be a file with one of the following extensions: html |

#### Usage

``` bash
# To compare a pdf and its html source:
song-converter compare -s songs.html -f songs.pdf

# To compare a pdf and its html source where the first several lines of text are meant to be the heading on a single line:
song-converter compare -s songs.html -f songs.pdf --join-lines 4
```

### create

Deals with creating files from the song Markdown files

#### book



##### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| c | cover-file | the markdown cover file to use | string |  | true | Should be a file with one of the following extensions: md |
| l | location | the specific book to recreate by filtering songs down to just that book location | string |  | true |  |
| o | output | the html file to write the output to | string |  | false | Should be a file with one of the following extensions: html |
| d | working-dir | the directory where the Markdown files are located | string |  | true | Should be a directory |

#### csv

How it works:
- Reads in all of the files in the specified folder.
- Sorts the files alphabetically
- Converts each file into a CSV row
- Writes the content to the specified source


##### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| o | output-file | the file to write the csv to | string |  | false | Should be a file with one of the following extensions: csv |
| d | working-dir | the directory where the Markdown files are located | string |  | true | Should be a directory |

##### Usage

``` bash
# To write the output of converting the files in the specified folder into a csv format to a file:
song-converter create csv -d working-dir -o churchSongs.csv

# To write the output of converting the files in the specified folder into a csv format to std out:
song-converter create csv -d working-dir
```

#### html

How it works:
- Reads in all of the files in the specified folder
- Sorts the files alphabetically
- Adds the cover to the start of the content after converting it to html
- Converts each file into html
- Writes the content to the specified source


##### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| c | cover-file | the markdown cover file to use | string |  | true | Should be a file with one of the following extensions: md |
| o | output | the html file to write the output to | string |  | false | Should be a file with one of the following extensions: html |
| v | version-type | the version descriptor for the type of songs to generate (generally just abridged or unabridged) | string |  | true |  |
| d | working-dir | the directory where the Markdown files are located | string |  | true | Should be a directory |

##### Usage

``` bash
# To write the output of converting the files in the specified folder to html to a file:
song-converter create html -d working-dir -c cover.md -o songs.html

# To write the output of converting the files in the specified folder to html to std out:
song-converter create html -d working-dir -s cover.md
```


