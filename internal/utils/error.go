package utils

import "errors"

func ErrListToError(errList []error) error {
	if len(errList) > 0 {
		errMsg := ""
		for _, err := range errList {
			errMsg += err.Error() + ";\n"
		}
		return errors.New(errMsg)
	}

	return nil
}
