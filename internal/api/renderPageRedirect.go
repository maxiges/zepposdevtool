package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/maxiges/ZeppOsDevTool/internal/models"
	"github.com/maxiges/ZeppOsDevTool/internal/storage"
	"github.com/maxiges/ZeppOsDevTool/static"

	"go.uber.org/zap"
)

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
func HandlerDisplayMainPage(w http.ResponseWriter, r *http.Request) {

	appNames := storage.GetAppList()

	// === RENDER MAIN PAGE HEADER ===
	// Parse and render the main page template (header, title, intro content)
	tmpl, err := template.New("renderPage").Parse(static.RenderMainPage)
	if err != nil {
		logger, _ := zap.NewDevelopment()
		logger.Sugar().With("error", err.Error()).Error("Parse file Error")

	}
	var b bytes.Buffer
	foo := io.Writer(&b)
	tmpl.Execute(foo, nil)
	dataBytes, _ := io.ReadAll(&b)
	w.Write(dataBytes)

	// === RENDER APPLICATION BUTTONS ===
	// Generate a clickable button for each registered application
	for _, appName := range appNames {

		tmpl, err := template.New("LinkButton").Parse(static.RenderAppButtonPage)
		if err != nil {
			logger, _ := zap.NewDevelopment()
			logger.Sugar().With("error", err.Error()).Error("Parse file Error")

		}
		var b bytes.Buffer
		foo := io.Writer(&b)
		tmpl.Execute(foo, &models.BasicAppValues{
			AppName: appName,
		})
		dataBytes, _ := io.ReadAll(&b)
		w.Write(dataBytes)
	}

	// === RENDER README SECTION WITH CLIENT-SIDE MARKDOWN CONVERSION ===
	// Define HTML structure for readme content and markdown conversion script
	// The Showdown.js library converts markdown to HTML on the client-side in real-time
	endVal := `
	 <div class="mt-5" id="content"></div>
 	 <script src="https://unpkg.com/showdown/dist/showdown.min.js"></script>
	 <script>
	 var converter = new showdown.Converter();
     document.getElementById('content').innerHTML =
     converter.makeHtml(` + "`%s" + "`" + `);
  	</script>
</'body></html>
`

	// Inject the readme markdown content and render the final section
	// The %s placeholder is replaced with the static README content
	w.Write([]byte(fmt.Sprintf(endVal, static.RenderReadme)))

}
