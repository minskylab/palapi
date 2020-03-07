package palapi

type Persistence interface {
	SaveWord(word Word) (*Word, error)
	GetWord(word string) (*Word, error)
	UpdateWord(word Word) (*Word, error)
	DeleteWord(word string) (*Word, error)
}
