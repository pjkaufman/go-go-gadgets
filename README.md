# Go Go Gadgets

This is a set of cli tools that I find useful for me. It helps me get things done more easily than I otherwise would
be able to.

## Installation

In general on a device you will want to run the following in order to install the CLI tools:

``` bash
make install
```

This will install the CLI programs and their completions.

### Termux (Mobile)

To install on mobile in Termux run the following command:

``` bash
make install-termux
```

As of right now, that just installs `ebook-lint` at this time.

## Available Gadgets

### [Ebook Linter](./ebook-lint/README.md)

A versatile ebook management tool that helps maintain and improve your digital library.

#### Key Features

- Supports CBR, CBZ, and EPUB formats
- Converts CBR to CBZ files
- Compresses images in CBZ and EPUB files
- Comprehensive EPUB linting and formatting
- String replacement and content fixes
- EPUB validation via W3C's Epubcheck program

#### Use Cases

Perfect for maintaining an ebook collection, preparing files for e-readers, and ensuring consistent formatting.

### Git Helper

Helps simplify some git commands and actions I normally forget.

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

### Cat ASCII

A tool to display nice little cat ASCII art.

## Key Dependencies

### Core Libraries
- [cobra](https://github.com/spf13/cobra) - CLI application framework
- [pflag](https://github.com/spf13/pflag) - Command line flag parsing
- [promptui](https://github.com/manifoldco/promptui) - Interactive prompt UI
- [colly](https://github.com/gocolly/colly) - Web scraping framework

### User Interface
- [bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) - TUI styling
- [color](https://github.com/fatih/color) - Terminal color output

### File Processing
- [archiver](https://github.com/mholt/archiver) - Archive file handling
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
- [dbus](https://github.com/godbus/dbus) - D-Bus interface
- [consolesize-go](https://github.com/nathan-fiscaletti/consolesize-go) - Console dimensions
