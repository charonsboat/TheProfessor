package main

import "strings"

//type POS int

type Result struct {
	Value  Word
	Weight byte
}

type Goal byte
type P_O_S string

const (
	FIND_SUBJECT  Goal = 0
	FIND_SUBJECTS Goal = 1
	FIND_ACTION   Goal = 2
	//For Professor 2.0
	// BEST_SITUATION  Goal = 3
	// WORST_SITUATION Goal = 4

	NOUN        P_O_S = "noun"
	VERB        P_O_S = "verb"
	ADVERB      P_O_S = "adverb"
	CONJUNCTION P_O_S = "conjunction" //Maybe just AND
	ADJECTIVE   P_O_S = "adjective"
)

type Word struct {
	Value        string
	PartOfSpeech []string //POS, for use later
	Capital      bool
}

type Perception struct {
	Value     string
	Reference *PerceptionEngine
	Weight    byte //I would use float but we only need 1-10 and float32 uses more memory
}

type PerceptionEngine struct {
	Sentence        string
	Words           []Word
	Perceptions     *[]Perception
	masterWeightMax byte
	masterWeightMin byte
	Challenge       Goal
	WordType        P_O_S
	SemiConclusion  []Result
}

func CreateEngine(sentence string, goal Goal) *PerceptionEngine {
	words := findWords(sentence)
	var wordtype P_O_S
	if goal == FIND_SUBJECT || goal == FIND_SUBJECTS {
		wordtype = NOUN
	} else if goal == FIND_ACTION {
		wordtype = VERB
	}
	return &PerceptionEngine{Sentence: sentence, Words: words, masterWeightMax: 255, masterWeightMin: 0, Challenge: goal, WordType: wordtype}
}

func (p *PerceptionEngine) Run() {
	NumPerceptions := make([]Perception, len(p.Words))
	DataSet := make([]Result, len(p.Words))

	for i := 0; i < len(p.Words); i++ {
		NumPerceptions[i] = NewPerception(p.Words[i])
		val := NumPerceptions[i].Run()
		if val > p.masterWeightMin && val < p.masterWeightMax {
			DataSet[i] = Result{p.Words[i], val}
		} else {
			DataSet[i] = 0
		}
	}

}

func (p *PerceptionEngine) Result() []string {
	values := make([]string, len(p.SemiConclusion))
	for i := 0; i < len(p.SemiConclusion); i++ {
		values[i] = p.SemiConclusion[i]
	}

	return values
}

//TODO!
func (p *PerceptionEngine) ReEval(result Result, correct bool) {

}

func (p *PerceptionEngine) NewPerception(target Word) *Perception {
	return &Perception{Value: target.Value, Reference: p} //Weight is set by Perception upon Run
}

func (p *Perception) Run() byte {
	if is(p.Value, p.Reference.WordType) {
		p.Weight++
		if p.Reference.WordType == NOUN {
			if isCap(p.value) {
				p.Weight++

				//TODO!
				//POSITION OF WORD IN SENTENCE NEXT
				//SURROUNDING/NEXT WORDS TYPES (FIND COMMON GROUPINGS)
			}
		}
	} else {
		p.Weight = 0
	}
	return p.Weight
}

//TODO
func isCap(value string) bool {
	return
}

//TODO
//Internal Data-Retrieval Functions
func is(value string, partofspeech P_O_S) bool {

}

//TODO
func checkDB()                   {}
func performLookup(value string) {} //Get info from web sources (mulitple)

func findWords(sentence string) []string {
	var lastIndex int = 0
	var numletters int8 = 0
	var numspaces int8 = 0
	var numwords int = 0
	var inWord bool = false

	for i := 0; i < len(sentence); i++ {
		if string(sentence[i]) == " " {
			if inWord {
				numspaces++
				numwords++
			}
			inWord = false
		} else {
			numletters++
			inWord = true
		}

	}

	words := make([]string, numwords+1)
	numwords = 0

	for i := 0; i < len(sentence); i++ {
		if string(sentence[i]) == " " {
			if inWord {
				numspaces++
				words[numwords] = sentence[lastIndex:i]
				numwords++
				lastIndex = i
			}
			inWord = false
		} else {
			numletters++
			inWord = true
		}

	}

	if lastIndex < len(sentence) {
		words[numwords] = sentence[lastIndex:len(sentence)]
	}

	for i := range words {
		if strings.Contains(words[i], ".") {
			words[i] = removeChar(words[i], ".")
		} else if strings.Contains(words[i], ",") {
			words[i] = removeChar(words[i], ",")
		} else if strings.Contains(words[i], "?") {
			words[i] = removeChar(words[i], "?")
		} else if strings.Contains(words[i], "!") {
			words[i] = removeChar(words[i], "!")
		}
	}

	return words
}
