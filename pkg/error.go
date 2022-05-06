package pkg

import (
// External
)

func CheckError(err error) {

	if err != nil {
		panic(err)
	}
}
