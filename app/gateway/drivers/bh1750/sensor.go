package bh1750

type BH1750 struct{}

func (s *BH1750) Name() string {
	return "BH1750"
}

func (s *BH1750) Read() (float64, error) {
	return 0, nil
}
