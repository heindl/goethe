// Copyright 2018 Parker Heindl. All rights reserved.
// Licensed under the MIT License. See LICENSE.md in the
// project root for information.
//
package pkg_compile_test

// Added unused function to test that it is replaced when package is converted to main.
func main() {
	RootCmd.Execute()
}
