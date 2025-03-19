package fiatshamir


const (
	Name = "FiatShamirAlgo"

	PhaseSetKey = "SET-KEY"
	Phase1      = "1"
	Phase2      = "2"
)

const randFrom = 1024
const randTo = 1024 * 1024

type Private struct {
	LastPhase string
	Q         int64
	P         int64
}

type Public struct {
	N         uint64 `json:"start_key"`
	PublicKey uint64 `json:"public_key"`
}