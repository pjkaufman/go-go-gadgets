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

## Todos

- Resize png test
