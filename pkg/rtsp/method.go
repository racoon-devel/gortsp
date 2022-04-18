package rtsp

// Method represents RTSP request method
type Method string

const (
	Options      Method = "OPTIONS"
	Describe     Method = "DESCRIBE"
	Announce     Method = "ANNOUNCE"
	Setup        Method = "SETUP"
	Play         Method = "PLAY"
	Pause        Method = "PAUSE"
	Teardown     Method = "TEARDOWN"
	GetParameter Method = "GET_PARAMETER"
	SetParameter Method = "SET_PARAMETER"
	Redirect     Method = "REDIRECT"
	Record       Method = "RECORD"
)

// IsValid returns true if the method value is known
func (m Method) IsValid() bool {
	methods := map[Method]bool{
		Options:      true,
		Describe:     true,
		Announce:     true,
		Setup:        true,
		Play:         true,
		Pause:        true,
		Teardown:     true,
		GetParameter: true,
		SetParameter: true,
		Redirect:     true,
		Record:       true,
	}

	_, ok := methods[m]
	return ok
}
