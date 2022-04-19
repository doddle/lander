package deployments

import (
	"fmt"
	"strings"
)

func ParseShowTagsFor(input string) (result []TagFilters) {
	commaList := strings.Split(input, ",")
	if len(commaList) > 0 {
		for _, item := range commaList {
			// fmt.Println("parsing item: " + item)
			splitAgain := strings.Split(item, "=")
			// fmt.Println(len(splitAgain))
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
