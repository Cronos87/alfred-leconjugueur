package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/raguay/goAlfred"
)

var positions = map[string]int{
	"present":                           0,
	"passe_compose":                     100,
	"imparfait":                         6,
	"plus_que_parfait":                  106,
	"passe_simple":                      12,
	"passe_anterieur":                   112,
	"futur_simple":                      18,
	"futur_anterieur":                   118,
	"subjonctif_present":                24,
	"subjonctif_passe":                  124,
	"subjonctif_imparfait":              30,
	"subjonctif_plus_que_parfait":       124,
	"conditionnel_present":              36,
	"conditionnel_passe_premiere_forme": 136,
	"conditionnel_passe_deuxieme_forme": 137,
	"imperatif_present":                 42,
	"imperatif_passe":                   142,
	"participe_present":                 46,
	"participe_passe":                   255,
	"infinitif_present":                 253,
	"infinitif_passe":                   254,
	"gerondif_present":                  251,
	"gerondif_passe":                    221,
}

func parse(verb string, position int) string {
	url := fmt.Sprintf("http://leconjugueur.lefigaro.fr/conjugaison/verbe/%s.html", verb)

	doc, err := goquery.NewDocument(url)

	if err != nil {
		return goAlfred.ToXML()
	}

	verbs, err := doc.Find(fmt.Sprintf("#temps%d + p", position)).Html()

	if err != nil {
		return goAlfred.ToXML()
	}

	if len(verbs) == 0 {
		verbs := doc.Find(".pt li a")
		nbVerbs := verbs.Length()

		if nbVerbs > 0 {
			goAlfred.AddResult("0", "0", "You mean:", "", "", "yes", "", "", 1)

			verbs.EachWithBreak(func(i int, s *goquery.Selection) bool {
				if s.Text() == "être" {
					return false
				}

				goAlfred.AddResult(strconv.Itoa(i), s.Text(), s.Text(), "", "", "yes", "", "", 1)

				return true
			})
		}

		return goAlfred.ToXML()
	}

	times := strings.Split(verbs, "<br/>")

	for i, time := range times {
		time = strings.Replace(time, "<b>", "", 1)
		time = strings.Replace(time, "</b>", "", 1)
		time = strings.Replace(time, "&#39;", "'", 1)

		goAlfred.AddResult(strconv.Itoa(i), time, time, "", "", "yes", "", "", 1)
	}

	return goAlfred.ToXML()
}

func main() {
	var position int
	var cmd, verb, xml string

	cmd = os.Args[1]
	verb = os.Args[2]

	position = positions[cmd]

	xml = parse(verb, position)

	fmt.Print(xml)
}
