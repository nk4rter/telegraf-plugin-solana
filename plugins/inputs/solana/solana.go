package solana

import (
	"context"
	// "fmt"
	"os"
	"os/signal"
	"syscall"

	// "fmt"
	// "math/big"
	// "time"

	"github.com/influxdata/telegraf"
	// "github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/inputs"

	// "github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type SolanaPlugin struct {
	Pubkey string `toml:"pubkey"`

	Log telegraf.Logger `toml:"-"`

	ctx    context.Context
	cancel context.CancelFunc

	client *rpc.Client
}

func init() {

	inputs.Add("solana", func() telegraf.Input {
		return &SolanaPlugin{}
	})
}

func (s *SolanaPlugin) Init() error {
	s.client = rpc.New(rpc.TestNet_RPC)
	return nil
}

func (s *SolanaPlugin) SampleConfig() string {
	return `
## Gathering info from the Solana blockchain
[[inputs.solana]]
  pubkey = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
`
}

func (s *SolanaPlugin) Description() string {
	return "Gathering info from the Solana blockchain"
}

func (s *SolanaPlugin) Gather(a telegraf.Accumulator) error {
	s.sendMetric(a)
	return nil
}

func (s *SolanaPlugin) Start(a telegraf.Accumulator) error {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1)
	go func() {
		for {
			<-sigs
			s.sendMetric(a)
		}
	}()

	return nil
}

func (s *SolanaPlugin) Stop() {
	s.cancel()
}

func (s *SolanaPlugin) sendMetric(a telegraf.Accumulator) {
	var pubKey = solana.MustPublicKeyFromBase58(s.Pubkey)

	{
		out, err := s.client.GetClusterNodes(s.ctx)
		if err != nil {
			s.Log.Error(err)
		} else {
			for _, v := range out {
				if v.Pubkey == pubKey {
					a.AddFields("clusterNode",
						map[string]interface{}{
							"pubkey":       v.Pubkey,
							"gossip":       v.Gossip,
							"tpu":          v.TPU,
							"tpuQuic":      v.TPUQUIC,
							"pubsub":       v.PubSub,
							"rpc":          v.RPC,
							"version":      v.Version,
							"featureSet":   v.FeatureSet,
							"shredVersion": v.ShredVersion,
						},
						nil,
					)
					break
				}
			}
		}
	}

	{
		out, err := s.client.GetVoteAccounts(s.ctx, nil)
		if err != nil {
			s.Log.Error(err)
		} else {
			for _, v := range out.Current {
				if v.NodePubkey == pubKey {
					a.AddFields("voteAccount",
						map[string]interface{}{
							"votePubkey":       v.VotePubkey.String(),
							"nodePubkey":       v.NodePubkey.String(),
							"activatedStake":   v.ActivatedStake,
							"epochVoteAccount": v.EpochVoteAccount,
							"commission":       v.Commission,
							"lastVote":         v.LastVote,
							"rootSlot":         v.RootSlot,
							"epochCredits":     v.EpochCredits,
						},
						nil,
					)
					break
				}
			}
		}
	}

	{
		out, err := s.client.GetBlockProduction(s.ctx)
		if err != nil {
			s.Log.Error(err)
		} else {
			v := out.Value.ByIdentity[pubKey]
			a.AddFields("blockProduction",
				map[string]interface{}{
					"leaderSlots":    v[0],
					"blocksProduces": v[1],
				},
				nil,
			)
		}
	}

	{
		out, err := s.client.GetBalance(
			s.ctx,
			pubKey,
			rpc.CommitmentFinalized,
		)
		if err != nil {
			s.Log.Error(err)
		} else {
			// lamportsOnAccount := new(big.Float).SetUint64(uint64(out.Value))
			// solBalance := new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))
			a.AddFields("balance",
				map[string]interface{}{
					"pubkey": pubKey.String(),
					"value":  out.Value,
				},
				nil,
			)
		}
	}
}
