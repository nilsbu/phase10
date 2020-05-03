package game

type SeqType int

const (
	Invalid SeqType = iota
	Straight
	Kind
	Ambiguous
)

type Sequence struct {
	Type SeqType
	N    int
}

func validate(cards Cards) Sequence {
	if len(cards) == 0 {
		return Sequence{Invalid, 0}
	}

	straight := isValidStraight(cards)
	kind := isValidKind(cards)
	if straight && kind {
		return Sequence{Ambiguous, len(cards)}
	} else if straight {
		return Sequence{Straight, len(cards)}
	} else if kind {
		return Sequence{Kind, len(cards)}
	}
	return Sequence{Invalid, len(cards)}
}

func isValidStraight(cards Cards) bool {
	offset := -2
	for i, card := range cards {
		if card != 13 {
			if offset == -2 {
				offset = int(card) - i
				if offset < 1 {
					return false
				}
			} else if i+offset != int(card) {
				return false
			}
		}
		if i+offset > 12 {
			return false
		}
	}

	return true
}

func isValidKind(cards Cards) bool {
	kind := -1
	for _, card := range cards {
		if card != 13 {
			if kind == -1 {
				kind = int(card)
			} else if kind != int(card) {
				return false
			}
		}
	}

	return true
}
