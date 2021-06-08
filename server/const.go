package server

// commands
var (
	CMD_GENERATE = []byte("GEN")
	CMD_POLYGON  = []byte("POL")
	CMD_CENTROID = []byte("CEN")
)

// responses
var (
	RESP_GEN      = []byte("GENERATED")
	RESP_POLYGON  = []byte("POLYGONS")
	RESP_CENTROID = []byte("CENTROIDS")

	RESP_NOGAME  = []byte("NOGAME")
	RESP_INVALID = []byte("INVALID")
)
