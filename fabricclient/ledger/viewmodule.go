package ledger

import (
	"fmt"
	"github.com/lxn/walk"
	"sort"
)

type MspView struct {
	Name string
	SignAlgo string
	HashAlgo string
	PubAlgo string
	PubCurve string
	CheckResult string
	Selected bool
}

type ViewModule struct {
	walk.TableModelBase
	walk.SorterBase
	sortOrder walk.SortOrder
	sortColumn int
	items []*MspView
}

func (m *ViewModule) RowCount() int {
	return len(m.items)
}

func (m *ViewModule) AddNewItem(item MspView) {
	m.items = append(m.items, &item)
}

func (m *ViewModule) DelAll() {
	m.items = nil
}

func (m *ViewModule) GetItem(i int) *MspView{
	return m.items[i]
}


func (m *ViewModule) ResetRows() {
	fmt.Println("ResetRows")
}

func (m *ViewModule)  Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.Name
	case 1:
		return item.SignAlgo
	case 2:
		return item.HashAlgo
	case 3:
		return item.PubAlgo
	case 4:
		return item.PubCurve
	case 5:
		return item.CheckResult
	default:
		return fmt.Sprintf("unexpect col %d", col)
	}
}

func (m *ViewModule) Sort(col int, order walk.SortOrder) error {
	m.sortColumn = col
	m.sortOrder = order
	sort.SliceStable(m.items, func(i, j int) bool {
		return i > j
	})

	return m.SorterBase.Sort(col, order)
}

func (m *ViewModule) SetChecked(row int, checked bool) error {
	m.items[row].Selected = checked
	return nil
}

func (m *ViewModule) Checked(row int) bool {
	return m.items[row].Selected
}
