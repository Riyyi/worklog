/*
 * Copyright (C) 2024 Riyyi
 *
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

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
