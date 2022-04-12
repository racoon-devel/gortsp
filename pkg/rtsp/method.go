package rtsp

type Method string

const (
	Options      Method = "OPTIONS"
	Describe     Method = "DESCRIBE"
	Setup        Method = "SETUP"
	Play         Method = "PLAY"
	Pause        Method = "PAUSE"
	Announce     Method = "ANNOUNCE"
	GetParameter Method = "GET_PARAMETER"
	SetParameter Method = "SET_PARAMETER"
	Record       Method = "RECORD"
)

func (m Method) IsValid() bool {
	methods := map[Method]bool{
		Options:      true,
		Describe:     true,
		Setup:        true,
		Play:         true,
		Pause:        true,
		Announce:     true,
		GetParameter: true,
		SetParameter: true,
		Record:       true,
	}

	_, ok := methods[m]
	return ok
}
