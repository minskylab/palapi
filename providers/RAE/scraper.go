package main

import (
	e "errors"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/minskylab/palapi"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func cleaner(def string) (*palapi.WordDefinition, string, error) {
	chunks := strings.Split(def, ".")
	valid := make([]string, 0)
	example := ""
	for _, chunk := range chunks {
		spaces := len(strings.Split(chunk, " "))
		// log.WithFields(log.Fields{
		// 	"chunk": chunk,
		// 	"spaces": spaces,
		// }).Info(def)

		if len(chunk) > 5 && spaces > 1 {
			valid = append(valid, chunk)
		}
	}

	if len(valid) == 0 {
		return nil, example, e.New("invalid definition, scr:28")
	}

	extractedDefinition := strings.TrimSpace(valid[0])

	if len(valid) > 1 {
		example = valid[1]
	}

	return &palapi.WordDefinition{
		PartOfSpeech: "",
		Definition:   extractedDefinition,
	}, example, nil
}

func scrapRAE(word string) ([]palapi.WordDefinition, []string, time.Duration, error) {
	c := colly.NewCollector()

	defs := make([]palapi.WordDefinition, 0)
	examples := make([]string, 0)
	var childError *error
	t1 := time.Now()
	c.OnHTML("div#resultados article p.j", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "p. us.") {
			return
		}
		w, example, err := cleaner(e.Text)
		if err != nil {
			tErr := errors.Wrap(err, "cleaning failed")
			childError = &tErr
			return
		}

		log.WithField("example", example).Info(w.Definition)
		defs = append(defs, *w)
		examples = append(examples, example)
	})

	// c.OnHTML("div#resultados article p", func(e *colly.HTMLElement) {
	// 	log.WithField("class", e.Attr("class")).Info(e.Text)
	// })

	err := c.Visit("https://dle.rae.es/" + strings.Trim(word, " \n\\/.',"))
	if err != nil {
		return nil, nil, 0, errors.Wrap(err, "visit on dle.rae done bad")
	}

	totalDur := time.Since(t1)
	if childError != nil {
		return nil, nil, 0, err
	}

	return defs, examples, totalDur, nil

}