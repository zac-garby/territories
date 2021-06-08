package server

// commands
var (
	CMD_GENERATE = []byte("GEN")
	CMD_POLYGON  = []byte("POL")
)

// responses
var (
	RESP_GEN     = []byte("GENERATED")
	RESP_POLYGON = []byte("POLYGONS")
	RESP_NOGAME  = []byte("NOGAME")
	RESP_INVALID = []byte("INVALID")
)
