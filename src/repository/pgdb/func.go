package pgdb

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func (g *GoPg) makeFilters(filters []string) (query string, args []interface{}) {
	var key, operations, value, ff, b1, b2 string
	args = make([]interface{}, 0)

	for _, v := range filters {
		slFilter := strings.Split(v, ":")
		key = slFilter[0]
		operations = slFilter[1]
		value = slFilter[2]
		ff = slFilter[3]

		slBetween := strings.Split(value, "|")
		if len(slBetween) == 2 {
			b1 = slBetween[0]
			b2 = slBetween[1]
		}

		if query != "" {
			switch ff {
			case "and", "And", "AND":
				query += " AND "
			case "or", "Or", "OR":
				query += " OR "
			default:
				query += " AND "
			}
		}

		switch operations {
		case "eq":
			query += key + " = ?"
			args = append(args, value)
			break
		case "ne":
			query += key + " <> ?"
			args = append(args, value)
			break
		case "gt":
			query += key + " > ?"
			args = append(args, value)
			break
		case "gte":
			query += key + " >= ?"
			args = append(args, value)
			break
		case "lt":
			query += key + " < ?"
			args = append(args, value)
			break
		case "lte":
			query += key + " <= ?"
			args = append(args, value)
			break
		case "like":
			query += key + " LIKE ?"
			args = append(args, "%"+value+"%")
			break
		case "ilike":
			query += key + " ILIKE ?"
			args = append(args, "%"+value+"%")
			break
		case "isNull":
			query += key + " IS NULL"
			args = append(args, "")
			break
		case "isNotNull":
			query += key + " <> ?"
			args = append(args, "")
			break
		case "btw":
			query += key + " BETWEEN ? AND ?"
			args = append(args, b1, b2)
			break
		case "in":
			query += key + " IN (?)"
			in := func() { strings.Split(value, "|") }
			args = append(args, in)
			break
		default:
			query += key + " = ?"
			args = append(args, value)
			break
		}
	}
	return query, args
}

//	func makeIn(val string) (result []string) {
//		result = strings.Split(val, "|")
//		return
//	}
func (g *GoPg) interfaceToSlice(slice interface{}) (res []interface{}) {
	s := reflect.Indirect(reflect.ValueOf(slice))
	if s.Kind() != reflect.Slice {
		log.Panic("InterfaceSlice() given a non-slice type")
	}
	if s.IsNil() {
		return nil
	}

	res = make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		res[i] = s.Index(i).Interface()
	}
	return res
}

func (g *GoPg) makeSorts(sorts []string) (strSort string) {
	for i, v := range sorts {
		slFilter := strings.Split(v, ":")
		if i > 0 {
			strSort += ", "
		}
		strSort += fmt.Sprintf("%s %s", slFilter[0], slFilter[1])
	}
	return
}
