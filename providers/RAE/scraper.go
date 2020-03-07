package rae

import (
	"path"
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
		chunk = strings.TrimSpace(chunk)
		parts := len(strings.Split(chunk, " "))
		// log.WithFields(log.Fields{
		// 	"chunk": len(chunk),
		// 	"parts": parts,
		// }).Info(chunk)

		if len(chunk) > 5 && parts > 1 {
			valid = append(valid, chunk)
		}
	}

	if len(valid) == 0 {
		return nil, example, ErrInvalidRAEDefinition
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

func (p *Provider) scraper(word string) ([]palapi.WordDefinition, []string, time.Duration, error) {
	c := colly.NewCollector()

	defs := make([]palapi.WordDefinition, 0)
	examples := make([]string, 0)
	var childError *error
	t1 := time.Now()

	c.OnHTML("div#resultados article p[class^=j]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "p. us.") {
			return
		}
		w, example, err := cleaner(e.Text)
		if err != nil {
			if errors.Is(err, ErrInvalidRAEDefinition) {
				return
			}
			tErr := errors.Wrap(err, "cleaning failed")
			childError = &tErr
			return
		}
		example = strings.Trim(example, " \n\\/.',")

		log.WithField("example", example).Info(w.Definition)
		defs = append(defs, *w)
		if example != "" {
			examples = append(examples, example)
		}
	})

	// c.OnHTML("div#resultados article p", func(e *colly.HTMLElement) {
	// 	log.WithField("class", e.Attr("class")).Info(e.Text)
	// })

	err := c.Visit(path.Join(p.baseURL, strings.Trim(word, " \n\\/.',")))
	if err != nil {
		return nil, nil, 0, errors.Wrap(err, "visit on dle.rae done bad")
	}

	totalDur := time.Since(t1)
	if childError != nil {
		return nil, nil, 0, *childError
	}

	return defs, examples, totalDur, nil

}