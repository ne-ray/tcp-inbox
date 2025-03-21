package hashcash

const (
	Name = "Hashcash"

	LenRandString = 255
)

type Public struct {
	Data string `json:"data"`
}
