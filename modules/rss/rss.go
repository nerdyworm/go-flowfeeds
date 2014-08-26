package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

func Read(reader io.Reader) (feed Feed, err error) {
	return feed, xml.NewDecoder(reader).Decode(&feed)
}

type Feed struct {
	Channel struct {
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Image       struct {
			HREF string `xml:"href,attr"`
			URL  string `xml:"url"`
		} `xml:"image"`

		Items []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Guid        string `xml:"guid"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Enclosure   struct {
		URL  string `xml:"url,attr"`
		Type string `xml:"type,attr"`
	} `xml:"enclosure"`
}

func (i Item) Published() time.Time {
	t, err := parseDate(i.PubDate)
	if err != nil {
		log.Printf("\nCould not parse `%s`\n", i.PubDate)
	}

	return t
}

func (f Feed) Title() string {
	return strings.TrimSpace(f.Channel.Title)
}

func (f Feed) Description() string {
	return f.Channel.Description
}

func (f Feed) Image() string {
	if f.Channel.Image.HREF == "" {
		return f.Channel.Image.URL
	} else {
		return f.Channel.Image.HREF
	}
}

func (f Feed) Items() []Item {
	return f.Channel.Items
}

var dateFormats = []string{
	"01-02-2006",
	"01/02/2006",
	"01/02/2006 - 15:04",
	"01/02/2006 15:04:05 MST",
	"01/02/2006 3:04 PM",
	"02-01-2006",
	"02/01/2006",
	"02.01.2006 -0700",
	"02/01/2006 - 15:04",
	"02.01.2006 15:04",
	"02/01/2006 15:04:05",
	"02.01.2006 15:04:05",
	"02-01-2006 15:04:05 MST",
	"02/01/2006 15:04 MST",
	"02 Jan 2006",
	"02 Jan 2006 15:04:05",
	"02 Jan 2006 15:04:05 -0700",
	"02 Jan 2006 15:04:05 MST",
	"02 Jan 2006 15:04:05 UT",
	"02 Jan 2006 15:04 MST",
	"02 Monday, Jan 2006 15:04",
	"06-1-2 15:04",
	"06/1/2 15:04",
	"1/2/2006",
	"1/2/2006 15:04:05 MST",
	"1/2/2006 3:04:05 PM",
	"1/2/2006 3:04:05 PM MST",
	"15:04 02.01.2006 -0700",
	"2006-01-02",
	"2006/01/02",
	"2006-01-02 00:00:00.0 15:04:05.0 -0700",
	"2006-01-02 15:04",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05-0700",
	"2006-01-02 15:04:05 MST",
	"2006-01-02 15:04:05Z",
	"2006-01-02 at 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05:00",
	"2006-01-02T15:04:05 -0700",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05-0700",
	"2006-01-02T15:04:05:-0700",
	"2006-01-02T15:04:05-07:00:00",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04-07:00",
	"2006-01-02T15:04Z",
	"2006-1-02T15:04:05Z",
	"2006-1-2",
	"2006-1-2 15:04:05",
	"2006-1-2T15:04:05Z",
	"2006 January 02",
	"2-1-2006",
	"2/1/2006",
	"2.1.2006 15:04:05",
	"2 Jan 2006",
	"2 Jan 2006 15:04:05 -0700",
	"2 Jan 2006 15:04:05 MST",
	"2 Jan 2006 15:04:05 Z",
	"2 January 2006",
	"2 January 2006 15:04:05 -0700",
	"2 January 2006 15:04:05 MST",
	"6-1-2 15:04",
	"6/1/2 15:04",
	"Jan 02, 2006",
	"Jan 02 2006 03:04:05PM",
	"Jan 2, 2006",
	"Jan 2, 2006 15:04:05 MST",
	"Jan 2, 2006 3:04:05 PM",
	"Jan 2, 2006 3:04:05 PM MST",
	"January 02, 2006",
	"January 02, 2006 03:04 PM",
	"January 02, 2006 15:04",
	"January 02, 2006 15:04:05 MST",
	"January 2, 2006",
	"January 2, 2006 03:04 PM",
	"January 2, 2006 15:04:05",
	"January 2, 2006 15:04:05 MST",
	"January 2, 2006, 3:04 p.m.",
	"January 2, 2006 3:04 PM",
	"Mon, 02 Jan 06 15:04:05 MST",
	"Mon, 02 Jan 2006",
	"Mon, 02 Jan 2006 15:04:05",
	"Mon, 02 Jan 2006 15:04:05 00",
	"Mon, 02 Jan 2006 15:04:05 -07",
	"Mon 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 --0700",
	"Mon, 02 Jan 2006 15:04:05 -07:00",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon,02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 GMT-0700",
	"Mon , 02 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05MST",
	"Mon, 02 Jan 2006, 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05 MST -0700",
	"Mon, 02 Jan 2006 15:04:05 MST-07:00",
	"Mon, 02 Jan 2006 15:04:05 UT",
	"Mon, 02 Jan 2006 15:04:05 Z",
	"Mon, 02 Jan 2006 15:04 -0700",
	"Mon, 02 Jan 2006 15:04 MST",
	"Mon,02 Jan 2006 15:04 MST",
	"Mon, 02 Jan 2006 15 -0700",
	"Mon, 02 Jan 2006 3:04:05 PM MST",
	"Mon, 02 January 2006",
	"Mon,02 January 2006 14:04:05 MST",
	"Mon, 2006-01-02 15:04",
	"Mon, 2 Jan 06 15:04:05 -0700",
	"Mon, 2 Jan 06 15:04:05 MST",
	"Mon, 2 Jan 15:04:05 MST",
	"Mon, 2 Jan 2006",
	"Mon,2 Jan 2006",
	"Mon, 2 Jan 2006 15:04",
	"Mon, 2 Jan 2006 15:04:05",
	"Mon, 2 January 2006 15:04:05",
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05-0700",
	"Mon, 2 Jan 2006 15:04:05 -0700 MST",
	"mon,2 Jan 2006 15:04:05 MST",
	"Mon 2 Jan 2006 15:04:05 MST",
	"Mon, 2 Jan 2006 15:04:05 MST",
	"Mon, 2 Jan 2006 15:04:05MST",
	"Mon, 2 Jan 2006 15:04:05 UT",
	"Mon, 2 Jan 2006 15:04 -0700",
	"Mon, 2 Jan 2006, 15:04 -0700",
	"Mon, 2 Jan 2006 15:04 MST",
	"Mon, 2, Jan 2006 15:4",
	"Mon, 2 Jan 2006 15:4:5 -0700 GMT",
	"Mon, 2 Jan 2006 15:4:5 MST",
	"Mon, 2 Jan 2006 3:04:05 PM -0700",
	"Mon, 2 January 2006",
	"Mon, 2 January 2006 15:04:05 -0700",
	"Mon, 2 January 2006 15:04:05 MST",
	"Mon, 2 January 2006, 15:04:05 MST",
	"Mon, 2 January 2006, 15:04 -0700",
	"Mon, 2 January 2006 15:04 MST",
	"Monday, 02 January 2006 15:04:05",
	"Monday, 02 January 2006 15:04:05 -0700",
	"Monday, 02 January 2006 15:04:05 MST",
	"Monday, 2 Jan 2006 15:04:05 -0700",
	"Monday, 2 Jan 2006 15:04:05 MST",
	"Monday, 2 January 2006 15:04:05 -0700",
	"Monday, 2 January 2006 15:04:05 MST",
	"Monday, January 02, 2006",
	"Monday, January 2, 2006",
	"Monday, January 2, 2006 03:04 PM",
	"Monday, January 2, 2006 15:04:05 MST",
	"Mon Jan 02 2006 15:04:05 -0700",
	"Mon, Jan 02,2006 15:04:05 MST",
	"Mon Jan 02, 2006 3:04 pm",
	"Mon Jan 2 15:04:05 2006 MST",
	"Mon Jan 2 15:04 2006",
	"Mon, Jan 2 2006 15:04:05 -0700",
	"Mon, Jan 2 2006 15:04:05 -700",
	"Mon, Jan 2, 2006 15:04:05 MST",
	"Mon, Jan 2 2006 15:04 MST",
	"Mon, Jan 2, 2006 15:04 MST",
	"Mon, January 02, 2006 15:04:05 MST",
	"Mon, January 02, 2006, 15:04:05 MST",
	"Mon, January 2 2006 15:04:05 -0700",
	time.ANSIC,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RubyDate,
	time.UnixDate,
	"Updated January 2, 2006",
}

func parseDate(d string) (t time.Time, err error) {
	d = strings.TrimSpace(d)

	for _, f := range dateFormats {
		if t, err = time.Parse(f, d); err == nil {
			return
		}
	}

	err = fmt.Errorf("could not parse date: %v", d)
	return
}
