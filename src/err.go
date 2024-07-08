package src

func assert(err error) {
    if err != nil {
        panic(err)
    }
}

func verify(condition bool, message string) {
	if !condition {
		panic(message)
	}
}
