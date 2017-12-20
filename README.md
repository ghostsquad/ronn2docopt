# ronn2docopt
Start with [ronn man page format](https://github.com/rtomayko/ronn), and parse out usage string for docopt

## Quick Start Example

1. `go run ronn2docopt.go ./examples/basic/docs/thingy.1.ronn`

This prints to stdout the [docopt usage string](http://docopt.org/)

From here, you can embed that usage string as an example by using [go-bindata](https://github.com/shuLhan/go-bindata):

```
go run main.go ./examples/basic/docs/thingy.1.ronn > ./examples/basic/docs/docopt.txt
go-bindata --prefix "$(pwd)/examples/basic" -o ./examples/basic/lib/bindata.go ./examples/basic/docs/docopt.txt
go run ./examples/basic/thingy.go --help
```

## Usage

### Rules

Basic Ronn/Ronn2Docopt Markdown rules. See [ronn format](https://rtomayko.github.io/ronn/ronn-format.7https://rtomayko.github.io/ronn/ronn-format.7) for ronn basics.

1. Ronn2Docopt only cares about the `## SYNOPSIS` and `## OPTIONS` sections:

   Sections start/terminate by H2 headers (`## Foo`) or end of file
   
2. The Synopsis section get's stripped of the backtick (`` ` ``) and `<br>`, but otherwise, is used verbatim.

     **Ronn Source**
     ```
     `naval_fate` `ship new <name>...`<br>
     `naval_fate` `ship <name> move <x> <y> [--speed=<kn>]`<br>
     ` naval_fate` `ship shoot <x> <y>`<br>
     `naval_fate` `mine (set|remove) <x> <y> [--moored|--drifting]`<br>
     `naval_fate` `-h | --help`<br>
     `naval_fate` `--version`<br>
     ```
    
     **Docopt Output**

     ```
     naval_fate ship new <name>...
     naval_fate ship <name> move <x> <y> [--speed=<kn>]
     naval_fate ship shoot <x> <y>
     naval_fate mine (set|remove) <x> <y> [--moored|--drifting]
     naval_fate -h | --help
     naval_fate --version
     ```

3. Options get "minified". First, there's a hard requirement that option lines start with 2 spaces and `*` and end with `:`
   They are stripped of the `*`, `:` and backticks `` ` ``.

    **Ronn Source**
    ```
      * `-speed=<kn>`:
    ```

    **Docopt Output**
    ```
      -speed=<kn>
    ```

4. Options terminate by other options, sections, "option sections", or the end of the file.
   An option section is simply text on a line that does not begin with spaces.
   Option sections are saved, and will be used in docopt output. Option sections can span multiple lines.

     **Ronn Source**
     ```
     Basic Options:

       * `--speed=<kn>`:

     Advanced Options:

       * `--acceleration=<kn_per_sec>`:
     ```

     **Docopt Output**
     ```
     Basic Options:

       --speed=<kn>

     Advanced Options:

       --acceleration=<kn_per_sec>
     ```

5. An options description is the first "sentence" following an option declaration. This short description is used in docopt. A "sentence" terminates by a period. If a period is not found on the first line, only the first line is used as the short description.

     **Ronn Source**
     ```
       * `--speed=<kn>`:
         Speed of your vessel. This determines how fast you'll get somewhere.
     ```

     **Docopt Output**

     ```
     --speed=<kn> Speed of your vessel
     ```

6. Within the option description, you can provide a default, which will be used by docopt. The default can occur anywhere in the description.

     **Ronn Source**
     ```
       * `--speed=<kn>`:
         Speed of your vessel. This determines how fast you'll get somewhere. [default: 30]
     ```

     **Docopt Output**

     ```
     --speed=<kn> Speed of your vessel [default: 30]
     ```

7. Option declarations (within an "option section") are compared and padded to provide nicely formatted docopt output.

     **Ronn Source**

     ```
       * `--speed=<kn>`:
         Speed of your vessel. This determines how fast you'll get somewhere.

       * `--acceleration=<kn_per_sec>`
         How quickly you will get up to speed. High acceleration is good, but usually requires other tradoffs.
     ```

     **Docopt Output**

     ```
     --speed=<kn>                 Speed of your vessel
     --acceleration=<kn_per_sec>  How quickly you will get up to speed.
     ```

That's It!

### Choosing between using manpages or docopt usage.

The various docopt implementations have a `help` argument ([python API](https://github.com/docopt/docopt#api)), that when set to false, will cause docopt to not automatically print help information and exit.
Instead it is up to your program to decide what to print. What you can do is include a verbose flag `-v`, `--verbose`, and specify that:

`--help` prints short/simple docopt usage and
`--help --verbose` drops you into a manpage within a pager (like `less`).

## Contributing

Make sure you have [glide](https://github.com/Masterminds/glide) installed.

```
glide install
./bin/test # or go test
```
