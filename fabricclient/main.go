package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/HanDaXia/BlockChainSafeTesting/fabricclient/ledger"
)

type MyWindow struct {
	ledgerUI ledger.LedgerUI
	mainWindow *walk.MainWindow
	mspCom *walk.Composite
}

func main(){
	mw := MyWindow{}
	mainWindow := MainWindow{
		AssignTo: &mw.mainWindow,
		Title: "安全检测工具",
		Size: Size{1000, 750},
		Layout:   HBox{MarginsZero: true},
		Children: []Widget{
			ScrollView{
				Visible: true,
				Layout: VBox{MarginsZero: true},
				Children: mw.ledgerUI.NewLedgerConfig(&mw.mainWindow),
			},

			VSpacer{},
		},
	}

	err := mainWindow.Create()
	if err != nil {
		fmt.Println("encounter error while mainWindow create, info:%s", err)
	}

	mw.mainWindow.Run()
}
