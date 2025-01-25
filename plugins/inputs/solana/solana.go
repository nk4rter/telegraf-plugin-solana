package solana

import (
	"context"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type SolanaPlugin struct {
	valueName string `toml:"value_name"`

	sampleFrequency config.Duration `toml:"sample_frequency"`
	ctx             context.Context
	cancel          context.CancelFunc

	Log telegraf.Logger `toml:"-"`
}

func init() {
	inputs.Add("solana", func() telegraf.Input {
		return &SolanaPlugin{
			valueName:       "value",
			sampleFrequency: config.Duration(1 * time.Second),
		}
	})
}

func (s *SolanaPlugin) Init() error {
	return nil
}

func (s *SolanaPlugin) SampleConfig() string {
	return `
## Gathering info from the Solana blockchain
[[inputs.solana]]
  # The name of the measurement to write out to.
  value_name = "value"
  sample_frequency = "1000ms"
`
}

func (s *SolanaPlugin) Description() string {
	return "Gathering info from the Solana blockchain"
}

func (s *SolanaPlugin) Gather(a telegraf.Accumulator) error {
	s.sendMetric(a)
	return nil
}

// func (s *SolanaPlugin) Start(a telegraf.Accumulator) error {
// 	s.Log.Info("Started as service")

// 	s.ctx, s.cancel = context.WithCancel(context.Background())
// 	go func() {
// 		t := time.NewTicker(time.Duration(s.sampleFrequency))
// 		for {
// 			select {
// 			case <-s.ctx.Done():
// 				t.Stop()
// 				return
// 			case <-t.C:
// 				s.sendMetric(a)
// 			}
// 		}
// 	}()

// 	return nil
// }

// func (s *SolanaPlugin) Stop() {
// 	s.cancel()
// }

func (s *SolanaPlugin) sendMetric(a telegraf.Accumulator) {
	a.AddFields("hello",
		map[string]interface{}{
			s.valueName: "world",
		},
		nil,
	)
}
