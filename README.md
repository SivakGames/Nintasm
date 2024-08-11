# Nintasm - A 6502 Assembler for the NES

This project is a 6502 assembler written entirely in Golang with the intended purpose of building NES games.

The syntax used is based on the original NESASM assembler. NESASM itself has undergone several iterations by different contributors over the years. Nintasm, however, seeks to take the general concept of NESASM and add onto it more features as well as rectify some issues the orignal program had.

## **\*\*Important Note About This Readme\*\***

This readme is **_NOT_** a guide on how to program in 6502 assembly nor does it contain information on how to develop games for the NES. Instead, it focuses on explaining how to use the assembler and highlights the available features.

There are, however, sample programs included in the repository that can be built with Nintasm and will run on the NES or in an NES emulator.

&nbsp;

---

## Sections

- [Nintasm - A 6502 Assembler for the NES](#nintasm---a-6502-assembler-for-the-nes)
	- [**\*\*Important Note About This Readme\*\***](#important-note-about-this-readme)
	- [Sections](#sections)
	- [Known Bugs/Issues/Todo](#known-bugsissuestodo)
	- [Major Issues Fixed/Rectified from NESASM](#major-issues-fixedrectified-from-nesasm)
	- [Where to get the executable file](#where-to-get-the-executable-file)
	- [Building a Nintasm executable from the source](#building-a-nintasm-executable-from-the-source)
	- [Running the Program](#running-the-program)
		- [Command Line Options](#command-line-options)
	- [General Assembler Syntax](#general-assembler-syntax)
		- [**Labels**](#labels)
		- [**Instructions**](#instructions)
		- [**Directives**](#directives)
		- [**Comments**](#comments)
		- [**Sample Program Snippet**](#sample-program-snippet)
	- [Working With Numbers](#working-with-numbers)
		- [Example of General Number Use](#example-of-general-number-use)
		- [Negative Numbers](#negative-numbers)
			- [Examples of using negative numbers](#examples-of-using-negative-numbers)
	- [Working With Mathematical Expressions](#working-with-mathematical-expressions)
		- [Supported Operations](#supported-operations)
			- [Examples of doing some basic math operations](#examples-of-doing-some-basic-math-operations)
		- [Using Parentheses](#using-parentheses)
			- [Examples with Parentheses](#examples-with-parentheses)
	- [Working with Strings](#working-with-strings)
		- [Example of Using Strings](#example-of-using-strings)
	- [Working with Arrays](#working-with-arrays)
	- [Assembler Constants](#assembler-constants)
		- [PPU Constants](#ppu-constants)
			- [PPU Control ($2000)](#ppu-control-2000)
			- [PPU Mask ($2001)](#ppu-mask-2001)
			- [PPU Status ($2002)](#ppu-status-2002)
			- [PPU Scroll ($2005)](#ppu-scroll-2005)
			- [PPU Address ($2006)](#ppu-address-2006)
			- [PPU Data ($2007)](#ppu-data-2007)
		- [OAM Attribute Constants](#oam-attribute-constants)
		- [Controller Button Constants](#controller-button-constants)
	- [User-Defined Symbols](#user-defined-symbols)
		- [Simple Assignment](#simple-assignment)
		- [Reserve Method](#reserve-method)
		- [Namespaces](#namespaces)
	- [Assembler Directives](#assembler-directives)
		- [iNES Header Directives](#ines-header-directives)
			- [Required iNES Directives](#required-ines-directives)
			- [Optional iNES Directives](#optional-ines-directives)
			- [PRG and CHR alternate definitions](#prg-and-chr-alternate-definitions)
		- [RomSegment Directives](#romsegment-directives)
		- [Bank and Org Directives](#bank-and-org-directives)
		- [File Inclusion Directives](#file-inclusion-directives)
		- [Raw Byte Directives](#raw-byte-directives)
			- [Mixed Byte Directives](#mixed-byte-directives)
		- [Repeat Directive](#repeat-directive)
		- [Range Directive](#range-directive)
		- [GNSI Directive](#gnsi-directive)
		- [Reassign Directive](#reassign-directive)
		- [Misc Directives](#misc-directives)
	- [Conditional Directives](#conditional-directives)
		- [If Statement Directives](#if-statement-directives)
		- [Switch Statement Directives](#switch-statement-directives)
	- [Assembler Functions](#assembler-functions)
		- [Number and Math Functions](#number-and-math-functions)
		- [String functions](#string-functions)
		- [Array functions](#array-functions)
		- [Symbol-Related Functions](#symbol-related-functions)
		- [Namespace functions](#namespace-functions)
		- [Bytes Under Label Counting Functions](#bytes-under-label-counting-functions)
		- [Examples of Using Functions](#examples-of-using-functions)
	- [User-Defined Functions](#user-defined-functions)
	- [Macros](#macros)
		- [Defining and using a simple macro](#defining-and-using-a-simple-macro)
		- [Using arguments in macros](#using-arguments-in-macros)
		- [Using anonymous arguments in macros](#using-anonymous-arguments-in-macros)
			- [Checking supplied argument number](#checking-supplied-argument-number)
		- [Exiting a Macro Early](#exiting-a-macro-early)
	- [Character Maps](#character-maps)
		- [Directives Used Within Character Maps](#directives-used-within-character-maps)
		- [Other Character Map Directives](#other-character-map-directives)
	- [Expression Maps](#expression-maps)
		- [Other Expression Map Directives](#other-expression-map-directives)
		- [Examples of Using Expression Maps](#examples-of-using-expression-maps)
	- [Dynamic Labels and Strings](#dynamic-labels-and-strings)
		- [Examples of Using Dynamic Labels](#examples-of-using-dynamic-labels)
	- [**Happy Assembly!**](#happy-assembly)

---

&nbsp;

## Known Bugs/Issues/Todo

- There may be situations where error output may be that from Golang involving unanticipated operations instead of the standard assembler error output.
- On the second pass, if there are errors, the line number display is sometimes off
- Currently there are no safeguards for circular `.include` statements, function calling, or macro invoking.

[] TODO: Allow iNES header version 2.0 format
[] TODO: Implement the ines battery flag
[] TODO: Add support for escaped characters in strings
[] TODO: (low priority) Support illegal 6502 opcodes
[] TODO: (low priority) Clean up the readability of the block stack handler
[] TODO: (low priority) Implement rom segment description parameter

## Major Issues Fixed/Rectified from NESASM

The following were some of the most significant limitations in the original NESASM that have been addressed and resolved in Nintasm:

- Various erroneous actions that assembled without any error or warning output will no longer work and actually show error output.
- Banks can now be sizes other than 8KB
- If a label is known to be zero page, then it will be recognized as such and is no longer necessary to write the `<` character before the operand. (It is, however, still possible to use this to explicitly state zero page).
- Label names can be much longer (the original limit was 32 characters).
- `.db`/`.byte` and `.dw`/`.word` blocks can accommodate a larger amount of data per line than before.
- Local labels within a parent label can be directly referenced. For example, if a parent label is `ABC` and it has a local label `.def`, it can be referenced in a JMP command as `JMP ABC.def`.

## Where to get the executable file

Look in the `/dist` folder of this repository to find a binary of the latest version.  Currently the binary is available for Windows only.

## Building a Nintasm executable from the source

To build a binary executable yourself you will need Golang installed on your system. From the root folder, run:

```text
go build ./nintasm.go
```

The final executable should show up in the same folder.

## Running the Program

Simply put the `nintasm.exe` file in the desired folder and type:

```text
nintasm <SOURCE_FILE> [-options]
```

### Command Line Options

```text
  -h  Show help about the command line
  -o  Change the output filename to something else (Simply type the desired name after -o)
  -r  Write a raw file without an iNES header
  -s  Show segment usage of each bank (i.e. how much of each bank has been used versus how much is available)
```

&nbsp;

---

&nbsp;

## General Assembler Syntax

Programs have these basic operations: labels, directives, and 6502 instructions. There can only be one operation per line and a line break will terminate the operation.

### **Labels**

- Labels are any text (up to the programmer what to call them) at the start of a line. (i.e. NO whitespace can precede a label). The general purpose is to give a name to a particular section in the ROM that can be referenced.
- Label names must start with a letter or underscore, followed by any combination of letters, underscores, or numbers.
- Label names must also end with a single colon.
- It's important to note that label names are case-sensitive. (e.g. `LDA Label` vs `LDA label` are **NOT** the same thing)
- It's possible to define local labels underneath a parent label. These follow the same rules and naming conventions as a label, however, they start with a period ( `.` ). Local labels can be referenced anywhere in the code by placing their parent label in front of them. However, while inside the local scope, the parent label's name can be omitted.

### **Instructions**

The standard set of 56 instructions offered by the 6502 processor are available.

- Instructions must have _at least_ one whitespace character preceding them.
- Instruction names are NOT case sensitive. For example, `lda` and `LDA` will have the same effect. However, it is generally recommended to consistently use a specific case throughout the program.

### **Directives**

Directives are operations to assist in building certain parts of the game. These operations, like instructions, need _at least_ one whitespace character before them, are case-insensitive, and start with a period immediately followed by the directive name. (e.g. `.byte` )

Main section: [Assembler Directives](#assembler-directives)

### **Comments**

- Comments are simply any text that follows a semicolon ( `;` ) on a line. Note that there's no way to write block comments.

### **Sample Program Snippet**

This simple program illustrates the usage of the basic building blocks of a program.

```text
;==============================================================
; This is a comment as it has a semicolon preceding it
; You can write as much text as you want, but
; be sure to keep adding semicolons at the start
; of new lines if your plan on writing multiple lines of comments.
;==============================================================

;Below here is a label called Label1 and below the label is a .BYTE directive which is used to store some bytes in the ROM.

Label1:
    .byte $80,%11000001,33 ;The numbers used appear in hexadecimal, binary, and decimal

;Below here is another label called Label2 containing some 6502 instructions.
; * These instructions don't really do much programmatically.

;Below here is another label called Label2 containing some 6502 instrutions as well as a local label called .loop

Label2:
    LDX #2
.loop:
    LDA Label1, X
    STA $0200, X
    DEX
    BPL .loop
    RTS
```

---

## Working With Numbers

While this concept is standard across 6502 assemblers, it's important to note that numbers can be expressed in decimal, hexadecimal, or binary formats.

- Decimal numbers are simply written as-is (10 is 10)
- Hexadecimal numbers are written with a dollar sign ( `$` ) in front of them (e.g. A hex value of `0x1a` would be written as `$1a`). The letters in the number are also case insensitive. In this case, $1a or $1A are acceptable.
- Binary numbers are written with a percent sign ( `%` ) in front of them (e.g. A binary value of `0b01111010` would be written as `%01111010` )
- _\*Note: Octal is not supported_

### Example of General Number Use

```text
;Load the decimal number 11 into A
    LDA #11
;Load the hexadecimal number 0x0b (which is also 11) into A
    LDA #$0b
;Load the binary number 0b00001011 (which is also 11) into A
    LDA #%00001011
```

### Negative Numbers

It's also possible to express values as negative numbers. After assembling, the negative number will be converted into its proper signed representation.

#### Examples of using negative numbers

```text
;Within a .BYTE directive
    .BYTE -1,-2,-3 ;Would become .BYTE 255,254,253 (or $ff,$fe,$fd)
;Within an instruction
    LDA #-4 ;Would become LDA #252 (or $fc)
;Using negative hex...
    LDA #-$0b
;Using negative binary...
    LDA #-%00001100
```

---

## Working With Mathematical Expressions

Sometimes you might want to use a mathematical expression to generate or manipulate a value.

### Supported Operations

These are all the supported unary/binary operations and precedence is the same as the standard order of operations.

```text
+, -, *, /, %, ^, &, |, ~, <<, >>
&&, ||, !=, ==, <, >, <=, >=
```

Ternary operations ( `?:` ) are also supported.

```text
VAR1 = 5
VAR2 = 6 > VAR1 ? 1 : 0 ;Would set to 1
```

#### Examples of doing some basic math operations

```text
ABC = 5
DEF = ABC + 1  ;Will be 6
GHI = ABC - 2  ;Will be 3
```

### Using Parentheses

In addition, parentheses can be used for specifying order of operations. (Multiple sets of parentheses may also be used)

#### Examples with Parentheses

```text
ABC = 5
DEF = (1 + ABC) * 2  ;Will be 12
GHI = 1 + ABC * 2    ;Will be 11
```

---

## Working with Strings

It's possible to use strings in `.db`/`.byte` directives or by assigning them to constants. Strings can be enclosed in either single quotes ( `'` ) or double quotes ( `"` ).

When working with strings, each character is represented by one byte and corresponds to the ASCII value of that character.

- Escaped characters (e.g. "\"") are NOT supported.
- Please note: Unicode is supported but may result in multiple bytes per character. A warning will be shown if this occurs.

### Example of Using Strings

```text
    .BYTE "ABC"  ;Would evaluate to .BYTE $41,$42,$43 (the ASCII values of each character)
    .BYTE 'DEF'  ;Would evaluate to $44,$45,$46
```

---

## Working with Arrays

It's also possible to assign multiple values to a constant using square brackets (`[` `]`)

If using this data directly, the array will be flattened and all bytes will come out sequentially

```text
MyArray = [1,2,3,4]
MyArray2 = [5,"ABC"]

;The following would evaluate to $01,$02,$03,$04,$05,$41,$42,$43
    .db MyArray, MyArray2
```

---

## Assembler Constants

The assembler comes with some built-in constants for commonly used memory addresses.

### PPU Constants

A lot of these contain the address name as well as individual bits associated with setting things related to the respective address.

#### PPU Control ($2000)

- `PPUCTRL` - $2000
- `PPUCTRL.nameTable0` - $00
- `PPUCTRL.nameTable1` - $01
- `PPUCTRL.nameTable2` - $02
- `PPUCTRL.nameTable3` - $03
- `PPUCTRL.drawDirection` - $04
- `PPUCTRL.spritePatternTable` - $08
- `PPUCTRL.bgPatternTable` - $10
- `PPUCTRL.use8x16Sprites` - $20
- `PPUCTRL.masterSlave` - $40
- `PPUCTRL.enableNMI` - $80

Use Example

```text
;Would set the PPU Control so that NMIs are enabled and that sprites are in the $1000 pattern table
;Evalutes to %10001000
  LDA #(PPUCTRL.enableNMI | PPUCTRL.spritePatternTable)
  STA PPUCTRL
```

#### PPU Mask ($2001)

- `PPUMASK` - $2001
- `PPUMASK.grayscale` - $01
- `PPUMASK.disableBgClip` - $02
- `PPUMASK.disableSpriteClip` - $04
- `PPUMASK.showBg` - $08
- `PPUMASK.showSprites` - $10
- `PPUMASK.emphasizeRed` - $20
- `PPUMASK.emphasizeGreen` - $40
- `PPUMASK.emphasizeBlue` - $80

Use Example

```text
;Would set the PPU Mask so that both sprites and BG are on and neither are clipped in the left column
;Evalutes to %00011110
  LDA #(PPUMASK.showSprites | PPUMASK.showBg | PPUMASK.disableBgClip | PPUMASK.disableSpriteClip)
  STA PPUMASK
```

#### PPU Status ($2002)

- `PPUSTATUS` - $2002

#### PPU Scroll ($2005)

- `PPUSCROLL` - $2005

#### PPU Address ($2006)

> Note: Do not use a leading zero for the operations that involve replacing XX with a number.

- `PPUADDR` - $2006
- `PPUADDR.nt0lineXX` - Will evaluate to `$2000 + (XX * $20)`. Replace XX with a value from 0-29. Useful for setting the address to a specific line of the background on Name Table 0.
- `PPUADDR.nt1lineXX` - `$2400 + (XX * $20)`
- `PPUADDR.nt2lineXX` - `$2800 + (XX * $20)`
- `PPUADDR.nt3lineXX` - `$2c00 + (XX * $20)`
- `PPUADDR.nt0attLineXX` - Will evaluate to `$23c0 + (XX * $08)`. Replace XX with a value from 0-7. Useful for setting the address to a specific line of the attribute table on Name Table 0.
- `PPUADDR.nt1attrLineXX` - `$27c0 + (XX * $08)`
- `PPUADDR.nt2attrLineXX` - `$2bc0 + (XX * $08)`
- `PPUADDR.nt3attrLineXX` - `$2fc0 + (XX * $08)`
- `PPUADDR.palBgXX` - Will evaluate to `$3f00 + (XX * $04)`. Replace XX with a value from 0-3. Useful for setting the background palette to a specific spot.
- `PPUADDR.palSpriteXX` - Will evaluate to `$3f10 + (XX * $04)`. Replace XX with a value from 0-3. Useful for setting the sprite palette to a specific spot.

Use Example

```text
;Would set the PPU address to $2100 (line 8, column 0 on name table 0)
  LDA #high(PPUADDR.nt0line8)
  STA PPUADDR
  LDA #low(PPUADDR.nt0line8)
  STA PPUADDR

;Would set the PPU address to $2232 (line 17, column 18 on name table 0)
  LDA #high(PPUADDR.nt0line17 + 18)
  STA PPUADDR
  LDA #low(PPUADDR.nt0line17 + 18)
  STA PPUADDR
```

#### PPU Data ($2007)

- `PPUDATA` - $2007

### OAM Attribute Constants

Like with PPU constants, these can be useful when writing a mask for sprite attributes.

- `OAMATTR.pal0` - $00
- `OAMATTR.pal1` - $01
- `OAMATTR.pal2` - $02
- `OAMATTR.pal3` - $03
- `OAMATTR.priority` - $20
- `OAMATTR.flipHoriz` - $40
- `OAMATTR.flipVert` - $80

### Controller Button Constants

While not necessarily NES-specific, most games will store a controller read into some kind of memory and this is based on the order that they are often stored into said memory.

- `CTRLBTN.right` - $01
- `CTRLBTN.left` - $02
- `CTRLBTN.down` - $04
- `CTRLBTN.up` - $08
- `CTRLBTN.start` - $10
- `CTRLBTN.select` - $20
- `CTRLBTN.b` - $40
- `CTRLBTN.a` - $80

---

## User-Defined Symbols

You can assign a meaningful name to a value that can be used in different parts of the program. There are a few ways to do this.

### Simple Assignment

The simplest way is to use either an equal sign ( `=` ) or the letters `EQU` between a label name and the desired value. (both are synonymous)

**IMPORTANT:** When assigning names to constants there can be **no whitespace** before the expression.

Basic examples:

```text
ABC = 1     ;Sets the name of ABC equal to 1
DEF EQU 2   ;Sets the name of DEF equal to 2
 GHI = 3    ;This is INVALID because of the space before the GHI
```

Symbols can then be used as operands in both directives and instructions:

```text
ABC = 1
DEF = 255

;Would put 3 bytes with a value of 1 in this part of the ROM
    .BYTE ABC, ABC, ABC

;Would load in a value of 255 into A
    LDA #DEF
```

### Reserve Method

This method is used for setting a symbol using an internal counter and specifying how much to increment the counter after setting it. This can be useful for specifying a block of memory that should be `n` bytes long, for instance.

To use this, you must first add an `.rsset` directive which is followed by a number indicating where the counter should start. From there, you simply put in a unique label followed by `.rs` followed by the desired value of which to increment the counter.

Example using `.rsset` and `.rs`

```text
 .rsset $0000    ;Set the reserve counter to 0
rsBlock1 .rs 10  ;rsBlock1 will be equal to 0 and increment the counter by 10
rsBlock2 .rs 20  ;rsBlock2 will be equal to 10 and increment the counter by 20
rsBlock3 .rs 5   ;rsBlock3 will be equal to 30 and increment the counter by 5

; ... etc. (The counter is currently 35 here in this example)
```

### Namespaces

It's also possible to define a series of constants within a namespace, which can be helpful for organizing and grouping related variables. This can be accomplished by using a label name and then following it with a `.namespace` directive.

Within the namespace, you can declare your desired symbols. These symbols **must** be preceded with a period, similar to local labels.

Once you have defined the symbols, you can conclude the namespace by using the same label name followed by the `.endnamespace` directive. Any symbol declared within the namespace can be referenced using the namespace's name, followed by a period, and then the symbol/label name.

Example demonstrating the usage of a namespace:

```text
MyVars .namespace
.a = 1
.b = 2
.c = 255
MyVars .endnamespace

;Doing this would load 255 into A
    LDA #MyVars.c
```

Strings or series of bytes can also be used in namespaces:

```text
MyVars2 .namespace
.d = "TEXT"
.e = $01,$02,$03,$04
MyVars2 .endnamespace

;Doing this would output the ASCII codes for the 4 letters that make up "TEXT"
    .BYTE MyVars.d
;Doing this would output the bytes 1,2,3,4
    .BYTE MyVars.e
```

Finally, namespaces can reference anything declared in them locally likeso:

```text
MyVars3 .namespace
.f = 1
.g = .f + 2      ;.g will equal 3
.h = "TEXT"
.i = "MY", .h    ;.i will contain the ASCII bytes of "MYTEXT"
.j = .g,.i       ;.j will contain a 3 followed by the byte representation in .i
MyVars3 .endnamespace
```

---

## Assembler Directives

Directives are assembler-defined operations that are not specific to the 6502 processor but assist in building and defining various aspects of the game.

- A directive is a statement that begins with _at least_ one whitespace character followed by a period ( `.` ) and then the name of the directive. (e.g. `.byte` )
- Directive names are NOT case sensitive. `.byte` and `.BYTE` will have the same effect. However, it is generally recommended to consistently use a specific case throughout the program.

### iNES Header Directives

When building an NES ROM, it's necessary to specify certain values to be written to the iNES header. These values tell emulators about what configuration the game uses.

#### Required iNES Directives

- `.inesMap` - What mapper the game uses (0-255)
- `.inesPrg` - How big PRG ROM is in terms of the number of 16kb pieces it should have. (e.g. a value of 2 = 32kb)

#### Optional iNES Directives

- `.inesChr` - How big CHR ROM is in terms of the number of 8kb pieces it should have (e.g. a value of 4 = 32kb)
  - If using CHR RAM, however, **do NOT** use this directive.
- `.inesMir` - What type of mirroring the game uses (0 for horizontal, 1 for vertical)
- `.inesBat` - Whether the game has battery backup or not (0 for no, 1 for yes)

#### PRG and CHR alternate definitions

Both PRG and CHR ROM can _also_ be declared the following ways:

- Expressing the sizes numerically as the total number of bytes (e.g. `.inesPrg $8000` for setting PRG to 32kb)
- Expressing the sizes as a pre-defined string alias (e.g. `.inesPrg "32kb"` for setting PRG to 32kb)
  - Currently defined string aliases are: `"1kb", "2kb", "4kb"` (keeps doubling up until `"512kb"`) as well as `"1mb"` and `"2mb"`

### RomSegment Directives

This will declare where a segment of ROM will begin and how big it is. Takes 2-3 arguments.

Example 1:

```text
    .romSegment $8000, $8000, "PRG"
```

Example 2:

```text
    .romSegment $10000, $4000, "PRG"
```

- The first argument states how big the total segment size is. In general, this will be equal to the size value of the section of ROM you are building (e.g. if the declared PRG ROM is 32kb, then this should be set to the same amount of 32kb which is $8000 in this case)
- The second argument specifies how big the individual bank size is.
  - In the first example above, it's also set to $8000 which is the same size as the PRG ROM.
  - In the second example, the total ROM segment is declared to be 64kg and will be divided into 16kb banks.
- The third argument is optional. It specifies a user-defined name for the segment that will show up if using the `-s` option to show segment use on the command line.

> Like with the iNES directives, it's also possible to declare rom segment sizes with pre-defined strings such as "32kb", "64kb" etc.

### Bank and Org Directives

- `.bank` - (Required) Specify a bank ID number. The size of which will be whatever was declared in the `.romSegment` . These must be incremental in their placement. (e.g. `.bank 2` cannot precede `.bank 1`)
- `.org` - Specify the origin or the location counter for the bank. Basically, this number will be added to each label in the bank. For instance, if a bank is 8kb in size ($2000 bytes) and the `.org` is set to $8000, then the bank will start at $8000 and end at $9fff

> In general there's no need to supply an org directive more than once per bank with the exception of setting NMI/reset/IRQ vectors

An example of defining a bare minimum, empty, 32kb NROM game might look like:

```text
    .inesPrg "32kb"
    .inesChr "8kb"
    .inesMap 0
    .inesMir 0

    .romSegment "32kb", "32kb", "PRG"
    .bank 0
    .org $8000

    .romSegment "8kb", "8kb", "CHR"
    .bank 0
    .org $0000
```

### File Inclusion Directives

- `.include` - Will include a new input file relative to the current point in assembly. The file path is relative to the input source file.
- `.incBin` - Includes a file as binary data. Will effectively take each byte of the file and convert it into a series of `.byte` directives for the respective byte. It's also possible to add up to 2 arguments at the end:
  - The first argument will be a `seek` value (i.e. how many bytes into the file to start from. Default is 0)
  - The second argument will be a `read` value and represents how many total bytes to read. (default is -1 or all bytes)

### Raw Byte Directives

- `.db` or `.byte` - (Both are synonymous) Define a series of bytes in the current location
- `.dw` or `.word` - (Both are synonymous) Define a series of words (16 bit values) in the current location. Will be stored in the proper little endian order
- `.dwbe` or `.wordBe` - (Both are synonymous) Just like with `.word`, will take a 16 bit value, however endianness will be big endian order instead.
- `.ds` or `.pad` - (Both are synonymous) Define a series the _same_ byte N times. Takes 1 or 2 arguments. The first argument is the number of times. The optional second argument determines what value to repeat. If the second argument is not set, then it will be 0. Effectively a shorthand way of writing `.byte 0,0,0` would be `.pad 3`

#### Mixed Byte Directives

To keep things on one line and save on space, it's also possible to mix bytes, words, and big endian words on a single line using a mixed byte directive.

To use, simply start with: `.d_` followed by a combination of `b`, `w`, or `e` (e is big endian word). Doing so will allow the combination specified and this combination may be repeated however many times one wishes.

Example:

```text
;Allow a byte and then a word
    .d_bw 1,$0203 ;Would evaluate to $01,$03,$02

;Allow 2 bytes and then a big endian word
    .d_bbe 1,2,$0304 ;Would evaluate to $01,$02,$03,$04

;Alternate between bytes and words
    .d_bw 1,$0203,4,$0506 ;Would evaluate to $01,$03,$02,$04,$06,$05
```

It's also possible to have it so the very last type specified can repeat indefinitely. This can be done by placing an underscore at the end. Whatever precedes the underscore will be repeated for however many operands there are at and following that position in the operands.

Example:

```text
;Allow a byte, word, and any number of bytes afterwards
    .d_bwb_ 1,$2000,3,4,5 ;Would evaluate to $01,$00,$20,$03,$04,$05
```

### Repeat Directive

As the name implies the `.repeat` directive will repeat a series of instructions however many times stated. Use an `.endrepeat` directive to terminate the block.

```text
    .repeat 4
    .byte 10
    .endrepeat

;Will do a .byte 10 and then repeat it to 4 total instances
;Effectively becomes:
    .byte 10
    .byte 10
    .byte 10
    .byte 10
```

In addition, it's possible to define an iterator which is a simple variable that can serve as an argument in each loop. The value will start at 0 and count up by 1 per iteration.

```text
    .repeat 4, i
    .byte 10 + i
    .endrepeat

;Will do a .byte 10 and then repeat it 4 additional times plus the value of the iterator.
;i will count from 0 to 3 in this case and stop at 4.
;Effectively becomes:
    .byte 10
    .byte 11
    .byte 12
    .byte 13
```

Repeats can also be nested. If using iterators, then subsequent iterators should use a different letter.

```text
    .repeat 2, i
    .byte 10 + i
        .repeat 2, j
        .byte 20 + j
        .endrepeat
    .endrepeat

;Effectively becomes:
    .byte 10 ;i = 0
    .byte 20 ;j = 0
    .byte 21 ;j = 1
    .byte 11 ;i = 1
    .byte 20 ;j = 0
    .byte 21 ;j = 1
```

### Range Directive

This is another form of repetition that will iterate over an array.  It takes 2 arguments representing the current index as well as the current value of the array element at that index

```text

ArrayOfEvenNumbers = [2,4,6,8]

    .range ArrayOfEvenNumbers, i, v
    .byte i,v
    .endrange

;Will output the index and value as bytes
;Effectively becomes:
    .byte 0, 2
    .byte 1, 4
    .byte 2, 6
    .byte 3, 8
```

### GNSI Directive

- `.gnsi` - "Generate NameSpace Indexes" from labels. There may be cases where data is a set length and referencing points within it by an index can be helpful. Basically this will take all local labels under the parent and calculate the distance between them in bytes and generate a namespace with the exact same local labels as its keys. `.gnsi` can also take a transformative function as a second argument which uses `\1` to represent the current number.

Simple Example:

```text
MyData:
.set0:
 .db 1,2,3
.set1:
 .db 4,5,6

;Will generate a namespace called MyIndexes and will have keys of:
; .set0 = 0
; .set1 = 3
MyIndexes .gnsi MyData
```

Example with function:

```text
MyData:
.set0:
 .db 1,2,3
.set1:
 .db 4,5,6

;With the transformative function, the final values will be divided by 3. This will generate a namespace called MyIndexes and will have keys of:
; .set0 = 0
; .set1 = 1
MyIndexes .gnsi MyData, \1 / 3
```

### Reassign Directive

Using the `.reassign` directive, a previously defined symbol can be changed.  This only works on symbols. Things like charmaps, exprmaps, macros, and namespaces CANNOT be reassigned.

```text
ABC = 1
    .db ABC ;Will be 1
    .reassign ABC, 2
    .db ABC ;Will be 2
```

### Misc Directives

- `.autoZP` or `.autoZeroPage` - (Both are synonymous) This is turned on by default. If an operand is determined to be an 8 bit value and the instruction allows zero page then it will automatically convert the instruction to the zero page version. Realistically, the only reason to turn this off is for some niche cases where one might want to use the full, absolute version of the instruction. Set this to `1` to turn on and `0` to turn off. When turned off, if you still want to use the zero page version of the instruction, the operand must be immediately preceded by a `<` symbol.

- `.emptyRomFill` - After being built, there may be unused space in ROM and something needs to fill it. This directive will specify the desired value to fill each empty byte with. If not set, all empty bytes will be set to `$ff` by default. (Generally best to set this at the beginning of a program)

- `.throw` - Takes a string. Will add a user-thrown error to the error output. This can be handy when using macros with differing arguments.

- `.log` - Takes a string. Will output whatever message you want at the given line.  Output follows error output.

---

## Conditional Directives

It's possible to achieve conditional assembly using either if or switch directives.

### If Statement Directives

Conditional assembly can be achieved using if statements. Follow these guidelines to work with them:

- Begin a conditional block with an `.if` directive.
- Use the `.elseif` directive to check an additional condition if the previous condition fails.
- Optionally, provide a default failure case using the `.else` directive.
- Finally, close the conditional block with an `.endif` directive.
- Note: It is possible to nest if statements as well

```text
;Would generate LDA #2
SYMBOL = 2

    .if SYMBOL == 1
        LDA #1
    .elseif SYMBOL == 2
        LDA #2
    .else
        LDA #3
    .endif
```

### Switch Statement Directives

Another way to achieve conditional assembly is with switch statements. Unlike if statements, these are more for checking and reacting to symbols that are a specific value. Follow these guidelines to work with them:

- Begin a conditional block with an `.switch` directive followed by the symbol you want to check.
- Use the `.case` to check if the symbol matches a specific value. Place your conditional assembly underneath each `.case` directive.
- Optionally, provide a default case using the `.default` directive for handling unmatched values. Like with `.case` , place your conditonial assembly underneath the `.default` directive.
- Finally, close the conditional block with an `.endswitch` directive.

```text
;Would generate LDA #3
SYMBOL = 99

    .switch SYMBOL
      .case 1
        LDA #1
      .case 2
        LDA #2
      .default
        LDA #3
    .endswitch
```

---

## Assembler Functions

There exist several functions built into the assembler.

### Number and Math Functions

- `low(NUMBER)` - Will take a value and return the low byte. (e.g. `low($1122)` would return $22).
- `high(NUMBER)` - Will take a value and return the high byte. (e.g. `high($1122)` would return $11).
  - If the value is 8 bit, this will simply return 0
- `round(NUMBER)` - Will always round a number to the nearest int
- `floor(NUMBER)` - Will always round a number down regardless of the decimal point value to an int
- `ceil(NUMBER)` - Will always round a number up regardless of the decimal point value to an int
- `modfDeci(NUMBER)` - Will take the decimal point value of a number and discard the int value
- `modfInt(NUMBER)` - Will take the int value of a number and discard the decimal point value
- `sin(NUMBER)` - Sine of the number (in radians)
- `sindeg(NUMBER)` - Sine of the number (in degrees)
- `cos(NUMBER)` - Cosine of the number (in radians)
- `cosdeg(NUMBER)` - Cosine of the number (in degrees)
- `tan(NUMBER)` - Tangent of the number (in radians)
- `tandeg(NUMBER)` - Tangent of the number (in degrees)

### String functions

- `reverseStr(STRING)` - Reverses the characters in a specified string
- `strLen(STRING)` - Returns the length of the specified string
- `subStr(STRING, START, END)` - Returns part of a string. START is the starting position in the string (value must be 0 or higher and less than the string length). If only START is specified, then it will take a single character from the string. If END is also specified, then it must be a value greater than START and less than the length of the string and will return however many characters are in the range.
- `toCharmap(STRING)` - This will replace each individual character in a string from its ASCII value to the value defined in a character map. It will be returned as an array of numbers.

### Array functions

- `arrayByteLen(ARRAY)` - If the array (including all contained arrays) were to be fully flattened down to individual bytes, this will count how many bytes it would need
- `arrayContains(ARRAY, value)` - Whether the array has an instance of the value. (currently only numbers are supported)
- `arrayIndex(ARRAY, value)` - The index the target value is at. (currently only numbers are supported)
- `arrayLen(ARRAY)` - How many elements are in the array
- `subArray(ARRAY, START, END)` - Returns part of an array as a new array. START is the starting index in the array. If END is provided, then it will be the ending index.  If omitted, it will just take from the start index to very end.

### Symbol-Related Functions

- `bank(LABEL)` - Will take a label and return what bank the label is located in. (useful for programs that utilize some form of bank switching)
- `defined(SYMBOL)` - Determine whether a symbol has been defined at the point of execution.

### Namespace functions

- `namespaceValuesToArray(NAMESPACE)` - Returns all the values of a namespace (in order of declaration) as one array of values.

### Bytes Under Label Counting Functions

- `bytesInCurrentLabel()` - Will count the total bytes in the current label up until the next label
- `bytesInCurrentLocal()` - Will count the total bytes in the current local label up until the next local label (or next parent label if none)
- `bytesInLabel(LABEL)` - Will count the total bytes in the provided label up until whatever label follows it

These can be handy for cases of having some data of N length in ROM and there is some counter of bytes preceding it.

Examples of use:

```text
DemoData:
.data1:
  .db bytesInCurrentLocal() ;Would evaluate to 5.
  .db $f8,$20,$00,$f0

.data2:
  .db bytesInCurrentLocal() ;Would evaluate to 9
  .db $f8,$30,$00,$00
  .db $f8,$38,$00,$08

; Assumes some other label would follow this...

```

### Examples of Using Functions

```text
SOME_LABEL:
    LDA #low(SOME_LABEL)   ;Load in the low byte value of the label SOME_LABEL into A
    LDX #high(SOME_LABEL)  ;Load in the high byte value of the label SOME_LABEL into X
    LDY #bank(SOME_LABEL)  ;Load in the bank that the label SOME_LABEL is located in into Y
```

---

## User-Defined Functions

It's possible to write your own functions to generate specific values. User-defined functions must have a label (with no whitespace before the label), followed by the `.func` directive, followed by the desired expression AND be contained within a single line.

Functions can take arguments starting from an ID of 1 and are preceded with a backslash ( `\` ) .

> If no arguments are supplied in a function, a warning will be displayed as it's basically pointless to have a function without arguments. A constant would make more sense in such cases.

```text
Add1Function .func \1 + 1 ;Would add 1 to the argument supplied

;Using the function:
    .BYTE Add1Function(2), Add1Function(3) ;Would evaluate to: .BYTE 3, 4
```

---

## Macros

Macros provide a user-defined way to reuse blocks of directives or instructions.

### Defining and using a simple macro

To create a simple macro, it must have a label followed by `.macro` directive. After this statement, insert your desired instructions. To terminate the macro, use the same label and then `.endm` directive.  As with other labeled operations, there can be no whitespace before the label

```text
MyMacro .macro
    LDA #$3F
    STA $2006
    LDA #$00
    STA $2006
MyMacro .endm
```

Once a simple macro is defined, it can be inserted into the code as if it were a regular instruction or directive and invoked accordingly.

```text
;Invoking MyMacro in code
    MyMacro  ;The instructions defined in MyMarco (the example above) would be inserted here.
```

### Using arguments in macros

It's also possible to give names to the arguments of a macro by adding a list after the `.macro` directive.  Each argument must be preceded by a backslash in this case.

Example of use:

```text
;Macro with 2 named arguments
MyMacroWithNamedArgs .macro \var1, \var2
    .db \var1, \var2
MyMacroWithNamedArgs .endm

;Calling the macro
    MyMacroWithNamedArgs 2, 3 ;Would generate .db 2, 3

```

### Using anonymous arguments in macros

It's also possible to omit the list of arguments. This allows for any amount of arguments to be supplied and to use them, a backslash followed by a number (starting at `\1` for the first argument) will signify the argument.

Example of use:

```text
MyMacroWithAnonymouseArgs .macro
    LDA \1
    STA $2006
    LDA \2
    STA $2006
MyMacroWithAnonymouseArgs .endm

;Using MyMacroWithArgs in code
    MyMacroWithAnonymouseArgs #$3F, #$00

;The instructions defined in MyMarcoWithArgs will be inserted here and #$3F will replace any instances of \1 and #$00 will replace any instances of \2
```

#### Checking supplied argument number

It's also possible to check the number of arguments that have been supplied to a macro with the `\#` operation.

```text
;Check number of arguments
MyMacroWithVariableArgs .macro
 .if \# == 1
    LDA \1
    STA $2006
 .elseif \# == 2
    LDA \1
    STA $2006
    LDA \2
    STA $2006
 .endif

MyMacroWithVariableArgs .endm

;This maco will do different things depending on how many arguments are supplied to it.
```

### Exiting a Macro Early

It's possible to stop a macro using the `.exitMacro` directive.

```text
ExitMacroDemo .macro \arg1
    .if \arg1 > 2
    .exitMacro ;Will exit here and do nothing if the first argument is greater than 2
    .endif

;Macro execution continues otherwise...

ExitMacroDemo .endm

```

---

## Character Maps

Character maps are used for replacing ASCII with specific number values. This is useful for when you want to map text to specific tiles, but don't necessarily want to use the actual ASCII values for the tiles in-game.

- To define a character map, simply use a unique label followed by a `.charmap` directive.
- To end a character map, close it using the same label followed by `.endcharmap`

> Note: The very first charmap defined will also be the default charmap.

### Directives Used Within Character Maps

These are used within the character map definition

- `.defchar` - When inside a `.charmap` operation, put a single character in quotes for the desired ASCII character to map followed by the desired value to map it to. Values can be numbers, arrays, or (not recommended) strings.
  - Arrays can be useful for text engines that might use multiple bytes for characters
- `.defcharRange` - When inside a `.charmap` operation, put a single character in quotes for the desired ASCII character to map followed by the desired value to map it to. Unlike `defChar`, the value must be a numeric.

### Other Character Map Directives

- `.setCharMap` - Will set the current character map to whatever operand is specified.
- `.resetCharMap` - Set the current character map to the default character map (i.e. whichever one was first-defined in code)

Examples of Defining:

```text
MyCharmap .charmap
    .defchar 'A', $01
    .defchar 'B', $02
    .defChar "Ä", [$80,$01]
MyCharmap .endcharmap

    .db "AB" ;Would evaluate to $41,$42 (i.e. the actual ASCII values)
    .db toCharmap("AB") ;Would evaluate to $01,$02
    .db toCharmap("Ä") ;Would evaluate to $80,$01
```

It's also possible to use a `.defcharRange` directive to group characters. Simply set the starting and ending character and the desired value to start from which will increment by 1.

```text
MyCharmap2 .charmap
    .defcharRange 'a', 'z', $2a
MyCharmap2 .endcharmap

    .db toCharmap("abc") ;Would evaluate to $2a,$2b,$2c
```

Examples of Using:

```text
;If MyCharmap from earlier is the default character map, then this would be $1a,$1b
    .db toCharmap("AB")
```

---

## Expression Maps

Similar to a charmap, but this will replace an entire string value within backticks (e.g. `` `value` ``) with a single value. Unlike contstants, these are strings and they can be mapped to any text value. (e.g. `` `c#` `` might represent C sharp in data for a music engine)

- To define an expression map, simply use be a unique label followed by an `.exprmap` directive.
- Within the block, each character must be defined using a `.defexpr` directive followed by the string and then the desired value to change that string to. For the string, **DO NOT** use backticks but either single or double quotes.
- Close using the same label followed by `.endexprmap`

> Note: The very first exprmap defined in game will also be the default exprmap.

### Other Expression Map Directives

- `.setExprMap` - Will set the current expression map to whatever operand is specified.
- `.resetExprMap` - Set the current expression map to the default expression map (i.e. whichever one was first-defined in code)

### Examples of Using Expression Maps

```text
MyExprmap .exprmap
    .defchar 'Abb', $1a
    .defchar 'Cdd', $1b
MyExprmap .endexprmap

MyExprmap2 .exprmap
    .defchar 'Abb', $2a
    .defchar 'Cdd', $2b
MyExprmap2 .endexprmap

    .db `Abb` ;Would evaluate to $1a
    .db `Cdd` ;Would evaluate to $1b

    .setExprMap MyExprmap2
   .db `Abb` ;Would evaluate to $2a

    .resetExprMap
   .db `Abb` ;Would once again evaluate to $1a
```

---

## Dynamic Labels and Strings

There may be cases where multiple labels are desired (typically in repeat statements or macros) where providing a unique name every time can be cumbersome.  Strings can also be dynamically generated.

Dynamic labels can be created by using a lowercase `l` immediately followed by an opening double quote. Strings can be created by using a lowercase `s` immediately followed by an opening double quote. Within the double quotes, curly braces can be used to evaluate expressions that will become part of the generated output.

Local labels can also be dynamic.  The period **must** be inside the quotes.

### Examples of Using Dynamic Labels

```text
abc = 1

;This label becomes Dynamic1
l"Dynamic{abc}":

;This label becomes Dynamic2
l"Dynamic{abc+1}":

;Creates 10 local labels called .repeatDynamic0, .repeatDynamic1, ..., .repeatDynamic9
    .repeat 10, i
l".repeatDynamic{i}":
    .endRepeat

;Using the 10 dynamic labels in a .dw statement
    .repeat 10, i
    .dw l".repeatDyanmic{i}"
    .endRepeat
```

;Create a dynamic string
DynamicString = s"Some text with {abc}" ;Becomes "Some text with 1"

---

## **Happy Assembly!**

EOF!
