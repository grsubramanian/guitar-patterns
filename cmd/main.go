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

type note struct {
    stepsAwayFromRootNote uint
}

func newNote(stepsAwayFromRootNote uint) *note {
    return &note{stepsAwayFromRootNote: stepsAwayFromRootNote}
}

func (n *note) equals(other *note) bool {
    if other == nil {
        return false
    }
    return n.stepsAwayFromRootNote == other.stepsAwayFromRootNote
}

type gStringPattern struct {
    notes []*note
}

func newGStringPattern() *gStringPattern {
    return &gStringPattern{notes: make([]*note, 0)}
}

func (gsp *gStringPattern) add(n *note) {
    gsp.notes = append(gsp.notes, n)
}

func (gsp *gStringPattern) leftAligned() bool {
    return len(gsp.notes) > 0 && gsp.notes[0] != nil
}

func (gsp *gStringPattern) rightAligned() bool {
    return len(gsp.notes) > 0 && gsp.notes[len(gsp.notes) - 1] != nil
}

func (gsp *gStringPattern) trailingEmptyFrets() int {
    out := 0
    for i := len(gsp.notes) - 1; i >= 0; i-- {
        if gsp.notes[i] != nil {
            break
        }
        out += 1
    }
    return out
}

type pattern struct {
    gStringPatterns []*gStringPattern
}

func newPattern() *pattern {
    return &pattern{gStringPatterns: make([]*gStringPattern, 0)}
}

func (p *pattern) add(gsp *gStringPattern) {
    p.gStringPatterns = append(p.gStringPatterns, gsp)
}

func (p *pattern) leftAligned() bool {
    for _, gsp := range p.gStringPatterns {
        if gsp.leftAligned() {
            return true
        }
    }
    return false
}

func (p *pattern) rightAligned() bool {
    for _, gsp := range p.gStringPatterns {
        if gsp.rightAligned() {
            return true
        }
    }
    return false
}

func (p *pattern) rtrim() *pattern {
    minTrailingEmptyFrets := -1
    for _, gsp := range p.gStringPatterns {
        trailingEmptyFrets := gsp.trailingEmptyFrets()
        if minTrailingEmptyFrets < 0 || trailingEmptyFrets < minTrailingEmptyFrets {
            minTrailingEmptyFrets = trailingEmptyFrets
        }
    }
    for _, gsp := range p.gStringPatterns {
        gsp.notes = gsp.notes[0:len(gsp.notes) - minTrailingEmptyFrets]
    }
    return p
}

func (p *pattern) subPatternOf(other *pattern) bool {
    if other == nil {
        return false
    }

    numGStrings := len(p.gStringPatterns)
    if numGStrings != len(other.gStringPatterns) {
        return false
    }

    if numGStrings == 0 {
        return true
    }

    numFretsInPattern := len(p.gStringPatterns[0].notes)
    numFretsInOtherPattern := len(other.gStringPatterns[0].notes)
    if numFretsInPattern > numFretsInOtherPattern {
        return false
    }

    // TODO: Implement KMP if really necessary.
    j := 0
    for i := 0; i <= numFretsInOtherPattern - numFretsInPattern; i++ {
        matchFound := true
        for j < numFretsInPattern {
            notesOnAllGStringsForSameFretMatch := true
            for k := 0; k < numGStrings; k++ {
                note := p.gStringPatterns[k].notes[j]
                othernote := other.gStringPatterns[k].notes[i + j]
                notesmatch := (note == nil && othernote == nil) || (note != nil && note.equals(othernote))
                if !notesmatch {
                    notesOnAllGStringsForSameFretMatch = false
                    break
                }
            }
            if !notesOnAllGStringsForSameFretMatch {
                matchFound = false
                break
            }
            j++
        }
        if matchFound {
            return true
        } else {
            j = 0
        }
    }
    return false
}

type patternPrinter interface {
    pprint(*pattern)
}

type asciiPatternPrinter struct {
    paddedNoteRepresentations stringslice
    emptyNoteRepresentation string
}

func newAsciiPatternPrinter(noteRepresentations stringslice) *asciiPatternPrinter {
    pf := &asciiPatternPrinter{}
    paddedNoteRepresentations := getPadded(noteRepresentations)
    emptyNoteRepresentation := strings.Repeat("-", len(paddedNoteRepresentations[0]))
    pf.paddedNoteRepresentations = paddedNoteRepresentations
    pf.emptyNoteRepresentation = emptyNoteRepresentation
    return pf
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

func (apf asciiPatternPrinter) pprint(p *pattern) {

    var sb strings.Builder

    // print out the highest frequency string first.
    for i := len(p.gStringPatterns) - 1; i >= 0; i-- {
        gStringPattern := p.gStringPatterns[i]
        sb.WriteString("|")
        for _, note := range gStringPattern.notes {
           sb.WriteString("-")
           if note != nil {
               sb.WriteString(apf.paddedNoteRepresentations[note.stepsAwayFromRootNote])
           } else {
               sb.WriteString(apf.emptyNoteRepresentation)
           }
           sb.WriteString("-|")
        }
        if i > 0 {
            sb.WriteString("\n")
        }
    }
    fmt.Printf("%s\n\n", sb.String())
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

func addMod(arr uintslice, toAdd uint, maxVal uint) {
    for i, _ := range arr {
        arr[i] += toAdd 
        arr[i] %= maxVal
    }
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

var noteRepresentations stringslice
var defaultNoteRepresentations = stringslice{"1", "b2", "2", "b3", "3", "4", "#4", "5", "b6", "6", "b7", "7"}

var stepsBetweenConsecutiveGStrings uintslice
var defaultStepsBetweenConsecutiveGStrings = uintslice{5, 5, 5, 4, 5}

var stepsBetweenConsecutiveNotesInSequence uintslice

var aliasedRoot uint

var numFretsPerPattern uint

func main() {

    flag.Var(&noteRepresentations, "n", "the textual representations of the notes as an ordered list of strings, starting from the representation of the un-aliased (i.e. absolute) root note. equals " + fmt.Sprintf("%v", defaultNoteRepresentations) + " if not specified.") 
    
    flag.Var(&stepsBetweenConsecutiveGStrings, "ss", "the tuning represented as an ordered list of non-negative integers. each value represents the step jump from previous frequency string. first value represents jump from the lowest frequency string. equals " + fmt.Sprintf("%v", defaultStepsBetweenConsecutiveGStrings) + " if not specified.")
    
    flag.Var(&stepsBetweenConsecutiveNotesInSequence, "s", "each value represents the step jumps from the previous note in the sequence (chord or scale). first value represents jump from the aliased root note. must specify explicitly.")
    
    flag.UintVar(&aliasedRoot, "r", uint(0), "the number of steps away from the absolute root to treat as the temporary root. by default, there is no aliasing.")

    flag.UintVar(&numFretsPerPattern, "frets", uint(4), "the number of frets per pattern")
    flag.Parse()

    if len(noteRepresentations) == 0 {
        noteRepresentations = defaultNoteRepresentations
    }
    numNotes := len(noteRepresentations)
    aliasedRoot = aliasedRoot % uint(numNotes)

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
    addMod(sequenceNotes, aliasedRoot, uint(numNotes))

    pf := newAsciiPatternPrinter(noteRepresentations)

    var lastAcceptedPattern *pattern = nil
    
    for referenceFretOffset := 0; referenceFretOffset < numNotes; referenceFretOffset++ {
        referenceNoteFretNumOnLowestFrequencyGString := int(numFretsPerPattern) - 1 - referenceFretOffset

        pattern := newPattern()
        for gStringNum := 0; gStringNum < numGStrings; gStringNum++ {
            gStringPattern := newGStringPattern()

            for fretNumOnCurrentString := uint(0); fretNumOnCurrentString < numFretsPerPattern; fretNumOnCurrentString++ {
                stepsAwayFromRootNote := (int(stepsAwayFromLowestFrequencyGString[gStringNum] + fretNumOnCurrentString + sequenceNotes[0]) - referenceNoteFretNumOnLowestFrequencyGString) % numNotes
                if stepsAwayFromRootNote < 0 {
                    stepsAwayFromRootNote = numNotes + stepsAwayFromRootNote
                }

                if search(sequenceNotes, uint(stepsAwayFromRootNote)) >= 0 {
                    gStringPattern.add(newNote(uint(stepsAwayFromRootNote)))
                } else {
                    gStringPattern.add(nil)
                }
            }
            pattern.add(gStringPattern)
        }

        // we'll only accept left aligned patterns.
        // we'll expect that they not be a subpattern of another pattern, but given how we iterate, we only need to
        // compare with the last accepted pattern.
        acceptPattern := pattern.leftAligned() && (pattern.rightAligned() || lastAcceptedPattern == nil || !pattern.rtrim().subPatternOf(lastAcceptedPattern))
        if acceptPattern {
            lastAcceptedPattern = pattern
            pf.pprint(pattern)
        }
    }
}
