package modules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/evilsocket/bettercap-ng/core"
	"github.com/evilsocket/bettercap-ng/net"
	"github.com/evilsocket/bettercap-ng/session"
)

const eventTimeFormat = "2006-01-02 15:04:05"

func (s EventsStream) viewLogEvent(e session.Event) {
	fmt.Printf("[%s] [%s] (%s) %s\n",
		e.Time.Format(eventTimeFormat),
		core.Green(e.Tag),
		e.Label(),
		e.Data.(session.LogMessage).Message)
}

func (s EventsStream) viewTargetEvent(e session.Event) {
	t := e.Data.(*net.Endpoint)
	fmt.Printf("[%s] [%s] %s\n",
		e.Time.Format(eventTimeFormat),
		core.Green(e.Tag),
		t)
}

func (s EventsStream) viewModuleEvent(e session.Event) {
	fmt.Printf("[%s] [%s] %s\n",
		e.Time.Format(eventTimeFormat),
		core.Green(e.Tag),
		e.Data)
}

func (s EventsStream) viewSnifferEvent(e session.Event) {
	se := e.Data.(SnifferEvent)

	fmt.Printf("[%s] [%s] %s > %s | ",
		e.Time.Format(eventTimeFormat),
		core.Green(e.Tag),
		se.Source,
		se.Destination)

	keys := make([]string, 0)
	for k, _ := range se.Data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s=%s ", core.Green(k), se.Data[k])
	}

	fmt.Println()
}

func (s *EventsStream) View(e session.Event, refresh bool) {
	if s.filter == "" || strings.Contains(e.Tag, s.filter) {
		if e.Tag == "sys.log" {
			s.viewLogEvent(e)
		} else if strings.HasPrefix(e.Tag, "target.") {
			s.viewTargetEvent(e)
		} else if strings.HasPrefix(e.Tag, "mod.") {
			s.viewModuleEvent(e)
		} else if strings.HasPrefix(e.Tag, "net.sniff.") {
			s.viewSnifferEvent(e)
		} else {
			fmt.Printf("[%s] [%s] %v\n", e.Time.Format(eventTimeFormat), core.Green(e.Tag), e)
		}

		if refresh {
			s.Session.Refresh()
		}
	}
}
