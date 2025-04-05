package gui

import (
	"fmt"
	"zepp-os-dev-tool/internal/models"
	"zepp-os-dev-tool/internal/storage"

	g "github.com/AllenDang/giu"
	"github.com/inhies/go-bytesize"
)

var (
	plot1Data PlotData
	plot2Data PlotData
	plot3Data PlotData

	dashboardData models.TemplateData
)

func GenerateDataForData(appName string) {
	dashboardData.AppName = appName
	data, _ := storage.GetDataForApp(appName)

	dashboardData = models.TemplateData{
		AppName: appName,
		Total:   data[len(data)-1].Memory.System.Total,
		TotalKB: data[len(data)-1].Memory.System.Total / 1024,
	}
	plot1Data.Clear()
	plot2Data.Clear()
	plot3Data.Clear()

	bSize := bytesize.New(float64(dashboardData.Total))
	usedVal := data[len(data)-1].Memory.System.Used
	usedData := bytesize.New(float64(usedVal))

	*MemoryLabel = *g.Label(fmt.Sprintf("Total memory: %d (%s)  Used: %s  (%.2f %%)", dashboardData.Total, bSize, usedData, float64(usedVal*100)/float64(dashboardData.Total)))

	for i, dataRecord := range data {
		// SYSTEM

		plot1Data.Add(float64(dataRecord.Memory.System.Used), "system-used")
		plot1Data.AddCustomLabel(i, dataRecord.Description)

		plot1Data.Add(float64(dataRecord.Memory.System.Total), "system-total")
		plot1Data.AddCustomLabel(i, dataRecord.Description)

		plot2Data.Add(float64(dataRecord.Memory.App[0].Used), "app-used")
		plot2Data.AddCustomLabel(i, dataRecord.Description)

		plot2Data.Add(float64(dataRecord.Memory.System.Total), "system-total")
		plot2Data.AddCustomLabel(i, dataRecord.Description)

		// Modules

		for _, moduleVal := range dataRecord.Memory.App[0].Modules {
			element := moduleVal.File + "(Used)"
			if i > 0 {
				packageLen, _ := plot3Data.IsPackageExist(element)
				if packageLen < i {
					for range i - packageLen {
						plot3Data.Add(0, element)
					}
				}
			}
			plot3Data.Add(float64(moduleVal.Used), element)
		}
		plot3Data.AddCustomLabel(i, dataRecord.Description)
	}
}

func TryRefreshUI(appName string) {
	go GenerateMenu()
	if !AutoRefreshIsOn {
		return
	}
	if currentMenu == MenuPlot && dashboardData.AppName == appName {
		DrawChartsForApp(appName)
	}

}
func DrawChartsForApp(appName string) {

	if currentMenu != MenuPlot {
		currentMenu = MenuPlot
		Plot1 = g.Plot("Memory Usage - by watch - total consumption")
		Plot2 = g.Plot("Memory Usage - by app")
		Plot3 = g.Plot("Memory Usage - by module")
		MemoryLabel = g.Label("Total memory: ---")

		MainUILayout = g.SplitLayout(g.DirectionHorizontal, &sashPos2,
			g.Layout{
				g.Button("Erase all data").OnClick(func() {
					storage.ClearAllDataForApp(appName)
				}),
				MemoryLabel,
				g.SplitLayout(g.DirectionVertical, &sashPos1,
					g.Layout{
						Plot1,
					},
					g.Layout{
						Plot2,
					},
				).SplitRefType(g.SplitRefRight),
			},
			g.Layout{
				Plot3,
			},
		).SplitRefType(g.SplitRefRight)

	}

	GenerateDataForData(appName)

	plot1Data.GenerateX()
	plot2Data.GenerateX()
	plot3Data.GenerateX()

	chartHmax := windowH/2 - 100

	Plot1.AxisLimits(0, float64(plot1Data.AxeXLen), plot1Data.Min-plot1Data.MinMaxDiff()*0.1, plot1Data.Max+plot1Data.MinMaxDiff()*0.1, g.ConditionOnce).XTicks(plot1Data.Ticks, false).Plots(plot1Data.GetPlots()...).Size(windowW/2, chartHmax)
	Plot2.AxisLimits(0, float64(plot2Data.AxeXLen), plot2Data.Min-plot2Data.MinMaxDiff()*0.1, plot2Data.Max+plot2Data.MinMaxDiff()*0.1, g.ConditionOnce).XTicks(plot2Data.Ticks, false).Plots(plot2Data.GetPlots()...).Size(windowW/2, chartHmax)
	Plot3.AxisLimits(0, float64(plot3Data.AxeXLen), plot3Data.Min-plot3Data.MinMaxDiff()*0.1, plot3Data.Max+plot3Data.MinMaxDiff()*0.1, g.ConditionOnce).XTicks(plot3Data.Ticks, false).Plots(plot3Data.GetPlots()...).Size(windowW, chartHmax)
}
