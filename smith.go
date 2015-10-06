package smith

import "strings"

type Handler interface {
	Process(msg MessageEvent, smith *Smith)
}

type Smith struct {
	CmdPrefix string
	Handlers  map[string]Handler
	Loggers   []Handler
	Rtm       *Rtm
}

func NewSmith() *Smith {
	s := &Smith{}
	s.CmdPrefix = "!"
	s.Handlers = make(map[string]Handler)
	return s
}

func (s *Smith) Serve() {
	for {
		select {
		case msg, ok := <-s.Rtm.InChan:
			if !ok {
				return
			}
			if msg.Text != "" {
				s.Parse(msg)
			}
		}
	}
}

func (s *Smith) JackIn(token string) error {
	url := "https://slack.com/api/rtm.start?token=" + token
	rtm, err := StartRtm(url)
	if err != nil {
		return err
	}
	s.Rtm = rtm
	return nil
}

func (s *Smith) Parse(msg MessageEvent) {
	msg.User = s.Rtm.GetUser(msg.UserId)
	for _, k := range s.Loggers {
		k.Process(msg, s)
	}

	if !strings.HasPrefix(msg.Text, s.CmdPrefix) {
		return
	}
	msg.Text = strings.TrimLeft(msg.Text, "!")
	fields := strings.Fields(msg.Text)
	if len(fields) > 0 {
		for i, k := range s.Handlers {
			if i == fields[0] {
				msg.Text = strings.TrimLeft(msg.Text, i)
				msg.Text = strings.TrimLeft(msg.Text, " ")
				k.Process(msg, s)
				return
			}
		}
	}
	return
}

func (s *Smith) Handle(command string, handler Handler) {
	s.Handlers[command] = handler
}

func (s *Smith) HandleAll(handler Handler) {
	s.Loggers = append(s.Loggers, handler)
}
