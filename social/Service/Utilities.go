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
