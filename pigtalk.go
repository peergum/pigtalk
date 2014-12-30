package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

const (
	MAXWORDLEN = 1
)

type PSeq struct {
	words    []*Word
	position int
	previous *Phrase
	next     *Phrase
	context  int
}

type WSeq struct {
	characters []*Character
	position   int
	previous   *Word
	next       *Word
	phrase     *Phrase
}

type CSeq struct {
	position int
	previous *Character
	word     *Word
}

type CSeqs []*CSeq

type CStat struct {
	char  *Character
	count int
}

type CStats []*CStat

type Character struct {
	value rune
	seqs  CSeqs
	next  [MAXWORDLEN]CStats
	count int
}

type Characters []*Character

type Word struct {
	chars Characters
	seqs  []*WSeq
	count int
}

type Words []*Word

type Phrase struct {
	words []*Word
	seqs  []*PSeq
}

type Phrases []*Phrase

type Current struct {
	value   string
	pChar   *Character
	pWord   *Word
	pPhrase *Phrase
}

var (
	Chars Characters

	chars   Characters
	words   Words
	phrases Phrases

	cposition = 0
	pChar     *Character
	ipCSeq    int = -1

	wposition = 0
	pWord     *Word
	cWord     *Word
	ipWSeq    int = -1

	pPhrase *Phrase
	cPhrase *Phrase
	ipPSeq  int = -1

	//cQueue    Characters = make(Characters, 0, MAXWORDLEN)
	//totalChar int        = 0
	space *Character

	debug *bool = flag.Bool("d", false, "Debugging")
)

func main() {

	flag.Parse()
	fmt.Println(*debug)
	args := flag.Args()
	if len(args) < 1 {
		panic("Missing filename")
		os.Exit(1)
	}
	filename := args[0]
	scan(filename)
	fmt.Println()
	charStats()
	pass2()
	wordStats()
}

func scan(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	eof := false
	for !eof {
		char, size, err := reader.ReadRune()
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			eof = true
		}
		if size > 0 {
			_, pChar, err = searchCSeq(char)
			Chars = append(Chars, pChar)
		}
	}
}

func (char *Character) addNext(nextChar *Character, position int) {
	found := false
	for _, cstat := range char.next[position] {
		if cstat.char == nextChar {
			cstat.count++
			found = true
			break
		}
	}
	if !found {
		if char.next[position] == nil {
			char.next[position] = make(CStats, 0)
		}
		char.next[position] = append(char.next[position], &CStat{
			char:  nextChar,
			count: 1,
		})
	}
}

func searchCSeq(char rune) (cseq *CSeq, cchar *Character, err error) {
	// search char
	for _, schar := range chars {
		if schar.value == char {
			schar.count++
			//logf("Found char %c\n", char)
			found := false
			for _, cseq := range schar.seqs {
				if cseq.previous == pChar {
					fmt.Printf("\x1b[30m%c\x1b[0m", char)
					found = true
					cseq.previous.addNext(schar, 0)
					break
				}
			}
			if !found {
				fmt.Printf("\x1b[41m%c\x1b[0m", char)
				cseq := CSeq{
					position: cposition,
					previous: pChar,
					word:     nil,
				}
				cseq.previous.addNext(schar, 0)
				schar.seqs = append(schar.seqs, &cseq)
			}
			cposition++
			return cseq, schar, nil
		}
	}
	fmt.Printf("\x1b[42m%c\x1b[0m", char)
	cseq = &CSeq{
		position: cposition,
		previous: pChar,
		word:     nil,
	}
	schar := &Character{
		value: char,
		seqs: CSeqs{
			cseq,
		},
		count: 1,
	}
	if pChar != nil {
		cseq.previous.addNext(schar, 0)
	}
	chars = append(chars, schar)
	logf("Added char %c [%d]\n", char, len(chars))
	cposition++
	return cseq, schar, nil
}

func (char *Character) String() string {
	if char == nil {
		return "ø"
	}
	return fmt.Sprintf("Char %c [%d] Previous: %s Next: %s\n", char.value, len(char.seqs), char.seqs, char.next)
}

func (cseq *CSeq) String() string {
	previous := 'ø'
	if cseq.previous != nil {
		previous = cseq.previous.value
	}
	return fmt.Sprintf("%c", previous)
}

func (word *Word) String() (result string) {
	if word == nil {
		return "ø"
	}
	for _, char := range word.chars {
		result = fmt.Sprintf("%s%c", result, char.value)
	}
	return result
}

func logf(fmt string, v ...interface{}) {
	if *debug {
		log.Printf(fmt, v...)
	}
}

func (cstat *CStat) String() string {
	if cstat == nil {
		return "."
	}
	return fmt.Sprintf("%0X-%c:%d", cstat.char.value, cstat.char.value, cstat.count)
}

// CSeqs sorting
func (cseqs CSeqs) Less(i, j int) bool {
	switch {
	case cseqs[i].previous == nil && cseqs[j].previous != nil:
		return true
	case cseqs[i].previous == nil || cseqs[j].previous == nil:
		return false
	}
	return cseqs[i].previous.value < cseqs[j].previous.value
}
func (cseqs CSeqs) Swap(i, j int) {
	cseqs[i], cseqs[j] = cseqs[j], cseqs[i]
}
func (cseqs CSeqs) Len() int {
	return len(cseqs)
}

// Characters sorting
func (chars Characters) Less(i, j int) bool {
	return chars[i].value < chars[j].value
}
func (chars Characters) Swap(i, j int) {
	chars[i], chars[j] = chars[j], chars[i]
}
func (chars Characters) Len() int {
	return len(chars)
}

// Characters sorting by count
type CharactersByFrequency Characters

func (chars CharactersByFrequency) Less(i, j int) bool {
	return chars[i].count < chars[j].count
}
func (chars CharactersByFrequency) Swap(i, j int) {
	chars[i], chars[j] = chars[j], chars[i]
}
func (chars CharactersByFrequency) Len() int {
	return len(chars)
}

// Stats sorting
func (stats CStats) Less(i, j int) bool {
	return stats[i].char.value < stats[j].char.value
}
func (stats CStats) Swap(i, j int) {
	stats[i], stats[j] = stats[j], stats[i]
}
func (stats CStats) Len() int {
	return len(stats)
}

// sequence frequency sorting
type BySeqFrequency Characters

func (chars BySeqFrequency) Len() int {
	return len(chars)
}
func (chars BySeqFrequency) Less(i, j int) bool {
	return len(chars[i].next[0])+len(chars[i].seqs) < len(chars[j].next[0])+len(chars[j].seqs)
}
func (chars BySeqFrequency) Swap(i, j int) {
	chars[i], chars[j] = chars[j], chars[i]
}

// Words sorting
func (words Words) Len() int {
	return len(words)
}
func (words Words) Less(i, j int) bool {
	for index, char := range words[i].chars {
		if index >= len(words[j].chars) {
			return false
		}
		if char.value < words[j].chars[index].value {
			return true
		}
		if char.value > words[j].chars[index].value {
			return false
		}
	}
	if len(words[i].chars) < len(words[j].chars) {
		return true
	}
	return false
}
func (words Words) Swap(i, j int) {
	words[i], words[j] = words[j], words[i]
}
func (words Words) Search(needle *Word) int {
	for index, word := range words {
		compare := Words{
			word,
			needle,
		}
		if !compare.Less(0, 1) {
			return index
		}
	}
	return len(words)
}

// Words sorting by frequency
type WordsByFrequency Words

func (words WordsByFrequency) Len() int {
	return len(words)
}
func (words WordsByFrequency) Less(i, j int) bool {
	return words[i].count < words[j].count
}
func (words WordsByFrequency) Swap(i, j int) {
	words[i], words[j] = words[j], words[i]
}

func charStats() {
	sort.Sort(chars)
	for _, char := range chars {
		fmt.Printf("%c[%0x]-(%d)", char.value, char.value, char.count)
		sort.Sort(char.seqs)
		fmt.Printf(" -> %s", char.seqs)
		sort.Sort(sort.Reverse(char.next[0]))
		fmt.Printf(" -> %s\n", char.next[0])
	}
	fmt.Println("Sorting by Sequence Frequency")
	sort.Sort(sort.Reverse(BySeqFrequency(chars)))
	for _, char := range chars {
		fmt.Printf("[%0X:%c-%d]\n", char.value, char.value, len(char.next[0])+len(char.seqs))
	}
	fmt.Println("Sorting by Appearance Frequency")
	sort.Sort(sort.Reverse(CharactersByFrequency(chars)))
	for _, char := range chars {
		fmt.Printf("%c[%0x]: %d\n", char.value, char.value, char.count)
	}
	space = chars[0]
	fmt.Printf("SPACE is '%c'\n", space.value)
}

func wordStats() {
	sort.Sort(words)
	for _, word := range words {
		fmt.Printf("%s (%d)\n", word, word.count)
		// sort.Sort(char.seqs)
		// fmt.Printf(" -> %s", char.seqs)
		// sort.Sort(sort.Reverse(char.next[0]))
		// fmt.Printf(" -> %s\n", char.next[0])
	}
	/*
		fmt.Println("Sorting by Sequence Frequency")
		sort.Sort(sort.Reverse(BySeqFrequency(chars)))
		for _, char := range chars {
			fmt.Printf("[%0X:%c-%d]\n", char.value, char.value, len(char.next[0])+len(char.seqs))
		}
		space = chars[0]
		fmt.Printf("SPACE is '%c'\n", space.value)
	*/
	fmt.Println("Sorting by Appearance Frequency")
	sort.Sort(WordsByFrequency(words))
	for _, word := range words {
		fmt.Printf("%03d: %s\n", word.count, word)
	}
}

func pass2() {
	logf("Pass2: words\n")
	//var buffer Characters
	//var pWord Word
	position := 0
	for i, char := range Chars {
		//fmt.Printf("%c", char.value)
		if char.value == space.value || char.value == '\r' || char.value == '\n' {
			if position == i {
				position++
				continue
			}
			word := &Word{
				chars: Chars[position:i],
				count: 1,
			}
			position = i + 1
			found := words.Search(word)
			if found == len(words) || word.String() != words[found].String() {
				words = append(words, word)
				logf("Added word: [%s]\n", word)
				sort.Sort(words)
				fmt.Printf("\x1b[41m%s\x1b[0m ", word)
			} else {
				words[found].count++
				logf("Found word: [%s]\n", word)
				fmt.Printf("%s ", word)
			}
		} else {

		}
	}
	if position < len(Chars) {
		word := &Word{
			chars: Chars[position:],
			count: 1,
		}
		found := words.Search(word)
		if found == len(words) || word.String() != words[found].String() {
			words = append(words, word)
			logf("Added word: [%s]\n", word)
			sort.Sort(words)
			fmt.Printf("\x1b[41m%s\x1b[0m ", word)
		} else {
			words[found].count++
			logf("Found word: [%s]\n", word)
			fmt.Printf("%s ", word)
		}
	}
	fmt.Println()
}
