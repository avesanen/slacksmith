package smith

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/websocket"
)

type Channel struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	IsChannel bool   `json:"is_channel"`
	IsGeneral bool   `json:"is_general"`
}

type Rtm struct {
	Ok   bool   `json:"ok"`
	Url  string `json:"url"`
	Self struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"self"`
	Team struct {
		Id   string `json:"id"`
		Name string `json:"name`
	} `json:"team"`
	Users    []User            `json:"users"`
	Channels []Channel         `json:"channels"`
	Groups   json.RawMessage   `json:"groups"`
	Ims      json.RawMessage   `json:"ims"`
	Bots     json.RawMessage   `json:"bots"`
	Error    string            `json:"error,omitempty"`
	Conn     *websocket.Conn   `json:"-"`
	OutChan  chan string       `json:"-"`
	InChan   chan MessageEvent `json:"-"`
}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Deleted  bool   `json:"deleted"`
	RealName string `json:"real_name"`
	IsAdmin  bool   `json:"is_admin"`
	IsOwned  bool   `json:"is_owner"`
	Presence string `json:"presence"`
	Profile  struct {
		FirstName          string `json:"first_name"`
		LastName           string `json:"last_name"`
		Phone              string `json:"phone"`
		RealNameNormalized string `json:"real_name_normalized"`
		Email              string `json:"email"`
	} `json:"profile"`
}

func (rtm *Rtm) Reader() {
	go func() {
		defer close(rtm.InChan)
		for {
			var b = make([]byte, 2048)
			var n int
			var err error
			n, err = rtm.Conn.Read(b)
			if err != nil {
				log.Println(err.Error())
				return
			}
			var msg Msg
			err = json.Unmarshal(b[:n], &msg)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			if msg.Type == "pong" {
				continue
			} else if msg.Type == "message" {
				var message MessageEvent
				err = json.Unmarshal(b[:n], &message)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				rtm.InChan <- message
			}
		}
	}()
}

func (rtm *Rtm) Writer() {
	go func() {
		for {
			select {
			case msg, ok := <-rtm.OutChan:
				if !ok {
					return
				}
				_, err := rtm.Conn.Write([]byte(msg))
				if err != nil {
					log.Println(err.Error())
					continue
				}
			}
		}
	}()
	go func() {
		for {
			<-time.After(time.Second * 1)
			rtm.OutChan <- "{\"id\":1234,\"type\":\"ping\"}"
		}
	}()
}

func (rtm *Rtm) Dial() error {
	u, err := url.Parse(rtm.Url)
	if err != nil {
		return err
	}
	u.Host += ":443"

	ws, err := websocket.Dial(u.String(), "", "http://localhost/")
	if err != nil {
		return err
	}

	rtm.Conn = ws
	return nil
}

func (rtm *Rtm) Say(s string, c string) error {
	m := &ChatMsg{
		Id:      1,
		Type:    "message",
		Channel: c,
		Text:    s,
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	rtm.OutChan <- string(b)
	return nil
}

func StartRtm(rtmUrl string) (*Rtm, error) {
	resp, err := http.Get(rtmUrl)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rtm Rtm
	if err := json.Unmarshal(body, &rtm); err != nil {
		return nil, err
	}

	if !rtm.Ok {
		return nil, errors.New(rtm.Error)
	}

	err = rtm.Dial()
	if err != nil {
		return nil, err
	}

	rtm.OutChan = make(chan string)
	rtm.InChan = make(chan MessageEvent)
	rtm.Reader()
	rtm.Writer()
	return &rtm, nil
}

func (rtm *Rtm) GetUser(id string) *User {
	for _, u := range rtm.Users {
		if u.Id == id {
			retU := u
			return &retU
		}
	}
	return nil
}

func (rtm *Rtm) GetChannel(id string) *Channel {
	for _, c := range rtm.Channels {
		if c.Id == id {
			retC := c
			return &retC
		}
	}
	return nil
}
