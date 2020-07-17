package static

import "fmt"

func (a *Ints) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]int64, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Int8s) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]int8, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Int16s) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]int16, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Int32s) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]int32, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Int64s) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]int64, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Uint8s) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]uint8, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Uint16s) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]uint16, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Uint32s) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]uint32, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Uint64s) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]uint64, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Bools) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]bool, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Timestamps) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]int64, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Floats) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]float64, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Float32s) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]float32, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}

func (a *Float64s) String() string {
	vs := a.Vs
	if a.Is != nil {
		vs = make([]float64, len(a.Is))
		for i, o := range a.Is {
			vs[i] = a.Vs[o]
		}
	}
	return fmt.Sprintf("%v", vs)
}
