# Go Go Gadgets

This repository is full of CLI tools that I find useful. They are meant to help me with various tasks that I do. Some are more useful than others. But they are all gadgets that I have put together here in this repository. 

## Installation

In general on a device you will want to run the following in order to install the CLI tools:

``` bash
make install
```

This will install the CLI programs and their completions.

### Termux (Mobile)

To install on mobile in Termux run the following command which only installs `epub-lint` at this time:

``` bash
make install-termux
```

## Uninstalling

If you need to uninstall the packages, the simplest way is to run the following make rule:

```bash
make clean
```

## Testing

To run the unit tests, you can go ahead and run the following:

``` bash
make test
```

This will run all of the unit tests and let you know if anything is broken.

## Documentation

The main documentation for the gadgets and this repository live in this README, each gadget's README, and the [CONTRIBUTING.md](CONTRIBUTING.md). A lot of the Gadget README files are generated, so keep that in mind if a change is needed.

## Available Gadgets

### [Epub Linter](./epub-lint/README.md)

A versatile ebook management tool that helps maintain and improve your epubs.

#### Key Features

- Compresses images in EPUB files
- Comprehensive EPUB linting and formatting
- String replacement and content fixes
- EPUB validation via W3C's EpubCheck program
- Moving author's notes to their own file at the end of the EPUB

#### Use Cases

Perfect for maintaining an ebook collection, preparing files for e-readers, and ensuring consistent formatting.

### [Song Converter](./song-converter/README.md)

A specialized tool for converting song collections between different formats to create songbooks.

#### Key Features

- Converts Markdown files with YAML frontmatter to HTML or CSV
- Alphabetical sorting of songs
- Support for cover pages
- Batch processing of multiple files
- Flexible output options (file or stdout)
- Version control for different song variants (abridged/unabridged)

#### Use Cases

Perfect for creating church songbooks, music collections, or any organized compilation of songs that needs to be formatted consistently.

### [Magnum](./magnum/README.md)

A light novel release tracker that monitors and manages information about book series from various publishers.

#### Key Features

- Tracks releases from multiple publishers (Yen Press, JNovel Club, Seven Seas Entertainment, etc.)
- Maintains series status (Ongoing/Completed/On Hold)
- Automated release date tracking (has to be manually triggered to look for new releases)
- Series management (add/remove/list)
- Customizable series information
- Shows upcoming releases in chronological order

#### Use Cases

Ideal for light novel enthusiasts who want to stay updated on release dates and manage their reading lists.

### [JP Processor](./jp-proc/README.md)

An image processing utility focused on JPEG and PNG optimization.

#### Key Features

- Image resizing capabilities
- EXIF data removal
- Quality adjustment for JPEG files
- File overwrite protection
- Comparable performance to imgp

#### Use Cases

Useful for batch processing images, preparing photos for web use, or cleaning up image metadata.

### [Cat ASCII](./cat-ascii/README.md)

A tool to display nice little cat ASCII art.

### [Versy](./versy/README.md)

A tool for getting the verse of the day and specified verses in two translations (often one being in English and the other being in Spanish).

## Contributing

As with most, if not all, software there are always improvements that can be made. If you believe that something can be improved or you would like to contribute, look at [the contributing docs](CONTRIBUTING.md).

## Key Dependencies

### Core Libraries
- [cobra](https://github.com/spf13/cobra) - CLI application framework
- [pflag](https://github.com/spf13/pflag) - Command line flag parsing
- [colly](https://github.com/gocolly/colly) - Web scraping framework

### User Interface
- [bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) - TUI styling and terminal color
- [promptui](https://github.com/manifoldco/promptui) - Interactive prompt UI

### File Processing
- [markdown](https://github.com/gomarkdown/markdown) - Markdown processing
- [frontmatter](https://github.com/adrg/frontmatter) - YAML frontmatter parsing
- [go-exif](https://github.com/dsoprea/go-exif) - EXIF data manipulation
- [jpegquality](https://github.com/liut/jpegquality) - JPEG quality management
- [go-jpeg-image-structure](https://github.com/dsoprea/go-jpeg-image-structure) - JPEG processing
- [go-png-image-structure](https://github.com/dsoprea/go-png-image-structure) - PNG processing

### Testing and Development
- [testify](https://github.com/stretchr/testify) - Testing toolkit
- [heredoc](https://github.com/MakeNowJust/heredoc) - Multiline string literals
- [diff](https://github.com/andreyvit/diff) - Text difference comparison

### System Integration
- [clipboard](https://github.com/atotto/clipboard) - Clipboard interaction
- [consolesize-go](https://github.com/nathan-fiscaletti/consolesize-go) - Console dimensions
