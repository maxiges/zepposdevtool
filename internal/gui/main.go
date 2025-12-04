package gui

import (
	"zepp-os-dev-tool/internal/api"
	"zepp-os-dev-tool/internal/storage"

	g "github.com/AllenDang/giu"
)

type MenuIndex int

const (
	MenuNull  MenuIndex = 0
	MenuHello MenuIndex = 1
	MenuPlot  MenuIndex = 2
)

var (
	currentMenu     MenuIndex
	currentMenuBar  *g.MenuBarWidget
	AutoRefreshIsOn = true

	MainUILayout g.Widget
	Plot1        *g.PlotCanvasWidget
	Plot2        *g.PlotCanvasWidget
	Plot3        *g.PlotCanvasWidget

	MemoryLabel       *g.LabelWidget
	MemorySizeUsedBar *g.ProgressBarWidget

	SelectMem1 *g.CheckboxWidget
	SelectMem2 *g.CheckboxWidget
	SelectMem3 *g.CheckboxWidget

	AutoMove *g.SliderIntWidget
)

var (
	windowW int = 1800
	windowH int = 900

	sashPos1 float32 = float32(windowW / 2)
	sashPos2 float32 = 400
	sashPos3 float32 = 50
	sashPos4 float32 = 100
)

var (
	markdown               = defaultMd
	splitLayoutPos float32 = 900
)

func firstLayout() {
	currentMenu = MenuHello
	MainUILayout = g.Layout{
		g.Markdown(markdown).
			Header(0, (g.Context.FontAtlas.GetDefaultFonts())[0].SetSize(28), true).
			Header(1, (g.Context.FontAtlas.GetDefaultFonts())[0].SetSize(26), false).
			Header(2, nil, true),
	}

}

func loop() {

	if currentMenu == MenuNull {
		firstLayout()
		GenerateMenu()
	}

	g.SingleWindowWithMenuBar().Layout(
		currentMenuBar,
		MainUILayout,
	)
}

func GenerateMenu() {
	autoRefreshText := "Auto refresh"
	if AutoRefreshIsOn {
		autoRefreshText += " (On)"
	} else {
		autoRefreshText += " (Off)"
	}

	appList := storage.GetAppList()
	menuButtons := []g.Widget{}

	for _, appName := range appList {
		menuElement := g.MenuItem(appName).OnClick(func() {
			DrawChartsForApp(appName)
		})
		menuButtons = append(menuButtons, menuElement)

	}

	currentMenuBar = g.MenuBar().Layout(
		g.Menu("Home").Layout(
			g.MenuItem("Home").OnClick(func() {
				firstLayout()
			}),
		),
		g.Menu("Apps").Layout(
			menuButtons...,
		),
		g.MenuItem(autoRefreshText).OnClick(func() {
			AutoRefreshIsOn = !AutoRefreshIsOn
			go GenerateMenu()
		}),
	)

}
func RunGui(wFlag, hFlag *int) {

	if wFlag != nil && *wFlag > 0 {
		windowW = *wFlag
		sashPos1 = float32(windowW / 2)
	}
	if hFlag != nil && *hFlag > 0 {
		windowH = *hFlag
	}
	api.RefreshFun = TryRefreshUI

	wnd := g.NewMasterWindow("Plot", windowW, windowH, g.MasterWindowFlagsNotResizable)
	wnd.Run(loop)
}
func DisableGui() {
	AutoRefreshIsOn = false
}
