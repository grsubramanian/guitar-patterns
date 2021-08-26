# Tool to generate chord and scale shapes for any uniformly fretted instrument for any musical system

## Building

```
$ go build -o bin/ cmd/main.go
```

## Use cases

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

**Can rename the 12 notes - e.g. augmented chord**

```
$ ./bin/main -s 4 -s 4 -n "R" -n "b2" -n "2" -n "b3" -n "3" -n "4" -n "#4" -n "5" -n "#5" -n "6" -n "b7" -n "7"

|----|----|----|-R--|
|-3--|----|----|----|
|-R--|----|----|----|
|----|-#5-|----|----|
|----|----|-3--|----|
|----|----|----|-R--|

...
```

**Can rename the 12 notes to a specific key - e.g. C sus 2 chord**

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
