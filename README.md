# go-notepad

Windows XP inspired notepad written in Go using GTK

![Go-Notepad](images/screenshot.png?raw=true "Go-Notepad")

## Building from Source

Prerequisits:

- Go 1.17+
- magefile: <https://magefile.org/>
- golangci-lint: <https://github.com/golangci/golangci-lint>

To build go-notepad from source:

```bash
git clone git@github.com:ajm113/go-notepad.git
cd go-notepad
mage build
```
This will create a binary in the `./build/` directory and may take a few minutes to compile depending on your system on for the first time.

You can also run go-notepad directly from mage:
```bash
mage run
```

## Config

Currently, a basic `.notepad.yml` config will look like this (more coming soon):

```yml

font:
  family: ariel
  size: 10
  wrap: false
statusbar:
  enable: true

```

Go-notepad will look for the `.notepad.yml` in the following dirs:

- `./.notepad.yml`
- `~/.notepad.yml`
- `~/config/go-notepad/notepad.yml`


## Current Features

- Simular UI layout of Win XP notepad
- Save and open files
- Word Wrap
- Status Bar
- Simple user config
- Drag & Drop!

## TODO or Citation Needed

- Improve error/dialog messages to match Win XP Notepad.
- Add Find/Replace
- Add print functionality.
- Emulate About dialog.
- Add undo/redo functionality.
- ~~Improve makefile for crossplatform support.~~
- Add GH workflows for automated deploys.
- Add app icon and about banner image.

## The MIT License (MIT)

Copyright © `2021-2024` `@ajm113`

Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the “Software”), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.
