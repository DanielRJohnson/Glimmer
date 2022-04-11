# Welcome to Glimmer ✨
Glimmer is a fully statically-typed scripting language which aims to have clean and simple syntax while drawing some inspiration from each of my favorite languages and topics.

This implementation of Glimmer comes with a complete Read-Eval-Print-Loop (REPL), Read-Parse-Print-Loop (RPPL), and Read-Lex-Print-Loop (RLPL). Of course, you can also execute source files directly.

# Features

## Number Arithmetic
 - numeric types supported are integer, float, and boolean
 - integers and floats are 64 bits, borrowing Go's typing
 - numeric types are defined over +, -, *, / with promotion

```
>> (1 + 1) # integer arithmetic
2
>> (true + true) # boolean promotion
2
>> (1 + true + 2.2) # integer promotion
4.2
```

## Strings & Their Operations
 - strings are defined over +, -, *, / with other strings
 - pythonic (string * number) is also defined

```
>> ("a" + "b") # concat
ab
>> ("aa" - "a") # remove first
a
>> ("ab" * "ac") # cross product
aaacbabc
>> ("aabbaaa" / "aa") # remove all
bba
>> ("a" * 4) # repeat N times
aaaa
```

## Arrays & Builtin Array Functions
 - Arrays are immutable objects with indexing as the only operation
 - Builtin functions are used to make working with arrays nicer

```
>> [1, 2, 3, 4][2]
3
>> head([1, 2, 3, 4])
1
>> push([1, 2, 3, 4], 5)
[1, 2, 3, 4, 5]
```

## Dictionaries
 - Dictionaries are objects of pairs indexed by strings
 - More functionality is planned in the future

```
>> {"a": 1, "b": 2}["a"]
1
>> key = "a"; {"a": 1, "b": 2}[key]
1
>> key = "a"; {key: 1, "b": 2}["a"]
1
```

## Variable Declaration and Assignment
 - Assignment binds an identifier to a value in an environment
 - Reassignment updates the value for the identifier
 - Values include integers, floats, booleans, strings, arrays, dictionaries, and functions

```
>> x = 5
>> x
5
>> myArr = [1, 2, 3, 4, 5]
>> myArr[2]
3
>> myDict = {"a": 1, "b": 2}
>> myDict["a"]
1
```

## First-Class Functions
 - Functions are first-class values that can be applied to parameters 
 - Functions are statically scoped, allow recursion, return the last statement if no explicit return has happened 
 - Note: Glimmer is whitespace-agnostic so while the examples shown are on one line, you may have any indentation/newlines you want in a file.

```
>> inc = fn(x: int) -> int { x + 1 }
>> applyTwice = fn(f: fn(int) -> int, x: int) -> int { f(f(x)) }
>> applyTwice(inc, 1)
3
>> fact = fn(n: int) -> int { if n == 0 { 1 } else { fact(n - 1) * n } }
>> fact(5)
120
```

## If Expressions
 - Ife's are expressions in Glimmer that evaluate to the last statement of which branch gets evaluated
 - The condition of an ife is also multi-statement and evaluates to the last statement
 - truthy values are not null, false, or zero
 - Any amount of "else ife" branches are allowed that are also have multi-statement conditions
 - Funcions are the only scope extenders, so the blocks of an ife operate in the same environment as its parent

```
>> ife (true) { 1 } else { 0 }
1
>> ife x = 5; x > 4 { 1 } else { 0 }
1
>> ife (false) { 1 } else ife (true) { 0 } else { 1 }
0
>> ife x = 5; x <= 4 { 1 } else ife x -= 1; x <= 4 { 0 } else { 1 }
0
```

## If Statements
 - Like if expressions, but without the burden of returning a value with strict typing
```
>> if true { print("hi") } else { return 1 }
hi
null # statements return null
```

## For Loops
 - Loop over collections with ease using a for-in structure
 - arrays give you value with one loop variable, index then value with two
 - dicts give you key with one loop variable, key then value with two

```
>> for v in range(5) { print(v * v) }
0
1
4
9
16
null # loops are statements, thus return null

>> for i, v in range(4) { print("i: ", i, " v^2: ", v*v) }
i: 0 v^2: 0
i: 1 v^2: 1
i: 2 v^2: 4
i: 3 v^2: 9
null # loops are statements, thus return null
```

## While Loops
 - Your standard while construct that we all need, as long as you're not an expert in lambda calculus
 - However, we still have the extra pizazz of multi-statement conditions
 - truthy values are not null, false, or zero

```
>> x = true; while x { print("PING"); x = false }
PING
null

>> x = 0; while x += 1; x < 5 { print(x) }
1
2
3
4
null
```

## Other Builtins
 - Builtin functions can be found in `evaluator/builtins.go`
 - Many more are planned in the future, as well as a library structure
 - However, we all know the most important one

```
>> print("Hello, World!")
Hello, World!
```

## Full Static Typing!!!
 - Static typing means that the language makes some concessions to determine the type of every object in the program before the program even runs. This leads to many less weird runtime errors, less crashes = good. These concessions include:
     - manually fixing fn arguments and return type
     - containers must hold only one type
     - all branches of an `ife` expression must match types

```
>> 1 + "string"
Static TypeError at [1,3]: infix operator for 'int + string' not found
>> push([ fn(x: int) -> int { x }, fn(x: int) -> int { x } ], "not fn")
Static TypeError at [1,5]: Argument 2 to push must be match Argument 1's held type: fn(int) -> int, got=string
```

# Usage
* To run a source file, run `glimmer <my source file>`
* To open the Glimmer REPL, run `glimmer`
* To open the Glimmer RPPL, run `glimmer -p`
* To open the Glimmer RLPL, run `glimmer -l`
* When evaluating and parsing, you can also use the flag `--dot` to generate a dotfile & image for the AST of your input.

# Changelog
* V0.0: Base Language Push
* V0.1: Added `for` construct as well as assignment and arithmetic assignment (i.e. +=)
* V0.2: Added line and col numbers for parser errors, multi-line `ife`'s, and deprecated let in favor of defining and updating assignment 
* V0.3: Added static typing, changing function syntax `fn() -> none { print("WOOHOO") }()`
* V0.4: Resigned `for`, and added `if` (non-valued if statements), `while`, and `range` 

# Possible Future Work
Near:
* async-finish blocks?
* OS interaction (exec, input, etc)
* `in` as an infix operator
* Imports & standard library/ more builtins
* More dict functionality

Far:
* Some semblance of objects/structs/data
* Bytecode interpreter (down the road)

# Credit
Much of the methodologies, code, and knowledge in the writing of this came from Thorsten Ball's book, Writing an Interpreter in Go. I wrote every line in this repo character by character without copy-pasting, changed methods where I saw fit, and added much on top of the code from this book. Reading this was a great inspiration, and I give my sincere thanks to Mr. Ball. Check out the book at https://interpreterbook.com/.
