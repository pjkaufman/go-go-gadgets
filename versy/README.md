<!-- This file is generated from  https://github.com/pjkaufman/go-go-gadgets/versy/README.md.tmpl. Please make any necessary changes there. -->

# Versy

This is a program that grabs the verse of the day or the specified verse(s) in two translations (either the default or the user specified ones).

## Commands

- [show](#show)

### show

Displays the specified verse reference in the two specified Bible versions

#### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| v | verbose | show more info about what is going on |  | false | false |  |
|  | verse | the Bible verse to get the two versions of | string |  | true |  |
|  | version-one | gets the first instance of the verse in the specified version (default is ESV) | string | ESV | false |  |
|  | version-two | gets the second instance of the verse in the specified version (default is NVI) | string | NVI | false |  |


