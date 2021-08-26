package main

import (
    "flag"
    "fmt"
    "os"
    "sort"
    "strconv"
    "strings"
)

type uintslice []uint

func (a uintslice) Len() int           { return len(a) }
func (a uintslice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a uintslice) Less(i, j int) bool { return a[i] < a[j] }

func (i *uintslice) String() string {
    return fmt.Sprintf("%v", *i)
}
 
func (i *uintslice) Set(value string) error {
    tmp, err := strconv.ParseUint(value, 10, 32)
    if err != nil {
        return err
    } else {
        *i = append(*i, uint(tmp))
    }
    return nil
}

type stringslice []string

func (s *stringslice) String() string {
    return fmt.Sprintf("%v", *s)
}
 
func (s *stringslice) Set(value string) error {
    *s = append(*s, value)
    return nil
}

 
func cumSum(arr uintslice) uintslice {
    out := make(uintslice, 0)
    out = append(out, 0)
    rsum := uint(0)
    for _, item := range arr {
        rsum += item
        out = append(out, rsum)
    }
    return out
}

func cumSumMod(arr uintslice, maxVal uint) uintslice {
    out := make(uintslice, 0)
    out = append(out, 0)
    rsum := uint(0)
    for _, item := range arr {
        rsum += item
        rsum %= maxVal
        out = append(out, rsum)
    }
    return out
}

func unique(arr uintslice) uintslice {
    if len(arr) == 0 {
        return arr
    }
    prevVal := arr[0]

    lPos := 1
    for rPos := 1; rPos < len(arr); rPos++ {
        if arr[rPos] != prevVal {
            arr[lPos] = arr[rPos]
            lPos++
            prevVal = arr[rPos]
        }
    }
    arr = arr[0:lPos]
    return arr
}

func search(arr uintslice, tgt uint) int {
    for pos, item := range arr {
        if item == tgt {
            return pos
        }
    }
    return -1
}

func reversed(arr stringslice) stringslice {
    l := len(arr)
    for i := 0; i < l / 2; i++ {
        temp := arr[i]
        arr[i] = arr[l - 1 - i]
        arr[l - 1 - i] = temp
    }
    return arr
}

func paddedNotes(notes stringslice) stringslice {

    maxNoteLen := 0
    for _, note := range notes {
        noteLen := len(note)
        if noteLen > maxNoteLen {
            maxNoteLen = noteLen
        }
    }

    out := make(stringslice, 0)
    for _, note := range notes {
        noteLen := len(note)
        var paddedNote string
        if noteLen == maxNoteLen {
            paddedNote = note
        } else {
            leftPadding := (maxNoteLen - noteLen) / 2
            rightPadding := maxNoteLen - noteLen - leftPadding
            paddedNote = strings.Repeat("-", leftPadding) + note + strings.Repeat("-", rightPadding) 
        }
        out = append(out, paddedNote)
    }
    return out
}

var notes stringslice
var defaultNotes = stringslice{"1", "b2", "2", "b3", "3", "4", "#4", "5", "b6", "6", "b7", "7"}

var stringSteps uintslice
var defaultStringSteps = uintslice{5, 5, 5, 4, 5}

var chordSteps uintslice

var numFretsPerPattern uint

func main() {

    flag.Var(&notes, "n", "the names of the notes represented as an ordered list of strings, starting from the name of the root note. equals " + fmt.Sprintf("%v", defaultNotes) + " if not specified.") 
    
    flag.Var(&stringSteps, "ss", "the tuning represented as an ordered list of non-negative integers. each value represents the step jump from previous frequency string. first value represents jump from the lowest frequency string. equals " + fmt.Sprintf("%v", defaultStringSteps) + " if not specified.")
    
    flag.Var(&chordSteps, "s", "each value represents the step jumps from the previous note in the chord. first value represents jump from the root note. must specify explicitly.")
    
    flag.UintVar(&numFretsPerPattern, "frets", uint(4), "the number of frets per pattern")
    flag.Parse()

    if len(notes) == 0 {
        notes = defaultNotes
    }
    paddedNotes := paddedNotes(notes)
    emptyFret := strings.Repeat("-", len(paddedNotes[0]))
    numNotes := len(paddedNotes)

    if len(stringSteps) == 0 {
        stringSteps = defaultStringSteps
    }
    cumStringSteps := cumSum(stringSteps)
    numStrings := len(cumStringSteps)

    if len(chordSteps) == 0 {
        flag.PrintDefaults()
        os.Exit(1)
    }
    cumChordSteps := cumSumMod(chordSteps, uint(numNotes))
    sort.Sort(cumChordSteps)
    cumChordSteps = unique(cumChordSteps)

    for _, chordStepOnTopString := range cumChordSteps {
        // fret on top string, but 1-indexed to avoid uint underflow.
        for fret1OnTopString := numFretsPerPattern; fret1OnTopString >= 1; fret1OnTopString-- {
            fretOnTopString := fret1OnTopString - 1
            stringPatterns := make([]string, 0)
            for stringNum := 0; stringNum < numStrings; stringNum++ {
                var sb strings.Builder
                sb.WriteString("|")

                // The following two nested for loops can in theory be optimized wrt runtime complexity,
                // but it's not worth it.
                for fretOnCurString := uint(0); fretOnCurString < numFretsPerPattern; fretOnCurString++ {
                    stepsAwayFromRoot := (int(cumStringSteps[stringNum] + fretOnCurString + chordStepOnTopString) - int(fretOnTopString)) % numNotes
                    if stepsAwayFromRoot < 0 {
                        stepsAwayFromRoot = numNotes + stepsAwayFromRoot
                    }
                    
                    sb.WriteString("-")
                    if search(cumChordSteps, uint(stepsAwayFromRoot)) >= 0 {
                        sb.WriteString(paddedNotes[stepsAwayFromRoot])
                    } else {
                        sb.WriteString(emptyFret)    
                    }
                    sb.WriteString("-|")
                }
                stringPatterns = append(stringPatterns, sb.String())
            }
            stringPatterns = reversed(stringPatterns)
            fmt.Printf("%s\n\n", strings.Join(stringPatterns, "\n"))
        }
    }
}
