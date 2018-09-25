// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package utilities

import (
	"fmt"
	"os"
)

// MakeRecommendation propagates notes to the user about ways to improve the project for Open Source publishing.
func MakeRecommendation(note string) {
	fmt.Fprint(os.Stdout, note)
}
