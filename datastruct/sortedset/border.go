package sortedset

import "strconv"

const (
	negativeInf int8 = -1
	positiveInf int8 = 1
)

type ScoreBorder struct {
	Inf     int8
	Value   float64
	Exclude bool
}

func (border *ScoreBorder) greater(value float64) bool {
	if border.Inf == negativeInf {
		return false
	} else if border.Inf == positiveInf {
		return true
	}
	if border.Exclude {
		return border.Value > value
	} else {
		return border.Value >= value
	}
}

func (border *ScoreBorder) less(value float64) bool {
	if border.Inf == negativeInf {
		return true
	} else if border.Inf == positiveInf {
		return false
	}
	if border.Exclude {
		return border.Value < value
	} else {
		return border.Value <= value
	}
}

var positiveInfBorder = &ScoreBorder{
	Inf: positiveInf,
}

var negativeInfBorder = &ScoreBorder{
	Inf: negativeInf,
}

func ParseScoreBorder(s string) (*ScoreBorder, error) {
	if s == "inf" || s == "+inf" {
		return positiveInfBorder, nil
	} else if s == "-inf" {
		return negativeInfBorder, nil
	}
	if s[0] == '(' {
		value, err := strconv.ParseFloat(s[1:], 64)
		if err != nil {
			return nil, err
		}
		return &ScoreBorder{
			Inf:     0,
			Value:   value,
			Exclude: true,
		}, nil
	}
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return &ScoreBorder{
		Inf:     0,
		Value:   value,
		Exclude: false,
	}, nil
}
