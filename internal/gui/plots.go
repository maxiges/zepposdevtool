package gui

import (
	"fmt"
	"sort"
	"sync"

	g "github.com/AllenDang/giu"
)

// lock synchronizes access to PlotData to ensure thread-safe concurrent updates.
var lock = sync.Mutex{}

// PlotData represents a collection of plot series with metadata for rendering.
//
// It manages multiple named data series (e.g., "system-used", "system-total"),
// tracks min/max values across all series, and maintains X-axis labels and ticks.
//
// Fields:
// - Data: map of series name to slice of float64 values
// - Ticks: list of X-axis tick marks with labels
// - TicksCustomLabels: map of data point indices to custom label strings
// - Min: minimum value across all data series
// - Max: maximum value across all data series
// - AxeXLen: total number of data points (length of X-axis)
type PlotData struct {
	Data              map[string][]float64
	Ticks             []g.PlotTicker
	TicksCustomLabels map[int]string
	Min               float64
	Max               float64
	AxeXLen           int
}

// Clear resets all plot data to an empty state.
//
// This reinitializes:
// - Min to a very large positive value
// - Max to a very large negative value
// - Data map, Ticks slice, and TicksCustomLabels map to empty collections
//
// Call this before populating a PlotData with new data to ensure
// a clean slate and accurate min/max calculations.
func (p *PlotData) Clear() {
	p.Min = 99999999999999
	p.Max = -999999999999999
	p.Data = map[string][]float64{}
	p.Ticks = []g.PlotTicker{}
	p.TicksCustomLabels = map[int]string{}
}

// MinMaxDiff returns the difference between max and min values.
//
// This is useful for calculating appropriate axis ranges with padding,
// e.g., to add a 10% margin around the data: max + MinMaxDiff()*0.1
//
// Returns: Max - Min
func (p *PlotData) MinMaxDiff() float64 {
	return p.Max - p.Min
}

// Add appends a value to a named data series and updates min/max tracking.
//
// If the series does not exist, it is created. The method is thread-safe
// via mutex locking.
//
// Parameters:
// - value: the data point to add (typically a memory metric in bytes)
// - plotName: the name of the series this value belongs to (e.g., "system-used")
//
// Thread-safety: This method is protected by a package-level mutex.
func (p *PlotData) Add(value float64, plotName string) {
	lock.Lock()
	defer lock.Unlock()

	// Update min/max boundaries
	if value < p.Min {
		p.Min = value
	}
	if value > p.Max {
		p.Max = value
	}

	// Get or create the series, then append the value
	val, exist := p.Data[plotName]
	if !exist {
		val = []float64{}
	}
	val = append(val, value)

	p.Data[plotName] = val

}

// IsPackageExist checks if a named series exists and returns its length.
//
// Parameters:
// - plotName: the name of the series to check
//
// Returns:
// - size: the number of data points in the series (0 if it doesn't exist)
// - isExist: true if the series exists, false otherwise
//
// Thread-safety: This method is protected by a package-level mutex.
func (p *PlotData) IsPackageExist(plotName string) (size int, isExist bool) {
	lock.Lock()
	defer lock.Unlock()

	val, exist := p.Data[plotName]
	return len(val), exist
}

// GenerateX builds X-axis ticks and labels for rendering.
//
// This function:
// 1. Determines the maximum data length across all series
// 2. Creates tick marks at regular intervals (every 10 data points)
// 3. Uses custom labels where provided, otherwise uses numeric indices
// 4. Sets AxeXLen to the total number of data points
//
// Call this after all data has been added but before rendering.
// The resulting Ticks are used to label the X-axis of the plot.
func (p *PlotData) GenerateX() {
	maxDataLen := 0
	for _, data := range p.Data {
		if len(data) > maxDataLen {
			maxDataLen = len(data)
		}
	}
	p.AxeXLen = maxDataLen
	labelIndex := maxDataLen / 30
	if labelIndex >= 0 {
		labelIndex = 1
	}

	for i := 0; i < maxDataLen; i += labelIndex {
		label := fmt.Sprintf("%d", i)
		if value, ok := p.TicksCustomLabels[i]; ok {
			label = value
		}
		p.Ticks = append(p.Ticks, g.PlotTicker{Position: float64(i), Label: label})
	}
}

// AddCustomLabel associates a custom label string with a specific data point index.
//
// These labels are used by GenerateX to display descriptive text on the X-axis
// instead of numeric indices. Typical use: storing descriptions or timestamps.
//
// Parameters:
// - index: the data point index (position in the series)
// - label: the custom label to display at that position
func (p *PlotData) AddCustomLabel(index int, label string) {
	p.TicksCustomLabels[index] = label
}

type sortedMap struct {
	key  string
	data []float64
}

// GetPlots returns all data series as a slice of giu plot widgets.
//
// Each series is converted to a g.Line widget for rendering in the GUI.
// The order of widgets in the returned slice may vary due to map iteration.
//
// Returns: slice of g.PlotWidget, one per named series in Data
func (p *PlotData) GetPlots() []g.PlotWidget {

	ret := make([]g.PlotWidget, 0, len(p.Data))

	sortedMapList := make([]sortedMap, 0, len(p.Data))

	for key, value := range p.Data {
		mapVal := sortedMap{
			key:  key,
			data: value,
		}
		sortedMapList = append(sortedMapList, mapVal)
	}

	sort.Slice(sortedMapList, func(i, j int) bool {
		return sortedMapList[i].key > sortedMapList[j].key

	})

	for _, value := range sortedMapList {
		ret = append(ret, g.Line(value.key, value.data))
	}

	return ret
}
