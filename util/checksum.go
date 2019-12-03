// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// CalcFileChecksum calculate the file md5 code.
func CalcFileChecksum(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		return ""
	}

	defer f.Close()

	m := md5.New()
	if _, err = io.Copy(m, f); err != nil {
		return ""
	}

	return hex.EncodeToString(m.Sum(nil))
}
