// Copyright (c) 2018 Parker Heindl. All rights reserved.
//
// Use of this source code is governed by the MIT License.
// Read LICENSE.md in the project root for information.

package pkg_compile_test

// Added unused function to test that it is replaced when package is converted to main.
func main() {
	RootCmd.Execute()
}
