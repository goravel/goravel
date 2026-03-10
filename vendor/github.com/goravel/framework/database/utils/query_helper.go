package utils

import "github.com/goravel/framework/errors"

func PrepareWhereOperatorAndValue(args ...any) (op any, value any, err error) {
	if len(args) == 0 || len(args) > 2 {
		return nil, nil, errors.DatabaseInvalidArgumentNumber.Args(len(args), "1 or 2")
	}

	if len(args) == 1 {
		op = "="
		value = args[0]
	} else {
		op = args[0]
		value = args[1]
	}

	return
}
