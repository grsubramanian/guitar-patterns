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

func getPadded(noteRepresentations stringslice) stringslice {

    maxNoteLen := 0
    for _, noteRepresentation := range noteRepresentations {
        noteRepresentationLen := len(noteRepresentation)
        if noteRepresentationLen > maxNoteLen {
            maxNoteLen = noteRepresentationLen
        }
    }

    out := make(stringslice, 0)
    for _, noteRepresentation := range noteRepresentations {
        noteRepresentationLen := len(noteRepresentation)
        var paddedNoteRepresentation string
        if noteRepresentationLen == maxNoteLen {
            paddedNoteRepresentation = noteRepresentation
        } else {
            leftPadding := (maxNoteLen - noteRepresentationLen) / 2
            rightPadding := maxNoteLen - noteRepresentationLen - leftPadding
            paddedNoteRepresentation = strings.Repeat("-", leftPadding) + noteRepresentation + strings.Repeat("-", rightPadding) 
        }
        out = append(out, paddedNoteRepresentation)
    }
    return out
}

var noteRepresentations stringslice
var defaultNoteRepresentations = stringslice{"1", "b2", "2", "b3", "3", "4", "#4", "5", "b6", "6", "b7", "7"}

var stepsBetweenConsecutiveGStrings uintslice
var defaultStepsBetweenConsecutiveGStrings = uintslice{5, 5, 5, 4, 5}

var stepsBetweenConsecutiveNotesInSequence uintslice

var numFretsPerPattern uint

func main() {

    flag.Var(&noteRepresentations, "n", "the textual representations of the notes as an ordered list of strings, starting from the representation of the root note. equals " + fmt.Sprintf("%v", defaultNoteRepresentations) + " if not specified.") 
    
    flag.Var(&stepsBetweenConsecutiveGStrings, "ss", "the tuning represented as an ordered list of non-negative integers. each value represents the step jump from previous frequency string. first value represents jump from the lowest frequency string. equals " + fmt.Sprintf("%v", defaultStepsBetweenConsecutiveGStrings) + " if not specified.")
    
    flag.Var(&stepsBetweenConsecutiveNotesInSequence, "s", "each value represents the step jumps from the previous note in the sequence (chord or scale). first value represents jump from the root note. must specify explicitly.")
    
    flag.UintVar(&numFretsPerPattern, "frets", uint(4), "the number of frets per pattern")
    flag.Parse()

    if len(noteRepresentations) == 0 {
        noteRepresentations = defaultNoteRepresentations
    }
    paddedNoteRepresentations := getPadded(noteRepresentations)
    emptyNoteRepresentation := strings.Repeat("-", len(paddedNoteRepresentations[0]))
    numNotes := len(paddedNoteRepresentations)

    if len(stepsBetweenConsecutiveGStrings) == 0 {
        stepsBetweenConsecutiveGStrings = defaultStepsBetweenConsecutiveGStrings
    }
    stepsAwayFromLowestFrequencyGString := cumSum(stepsBetweenConsecutiveGStrings)
    numGStrings := len(stepsAwayFromLowestFrequencyGString)

    if len(stepsBetweenConsecutiveNotesInSequence) == 0 {
        flag.PrintDefaults()
        os.Exit(1)
    }
    sequenceNotes := cumSumMod(stepsBetweenConsecutiveNotesInSequence, uint(numNotes))
    sort.Sort(sequenceNotes)
    sequenceNotes = unique(sequenceNotes)

    for _, sequenceNoteOnLowestFrequencyGString := range sequenceNotes {
        // fret on top string, but 1-indexed to avoid uint underflow.
        for referenceFretNumOnLowestFrequencyGString_1 := numFretsPerPattern; referenceFretNumOnLowestFrequencyGString_1 >= 1; referenceFretNumOnLowestFrequencyGString_1-- {
            referenceFretNumOnLowestFrequencyGString := referenceFretNumOnLowestFrequencyGString_1 - 1
            gStringPatterns := make([]string, 0)
            for gStringNum := 0; gStringNum < numGStrings; gStringNum++ {
                var sb strings.Builder
                sb.WriteString("|")

                // The following two nested for loops can in theory be optimized wrt runtime complexity,
                // but it's not worth it.
                for fretNumOnCurrentString := uint(0); fretNumOnCurrentString < numFretsPerPattern; fretNumOnCurrentString++ {
                    stepsAwayFromRootNote := (int(stepsAwayFromLowestFrequencyGString[gStringNum] + fretNumOnCurrentString + sequenceNoteOnLowestFrequencyGString) - int(referenceFretNumOnLowestFrequencyGString)) % numNotes
                    if stepsAwayFromRootNote < 0 {
                        stepsAwayFromRootNote = numNotes + stepsAwayFromRootNote
                    }
                    
                    sb.WriteString("-")
                    if search(sequenceNotes, uint(stepsAwayFromRootNote)) >= 0 {
                        sb.WriteString(paddedNoteRepresentations[stepsAwayFromRootNote])
                    } else {
                        sb.WriteString(emptyNoteRepresentation)    
                    }
                    sb.WriteString("-|")
                }
                gStringPatterns = append(gStringPatterns, sb.String())
            }
            gStringPatterns = reversed(gStringPatterns)
            fmt.Printf("%s\n\n", strings.Join(gStringPatterns, "\n"))
        }
    }
}
