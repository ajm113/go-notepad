# go-notepad

Windows XP inspired notepad written in Go using GTK

## Building from Source

Prerequisits:

**Anything marked with * is required.**

- *Go 1.17+
- *GTK3 Dev Libraries Installed: <https://github.com/gotk3/gotk3>
- golangci-lint: <https://github.com/golangci/golangci-lint>

After Go and needed libraries are installed. You should be able to run `$ make` and see a executable.
If not running `go build -o notepad -tags pango_1_42,gtk_3_22 .` will result the same output if you do not have Make installed for your targeted operating system.

Once the executable is built you should be able to run it via `./notepad` or if you are on Windows ... `notepad.exe`.

## Config

Currently, a basic `.notepad.yml` config will look like this (more coming soon):

```yml

font:
  family: ariel
  size: 10
  style: regular
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

- ~~Add time and date insert.~~
- Add Font selection dialog.
- ~~Get default font size/family/type from legacy Notepad on Win XP.~~
- Improve error/dialog messages to match Win XP Notepad.
- Add Find/Replace
- ~~Add Go To Line.~~
- Add print functionality.
- ~~Drag & drop files.~~
- Emulate About dialog.
- Improve makefile for crossplatform support.
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
