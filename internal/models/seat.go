package models

// SeatPosition 座位位置信息
type SeatPosition struct {
	ID       string  `json:"id"`
	Label    string  `json:"label"`
	Type     string  `json:"type"`
	Left     float64 `json:"left"`
	Top      float64 `json:"top"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Rotation float64 `json:"rotation"`
}

// SeatLayout 座位布局信息
type SeatLayout struct {
	DJ         []SeatPosition `json:"dj"`
	Tables     []SeatPosition `json:"tables"`
	Cards      []SeatPosition `json:"cards"`
	VIPs       []SeatPosition `json:"vips"`
	FirstClass []SeatPosition `json:"firstClass"`
}

