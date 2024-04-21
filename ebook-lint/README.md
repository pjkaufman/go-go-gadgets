# Ebook Lint

This is a program that helps lint and make updates to ebooks.

## Supported File Types

- Epub
- Cbr
- Cbz

## TODOs

- See about removing unused files and images when running epub linting
- See about removing requirement of one reference per file in the opf/ncx/toc
- See about validating that all links exist (but that could be handled potentially by the epub checker)
- See about having a portion of logic that uses the epub checker logic that allows for checking for issues in epubs
  - This would need to either be optional or only fail on errors as they are the only ones that I believe should halt loading an epub 
