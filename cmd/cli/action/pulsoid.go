package action

import (
	"github.com/fasthttp/websocket"
	"github.com/gizmo-ds/pulsoid-vrchat-osc/internal/global"
	"github.com/gizmo-ds/pulsoid-vrchat-osc/pkg/pulsoid"
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
	d := osc.NewStandardDispatcher()
	_ = d.AddMsgHandler("/avatar/change", func(msg *osc.Message) {
		if id, ok := msg.Arguments[0].(string); ok {
			log.Debug().Str("AvatarID", id).Msg("Avatar changed")
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
			log.Error().Err(err).Msg("Could not read from pulsoid")
			continue
		}
		if p.enabled {
			msg := osc.NewMessage("/avatar/parameters/OSC_HeartRate")
			msg.Append(float32(result.Data.HeartRate/(220/2) - 1))
			if err = p.client.Send(msg); err != nil {
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
