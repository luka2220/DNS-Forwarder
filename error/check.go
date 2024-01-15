package error

import "log"

func Check(err error) {
	if err != nil {

		log.Fatalf("An error occured: %v", err)
	}
}
