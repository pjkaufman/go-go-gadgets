<!-- This file is generated from  https://github.com/pjkaufman/go-go-gadgets/magnum/README.md.tmpl. Please make any necessary changes there. -->

# Magnum

Magnum is a program that checks if the list of specified light novels has any updates and notes the release dates of any new entries.

## Supported Publishers
- Yen Press
- JNovel Club
- Seven Seas Entertainment (uses Google Cache)
- One Peace Books (uses Wikipedia)
- Viz Media
- Hanashi Media (uses Wikipedia)

## Light Novels to Account for
- Daily Life of the Immortal King - Novel Updates?
- Eighth Son - Novel Updates

## TODOs
- Add more unit tests and validation for commands and parsing logic to make sure it works as intended and is easier to refactor down the road since breaking changes should be easier to catch

## Commands

- [add](#add)
- [get-info](#get-info)
- [list](#list)
- [remove](#remove)
- [set-status](#set-status)
- [show-info](#show-info)

### add

Adds the provided series info to the list of series to keep track of

#### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| n | name | the name of the series | string |  | true |  |
| p | publisher | the publisher of the series | string |  | false |  |
| s | slug | the slug for the series to use instead of the one based on the series name | string |  | false |  |
| r | status | the status of the series (defaults to Ongoing) | string | O | false |  |
| t | type | the series type | string |  | false |  |
| o | wikipedia-table-parse-override | the amount of tables that should parsed in the light novels section of the wikipedia page if it should not be all of them | int | 0 | false |  |

#### Usage

``` bash
# To add a series with just a name and other information to be filled out:
magnum add -n "Lady and the Tramp"
Note: that the other fields will be filled in via prompts except the series status which is assumed to be ongoing

# To add a series with a special URL slug that does not follow the normal pattern for the publisher in question or is on its own page:
magnum add -n "Re:ZERO -Starting Life in Another World" -s "re-starting-life-in-another-world"

# To add a series that is not ongoing (for example Completed):
magnum add -n "Demon Slayer" -r "C"
```

### get-info

Gets the book release info for books that have been added to the list of series to track

#### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| c | include-completed | get info for completed series |  | false | false |  |
| s | series | get info for just the specified series | string |  | false |  |
| v | verbose | show more info about what is going on |  | false | false |  |

#### Usage

``` bash
# To get all of the release data for non-completed series:
magnum get-info
```

### list

Lists the names of each of the series that is currently being tracked

#### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| p | publisher | show series with the specified publisher | string |  | false |  |
| r | status | show series with the specified status | string |  | false |  |
| t | type | show series with the specified type | string |  | false |  |
| v | verbose | show the publisher and other info about the series |  | false | false |  |

#### Usage

``` bash
# To show a list of all series names that are being tracked:
magnum list

# To include information like publisher, status, series, etc.:
magnum list -v
```

### remove

Removes the provided series from the list of series to keep track of

#### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| n | name | the name of the series | string |  | true |  |

#### Usage

``` bash
# To remove a series use the following command:
magnum remove -n "Lady and the Tramp"
```

### set-status

Sets the status of the provided/selected book name

#### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| c | include-completed | include completed series in the books to search |  | false | false |  |
| n | name | name of the book to set the status for | string |  | false |  |
| s | status | status to set for the selected book (O/H/C) | string |  | false |  |

#### Usage

``` bash
# To set the status of a book you know the name of:
magnum set-status -n "book_name"
This will result in being prompted for a status for that book.

# To set the status of a book you know the name and status of:
magnum set-status -n "book_name" -s C

# To set the status of a book by using the cli selection options:
magnum set-status

# To set the status of a book and include the completed series:
magnum set-status -c
```

### show-info

Shows each series that has upcoming releases along with when the releases are in the order they are going to be released

#### Usage

``` bash
# To show upcoming releases in order of when they are releasing:
magnum show-info
```


