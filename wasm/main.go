/*
Copyright © 2020 Alessandro Segala (@ItalyPaleAle)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package main

/*
Build with:
GOOS=js GOARCH=wasm go build -o  ../ui/dist/app.wasm
brotli -9k ../ui/dist/app.wasm

The Go WebAssembly runtime is at:
$GOROOT/misc/wasm/wasm_exec.js
*/

import (
	"syscall/js"
)

const MaxSafeInteger = 9007199254740991

func main() {
	// Export a "Prvt" global object that contains our functions
	js.Global().Set("Prvt", map[string]interface{}{
		"decryptRequest": DecryptRequest(),
		"getIndex":       GetIndex(),
	})

	// Prevent the function from returning, which is required in a wasm module
	select {}
}
