package model

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Int64toString []int64

func (is Int64toString) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	b.Grow(64)
	b.WriteByte('[')
	for i, v := range is {
		b.WriteByte('"')
		b.Write([]byte(strconv.FormatInt(v, 10)))
		b.WriteByte('"')
		if i+1 != len(is) {
			b.WriteByte(',')
		}
	}
	b.WriteByte(']')
	return b.Bytes(), nil
}

func (is *Int64toString) UnmarshalJSON(data []byte) error {
	str := strings.Split(strings.Trim(string(data), "[]"), ",")
	for _, v := range str {
		num, err := strconv.ParseInt(strings.Trim(v, " \""), 10, 64)
		if err != nil {
			fmt.Println(err)
			return err
		}
		*is = append(*is, num)
	}
	return nil
}

// Scan 方法实现了 sql.Scanner 接口
func (is *Int64toString) Scan(v interface{}) error {
	if s, ok := v.([]uint8); !ok {
		return fmt.Errorf("断言失败")
	} else {
		str := strings.Split(string(s), ",")
		for _, v := range str {
			num, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			*is = append(*is, num)
		}
	}
	return nil
}

func (is Int64toString) Value() (driver.Value, error) {
	if len(is) == 0 {
		return nil, nil
	}
	s := make([]string, 0, 4)
	for _, v := range is {
		s = append(s, strconv.FormatInt(v, 10))
	}
	return strings.Join(s, ","), nil
}
