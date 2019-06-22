package social

import "errors"

func AppendIfNotExists(appendto []string, appendelement string) ([]string, error) {
	for _, el := range appendto {
		if el == appendelement {
			return nil, errors.New("Element already exist")
		}
	}

	return append(appendto, appendelement), nil
}

func RemoveIfExists(removefrom []string, removeelement string) ([]string, error) {
	var returnVal []string
	var elementFound bool = false
	for _, el := range removefrom {
		if el == removeelement {
			elementFound = true
			continue
		}
		returnVal = append(returnVal, el)
	}

	if elementFound == false {
		return nil, errors.New("Elemement not found")
	}
	return returnVal, nil
}
