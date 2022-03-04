# Welcome to Glimmer ✨
Glimmer is a dynamically-typed scripting language made by me, Daniel Johnson, which aims to have clean and simple syntax while drawing some inspiration from each of my favorite languages and topics.

This implementation of Glimmer comes with a complete Read-Eval-Print-Loop (REPL), Read-Parse-Print-Loop (RPPL), and Read-Lex-Print-Loop (RLPL). Of course, you can also execute source files directly.

# Features

### Number Arithmetic
<img src="https://github.com/DanielRJohnson/Glimmer/blob/main/_examles/terminal%20images/numbers.png" alt="number type arithmetic" width="66%"/>

### Strings & Their Operations
<img src="https://github.com/DanielRJohnson/Glimmer/blob/main/_examles/terminal%20images/strings.png" alt="string arithmetic" width="66%"/>

### Arrays & Builtin Array Functions
<img src="https://github.com/DanielRJohnson/Glimmer/blob/main/_examles/terminal%20images/arrays.png" alt="array operations" width="66%"/>

### Dictionaries
<img src="https://github.com/DanielRJohnson/Glimmer/blob/main/_examles/terminal%20images/dicts.png" alt="Dictionary example" width="66%"/>

### First-Class Functions
<img src="https://github.com/DanielRJohnson/Glimmer/blob/main/_examles/terminal%20images/first-class%20functions.png" alt="first-class function example" width="66%"/>

### Statically-Scoped Variables
<img src="https://github.com/DanielRJohnson/Glimmer/blob/main/_examles/terminal%20images/static%20scoping.png" alt="static variable example" width="66%"/>

### Recursion
<img src="https://github.com/DanielRJohnson/Glimmer/blob/main/_examles/terminal%20images/recursion.png" alt="recursion example" width="66%"/>

# Usage
* To run a source file, run `glimmer <my source file>`
* To open the Glimmer REPL, run `glimmer`
* To open the Glimmer RPPL, run `glimmer -p`
* To open the Glimmer RLPL, run `glimmer -l`
* When evaluating and parsing, you can also use the flag --dot to generate a dotfile & image for the AST of your input.

# Changelog
* V1.0: Base Language Push
* V1.1: Added "For" construct as well as assignment and arithmetic assignment (i.e. +=)

# TODO
* Line and Col in parser error messages
* OS interaction (exec, input, etc)
* Imports & standard library
* More builtins (casts, etc.)
* Static typing (maybe separate fork or something)
* Some semblance of objects/structs/data
* Bytecode interpreter (down the road)

# Credit
Much of the methodologies, code, and knowledge in the writing of this came from Thorsten Ball's book, Writing an Interpreter in Go. I wrote every line in this repo character by character without copy-pasting, changed methods where I saw fit, and added much on top of the code from this book. Reading this was a great inspiration, and I give my sincere thanks to Mr. Ball. Check out the book at https://interpreterbook.com/.
