package action

import (
	"errors"
	"strings"
	"time"

	"github.com/gizmo-ds/pulsoid-vrchat-osc/internal/global"
	"github.com/gizmo-ds/pulsoid-vrchat-osc/pkg/pulsoid"
	"github.com/gorilla/websocket"
	"github.com/hypebeast/go-osc/osc"
	"github.com/rs/zerolog/log"
)

type Pulsoid struct {
	RamielUrl string
	WidgetID  string
	client    *osc.Client
	enabled   bool
}

func NewPulsoid() *Pulsoid {
	p := &Pulsoid{
		WidgetID: global.Config.WidgetID,
		client:   osc.NewClient("127.0.0.1", global.Config.VRChat.Port),
		enabled:  true,
	}
	return p
}

func (p *Pulsoid) startOscServer() {
	if global.Config.Address == "" || len(global.Config.EnableAvatars) == 0 {
		log.Info().Msg("OSC server disabled")
		return
	}
	d := osc.NewStandardDispatcher()
	_ = d.AddMsgHandler("*", func(msg *osc.Message) {
		if msg.Address == "/avatar/change" {
			if id, ok := msg.Arguments[0].(string); ok {
				log.Debug().Str("AvatarID", id).Msg("Avatar changed")
				enabled := false
				for _, eid := range global.Config.EnableAvatars {
					if eid == id || (eid == "local" && strings.HasPrefix(id, "local:")) {
						enabled = true
						break
					}
				}
				log.Info().
					Bool("Enabled", enabled).
					Msg("OSC enabled")
				p.enabled = enabled
			}
		}
	})
	server := &osc.Server{
		Addr:       global.Config.Address,
		Dispatcher: d,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start OSC server")
	}
}

func (p *Pulsoid) Start() {
	go p.startOscServer()

	p.GetRamielUrl()
	conn, _, err := websocket.DefaultDialer.Dial(p.RamielUrl, nil)
	if err != nil {
		log.Fatal().Err(err).Caller().Msg("Could not connect to pulsoid")
	}
	defer conn.Close()

	log.Info().Msg("Pulsoid connected")

	for {
		var result pulsoid.WebSocketResult
		err = conn.ReadJSON(&result)
		if err != nil {
			var e *websocket.CloseError
			if errors.As(err, &e) {
				log.Error().Err(e).Msg("Pulsoid disconnected")
				return
			}
			log.Error().Err(err).Msg("Could not read from pulsoid")
			continue
		}
		if p.enabled {
			if global.Config.FloatParameterName == "" || global.Config.IntParameterName == "" {
				log.Warn().Int("HeartRate", result.Data.HeartRate).Msg("HeartRate not sent")
				continue
			}

			bundle := osc.NewBundle(time.Now())

			if global.Config.FloatParameterName != "" {
				floatMsg := osc.NewMessage("/avatar/parameters/" + global.Config.FloatParameterName)
				floatMsg.Append(float32(clampIntRange(result.Data.HeartRate, 0, 254))/127 - 1)
				bundle.Append(floatMsg)
			}

			if global.Config.IntParameterName != "" {
				intMsg := osc.NewMessage("/avatar/parameters/" + global.Config.IntParameterName)
				intMsg.Append(clampIntRange(result.Data.HeartRate, 0, 255))
				bundle.Append(intMsg)
			}

			if err = p.client.Send(bundle); err != nil {
				log.Error().Err(err).Msg("Could not send OSC message")
				continue
			}
			log.Info().Int("HeartRate", result.Data.HeartRate).Msg("HeartRate sent")
		}
	}
}

func (p *Pulsoid) GetRamielUrl() {
	u, err := pulsoid.GetRamielUrl(p.WidgetID)
	if err != nil {
		log.Fatal().Err(err).Caller().Msg("GetRamielUrl")
	} else {
		p.RamielUrl = u
	}
}

func clampIntRange(number int, min int, max int) int {
	if number < min {
		number = min
	}
	if number > max {
		number = max
	}
	return number
}
