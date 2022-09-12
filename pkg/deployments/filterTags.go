package deployments

import (
	"fmt"
	"strings"
)

func ParseShowTagsFor(input string) (result []TagFilters) {
	commaList := strings.Split(input, ",")
	if len(commaList) > 0 {
		for _, item := range commaList {
			splitAgain := strings.Split(item, "=")
			if len(splitAgain) == 1 {
				fmt.Printf("couldn't split '%s' by '='\n", item)
				continue
			}
			foo := TagFilters{
				Registry: splitAgain[0],
				Name:     splitAgain[1],
			}
			result = append(result, foo)
		}
	}
	return result
}
