package fiatshamir

const (
	Name = "FiatShamirAlgo"

	PhaseSetKey = "SET-KEY"
	PhaseProof  = "PROOF"
	PhaseCall   = "CALL"
	PhaseAnswer = "ANSWER"
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
