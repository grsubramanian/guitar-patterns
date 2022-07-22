# Tool to generate chord and scale shapes for any uniformly fretted instrument for any musical system

<img src="./resources/minor_chord_example.svg">

## Building the program

1. If not already done, install git on your system. See https://git-scm.com/book/en/v2/Getting-Started-Installing-Git.

2. If not already done, install the Go programming language on your system. See https://golang.org/doc/install.

3. If not already done, install the `goimports` tool by running `go install golang.org/x/tools/cmd/goimports@latest`

4. If not already done, install the `golint` tool by running `go install golang.org/x/lint/golint@latest`.

5. If not already done, install `make`. The instructions can vary by OS and distribution.

6. Clone this repository. See https://docs.github.com/en/github/creating-cloning-and-archiving-repositories/cloning-a-repository-from-github/cloning-a-repository.

7. Go into the cloned directory, for example in Unix / Linux systems by running `cd guitar-patterns`.

8. Build the program by running `make`. This will generate the binary at `./bin/main`.

## Running the program

### Basics

To run the program, just type `./bin/main` with the appropriate command line arguments as will be shown below. This needs to be done from the command line terminal. If you have only used graphical interfaces in the past, this might sound daunting, but it is easy. There are several resources online.

### Available printers

The program will print output to STDOUT. In simple words, this means that the output will print to the screen.

There are two output options available.

 * ASCII (default) - lightweight for those who prefer to stay on the terminal.
 * SVG (when run with `-svg`) - convenient because the SVG output can be converted to JPG or PNG using free converter tools available online.

If you want the output to go to a file (which you will likely need when using the SVG option), you will need to redirect the output to the file such as `./bin/main -svg {some args} > somefileonthecomputer.svg`.

### Accessibility options

Run with the `-left` argument if you are left handed.

### Use cases

*All examples shown here print in ASCII and use right-handedness.*

**Simple chord - e.g. minor chord**

There are 3 steps between the root and the next note, and then 4 steps from that one to the next.

```
$ ./bin/main -s 3 -s 4

|----|----|----|-1--|
|----|----|----|-5--|
|-1--|----|----|-b3-|
|-5--|----|----|----|
|----|-b3-|----|----|
|----|----|----|-1--|

...

```

**Chords can have more than 3 notes - e.g. dominant 7th chord**

```
$ ./bin/main -s 4 -s 3 -s 3

|----|-b7-|----|-1--|
|-3--|----|----|-5--|
|-1--|----|----|----|
|-5--|----|----|-b7-|
|----|----|-3--|----|
|----|-b7-|----|-1--|

...
```

**Chords can wrap around - e.g. major 7th add 9 chord**

```
$ ./bin/main -s 4 -s 3 -s 4 -s 3

|----|----|-7--|-1--|
|-3--|----|----|-5--|
|-1--|----|-2--|----|
|-5--|----|----|----|
|-2--|----|-3--|----|
|----|----|-7--|-1--|

...
```

**Can handle scales - e.g. harmonic minor scale**

```
$ ./bin/main -s 2 -s 1 -s 2 -s 2 -s 1 -s 3

|----|----|-7--|-1--|
|----|-4--|----|-5--|
|-1--|----|-2--|-b3-|
|-5--|-b6-|----|----|
|-2--|-b3-|----|-4--|
|----|----|-7--|-1--|

...
```

**Can change the number of frets per pattern - e.g. minor pentatonic scale with 5 frets per pattern**

```
$ ./bin/main -s 3 -s 2 -s 2 -s 3 -frets 5

|----|----|-b7-|----|-1--|
|-b3-|----|-4--|----|-5--|
|----|-1--|----|----|-b3-|
|----|-5--|----|----|-b7-|
|----|----|-b3-|----|-4--|
|----|----|-b7-|----|-1--|

...
```

**Can view how the pattern looks across all 12 frets e.g. sus4 chord sweep**

```
$ ./bin/main -s 5 -s 2 -frets 12

|----|----|----|-4--|----|-5--|----|----|----|----|-1--|----|
|----|----|----|-1--|----|----|----|----|-4--|----|-5--|----|
|-4--|----|-5--|----|----|----|----|-1--|----|----|----|----|
|-1--|----|----|----|----|-4--|----|-5--|----|----|----|----|
|-5--|----|----|----|----|-1--|----|----|----|----|-4--|----|
|----|----|----|-4--|----|-5--|----|----|----|----|-1--|----|

...
```

**Can change guitar tuning - e.g. diminished chord with open tuning**

```
$ ./bin/main -s 3 -s 3 -ss 7 -ss 5 -ss 4 -ss 3 -ss 5

|----|-1--|----|----|
|-#4-|----|----|----|
|-b3-|----|----|-#4-|
|----|-1--|----|----|
|-#4-|----|----|----|
|----|-1--|----|----|

...
```

**Can extend to any uniformly fretted string instrument with any number of strings - e.g. major 6th chord on 5 string banjo**

```
$ ./bin/main -s 4 -s 3 -s 2 -ss 7 -ss 5 -ss 4 -ss 3

...

|-6--|----|----|-1--|
|----|-5--|----|-6--|
|----|----|-3--|----|
|-6--|----|----|-1--|
|----|----|-3--|----|

...
```

**Can rename the 12 notes to suit typical western notation - e.g. C sus 2 chord**

```
$ ./bin/main -s 2 -s 5 -n "C" -n "Db" -n "D" -n "Eb" -n "E" -n "F" -n "F#" -n "G" -n "Ab" -n "A" -n "Bb" -n "B"

|----|----|----|-C--|
|----|----|----|-G--|
|-C--|----|-D--|----|
|-G--|----|----|----|
|-D--|----|----|----|
|----|----|----|-C--|

...
```

**Can pick a specific key - e.g. F augmented chord**

```
$ ./bin/main -r 5 -s 4 -s 4 -n "C" -n "Db" -n "D" -n "Eb" -n "E" -n "F" -n "F#" -n "G" -n "Ab" -n "A" -n "Bb" -n "B"

|----|----|----|-F--|
|-A--|----|----|----|
|-F--|----|----|----|
|----|-Db-|----|----|
|----|----|-A--|----|
|----|----|----|-F--|

...
```

**Can work with any music system - e.g. hypothetical music system with only 7 notes**

```
$ ./bin/main -s 2 -s 3 -n "R" -n "Sh" -n "Tr" -n "Qw" -n "Qb" -n "Si" -n "Se"

|-R--|----|-Tr-|----|
|-Tr-|----|----|-Si-|
|-Si-|----|-R--|----|
|-R--|----|-Tr-|----|
|-Tr-|----|----|-Si-|
|----|-Si-|----|-R--|

...
```
