package ledger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	serviceUrl string
)

func init() {
	serviceUrl = os.Getenv("FABRIC_SERVER_URL")
	if len(serviceUrl) == 0 {
		serviceUrl = "http://172.16.0.250:8080/Check"
	}
}

type LedgerUI struct {
	mainWindow **walk.MainWindow
	com *walk.Composite
	randomCom *walk.Composite
	selectPathBtn *walk.PushButton
	selectRandom *walk.PushButton
	pathEdit *walk.TextEdit
	randonEdit *walk.TextEdit
	startCheckBtn *walk.PushButton
	status *walk.Label
	result *walk.Label
	blockTb *walk.TableView
	blockMd ViewModule
	txTb *walk.TableView
	txMd ViewModule
	endorseTb *walk.TableView
	endorseMd ViewModule
	ledgerData *LedgerData
	ledgerLabel[6] *walk.Label
	blockLabel *walk.Label
	txLabel *walk.Label
	endorseLabel *walk.Label
	randomLabel *walk.Label
	randomCheck[19] *walk.CheckBox
	randomResult *walk.Label
}

func (m *LedgerUI)NewLedgerConfig(mainWindow **walk.MainWindow) []Widget{
	m.mainWindow = mainWindow
	return []Widget{
		Composite{
			MaxSize: Size{0, 40},
			Layout:  HBox{},
			Children: []Widget{
				PushButton{
					Text:      "选择账本目录",
					AssignTo:  &m.selectPathBtn,
					OnClicked: m.selectPathClicked},
				TextEdit{AssignTo: &m.pathEdit},
				//HSpacer{},
			},
		},

		Composite{
			MaxSize: Size{0, 40},
			Layout:  HBox{},
			Children: []Widget{
				PushButton{
					Text:      "选择随机数文件",
					AssignTo:  &m.selectRandom,
					OnClicked: m.selectFileClicked},
				TextEdit{AssignTo: &m.randonEdit},
				//HSpacer{},
			},
		},

		Composite{
			//MaxSize:Size{0, 40},
			Layout: HBox{},
			Children: []Widget{
				PushButton{Text: "启动检测", AssignTo: &m.startCheckBtn, OnClicked: m.startBtnClicked},
				Label{Text: "", AssignTo: &m.status},
				HSpacer{},
			},
		},

		Composite{
			AssignTo: &m.com,
			Visible: false,
			Layout: VBox{},
			Children: []Widget{
				HSplitter{
					Children:[]Widget{
						Label{AssignTo: &m.ledgerLabel[0], TextAlignment:AlignNear, MaxSize:Size{100,0}},
						Label{AssignTo: &m.ledgerLabel[1], TextAlignment:AlignNear, MaxSize:Size{100,0}},
						Label{AssignTo: &m.ledgerLabel[2], TextAlignment:AlignNear, MaxSize:Size{100,0}},
						Label{AssignTo: &m.ledgerLabel[3], TextAlignment:AlignNear, MaxSize:Size{100,0}},
						Label{AssignTo: &m.ledgerLabel[4], TextAlignment:AlignNear, MaxSize:Size{100,0}},
						Label{AssignTo: &m.ledgerLabel[5], TextAlignment:AlignNear, MaxSize:Size{100,0}},
					},
				},

				Composite{
					Layout: HBox{},
					Children: []Widget {
						Label{Text: "区块检测数据", AssignTo: &m.blockLabel, TextAlignment:AlignNear},
						HSpacer{},
					},
				},
				TableView{
					AssignTo:                 &m.blockTb,
					CheckBoxes:               false,
					ColumnsOrderable:         false,
					MultiSelection:           true,
					NotSortableByHeaderClick: true,
					MinSize:Size{100,70},
					AlternatingRowBGColor:    walk.RGB(240, 240, 240),
					Model:                    &m.blockMd,
					Columns: []TableViewColumn{
						{Title: "区块序号"},
						{Title: "签名算法"},
						{Title: "Hash算法"},
						{Title: "公钥算法"},
						{Title: "公钥曲线"},
						{Title: "验证结果"},
					},
					OnSelectedIndexesChanged: m.selectIndexChangeBlock,
				},

				Composite{
					Layout: HBox{},
					Children: []Widget {
						Label{AssignTo: &m.txLabel, TextAlignment:AlignNear},
						HSpacer{},
					},
				},
				TableView{
					AssignTo:                 &m.txTb,
					CheckBoxes:               false,
					ColumnsOrderable:         false,
					MultiSelection:           true,
					NotSortableByHeaderClick: true,
					MinSize:Size{100,70},
					AlternatingRowBGColor:    walk.RGB(240, 240, 240),
					Model:                    &m.txMd,
					Columns: []TableViewColumn{
						{Title: "交易ID"},
						{Title: "签名算法"},
						{Title: "Hash算法"},
						{Title: "公钥算法"},
						{Title: "公钥曲线"},
						{Title: "验证结果"},
					},
					OnSelectedIndexesChanged: m.selectIndexChangeTx,
				},

				Composite{
					Layout: HBox{},
					Children: []Widget {
						Label{AssignTo: &m.endorseLabel, TextAlignment:AlignNear},
						HSpacer{},
					},
				},
				TableView{
					AssignTo:                 &m.endorseTb,
					CheckBoxes:               false,
					ColumnsOrderable:         false,
					MultiSelection:           false,
					NotSortableByHeaderClick: true,
					MinSize:Size{100,60},
					AlternatingRowBGColor:    walk.RGB(240, 240, 240),
					Model:                    &m.endorseMd,
					Columns: []TableViewColumn{
						{Title: "背书节点"},
						{Title: "签名算法"},
						{Title: "Hash算法"},
						{Title: "公钥算法"},
						{Title: "公钥曲线"},
						{Title: "验证结果"},
					},
				},

				Composite{
					Layout: HBox{},
					Children: []Widget {
						Label{
							Text:          "",
							AssignTo:      &m.result,
							TextAlignment: AlignNear,
							//Enabled: false,
							TextColor: walk.RGB(0, 0, 0),
						},
						HSpacer{},
					},
				},
			},
		},

		Composite{
			AssignTo: &m.randomCom,
			Layout:   VBox{},
			Visible: false,
			Children: []Widget {
				Composite{
					Layout: HBox{},
					Children: []Widget {
						Label{
							MinSize:Size{500, 0},
							Text:"",
							AssignTo:      &m.randomResult,
							TextAlignment: AlignNear,
						},
						HSpacer{},
					},
				},
				HSplitter{
					Children: []Widget{
						VSplitter{
							MinSize:Size{120,0},
							Children: []Widget{
								CheckBox{AssignTo: &m.randomCheck[0], OnClicked: m.onStateChange(0)},
								CheckBox{AssignTo: &m.randomCheck[1], OnClicked: m.onStateChange(1)},
								CheckBox{AssignTo: &m.randomCheck[2], OnClicked: m.onStateChange(2)},
							},
						},
						VSplitter{
							MinSize:Size{120,0},
							Children: []Widget{
								CheckBox{AssignTo: &m.randomCheck[3], OnClicked: m.onStateChange(3)},
								CheckBox{AssignTo: &m.randomCheck[4], OnClicked: m.onStateChange(4)},
								CheckBox{AssignTo: &m.randomCheck[5], OnClicked: m.onStateChange(5)},
							},
						},
						VSplitter{
							MinSize:Size{140,0},
							Children: []Widget{
								CheckBox{AssignTo: &m.randomCheck[6], OnClicked: m.onStateChange(6)},
								CheckBox{AssignTo: &m.randomCheck[7], OnClicked: m.onStateChange(7)},
								CheckBox{AssignTo: &m.randomCheck[8], OnClicked: m.onStateChange(8)},
							},
						},
						VSplitter{
							MinSize:Size{140,0},
							Children: []Widget{
								CheckBox{AssignTo: &m.randomCheck[9], OnClicked: m.onStateChange(9)},
								CheckBox{AssignTo: &m.randomCheck[10], OnClicked: m.onStateChange(10)},
								CheckBox{AssignTo: &m.randomCheck[11], OnClicked: m.onStateChange(11)},
							},
						},
						VSplitter{
							MinSize:Size{140,0},
							Children: []Widget{
								CheckBox{AssignTo: &m.randomCheck[12], OnClicked: m.onStateChange(12)},
								CheckBox{AssignTo: &m.randomCheck[13], OnClicked: m.onStateChange(13)},
								CheckBox{AssignTo: &m.randomCheck[14], OnClicked: m.onStateChange(14)},
							},
						},
						VSplitter{
							MinSize:Size{140,0},
							Children: []Widget{
								CheckBox{AssignTo: &m.randomCheck[15], OnClicked: m.onStateChange(15)},
								CheckBox{AssignTo: &m.randomCheck[16], OnClicked: m.onStateChange(16)},
							},
						},
						VSplitter{
							MinSize:Size{120,0},
							Children: []Widget{
								CheckBox{AssignTo: &m.randomCheck[17], OnClicked: m.onStateChange(17)},
								CheckBox{AssignTo: &m.randomCheck[18], OnClicked: m.onStateChange(18)},
							},
						},
					},
				},
			},
		},
		VSpacer{},
	}
}

func (m *LedgerUI)onStateChange(index int) (func()) {
	return func() {
		if !m.randomCom.Visible() {
			return
		}

		m.randomCheck[index].SetChecked(!m.randomCheck[index].Checked())
	}
}

func (m *LedgerUI)selectIndexChangeBlock() {
	m.endorseMd.DelAll()
	m.endorseMd.PublishRowsReset()
	if len(m.blockTb.SelectedIndexes()) == 0 {
		return
	}

	curIndex := m.blockTb.SelectedIndexes()[0]
	m.txLabel.SetText(fmt.Sprintf("区块[%d]交易检测数据", curIndex+1))
	m.endorseLabel.SetText("背书结果验证")

	m.txMd.DelAll()
	for _, tx := range m.ledgerData.Blocks[curIndex].Txs {
		result := "通过"
		if !tx.VerifyOk {
			result = "未通过"
		}
		item := MspView{tx.TxID, tx.SignAlgo, tx.HashAlgo, tx.PubAlgo,
			tx.PubCurve, result, false}
		m.txMd.AddNewItem(item)
	}
	m.txMd.PublishRowsReset()
}

func (m *LedgerUI)selectIndexChangeTx() {
	if len(m.blockTb.SelectedIndexes()) == 0 || len(m.txTb.SelectedIndexes()) == 0 {
		return
	}
	blockIndex := m.blockTb.SelectedIndexes()[0]
	txIndex := m.txTb.SelectedIndexes()[0]
	m.endorseMd.DelAll()
	for _, endorse := range m.ledgerData.Blocks[blockIndex].Txs[txIndex].Endorsers {
		result := "通过"
		if !endorse.VerifyOk {
			result = "未通过"
		}
		item := MspView{endorse.Name, endorse.SignAlgo, endorse.HashAlgo,
			endorse.PubAlgo, endorse.PubCurve, result, false}
		m.endorseMd.AddNewItem(item)
	}
	m.endorseMd.PublishRowsReset()
	m.endorseLabel.SetText(fmt.Sprintf("交易[%s]背书结果验证", m.ledgerData.Blocks[blockIndex].Txs[txIndex].TxID))
}

func (m *LedgerUI)selectPathClicked() {
	dlg := new(walk.FileDialog)
	dlg.Title = "选择账本目录"

	if _, err := dlg.ShowBrowseFolder(*m.mainWindow); err != nil {
		m.pathEdit.SetText("Error : File Open\r\n")
		return
	}
	m.pathEdit.SetText(dlg.FilePath)
}

func (m *LedgerUI)selectFileClicked() {
	dlg := new(walk.FileDialog)
	dlg.Title = "选择随机数文件"

	if _, err := dlg.ShowOpen(*m.mainWindow); err != nil {
		m.randonEdit.SetText("Error : File Open\r\n")
		return
	}
	m.randonEdit.SetText(dlg.FilePath)
}

func (m *LedgerUI) startBtnClicked() {
	m.com.SetVisible(false)
	m.randomCom.SetVisible(false)

	randomPath := m.randonEdit.Text()
	ledgerPath := m.pathEdit.Text()
	if len(randomPath) == 0 && len(ledgerPath) == 0 {
		m.status.SetText("至少设置一个路径！")
		return
	}

	if len(randomPath) > 0 {
		randomPathInfo, err := os.Stat(randomPath)
		if err != nil {
			m.status.SetText("随机数路径不存在，请重新设置！")
			return
		}
		if randomPathInfo.IsDir() {
			m.status.SetText("随机数路径错误，请重新设置文件！")
			return
		}
	}
	if len(ledgerPath) >0 {
		ledgerPathInfo, err := os.Stat(ledgerPath)
		if err != nil {
			m.status.SetText("账本路径不存在，请重新设置！")
			return
		}
		if !ledgerPathInfo.IsDir() {
			m.status.SetText("账本路径错误，请重新设置文件目录！")
			return
		}
	}

	m.startCheckBtn.SetEnabled(false)
	m.status.SetText("正在检测账本，请稍等。。。")
	go m.startCheckLedger(ledgerPath, randomPath)
}

func (m *LedgerUI) parseLedgrResult() {
	if len(m.ledgerData.Err) != 0 {
		m.status.SetTextColor(walk.RGB(250,0, 0))
		_ = m.status.SetText("解析账本出错" + m.ledgerData.Err)
		return
	}

	m.txMd.DelAll()
	m.txMd.PublishRowsReset()
	m.endorseMd.DelAll()
	m.endorseMd.PublishRowsReset()

	m.blockMd.DelAll()
	for _, block := range m.ledgerData.Blocks {
		result := "通过"
		if !block.VerifyOk {
			result = "未通过"
		}
		item := MspView{strconv.Itoa(int(block.Height)), block.SignAlgo, block.HashAlgo,
			block.PubAlgo, block.PubCurve, result, false}
		m.blockMd.AddNewItem(item)
	}
	m.blockMd.PublishRowsReset()
	m.status.SetTextColor(walk.RGB(0,250,0))
	_ = m.status.SetText("检测完成，请查看结果")
	_ = m.txLabel.SetText("交易检测数据")
	_ = m.endorseLabel.SetText("背书结果验证")
	_ = m.ledgerLabel[0].SetText(fmt.Sprintf("共识算法\n%s", m.ledgerData.ConsusType))
	_ = m.ledgerLabel[1].SetText(fmt.Sprintf("交易最大尺寸\n%s ",getSizeString(float64(m.ledgerData.TxMaxSize))))
	_ = m.ledgerLabel[2].SetText(fmt.Sprintf("打包大小\n%d",	m.ledgerData.TxCount))
	_ = m.ledgerLabel[3].SetText(fmt.Sprintf("区块最大尺寸\n%s", getSizeString(float64(m.ledgerData.BlockSize))))
	_ = m.ledgerLabel[4].SetText(fmt.Sprintf("超时时间\n%s", m.ledgerData.BlockTime))
	_ = m.ledgerLabel[5].SetText(fmt.Sprintf("orderer地址\n%s ",m.ledgerData.OrdererAddr))

	unpassItems := strings.Builder{}
	for _, block := range m.ledgerData.Blocks {
		if !block.VerifyOk {
			unpassItems.WriteString(fmt.Sprintf("\r\n区块[%d]", block.Height))
		}
		for _, tx := range block.Txs {
			if !tx.VerifyOk {
				unpassItems.WriteString(fmt.Sprintf("\r\n区块[%d]-交易[%s]", block.Height, tx.TxID))
			}
			for _, endorse := range tx.Endorsers {
				if !endorse.VerifyOk {
					unpassItems.WriteString(fmt.Sprintf("\r\n区块[%d]-交易[%s]-endorse[%s]", block.Height, tx.TxID, endorse.Name))
				}
			}
		}
	}
	if unpassItems.Len() == 0 {
		m.result.SetTextColor(walk.RGB(0, 250, 0))
		_ = m.result.SetText("账本检测结果： 通过")
	} else {
		m.result.SetTextColor(walk.RGB(250, 0, 0))
		_ = m.result.SetText("账本检测结果： 检测未通过\r\n未通过项："+unpassItems.String())
	}

	m.com.SetVisible(true)
}


func (m *LedgerUI) parseRandomResult(randomResult *RandomResponse) {
	checkItems := []string{
		"频率检验", "块内频数检验","累加和检验","游程检验","块内最长游程检验","二元矩阵秩检验", "离散傅里叶变换检验",
		"非重叠模块匹配检验","重叠模块匹配检验","Maurer的通用统计检验","近似熵检验","随机游动检验", "随机游动状态频数检验",
		"序列检验", "线性复杂度检验", "游程分布检验", "Poker检验", "自相关检验","二元推导检验",}

	if len(m.randomCheck[0].Text()) == 0 {
		for i:=0; i<19; i++ {
			_ = m.randomCheck[i].SetText(checkItems[i])
		}
	}

	m.randomCheck[0].Enabled()
	fmt.Println(randomResult.Result)
	m.randomResult.SetTextColor(walk.RGB(250,0, 0))
	splits := strings.Split(randomResult.Result, "&")
	if len(splits) != 19 {
		m.randomResult.SetText("随机数检测结果： 未通过\n未通过原因： " + randomResult.Result)
	} else {
		unpassItems := strings.Builder{}
		for index, result := range splits {
			rt := strings.Split(result, "=")
			if len(rt) != 2 {
				m.randomCheck[index].SetChecked(false)
				if unpassItems.Len() == 0 {
					unpassItems.WriteString(result)
				} else {
					unpassItems.WriteString(", " + result)
				}
			} else {
				if rt[1] == "0" {
					m.randomCheck[index].SetChecked(false)
					if unpassItems.Len() == 0 {
						unpassItems.WriteString(checkItems[index])
					} else {
						unpassItems.WriteString(", " + checkItems[index])
					}
				} else {
					m.randomCheck[index].SetChecked(true)
				}
			}

		}
		if unpassItems.Len() > 0 {
			m.randomResult.SetTextColor(walk.RGB(0,250, 0))
			_ = m.randomResult.SetText("随机数检测结果： 未通过 ")
		} else {
			m.randomResult.SetTextColor(walk.RGB(0,250, 0))
			_ = m.randomResult.SetText("随机数检测结果： 通过")
		}
	}

	m.randomCom.SetVisible(true)
}

func (m *LedgerUI) startCheckLedger(rootPath string, randomPath string) {
	defer m.startCheckBtn.SetEnabled(true)
	var err error
	var ledgerBytes []byte
	var randomBytes []byte
	if len(rootPath) != 0 {
		ledgerBytes, err = getLedgerBytes(rootPath)
		if err != nil {
			_ = m.status.SetText("读取账本数据错误！")
			return
		}
	}
	if len(randomPath) != 0 {
		randomBytes, err = getRandomBytes(randomPath)
		if err != nil {
			_ = m.status.SetText("读取随机数文件错误！")
			return
		}
	}

	var randomData *RandomResponse
	m.ledgerData, randomData, err = RequestAnalyze(ledgerBytes, randomBytes)
	if err != nil {
		_ = m.status.SetText(err.Error())
		return
	}

	if m.ledgerData != nil {
		m.parseLedgrResult()
	}

	if randomData != nil {
		m.parseRandomResult(randomData)
	}
}

func getLedgerBytes(ledgerPath string) ([]byte, error) {
	data, err := ioutil.ReadFile(ledgerPath + string(os.PathSeparator) + "blockfile_000000")
	return data, err
}

func getRandomBytes(randomPath string) ([]byte, error) {
	data, err := ioutil.ReadFile(randomPath)
	if err != nil {
		return nil, err
	}
	if len(data) < 1024*1024 {
		return nil, errors.New("随机数文件大小必须大于1M")
	}
	return data, nil
}

func getSizeString(len float64) string {
	unit := []string{"Byte", "KB", "MB", "GB", "TB"}
	unitIndex := 0
	for {
		if unitIndex == 4 {
			break
		}
		if len > 1024 {
			len = len / 1024
			unitIndex ++
		} else {
			break
		}
	}
	return fmt.Sprintf("%.2f%s", len, unit[unitIndex])
}

func PostBytes(url string, data []byte) ([]byte, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, body)
	if err != nil{
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Connection", "Keep-Alive")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("bad response status %d", resp.StatusCode))
	}

	return ioutil.ReadAll(resp.Body)
}

func RequestAnalyze(ledgerBytes, randomBytes []byte) (*LedgerData, *RandomResponse, error) {
	bytesData := SourceData{LedgerData: ledgerBytes, RandomData: randomBytes}
	cr := CheckRequest{CompanyName:"cetc", CompanyID:"cetc", CheckType:0}
	cr.Data, _ = json.Marshal(bytesData)
	postData, err := json.Marshal(&cr)
	if err != nil {
		return nil, nil, err
	}
	response, err := PostBytes(serviceUrl, postData)
	if err != nil {
		return nil, nil, errors.Wrap(err, "发送请求到服务器出错")
	}
	disRes := DistributeResult{}
	err = json.Unmarshal(response, &disRes)
	if err != nil {
		return nil, nil, errors.New("服务器返回数据格式错误: " + string(response))
	}

	checkResp := CheckResponse{}
	err = json.Unmarshal(disRes.Result, &checkResp)
	if err != nil {
		fmt.Println(string(response))
		return nil, nil, errors.New("服务器返回数据格式错误: " + string(disRes.Result))
	}

	var ld *LedgerData
	var rr *RandomResponse
	fbResp := &FabricResp{}
	err = json.Unmarshal(checkResp.Result.OtherResult, fbResp)
	if err != nil {
		return nil, nil, errors.New("服务器返回数据格式错误: " + string(checkResp.Result.OtherResult))
	}
	if len(fbResp.LedgerDetail) > 0 {
		ld = &LedgerData{}
		err = json.Unmarshal(fbResp.LedgerDetail, ld)
		if err != nil {
			return nil, nil, err
		}
	}
	if len(fbResp.RandomDetail) > 0 {
		rr = &RandomResponse{}
		err = json.Unmarshal(fbResp.RandomDetail, &rr)
		if err != nil {
			return nil, nil, errors.New("解析随机数数据错误："+string(fbResp.RandomDetail))
		}
	}
	return ld, rr, nil
}

