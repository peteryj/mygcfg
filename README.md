mygcfg
======

My Go Config File parser, based on .ini format, for example
[SectionName1]
name = value
[SectionName2]
name2 = value2

This software just parse the give file, and store them into a map.
Currently it is only a prototype, and easy to crash for misformatted file.

TODO
======
1. use LR(1) parser to rewrite the syntax parsing
2. refactoring
