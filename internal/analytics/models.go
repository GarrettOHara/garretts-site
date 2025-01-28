package analytics

type TimeSeriesData struct {
    Time  string
    Count int
}

type Analytics struct {
    DeviceStats      []DeviceStat
    BrowserStats     []BrowserStat
    PlatformStats    []PlatformStat
    RequestsOverTime []TimeSeriesData
    DistinctIPCount  int
}

type DeviceStat struct {
    DeviceType string
    Count      int
    Percentage float64
}

type BrowserStat struct {
    Browser    string
    Count      int
    Percentage float64
}

type PlatformStat struct {
    Platform   string
    Count      int
    Percentage float64
}
