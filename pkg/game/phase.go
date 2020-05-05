package game

import "errors"

func isPhaseFulfilled(seqs []Sequence, phase int) bool {
	expected, err := GetPhaseSequences(phase)
	if err != nil {
		return false // TODO: should this be an error
	}

	for _, seq := range seqs {
		found := false
		for j, ex := range expected {
			if seq.Fulfills(ex) {
				found = true
				expected = append(expected[:j], expected[j+1:]...)
				break
			}
		}
		if !found {
			return false
		}
	}
	return len(expected) == 0
}

func GetPhaseSequences(phase int) ([]Sequence, error) {
	switch phase {
	case 1:
		return []Sequence{{Kind, 3}, {Kind, 3}}, nil
	case 2:
		return []Sequence{{Kind, 3}, {Straight, 4}}, nil
	case 3:
		return []Sequence{{Kind, 4}, {Straight, 4}}, nil
	case 4:
		return []Sequence{{Straight, 7}}, nil
	case 5:
		return []Sequence{{Straight, 8}}, nil
	case 6:
		return []Sequence{{Straight, 9}}, nil
	case 7:
		return []Sequence{{Kind, 4}, {Kind, 4}}, nil
	case 8:
		return []Sequence{{Kind, 5}, {Kind, 2}}, nil
	case 9:
		return []Sequence{{Kind, 5}, {Kind, 3}}, nil
	case 10:
		return []Sequence{{Kind, 5}, {Straight, 5}}, nil
	default:
		return []Sequence{}, errors.New("invalid phase")
	}
}
