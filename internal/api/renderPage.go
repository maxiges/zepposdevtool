package api

import (
	"bytes"
	"fmt"
	"io"
	"my-chart-app/internal/models"
	"my-chart-app/internal/storage"
	"my-chart-app/static"
	"net/http"
	"text/template"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"go.uber.org/zap"
)

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

	dashboardData := models.TemplateData{
		AppName: appName,
		Total:   data[len(data)-1].Memory.System.Total,
		TotalKB: data[len(data)-1].Memory.System.Total / 1024,
	}

	diagramsInTime := map[string][]opts.LineData{}
	diagramsInTimePer := map[string][]opts.LineData{}
	diagramsInTime2 := map[string][]opts.LineData{}

	xLabels := make([]int, 0, len(data))
	for i, dataRecord := range data {
		xLabels = append(xLabels, i)
		dashboardData.Total = dataRecord.Memory.System.Total

		// SYSTEM
		val := diagramsInTime["system-used"]
		if len(val) == 0 {
			val = []opts.LineData{}
		}
		val = append(val, opts.LineData{Value: dataRecord.Memory.System.Used, XAxisIndex: i, Name: dataRecord.Description})
		diagramsInTime["system-used"] = val

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

		// APP
		valAppTotal := diagramsInTime2["Total_used_by_app"]
		if len(valAppTotal) == 0 {
			valAppTotal = []opts.LineData{}
		}
		valAppTotal = append(valAppTotal, opts.LineData{Value: dataRecord.Memory.App[0].Used, XAxisIndex: i, Name: dataRecord.Description})
		diagramsInTime2["Total_used_by_app"] = valAppTotal

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

	line.SetXAxis(xLabels)
	// Put data into instance
	line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))

	line.Render(w)

}
