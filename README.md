Dutycycle
=========

A simple library to perform duty cycle calculations. States are stored as bits in a 64-bit word. If you need 32-bit support, make an issue in the tracker and we'll see what can be done about it. See godoc for (a bit) more info. Example usage:

```go
// it's a good idea to use length that's multiple of 256 for performance reasons.
dc := dutycycle.NewDutyCycle(1024)

dc.SetOn()
dc.SetOff()
dc.SetOn()
dc.SetOn()
dc.DutyCycle()	// 0.75
```

Guntars