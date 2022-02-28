# Welcome to Glimmer ✨
Glimmer is a dynamically-typed scripting language with support for first-class functions with closures, basic types, arrays, dictionaries, and much more.

This implementation of Glimmer comes with a complete Read-Eval-Print-Loop (REPL), Read-Parse-Print-Loop (RPPL), and Read-Lex-Print-Loop (RLPL). Of course, you can also execute source files directly.

# Features
TODO README SECTION

# Usage
* To run a source file, run `glimmer <my source file>`
* To open the Glimmer REPL, run `glimmer`
* To open the Glimmer RPPL, run `glimmer -p`
* To open the Glimmer RLPL, run `glimmer -l`
* When evaluating and parsing, you can also use the flag --dot to generate a dotfile & image for the AST of your input.

# TODO
* Imports & standard library
* More builtins (casts, etc.)
* Static typing (maybe separate fork or something)
* Some semblance of objects/structs/data
* Bytecode interpreter (wayyyyyy down the road)

# Credit
Much of the methodologies, code, and knowledge in the writing of this came from Thorsten Ball's book, Writing an Interpreter in Go. I wrote every line in this repo character by character without copying, changed methods where I saw fit, and added much on top of the code from this book. Reading this was a great inspiration, and I give my sincere thanks to Mr. Ball. Check out the book at https://interpreterbook.com/.
