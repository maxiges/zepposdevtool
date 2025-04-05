package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"zepp-os-dev-tool/internal/models"
	"zepp-os-dev-tool/internal/storage"
	"zepp-os-dev-tool/static"

	"go.uber.org/zap"
)

func HandlerDisplayMainPage(w http.ResponseWriter, r *http.Request) {

	appNames := storage.GetAppList()

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

	//END
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

	w.Write([]byte(fmt.Sprintf(endVal, static.RenderReadme)))

}
