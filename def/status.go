// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package def

type STATUS uint8

const (
	NOTSTART  STATUS = 0
	SUCCESS   STATUS = 1
	ONGOING   STATUS = 2
	PENDING   STATUS = 3
	FAILURE   STATUS = 4
	CANCEL    STATUS = 5
	INTERRUPT STATUS = 6
)

func IsCompleted(status STATUS) bool {
	if status == SUCCESS || status == FAILURE || status == CANCEL {
		return true
	}

	return false
}
