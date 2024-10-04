<!-- This file is generated from  https://github.com/pjkaufman/go-go-gadgets/jp-proc/README.md.tmpl. Please make any necessary changes there. -->

# Jpeg and Png Processor

This is meant to be a replacement for my usage of imgp.

Currently I use imgp for the following things:
- image resizing
- exif data removal
- image quality setting

Given how this works, I find it easier to just go ahead and do a simple program in Go to see how things stack up and not be so reliant on Python. This also helps me learn some more about imaging processing as well. So a win-win in my book.

## How does this program compare with imgp?

| Operation | Original Size | New Size (imgp) | New Size (imgp with optimize flag) | New Size (jp-proc) |
| --------- | ------------- | --------------- | ---------------------------------- | ------------------ |
| Resize jpeg to 800x600 and remove exif data | 3.4M | 57KB | 56KB | 68KB |
| Resize jpeg to 800x600 and remove exif data and set quality to 40 | 3.4M | 32KB | 28KB | 37KB |

## Commands

- [proc](#proc)

### proc

Processes the provided image in the specified ways

#### Flags

| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |
| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |
| f | file | the image file to operate on | string |  | true | Should be a file with one of the following extensions: png, jpg, jpeg |
| m | mute | whether or not to keep from printing out values to standard out |  | false | false |  |
| o | overwrite | whether or not to overwrite the original file when done |  | false | false |  |
| q | quality | the quality of the jpeg to use when encoding the image (default is 75) | int | 75 | false |  |
| e | remove-exif | whether or not to remove exif data from the image |  | false | false |  |
| w | width | the width of the image to use when the image is resized (leave blank to keep original) | int | 0 | false |  |


