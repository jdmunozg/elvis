package linq

import (
	"reflect"
	"strings"

	"github.com/cgalvisleon/elvis/console"
	e "github.com/cgalvisleon/elvis/json"
	"github.com/cgalvisleon/elvis/strs"
)

func (s *Linq) strToCols(str string) []*Column {
	var result []*Column = []*Column{}
	str = strs.ReplaceAll(str, []string{" "}, "")
	cols := strings.Split(str, ",")

	for _, n := range cols {
		c := s.GetCol(n)
		if c != nil {
			result = append(result, c)
		}
	}

	return result
}

func (s *Linq) selCols(sel ...any) *Linq {
	var cols []*Column = []*Column{}
	for _, col := range sel {
		switch v := col.(type) {
		case Column:
			cols = append(cols, &v)
		case *Column:
			cols = append(cols, v)
		case []string:
			for _, n := range v {
				c := s.GetCol(n)
				if c != nil {
					cols = append(cols, c)
				}
			}
		case []*Column:
			cols = v
		case string:
			cols2 := s.strToCols(v)
			if len(cols2) == 0 {
				c := s.GetCol(v)
				console.Log(v, ", ", v)
				if c != nil {
					cols = append(cols, c)
				}
			} else {
				cols = append(cols, cols2...)
			}
		default:
			console.ErrorF("Linq select type (%v) value:%v", reflect.TypeOf(v), v)
		}
	}

	s._select = cols

	return s
}

func (s *Linq) Data(sel ...any) *Linq {
	s.SetTp(TpData)
	return s.selCols(sel...)
}

func (s *Linq) Select(sel ...any) *Linq {
	s.SetTp(TpSelect)
	return s.selCols(sel...)
}

/**
*
**/
func (s *Linq) Find() (e.Items, error) {
	s.SqlSelect()

	s.sql = strs.Format(`%s;`, s.sql)

	items, err := s.Query()
	if err != nil {
		return e.Items{}, err
	}

	for _, data := range items.Result {
		s.Details(&data)
	}

	return items, nil
}

func (s *Linq) All() (e.Items, error) {
	s.sql = s.SqlAll()

	items, err := s.Query()
	if err != nil {
		return e.Items{}, err
	}

	for _, data := range items.Result {
		s.Details(&data)
	}

	return items, nil
}

func (s *Linq) First() (e.Item, error) {
	s.sql = s.SqlLimit(1)

	item, err := s.QueryOne()
	if err != nil {
		return e.Item{}, err
	}

	s.Details(&item.Result)

	return item, nil
}

func (s *Linq) Limit(limit int) (e.Items, error) {
	s.sql = s.SqlLimit(limit)

	items, err := s.Query()
	if err != nil {
		return e.Items{}, err
	}

	for _, data := range items.Result {
		s.Details(&data)
	}

	return items, nil
}

func (s *Linq) Page(page, rows int) (e.Items, error) {
	offset := (page - 1) * rows
	s.sql = s.SqlOffset(rows, offset)

	items, err := s.Query()
	if err != nil {
		return e.Items{}, err
	}

	for _, data := range items.Result {
		s.Details(&data)
	}

	return items, nil
}

func (s *Linq) Count() int {
	s.sql = s.SqlCount()

	return s.QueryCount()
}

func (s *Linq) List(page, rows int) (e.List, error) {
	all := s.Count()

	items, err := s.Page(page, rows)
	if err != nil {
		return e.List{}, err
	}

	return items.ToList(all, page, rows), nil
}
