package utils

import "strconv"

func (env *Env) string2float64(s string) (f64 float64) {

	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Errorf("error converting string %v\t to float64: err: %v\n", s, err)
		return
	}
	return i
}
