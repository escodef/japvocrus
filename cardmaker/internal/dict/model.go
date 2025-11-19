package dict

type Sense struct {
	Ru       string
	Notes    string
	Examples []Example
}

type Example struct {
	Ja string
	Ru string
}

type Translation struct {
	Word    string
	Reading string
	Senses  []Sense
}
