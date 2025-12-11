package game

type Score int32

func (s *Score) Add(d int32) {
	*s += Score(d)
}

func (s *Score) Sub(d int32) {
	*s -= Score(d)
	if *s < 0 {
		*s = 0
	}
}
