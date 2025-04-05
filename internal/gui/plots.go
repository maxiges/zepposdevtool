package gui

import (
	"fmt"
	"sync"

	g "github.com/AllenDang/giu"
)

var lock = sync.Mutex{}

type PlotData struct {
	Data              map[string][]float64
	Ticks             []g.PlotTicker
	TicksCustomLabels map[int]string
	Min               float64
	Max               float64
	AxeXLen           int
}

func (p *PlotData) Clear() {
	p.Min = 99999999999999
	p.Max = -999999999999999
	p.Data = map[string][]float64{}
	p.Ticks = []g.PlotTicker{}
	p.TicksCustomLabels = map[int]string{}
}
func (p *PlotData) MinMaxDiff() float64 {
	return p.Max - p.Min
}

func (p *PlotData) Add(value float64, plotName string) {
	lock.Lock()
	defer lock.Unlock()

	if value < p.Min {
		p.Min = value
	}
	if value > p.Max {
		p.Max = value
	}

	val, exist := p.Data[plotName]
	if !exist {
		val = []float64{}
	}
	val = append(val, value)

	p.Data[plotName] = val

}

func (p *PlotData) IsPackageExist(plotName string) (size int, isExist bool) {
	lock.Lock()
	defer lock.Unlock()

	val, exist := p.Data[plotName]
	return len(val), exist
}

func (p *PlotData) GenerateX() {
	maxDataLen := 0
	for _, data := range p.Data {
		if len(data) > maxDataLen {
			maxDataLen = len(data)
		}
	}
	p.AxeXLen = maxDataLen

	for i := 0; i < maxDataLen; i += 10 {
		label := fmt.Sprintf("%d", i)
		if value, ok := p.TicksCustomLabels[i]; ok {
			label = value
		}
		p.Ticks = append(p.Ticks, g.PlotTicker{Position: float64(i), Label: label})
	}
}
func (p *PlotData) AddCustomLabel(index int, label string) {
	p.TicksCustomLabels[index] = label
}
func (p *PlotData) GetPlots() []g.PlotWidget {

	ret := make([]g.PlotWidget, 0, len(p.Data))
	for key, value := range p.Data {
		ret = append(ret, g.Line(key, value))
	}
	return ret
}
