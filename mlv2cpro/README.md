# mlv2cpro

Batch convert Magic Lantern RAW "MLV" files to Quicktime-wrapped
CineformRAW files.

## Usage

```
Usage of ./mlv2cpro:
  -extension="mov": File extension
  -keepfiles=true: Keep source files after transcoding
  -mlvdump="./mlv_dump": Path to mlv_dump binary
  -outdir=".": Output directory
  -raw2gpcf="./raw2gpcf": Path to raw2gpcf binary
  -threading=false: Use multi-threading
```

If no additional arguments are given, the current working directory will be
processed, otherwise all listed directories will be scanned and all .MLV
files will be processed.

The threading option allows parallelism (more than one file to be encoded
at once). It is disabled by default.

## Dependencies

 * [Go](http://golang.org) - for compilation
 * [raw2gpcf](http://www.magiclantern.fm/forum/index.php?topic=5479.msg41378#msg41378) - GoPro Studio Professional component for transcoding to CineformRAW from the command line.
 * [mlv_dump](http://www.magiclantern.fm/forum/index.php?topic=7122.msg88389#msg88389) - Magic Lantern component for converting between MLV and RAW (among other formats)

## Building

```
go build
```

To cross-compile, make sure that you have installed and setup
``golang-crosscompile`` first. Then run

```
make
```

