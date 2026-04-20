<!-- This file is generated from  https://github.com/pjkaufman/go-go-gadgets/versy/README.md.tmpl. Please make any necessary changes there. -->

# Versy

This is a program that grabs the verse of the day or the specified verse(s) in two translations (either the default or the user specified ones).

## Commands

- [versy](#versy-base-command)
  - [show](#show)

### versy (base command)

A verse of the day retriever for two translations

#### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
|  | translation-a | gets the verse reference specified in this translation first (default is ESV) | string | ESV | false |  |
|  | translation-b | gets the verse reference specified in this translation second (default is NVI) | string | NVI | false |  |
| v | verbose | show more info about what is going on |  | false | false |  |

#### show

Displays the specified verse reference in the two specified Bible versions

##### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
|  | translation-a | gets the verse reference specified in this translation first (default is ESV) | string | ESV | false |  |
|  | translation-b | gets the verse reference specified in this translation second (default is NVI) | string | NVI | false |  |
| v | verbose | show more info about what is going on |  | false | false |  |
|  | verse | the Bible verse to get the two versions of | string |  | true |  |


