package main

import "fmt"

func is_nigiri(ct int) bool {
	return ct == NIGIRI_1 || ct == NIGIRI_2 || ct == NIGIRI_3
}

// transform a NIGIRI_n into a NIGIRI_n_ON_WASABI
func wasabiify(ct int) (int, error) {
	switch ct {
	case NIGIRI_1:
		return NIGIRI_1_ON_WASABI, nil
	case NIGIRI_2:
		return NIGIRI_2_ON_WASABI, nil
	case NIGIRI_3:
		return NIGIRI_3_ON_WASABI, nil
	default:
		return -1, fmt.Errorf("wasabiify received non-nigiri card type %v (expected one of %v, %v, %v)", ct, NIGIRI_1, NIGIRI_2, NIGIRI_3)
	}
}
