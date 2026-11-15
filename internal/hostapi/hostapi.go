package hostapi

var HostCall func(method string, payload []byte) ([]byte, error)
