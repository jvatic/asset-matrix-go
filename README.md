# asset-matrix-go

**NOTE: This is currently in the initial stages of development and is not ready for production use.**

If you're tired of those tedious asset pipelines destroying your mood in the morning, fear not, help is here.

## Features

Compile, concatenate, and compress your assets with ease.

### Concatenation

Use comment directives to concatenate files. Assuming `//` is the syntax for a single-line comment, the following directives are available:

Directive                             | Description
------------------------------------- | -----------
`//= require ./path/to/file`          | Insert compatible file at path relative to the current file (path may omit the file extension).
`//= require path/to/file`            | Insert compatible file at path relative to any asset root (path may omit the file extension).
`//= require http://example.com/file` | Insert compatible file at given URL.
`//= require_tree ./path/to/dir`      | Insert all compatible files in dir at path relative to the current file's dir.
`//= require_tree path/to/dir`        | Insert all compatible files in dir at path relative to any asset root.
`//= require_self`                    | Insert the current file contents at given position relative to the other directives. The default is at the end.

All directives

- must appear at the top of a file (excluding whitespace)
- must follow the syntax `single-line-comment-sequence + = + directive-name + directive-param`
- raise an error when pointing to an incompatible file (e.g. attempting to concatenate HTML with JavaScript when there is no handler to convert HTML into JavaScript or vice versa).
- emit a warning when pointing to multiple files where a subset are incompatible (the previous point applies when none are).

Before any handlers are run, all directives are extracted from files allowing handlers to manipulate data without knowing about concatenation. Once all files to be concatenated share a common output format, they are joined into one.

### Asset Handler

Handlers may be used for compiling assets (e.g. Markdown -> HTML), crating mutations of the original (e.g. a handler for gzip compression may output the original plus the compressed version), or compiling a list of all assets they handled (e.g. an asset manifest.json).

All handlers

- must be deterministic (i.e. have a consistent output for any given input) as the output of a handler is cached with the digest of the input as the lookup key, ensuring handlers are only run for mutated inputs.
- must register for input format patterns (e.g. `md` and `markdown` to handle multiple markdown file extensions, or `*` to handle everything)
- may register output formats (e.g. `html` and `json`)
- may have dynamic output formats (e.g. `{input_format}.gz` or `min.{input_format}`)
- may register to take simultaneous inputs (all files matching the input pattern, e.g. manifest handler)
- may access metadata for the input (e.g. digest, file name/path/extensions, size in bytes, etc.)
- may append suffixes to the file name (joined with `-` before writing to disk)
- may use any helper (e.g. relative path to another asset)

Format patterns are plain text with the exception of the `*` wildcard.

Metadata available to handlers:

- digest
- file name
- file name suffixes
- file output name (name with suffixes applied)
- file path
- file output path (path with name suffixes applied)
- file extension
- file size in bytes

Helpers available to handlers:

Name   | Description
----   | -----------
URL    | Absolute URL of an asset (requires the base URL be configured)
Path   | Relative URL of an asset (no configuration required)
Base64 | Base64 encoded asset (e.g. embedding a font or image in css)
Data   | Asset data (e.g. inserting a template partial)

#### Core Handlers

Some handlers are included as part of the core package as they are commonly used:

Name                   | Input(s)        | Output(s)                                               | Description
---------------------- | --------------- | ------------------------------------------------------- | -----------
Manifest               | `*`[]           | {verbatim input formats}[], `manifest.json`             | Manifest generator.
GZip                   | `*`             | `{verbatim input format}`, `{verbatim input format}.gz` | GZip compression.
JavaScriptMinification | `js`            | `js`, `min.js`                                          | JavaScript minification.
CSSMinification        | `css`           | `css`, `min.css`                                        | CSS minification.
CoffeeScript           | `coffee`        | `js`                                                    | CoffeeScript compiler.
SASS                   | `sass`, `scss`  | `css`                                                   | SASS/SCSS compiler.

It would be nice to also have image compression as a core handler.

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## License

This software is provided under a BSD license as described in the [LICENSE file](https://github.com/tent/asset-matrix-go/blob/master/LICENSE).
