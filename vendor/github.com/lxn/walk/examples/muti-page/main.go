package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	DirStructError = "错误，此组织目录结构不符合标准规范。"
	NanoToSecond = 1000000000
)

type MyWindow struct {
	mainWindow *walk.MainWindow
	mspTable *walk.TableView
	selectPathBtn *walk.PushButton
	mspRadio *walk.RadioButton
	ledgerRadio *walk.RadioButton
	pathEdit *walk.TextEdit
	ledgerResultEdit *walk.Label
	startCheckBtn *walk.PushButton
	pageCom *walk.Composite
	status *walk.Label
}


func newPageAction2(title, image string, triggerStr string) (*walk.Action) {
	img, err := walk.Resources.Bitmap(image)
	if err != nil {
		fmt.Println("err")
		return nil
	}

	action := walk.NewAction()
	action.SetCheckable(true)
	action.SetExclusive(true)
	action.SetImage(img)
	action.SetText(title)

	action.Triggered().Attach(func() {
		fmt.Println(triggerStr)
	})
	return action
}


func main(){
	var ttt *walk.ToolBar
	var aaa *walk.Action
	//br := SolidColorBrush{walk.RGB(0,0,0)}
	mw := MyWindow{}

	mainWindow := MainWindow{
		AssignTo: &mw.mainWindow,
		Title: "sm tableview",
		Size: Size{800, 600},
		Background: SolidColorBrush{Color: walk.RGB(243,243,243)},
		Layout:   HBox{},
		ToolBar: ToolBar{
			Font:Font{PointSize:9},
			ButtonStyle: ToolBarButtonImageBeforeText,
			Items: []MenuItem{
				Action{
					AssignTo: &aaa,
					Text:  "MSP",
					Image: "./img/document-new.png",
					OnTriggered: func () {
						aaa.SetChecked(true)
						bg, err := walk.NewSolidColorBrush(walk.RGB(255, 255, 0))
						if err != nil {
							return
						}
						ttt.SetBackground(bg)
						mw.ledgerResultEdit.SetText("aaaa")
					},
				},
				Separator{},
				Action{
					Text:  "账本",
					Image: "./img/document-properties.png",
					OnTriggered: func () {
						mw.ledgerResultEdit.SetText("bbb")
					},
				},
				Separator{},
			},
		},
		Children: []Widget{
			ScrollView{
				Background: SolidColorBrush{walk.RGB(250, 250, 250)},
				HorizontalFixed: true,
				Layout:  VBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					Composite{
						Layout: VBox{MarginsZero: true},
						Children: []Widget{
							ToolBar{
								AssignTo:    &ttt,
								Orientation: Vertical,
								ButtonStyle: ToolBarButtonImageAboveText,
								MaxTextRows: 2,
							},
						},
					},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					TableView{
						Visible: true,
						AssignTo: &mw.mspTable,
						CheckBoxes: false,
						ColumnsOrderable: false,
						MultiSelection: false,
						NotSortableByHeaderClick: true,
						AlternatingRowBGColor: walk.RGB(239, 239, 239),
						Columns: []TableViewColumn{
							{Title: "组织"},
							{Title: "成员"},
							{Title: "一级目录名"},
							{Title: "证书目录名"},
							{Title: "证书名"},
							{Title: "检测结果", Width:200},
						},
						OnSelectedIndexesChanged: func() {
							fmt.Println("SelectedIndexes")
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget {
							Label{
								Text:          "hello world",
								Visible:       true,
								AssignTo:      &mw.ledgerResultEdit,
								TextAlignment: AlignNear,
								//Enabled: false,
								TextColor: walk.RGB(0, 0, 0),
							},
							HSpacer{},
						},
					},
				},
			},
			//VSpacer{},
		},
	}

	//code, err := mainWindow.Run()
	//if err != nil {
	//	fmt.Println("encounter error code:%d, info:%s", code, err)
	//}
	err := mainWindow.Create()
	if err != nil {
		fmt.Println("encounter error, info:%s", err)
	}

	actions := ttt.Actions()
	actions.Add(newPageAction2("aaa", "./img/document-new.png", "eeeee"))
	actions.Add(newPageAction2("fff", "./img/document-new.png", "fffff"))
	//mw.pageCom.

	code := mw.mainWindow.Run()
	if code != 0 {
		fmt.Println("encounter error, code:%d", code)
	}
}
