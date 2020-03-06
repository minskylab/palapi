package palapi



type Persistence interface {
	SaveWord(word Word) (Word, error)

}
