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
	MAXWORDLEN = 100
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
	previous *Character
	next     *Character
	count    int
	words    Words
}

type CSeqs []*CSeq

type CStat struct {
	char  *Character
	count int
}

type CStats []*CStat

type Character struct {
	value rune
	seqs  map[int]CSeqs
	count int
}

type Characters []*Character

type Prefix struct {
	value string
	next  CStats
}

type Prefixes []*Prefix

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

	chars    Characters
	prefixes Prefixes
	words    Words
	phrases  Phrases

	cposition   = 0
	pChar       *Character
	ipCSeq      int = -1
	cSeqCounter     = 0
	cSeqQueue   [MAXWORDLEN]*Character

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
	fmt.Printf("Text: ")
	input := bufio.NewReader(os.Stdin)
	for {
		var (
			s   string
			err error
			//size int
		)
		for {
			s, err = input.ReadString('\n')
			if err == nil {
				break
			}
		}
		s = s[:len(s)-1]
		fmt.Printf("[%s]", s)
		prefix := searchPrefix(s)
		for prefix != nil {
			fmt.Printf("Next: %s\n", prefix.next)
			s = fmt.Sprintf("%s%c", s, prefix.next[0].char.value)
			prefix = searchPrefix(s)
			fmt.Println(s)
		}
	}
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
			cChar := searchChar(char)
			Chars = append(Chars, cChar)
		}
	}
}

func searchChar(char rune) *Character {
	for _, schar := range chars {
		if schar.value == char {
			schar.count++
			return schar
		}
	}
	schar := &Character{
		value: char,
		seqs:  make(map[int]CSeqs, 0),
		count: 1,
	}
	chars = append(chars, schar)
	return schar
}

func updateCSeq(word *Word) {
	for i, char := range word.chars {
		var previous *Character
		var next *Character
		if i > 0 {
			previous = word.chars[i-1]
		}
		if i < len(word.chars)-1 {
			next = word.chars[i+1]
		}
		searchCSeq(i, char, previous, next, word)
	}
}

func updatePrefix(text string, char *Character) {
	prefix := searchPrefix(text)
	if prefix == nil {
		prefix = &Prefix{
			value: text,
			next:  CStats{},
		}
		prefixes = append(prefixes, prefix)
	}
	for _, next := range prefix.next {
		if next.char.value == char.value {
			next.count++
			sort.Sort(sort.Reverse(prefix.next))
			return
		}
	}
	cstat := &CStat{
		char:  char,
		count: 1,
	}
	prefix.next = append(prefix.next, cstat)
	sort.Sort(sort.Reverse(prefix.next))
}

func searchPrefix(text string) *Prefix {
	for _, prefix := range prefixes {
		if prefix.value == text {
			return prefix
		}
	}
	return nil
}

func searchCSeq(position int, char, previous, next *Character, word *Word) {
	found := false
	for _, cseq := range char.seqs[position] {
		if cseq.previous == previous && cseq.next == next {
			found = true
			cseq.count++
			wIndex := cseq.words.Search(word)
			if wIndex == len(cseq.words) || cseq.words[wIndex] != word {
				cseq.words = append(cseq.words, word)
				sort.Sort(cseq.words)
			}
			sort.Sort(sort.Reverse(char.seqs[position]))
			break
		}
	}
	if !found {
		cseq := &CSeq{
			previous: previous,
			next:     next,
			words:    Words{word},
			count:    1,
		}
		char.seqs[position] = append(char.seqs[position], cseq)
	}
}

func charStats() {
	sort.Sort(sort.Reverse(CharactersByFrequency(chars)))
	for _, char := range chars {
		fmt.Printf("%c[%0x]-(%d)\n", char.value, char.value, char.count)
		fmt.Printf(" -> %s\n", char.seqs)
	}
	/*
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
	*/
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
	fmt.Println("Pass2: word analysis")
	//var buffer Characters
	//var pWord Word
	position := 0
	total := len(Chars)
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
			prefix := Chars[position:i]
			updatePrefix(prefix.String(), &Character{
				value: 0,
				count: 0,
			})
			position = i + 1
			found := words.Search(word)
			if found == len(words) || word.String() != words[found].String() {
				words = append(words, word)
				logf("Added word: [%s]\n", word)
				sort.Sort(words)
				updateCSeq(word)
			} else {
				words[found].count++
				logf("Found word: [%s]\n", word)
			}
		} else {
			prefix := Chars[position:i]
			updatePrefix(prefix.String(), char)
		}
		fmt.Printf("\r%3d%% ", 100*i/total)
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
			updateCSeq(word)
		} else {
			words[found].count++
			logf("Found word: [%s]\n", word)
		}
	}
	fmt.Println("\r100%")
}

func (char *Character) String() string {
	if char == nil {
		return "ø"
	}
	return fmt.Sprintf("%c", char.value)
}

func (chars Characters) String() (result string) {
	for _, char := range chars {
		result = fmt.Sprintf("%s%c", result, char.value)
	}
	return result
}

func (cseq *CSeq) String() string {
	previous := 'ø'
	if cseq.previous != nil {
		previous = cseq.previous.value
	}
	next := 'ø'
	if cseq.next != nil {
		next = cseq.next.value
	}
	return fmt.Sprintf("%c|%c(%d)", previous, next, cseq.count)
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

func (cstat *CStat) String() string {
	if cstat == nil {
		return "."
	}
	return fmt.Sprintf("%0X-%c:%d", cstat.char.value, cstat.char.value, cstat.count)
}

func logf(fmt string, v ...interface{}) {
	if *debug {
		log.Printf(fmt, v...)
	}
}

// CSeqs sorting
func (cseqs CSeqs) Less(i, j int) bool {
	return cseqs[i].count < cseqs[j].count
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
	return stats[i].count < stats[j].count
}
func (stats CStats) Swap(i, j int) {
	stats[i], stats[j] = stats[j], stats[i]
}
func (stats CStats) Len() int {
	return len(stats)
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
