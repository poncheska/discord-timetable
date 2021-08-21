package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"net/http"
	"regexp"
)

var regexNotSpace = regexp.MustCompile("(\\S+\\s?)+[^\\s]")

type Lesson struct {
	Time    string
	Type    string
	Place   string
	Teacher string
}

type Day struct {
	Date    string
	Lessons []Lesson
}

type Timetable struct {
	Days []Day
}

func ParseTimetable(link string) (*Timetable, error) {
	c := colly.NewCollector()

	tt := &Timetable{make([]Day, 0, 0)}
	correct := true

	c.OnHTML("div.panel-group div.panel-default", func(e *colly.HTMLElement) {
		date := regexNotSpace.FindString(e.DOM.Find("div.panel-default > div.panel-heading").Text())
		d := Day{date, make([]Lesson, 0, 0)}
		e.DOM.Find(
			"div.panel-default > ul.panel-collapse > li.row",
		).Each(func(i int, selection *goquery.Selection) {
			time := regexNotSpace.FindString(selection.Find("li.row > div:nth-child(1) div:nth-child(2)").Text())
			typ := regexNotSpace.FindString(selection.Find("li.row > div:nth-child(2) div:nth-child(2)").Text())
			place := regexNotSpace.FindString(selection.Find("li.row > div:nth-child(3) div:nth-child(2)").Text())
			teacher := regexNotSpace.FindString(selection.Find("li.row > div:nth-child(4) div:nth-child(2)").Text())
			l := Lesson{time, typ, place, teacher}
			d.Lessons = append(d.Lessons, l)
		})
		tt.Days = append(tt.Days, d)
	})

	header := http.Header{"User-Agent": []string{c.UserAgent}}
	header.Set("Accept-Language", "ru")
	err := c.Request("GET", link, nil, nil, header)

	if !correct {
		return nil, errors.New("invalid timetable")
	}

	return tt, err
}

func (tt *Timetable) GetString() string {
	buf := bytes.Buffer{}
	buf.WriteString("Расписание на неделю:\n")
	for _, day := range tt.Days {
		ss := day.GetString()
		for _, s := range ss {
			buf.WriteString(s)
		}
	}
	return buf.String()
}

func (d Day) GetString() []string {
	res := make([]string, 0, 0)
	buf := bytes.NewBuffer([]byte{})
	counter := 0
	buf.WriteString("__***" + d.Date + "***__\n")
	for _, les := range d.Lessons {
		if counter == 5 {
			res = append(res, buf.String())
			buf = bytes.NewBuffer([]byte{})
			counter = 0
		}
		buf.WriteString(fmt.Sprintf("***(%v)***\t**%v**"+
			"\n*Место*: %v"+
			"\n*Препод.*: %v\n", les.Time, les.Type, les.Place, les.Teacher))
		counter++
	}
	res = append(res, buf.String())
	return res
}
