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

func (gsp *gStringPattern) mirrored() *gStringPattern {

    out := newGStringPattern()
    for i := len(gsp.notes) - 1; i >= 0; i -- {
        out.add(gsp.notes[i])
    }
    return out
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

func (p *pattern) mirrored() *pattern {
    out := newPattern()
    for _, gsp := range p.gStringPatterns {
        out.add(gsp.mirrored()) 
    }
    return out
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
    accept(*pattern)
    pprint()
    reset()
}

type asciiPatternPrinter struct {
    paddedNoteRepresentations stringslice
    emptyNoteRepresentation string
    asciiStringBuilder strings.Builder
}

func newAsciiPatternPrinter(noteRepresentations stringslice) patternPrinter {
    pp := &asciiPatternPrinter{}
    paddedNoteRepresentations := getPadded(noteRepresentations)
    emptyNoteRepresentation := strings.Repeat("-", len(paddedNoteRepresentations[0]))
    pp.paddedNoteRepresentations = paddedNoteRepresentations
    pp.emptyNoteRepresentation = emptyNoteRepresentation
    return pp
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

func (app *asciiPatternPrinter) accept(p *pattern) {

    // print out the highest frequency string first.
    for i := len(p.gStringPatterns) - 1; i >= 0; i-- {
        gStringPattern := p.gStringPatterns[i]
        app.asciiStringBuilder.WriteString("|")
        for _, note := range gStringPattern.notes {
           app.asciiStringBuilder.WriteString("-")
           if note != nil {
               app.asciiStringBuilder.WriteString(app.paddedNoteRepresentations[note.stepsAwayFromRootNote])
           } else {
               app.asciiStringBuilder.WriteString(app.emptyNoteRepresentation)
           }
           app.asciiStringBuilder.WriteString("-|")
        }
        if i > 0 {
            app.asciiStringBuilder.WriteString("\n")
        }
    }
    app.asciiStringBuilder.WriteString("\n\n")
}

func (app *asciiPatternPrinter) pprint() {
    fmt.Printf("%s", app.asciiStringBuilder.String())
}

func (app *asciiPatternPrinter) reset() {
    app.asciiStringBuilder.Reset() 
}

type svgPatternPrinter struct {
    noteRepresentations stringslice
    maxFretsInPattern uint
    svgStringBuilder strings.Builder
    numPatterns uint
    numGStringGaps uint
}

func newSvgPatternPrinter(noteRepresentations stringslice, maxFretsInPattern uint) patternPrinter {
    return &svgPatternPrinter{noteRepresentations: noteRepresentations, maxFretsInPattern: maxFretsInPattern} 
}

func (spp *svgPatternPrinter) accept(p *pattern) {

    // Lay out the gstrings.
    numGStringsInPattern := len(p.gStringPatterns)
    for i1 := numGStringsInPattern; i1 >= 1; i1-- {
        x1 := uint(10)
        y1 := (spp.numPatterns + 1) * uint(20) + (spp.numGStringGaps + uint(numGStringsInPattern - i1)) * uint(10)
        x2 := x1 + spp.maxFretsInPattern * uint(15)
        y2 := y1
        spp.svgStringBuilder.WriteString(fmt.Sprintf("  <line x1=\"%d\" y1=\"%d\" x2=\"%d\" y2=\"%d\" style=\"stroke:rgb(190,190,190);stroke-width:1\"/>\n", x1, y1, x2, y2))
    }

    // Lay out the frets.
    for f := uint(0); f <= spp.maxFretsInPattern; f++ {
        var width uint
        if f == uint(0) || f == spp.maxFretsInPattern {
            width = uint(3)
        } else {
            width = uint(2)
        }

        x1 := uint(10) + f * uint(15)
        y1 := (spp.numPatterns + 1) * uint(20) + spp.numGStringGaps * uint(10)
        x2 := x1
        y2 := y1 + uint(numGStringsInPattern - 1) * uint(10)
        spp.svgStringBuilder.WriteString(fmt.Sprintf("  <line x1=\"%d\" y1=\"%d\" x2=\"%d\" y2=\"%d\" style=\"stroke:rgb(184,115,51);stroke-width:%d\"/>\n", x1, y1, x2, y2, width))
    }

    // Lay out the notes.
    for i1 := numGStringsInPattern; i1 >= 1; i1-- {
        for f, note := range p.gStringPatterns[i1 - 1].notes {
            if note != nil {
                noteRepresentation := spp.noteRepresentations[note.stepsAwayFromRootNote]
                width := uint(4) * uint(len(noteRepresentation))
                height := uint(8)
                x := uint(10) + uint(f) * uint(15) + uint(15) / uint(2)
                y := (spp.numPatterns + 1) * uint(20) + (spp.numGStringGaps + uint(numGStringsInPattern - i1)) * uint(10)
                spp.svgStringBuilder.WriteString(fmt.Sprintf("  <text x=\"%d\" y=\"%d\" textLength=\"%d\" lengthAdjust=\"spacingAndGlyphs\" alignment-baseline=\"middle\" style=\"fill:rgb(71,140,204);text-anchor:middle;font-size:%d\">%s</text>\n", x, y, width, height, noteRepresentation))
            }
        }
    }

    spp.numGStringGaps += uint(numGStringsInPattern - 1)
    spp.numPatterns++
}

func (spp *svgPatternPrinter) pprint() {
    width := 2 * uint(10) + spp.maxFretsInPattern * uint(15)
    height := (spp.numPatterns + 1) * uint(20) + spp.numGStringGaps * uint(10)

    fmt.Printf("<svg width=\"%d\" height=\"%d\">\n", width, height)
    fmt.Printf(spp.svgStringBuilder.String())
    fmt.Printf("</svg>\n")
}

func (spp *svgPatternPrinter) reset() {
    spp.svgStringBuilder.Reset()
    spp.numPatterns = 0
    spp.numGStringGaps = 0
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

var leftHanded bool

var printSvg bool

func main() {

    flag.Var(&noteRepresentations, "n", "the textual representations of the notes as an ordered list of strings, starting from the representation of the un-aliased (i.e. absolute) root note. equals " + fmt.Sprintf("%v", defaultNoteRepresentations) + " if not specified.") 
    
    flag.Var(&stepsBetweenConsecutiveGStrings, "ss", "the tuning represented as an ordered list of non-negative integers. each value represents the step jump from previous frequency string. first value represents jump from the lowest frequency string. equals " + fmt.Sprintf("%v", defaultStepsBetweenConsecutiveGStrings) + " if not specified.")
    
    flag.Var(&stepsBetweenConsecutiveNotesInSequence, "s", "each value represents the step jumps from the previous note in the sequence (chord or scale). first value represents jump from the aliased root note. must specify explicitly.")
    
    flag.UintVar(&aliasedRoot, "r", uint(0), "the number of steps away from the absolute root to treat as the temporary root. by default, there is no aliasing.")

    flag.UintVar(&numFretsPerPattern, "frets", uint(4), "the number of frets per pattern")

    flag.BoolVar(&leftHanded, "left", false, "format for lefties")

    flag.BoolVar(&printSvg, "svg", false, "whether to print in SVG format")

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

    var pp patternPrinter
    if printSvg {
        pp = newSvgPatternPrinter(noteRepresentations, numFretsPerPattern)
    } else {
        pp = newAsciiPatternPrinter(noteRepresentations)
    }
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
            if leftHanded {
                pp.accept(pattern.mirrored())
            } else {
                pp.accept(pattern)
            }
        }
    }
    pp.pprint()
    pp.reset()
}
