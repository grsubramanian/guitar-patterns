# Chord and scale shapes for any uniformly fretted instrument

For diagrams that look like this:

```
|----|-R--|----|----|
|----|-5--|----|----|
|----|----|-3--|----|
|----|----|----|-R--|
|-3--|----|----|-5--|
|----|-R--|----|----|
```

## Building

`go build -o bin/ cmd/main.go`

## Use cases

**Simple chord - e.g. minor chord**

There are 3 steps between the root and the next note, and then 4 steps from that one to the next.

`./bin/main -s 3 -s 4`

**Chords can have more than 3 notes - e.g. dominant 7th chord**

`./bin/main -s 4 -s 3 -s 3`

**Chords can wrap around - e.g. major 7th add 9 chord**

`./bin/main -s 4 -s 3 -s 4 -s 3`

**Can handle scales - e.g. harmonic minor scale**

`./bin/main -s 2 -s 1 -s 2 -s 2 -s 1 -s 3`

**Can change the number of frets per pattern - e.g. minor pentatonic scale with 5 frets per pattern**

`./bin/main -s 3 -s 2 -s 2 -s 3 -frets 5`

**Can change guitar tuning - e.g. diminished chord with open tuning**

`./bin/main -s 3 -s 3 -ss 7 -ss 5 -ss 4 -ss 3 -ss 5`

**Can extend to any uniformly fretted string instrument with any number of strings - e.g. major 6th chord on 5 string banjo**

`./bin/main -s 4 -s 3 -s 2 -ss 7 -ss 5 -ss 4 -ss 3`

**Can rename the 12 notes - e.g. augmented chord**

`./bin/main -s 4 -s 4 -n "R" -n "b2" -n "2" -n "b3" -n "3" -n "4" -n "#4" -n "5" -n "#5" -n "6" -n "b7" -n "7"`

**Can work with any music system - e.g. hypothetical music system with only 7 notes**

`./bin/main -s 2 -s 3 -n "R" -n "Sh" -n "Tr" -n "Qw" -n "Qb" -n "Si" -n "Se"`
