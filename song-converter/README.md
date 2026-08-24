<!-- This file is generated from  https://github.com/pjkaufman/go-go-gadgets/song-converter/README.md.tmpl. Please make any necessary changes there. -->

# Song Converter

This is a program that helps converter some Markdown files with YAML frontmatter into html or csv to help with creating a song book.

## YAML Song Metadata

| Property | Usage | Used in Digital PDF (Y/N) | Used in Digital Book PDF (Y/N) |
| -------- | ----- | ------------------------- | ------------------------------ |
| `melody` | The melody that the song is to be sung or played to. | Y | Y |
| `key` | The key or keys the song can be played in.  | Y | Y |
| `authors` | The author or authors of the song. | Y | Y |
| `in-church` | Whether or not the author is a part of GMI. If they are, their name is bolded in the file. | Y | Y |
| `verse` | The Bible verse or verses referenced or used as inspiration for the song. | Y | Y |
| `location` | All the places where the song is found in the source books which can be multiple times in a single book. | Y | N |
| `copyright` | Who owns the copyright to the song if it is not something from someone in GMI. | N | N |
| `book-title` | The title of the song to use for the song instead of the song's title in the markdown. _Note: if there is an alternate title it may only be used for the first entry of the song in the Table of Contents._ | N | Y |
| `alternate-title` | The other name by which the song is known which should be used in the Table of Contents. This is typically used when you have two names for a song that shows up one time, one time in a primary and secondary location, or multiple times in the same location where you want the title to differ for entries or just be shown under both names if it only has one location. | N | Y |
| `no-break` | Whether or not the book version of the song should skip the line break after the song if it doesn't end a page. This is meant to help with keeping a set of songs on the same page. | N | Y |
| `skip-book` | Whether or not to skip the song when dealing with book PDF generation. This is meant for situations where a song exists and is the equivalent of something in the book, but should not actually be used. | N | Y |

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
|  | ignore-page-numbers | whether to ignore table of contents page numbers (this is for when the HTML or PDF will not have line numbers in the table of contents, but the other will) |  | false | false |  |
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
|  | secondary-location | a second book to include in the table of contents for the book to create | string |  | false |  |
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
| o | output | the file to write the csv to | string |  | false | Should be a file with one of the following extensions: csv |
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
|  | format | the version descriptor for the type of songs to generate (Abridged or Unabridged) | string |  | true | Should be a one of the following: Abridged, Unabridged |
| o | output | the html file to write the output to | string |  | false | Should be a file with one of the following extensions: html |
| d | working-dir | the directory where the Markdown files are located | string |  | true | Should be a directory |

##### Usage

``` bash
# To write the output of converting the files in the specified folder to html to a file:
song-converter create html -d working-dir -c cover.md -o songs.html

# To write the output of converting the files in the specified folder to html to std out:
song-converter create html -d working-dir -c cover.md
```


