# ibooks-notes-exporter

People remember less than 20% of what they read if they are not taking notes.
That's why it's always a good idea to take notes while you are reading.
But if you are reading in iBooks and taking notes there, it's very hard to get them out of that app.
You have to manually copy-paste quotes and your thoughts. ibooks-notes-exporter  solves
that problem. it's a little program that gets all your notes and highlights from iBooks and
exports them into markdown.

Note! **ibooks-notes-exporter** can only extract notes from EPUB files.


## Installation

ibooks-notes-exporter is available on OS X (both Intel and M-series processors).
It's distributed via a [homebrew](https://brew.sh/) package manager.

Run these commands in your terminal 

```shell

> brew tap 57uff3r/mac-apps
> brew install 57uff3r/mac-apps/ibooks_notes_exporter

```


## Usage

First of all, you have to get a list of all your books with notes and highlights.

```shell
❯ ibooks_notes_exporter books
+----------------------------------+-----------------+----------------------------------------------------------------------------------+
| SINGLEBOOK ID                    | NUMBER OF NOTES | TITLE AND AUTHOR                                                                 |
+----------------------------------+-----------------+----------------------------------------------------------------------------------+
| 4BAE5DA3C95788753173EAE8C63E6034 |               1 | Lorem impsum — John Doe                                                          |
| 7C3FA4F94689D97444BB4E0FD97D7197 |              54 | Hamlet — william shakespeare                                                     |
+----------------------------------+-----------------+----------------------------------------------------------------------------------+
```

And then you can export all your notes from the book into a markdown file.


```shell
❯ ibooks_notes_exporter export 4BAE5DA3C95788753173EAE8C63E6034 > ./LoremImpsum.md
```


## Feedback and contribution

Your feedback and pull requests are much appreciated.
Feel free to send your comments and thoughts to [me@akorchak.software](mailto:me@akorchak.software)


## Changelog

**0.0.4**

Fix for long titles made by @[NSBum](https://github.com/NSBum)

**0.0.3**

Bug fix: worng  order of notes and highlights. Quick and dirty fix, better solution requires to implement a parser 
for ISO/IEC 23736-6:2020 standard (EPUB Canonical Fragment Identifier or epubcfi) and this will be done in next 
versions.

**0.0.2**

Markdown fix: missing line break


**0.0.1**

Initial release


