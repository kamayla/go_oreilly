package theaurus

type Thesaurus interface {
	Synonyms(term string) ([]string, error)
}
