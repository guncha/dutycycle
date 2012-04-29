// Package dutycycle implements simple, but efficient duty cycle calculator. States are stored as
// bits in a machine word. Example usage:
//
//	dc := dutycycle.NewDutyCycle(1024)
//	dc.SetOn()
//	dc.SetOff()
//	dc.SetOn()
//	dc.SetOn()
//	dc.DutyCycle()	// 0.75
//
// DutyCycle() counts 4 64-bit words at a time so it's best to set length as a multiple of 256.
package dutycycle

type DutyCycle struct {
	bits []uint64
	idx  uint
	len  uint
	cap  uint
	cache float64
	cache_valid bool
}

func (dc *DutyCycle) SetOn() {
	dc.bits[dc.idx/64] |= 1 << (63 - dc.idx%64)
	dc.step()
}

func (dc *DutyCycle) SetOff() {
	dc.bits[dc.idx/64] &= ^(1 << (63 - dc.idx%64))
	dc.step()
}

func (dc *DutyCycle) step() {
	dc.idx++
	dc.idx = dc.idx % dc.cap
	if dc.len != dc.cap {
		dc.len++
	}
	dc.cache_valid = false
}

const c1 = 0x5555555555555555
const c2 = 0x3333333333333333
const c3 = 0x0F0F0F0F0F0F0F0F
const c4 = 0x000000FF000000FF

func pop(x uint64) int {
	
  x = x - ((x >> 1) & c1)
  x = (x & c2) + ((x >> 2) & c2)
  x = (x & c3) + ((x >> 4) & c3)

  x = x + (x >> 8)
  x = x + (x >> 16)
	x = x + (x >> 32)
  return int(x & 0xFF);	
}

// Re http://stackoverflow.com/a/1511920/357978
func pop4(x, y, z, w uint64) int {
		
    x = x - ((x >> 1) & c1)
    y = y - ((y >> 1) & c1)
		z = z - ((z >> 1) & c1)
		w = w - ((w >> 1) & c1)
		
    x = (x & c2) + ((x >> 2) & c2)
    y = (y & c2) + ((y >> 2) & c2)
		z = (z & c2) + ((z >> 2) & c2)
		w = (w & c2) + ((w >> 2) & c2)

		x = x + y
		z = z + w
		
    x = (x & c3) + ((x >> 4) & c3)
		z = (z & c3) + ((z >> 4) & c3)
		
		x = x + z
		
    x = x + (x >> 8)
    x = x + (x >> 16)
		x = x & c4
		x = x + (x >> 32)
    return int(x & 0x1FF);
}


// Returns duty cycle as a float64 variable in the range [0;1]
func (dc *DutyCycle) DutyCycle() float64 {
	
	if dc.cache_valid {
		return dc.cache
	}
	
	cnt := 0
	
	// do 4 words at a time
	l := (len(dc.bits) / 4) * 4
	for i := 0; i < l; i += 4 {
		cnt += pop4(dc.bits[i], dc.bits[i+1], dc.bits[i+2], dc.bits[i+3])
	}
	
	// reminder
	for i := l; i < len(dc.bits); i++ {
		cnt += pop(dc.bits[i])
	}
	
	dc.cache = float64(cnt) / float64(dc.len)
	dc.cache_valid = true
	return dc.cache
}

// Creates a new DutyCycle object and returns a pointer to it. Parameter length defines how many
// states should be remembered. This is essentially a binary moving average length.
func NewDutyCycle(length int) *DutyCycle {
	return &DutyCycle{
		bits: make([]uint64, (length+63)/64),
		idx:  0,
		len:  0,
		cap:  uint(length),
	}
}
