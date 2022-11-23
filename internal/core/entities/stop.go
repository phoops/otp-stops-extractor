package entities

type Stop struct {
	GtfsIDs  []string
	Agencies []string
	Code     string
	Name     string
	Lat      float32
	Lon      float32
}
