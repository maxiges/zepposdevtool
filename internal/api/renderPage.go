package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"text/template"
	"zepp-os-dev-tool/internal/models"
	"zepp-os-dev-tool/internal/storage"
	"zepp-os-dev-tool/static"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"go.uber.org/zap"
)

// HandlerDisplayData renders a comprehensive memory profiling dashboard for a specific application.
//
// This handler generates an HTML page with interactive charts displaying:
// 1. System memory usage over time (absolute values)
// 2. System memory usage as a percentage
// 3. Per-module memory usage breakdown for the application
//
// The handler retrieves stored memory profiling data for the specified application and
// visualizes it using line charts. If no data exists for the application, it renders
// a message indicating no data was found.
//
// Request Path Parameters:
// - appName: the identifier of the application to display data for
//
// Response: HTML page with embedded ECharts visualizations
//
// Example GET: /api/display/my-app
func HandlerDisplayData(w http.ResponseWriter, r *http.Request) {

	appName := r.PathValue("appName")
	data, exist := storage.GetDataForApp(appName)
	if !exist {
		line := charts.NewLine()
		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
			charts.WithTitleOpts(opts.Title{
				Title: "No data found for app : " + appName,
			}))
		line.Render(w)
		return
	}

	// Initialize dashboard template data with app name and total memory info
	dashboardData := models.TemplateData{
		AppName: appName,
		Total:   data[len(data)-1].Memory.System.Total,
		TotalKB: data[len(data)-1].Memory.System.Total / 1024,
	}

	// Initialize data structures to collect chart data series:
	// diagramsInTime: absolute memory values (bytes)
	// diagramsInTimePer: memory values as percentages
	// diagramsInTime2: per-module memory usage breakdown
	diagramsInTime := map[string][]opts.LineData{}
	diagramsInTimePer := map[string][]opts.LineData{}
	diagramsInTime2 := map[string][]opts.LineData{}

	xLabels := make([]int, 0, len(data))
	for i, dataRecord := range data {
		xLabels = append(xLabels, i)
		dashboardData.Total = dataRecord.Memory.System.Total

		// === SYSTEM MEMORY DATA ===
		// Track absolute system memory usage (Used)
		val := diagramsInTime["system-used"]
		if len(val) == 0 {
			val = []opts.LineData{}
		}
		val = append(val, opts.LineData{Value: dataRecord.Memory.System.Used, XAxisIndex: i, Name: dataRecord.Description})
		diagramsInTime["system-used"] = val

		// Track absolute system memory total capacity
		val2 := diagramsInTime["system-total"]
		if len(val2) == 0 {
			val2 = []opts.LineData{}
		}
		val2 = append(val2, opts.LineData{Value: dataRecord.Memory.System.Total, XAxisIndex: i, Name: dataRecord.Description})
		diagramsInTime["system-total"] = val2

		val3 := diagramsInTimePer["system-used-%"]
		if len(val3) == 0 {
			val3 = []opts.LineData{}
		}
		val3 = append(val3, opts.LineData{Value: dataRecord.Memory.System.Used * 100 / dataRecord.Memory.System.Total, XAxisIndex: i, Symbol: "roundRect", Name: dataRecord.Description})
		diagramsInTimePer["system-used-%"] = val3

		// === APPLICATION MODULE DATA ===
		// Track total memory used by the application
		valAppTotal := diagramsInTime2["Total_used_by_app"]
		if len(valAppTotal) == 0 {
			valAppTotal = []opts.LineData{}
		}
		valAppTotal = append(valAppTotal, opts.LineData{Value: dataRecord.Memory.App[0].Used, XAxisIndex: i, Name: dataRecord.Description})
		diagramsInTime2["Total_used_by_app"] = valAppTotal

		// Track per-module memory usage for each module in the application
		for _, moduleVal := range dataRecord.Memory.App[0].Modules {
			valAppModule := diagramsInTime2[moduleVal.File+"(Used)"]
			if len(valAppModule) == 0 {
				valAppModule = []opts.LineData{}
			}
			if len(valAppModule) < i {
				for _ = range i - len(valAppModule) {
					valAppModule = append(valAppModule, opts.LineData{Name: dataRecord.Description})
				}
			}
			valAppModule = append(valAppModule, opts.LineData{Value: moduleVal.Used, XAxisIndex: i, Name: dataRecord.Description})
			diagramsInTime2[moduleVal.File+"(Used)"] = valAppModule
		}

	}

	// tmpl, err := template.ParseFiles("./static/renderPageTop.htmx")
	tmpl, err := template.New("renderPage").Parse(static.RenderPageTop)
	if err != nil {
		logger, _ := zap.NewDevelopment()
		logger.Sugar().With("error", err.Error()).Error("Parse file Error")

	}
	var b bytes.Buffer
	foo := io.Writer(&b)
	tmpl.Execute(foo, &dashboardData)
	dataBytes, _ := io.ReadAll(&b)
	w.Write(dataBytes)

	w.Write([]byte("<div class=\"charts\">"))
	w.Write([]byte("<div class=\"row\">"))
	w.Write([]byte("<div class=\"col col-xl-6\">"))
	createDiagram(w, diagramsInTime, xLabels, "System", "SYSTEM")
	w.Write([]byte("</div>"))
	w.Write([]byte("<div class=\"col col-xl-6\">"))
	createDiagram(w, diagramsInTimePer, xLabels, "System by %", "SYSTEM_PER")
	w.Write([]byte("</div>"))
	w.Write([]byte("<div class=\"col col-xl-6\">"))
	createDiagram(w, diagramsInTime2, xLabels, "Modules", "MODULES")
	w.Write([]byte("</div>"))
	w.Write([]byte("</div></div>"))

}

// createDiagram generates and renders an ECharts line chart with the provided data.
//
// This helper function:
// 1. Writes a formatted title for the chart
// 2. Creates a line chart with optimized rendering settings (animations disabled)
// 3. Adds all data series to the chart
// 4. Sorts series alphabetically by name for consistent display
// 5. Configures smooth line interpolation
// 6. Renders the complete chart to the HTTP response
//
// Parameters:
// - w: HTTP response writer to write chart HTML to
// - data: map of series names to LineData arrays containing the chart data points
// - xLabels: array of x-axis labels (typically data point indices or timestamps)
// - chartText: display title for the chart
// - chartID: unique identifier for the ECharts instance
//
// The chart is configured with:
// - Animation disabled (AnimationDuration: 1, Animation: false)
// - Smooth line curves for better visualization
// - Macarons theme for consistent styling
func createDiagram(w http.ResponseWriter, data map[string][]opts.LineData, xLabels []int, chartText string, chartID string) {
	w.Write([]byte(fmt.Sprintf("<p class=\"text-center fs-4 fw-bold text-primary\">%s</p>", chartText)))
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeMacarons, ChartID: chartID}),
	)

	for key, val := range data {
		line.AddSeries(key, val, charts.WithSeriesOpts(func(s *charts.SingleSeries) {
			// default is 1000
			s.AnimationDuration = 1
			s.AnimationDurationUpdate = 0
			falseVal := false
			s.Animation = &falseVal
		}),
		)
	}

	sort.Slice(line.MultiSeries, func(i, j int) bool {
		return line.MultiSeries[i].Name < line.MultiSeries[j].Name
	})

	line.SetXAxis(xLabels)
	// Put data into instance
	line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))

	line.Render(w)

}
