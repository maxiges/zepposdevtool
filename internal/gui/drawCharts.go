package gui

import (
	"fmt"
	"sort"
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

	valueUnitDivider = 1024

	autoDiagramLast30s = false
)

// TryRefreshUI triggers a GUI update when appropriate.
//
// It regenerates the menu in the background and, if auto-refresh is enabled
// and the plots are currently visible for the same app, redraws the charts.
func TryRefreshUI(appName string) {
	go GenerateMenu()
	if !AutoRefreshIsOn {
		return
	}
	if currentMenu == MenuPlot && dashboardData.AppName == appName {
		DrawChartsForApp(appName)
	}
}

func changeUnits(divider int, appName string) {
	valueUnitDivider = divider
	TryRefreshUI(appName)
}

func DrawChartsForApp(appName string) {

	if currentMenu != MenuPlot {
		currentMenu = MenuPlot
		Plot1 = g.Plot("Memory Usage - by watch - total consumption")
		Plot2 = g.Plot("Memory Usage - by app")
		Plot3 = g.Plot("Memory Usage - by module")

		Plot1.XAxeFlags(g.PlotAxisFlagsAutoFit)
		Plot2.XAxeFlags(g.PlotAxisFlagsAutoFit)
		Plot3.XAxeFlags(g.PlotAxisFlagsAutoFit)

		Plot1.YAxeFlags(g.PlotAxisFlagsAutoFit, g.PlotAxisFlagsNone, g.PlotAxisFlagsNone)
		Plot2.YAxeFlags(g.PlotAxisFlagsAutoFit, g.PlotAxisFlagsNone, g.PlotAxisFlagsNone)
		Plot3.YAxeFlags(g.PlotAxisFlagsAutoFit, g.PlotAxisFlagsNone, g.PlotAxisFlagsNone)

		checked1 := valueUnitDivider == 1
		checked2 := valueUnitDivider == 1024
		checked3 := valueUnitDivider == 1024*1024
		SelectMem1 = g.Checkbox("Baits", &checked1).OnChange(
			func() {
				checked1 = true
				checked2 = false
				checked3 = false

				changeUnits(1, appName)
			},
		)

		SelectMem2 = g.Checkbox("Kilo Baits", &checked2).OnChange(
			func() {
				checked1 = false
				checked2 = true
				checked3 = false
				changeUnits(1024, appName)
			},
		)

		SelectMem3 = g.Checkbox("Mega Baits", &checked3).OnChange(
			func() {
				checked1 = false
				checked2 = false
				checked3 = true
				changeUnits(1024*1024, appName)
			},
		)

		isChecked := autoDiagramLast30s

		AutoMove = g.Checkbox("Auto Move", &isChecked).OnChange(
			func() {
				if autoDiagramLast30s {
					autoDiagramLast30s = false
					isChecked = false
				} else {
					autoDiagramLast30s = true
					isChecked = true
				}

				TryRefreshUI(appName)
			},
		)

		MemoryLabel = g.Label("Total memory: ---")
		MemorySizeUsedBar = g.ProgressBar(0)

		MainUILayout = g.SplitLayout(g.DirectionHorizontal, &sashPos2,
			g.Layout{
				g.Button("Erase all data").OnClick(func() {
					storage.ClearAllDataForApp(appName)
				}),
				g.Row(MemoryLabel, MemorySizeUsedBar),
				g.Row(SelectMem1, SelectMem2, SelectMem3, AutoMove),
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

	if autoDiagramLast30s {

		minX := 0
		if float64(plot1Data.AxeXLen) > 30 {
			minX = plot1Data.AxeXLen - 30
		}

		Plot1.AxisLimits(float64(minX), float64(plot1Data.AxeXLen), plot1Data.Min, plot1Data.Max, g.ConditionFirstUseEver)
	}
	chartHMax := windowH/2 - 100

	Plot1.AxisLimits(0, float64(plot1Data.AxeXLen), plot1Data.Min-plot1Data.MinMaxDiff()*0.1, plot1Data.Max+plot1Data.MinMaxDiff()*0.1, g.ConditionOnce).XTicks(plot1Data.Ticks, false).Plots(plot1Data.GetPlots()...).Size(windowW/2, chartHMax)
	Plot2.AxisLimits(0, float64(plot2Data.AxeXLen), plot2Data.Min-plot2Data.MinMaxDiff()*0.1, plot2Data.Max+plot2Data.MinMaxDiff()*0.1, g.ConditionOnce).XTicks(plot2Data.Ticks, false).Plots(plot2Data.GetPlots()...).Size(windowW/2, chartHMax)
	Plot3.AxisLimits(0, float64(plot3Data.AxeXLen), plot3Data.Min-plot3Data.MinMaxDiff()*0.1, plot3Data.Max+plot3Data.MinMaxDiff()*0.1, g.ConditionOnce).XTicks(plot3Data.Ticks, false).Plots(plot3Data.GetPlots()...).Size(windowW, chartHMax)
}

// GenerateDataForData builds plot data for a given application.
//
// It reads all stored samples for `appName`, computes summary values,
// and populates the plot data structures used by the GUI.
//
// Notes:
//   - If no data is available for `appName`, the function clears plots and returns.
//   - This function does not perform any rendering itself; it only mutates the
//     `plot*Data` variables and updates `dashboardData` and `MemoryLabel`.
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

	*MemorySizeUsedBar = *g.ProgressBar(float32(usedVal) / float32(dashboardData.Total))

	for i, dataRecord := range data {
		// SYSTEM

		plot1Data.Add(float64(dataRecord.Memory.System.Used)/float64(valueUnitDivider), "system-used")
		plot1Data.AddCustomLabel(i, dataRecord.Description)

		plot1Data.Add(float64(dataRecord.Memory.System.Total)/float64(valueUnitDivider), "system-total")
		plot1Data.AddCustomLabel(i, dataRecord.Description)

		plot2Data.Add(float64(dataRecord.Memory.App[0].Used)/float64(valueUnitDivider), "app-used")
		plot2Data.AddCustomLabel(i, dataRecord.Description)

		plot2Data.Add(float64(dataRecord.Memory.System.Total)/float64(valueUnitDivider), "system-total")
		plot2Data.AddCustomLabel(i, dataRecord.Description)

		// Modules

		sort.Slice(dataRecord.Memory.App[0].Modules, func(i, j int) bool {
			return dataRecord.Memory.App[0].Modules[i].File > dataRecord.Memory.App[0].Modules[j].File
		})

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
			plot3Data.Add(float64(moduleVal.Used)/float64(valueUnitDivider), element)
		}
		plot3Data.AddCustomLabel(i, dataRecord.Description)
	}
}
