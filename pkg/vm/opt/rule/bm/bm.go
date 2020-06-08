package bm

import "fmt"

func Gen(bs []Bm) string {
	switch {
	case len(bs) == 0:
		return ""
	case len(bs) == 1:
		return genSingal(bs[0])
	case bs[0].IsOr:
		return genOr(bs)
	default:
		return genAnd(bs)
	}
}

func genSingal(b Bm) string {
	if len(b.Bs) == 0 {
		return b.Name
	}
	return Gen(b.Bs)
}

func genOr(bs []Bm) string {
	if len(bs) == 2 {
		return fmt.Sprintf("bitmapOr(%s, %s)", genSingal(bs[0]), genSingal(bs[1]))
	}
	if bs[1].IsOr {
		return fmt.Sprintf("bitmapOr(bitmapOr(%s, %s), %s)", genSingal(bs[0]), genSingal(bs[1]), Gen(bs[2:]))
	}
	return fmt.Sprintf("bitmapAnd(bitmapOr(%s, %s), %s)", genSingal(bs[0]), genSingal(bs[1]), Gen(bs[2:]))
}

func genAnd(bs []Bm) string {
	if len(bs) == 2 {
		return fmt.Sprintf("bitmapAnd(%s, %s)", genSingal(bs[0]), genSingal(bs[1]))
	}
	if bs[1].IsOr {
		return fmt.Sprintf("bitmapOr(bitmapAnd(%s, %s), %s)", genSingal(bs[0]), genSingal(bs[1]), Gen(bs[2:]))
	}
	return fmt.Sprintf("bitmapAnd(bitmapAnd(%s, %s), %s)", genSingal(bs[0]), genSingal(bs[1]), Gen(bs[2:]))
}
