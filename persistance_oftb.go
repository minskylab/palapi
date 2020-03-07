package palapi

import "github.com/asdine/storm/v3"


type DefaultPersistence struct {
	filepath string
	db *storm.DB
}


type WordForStorm struct {
	ID         WordID `storm:"id"`
	Source     SourceID
	Definitions []WordDefinition
	Synonyms   map[int64]WordID
	Antonyms   map[int64]WordID
	Examples   []Sentence
	Frequency  WordFrequency
}

func NewDefaultPersistence(filepath string) (Persistence, error){
	db, err := storm.Open(filepath)
	if err != nil {
		return nil, err
	}

	return &DefaultPersistence{
		filepath: filepath,
		db:       db,
	}, nil
}

func (p *DefaultPersistence) SaveWord(word Word) (*Word, error) {
	w := WordForStorm(word)
	if err := p.db.Save(&w); err != nil {
		return nil, err
	}

	word = Word(w)
	return &word, nil
}

func (p *DefaultPersistence) GetWord(word string) (*Word, error) {
	w := new(WordForStorm)
	if err := p.db.Find("ID", word, w); err != nil {
		return nil, err
	}

	wNative := Word(*w)
	return &wNative, nil
}

func (p *DefaultPersistence) UpdateWord(word Word) (*Word, error) {
	w := WordForStorm(word)
	if err := p.db.Update(&w); err != nil {
		return nil, err
	}

	return p.GetWord(string(w.ID))
}

func (p *DefaultPersistence) DeleteWord(word string) (*Word, error) {
	w, err := p.GetWord(word)
	if err != nil {
		return nil, err
	}
	var backup = *w

	if err := p.db.DeleteStruct(w); err != nil {
		return nil, err
	}

	return &backup, nil
}
