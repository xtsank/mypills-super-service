package req

type DosageRuleDto struct {
	ValueFrom           int
	ValueTo             int
	Type                DosageType
	DosageValue         float32
	NumberOfDosesPerDay int
}

type DosageType string

const (
	ByWeight DosageType = "weight"
	ByAge    DosageType = "age"
)
