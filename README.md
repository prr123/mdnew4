# github markdown parser mdnew4

why? goldmark lacks documentation 
gomarkdown does not have extendible inline parsers

parser does two parses

1. blocks
2. each block is parsed for inline chars (inline parsers)

- creates an ast tree
- uses block parsers for different blocks
- all parsers are extendable
- block parsers are stored in a go map triggered by the first letter of a line
- inline parsers are stored in the inline map

changes from mdnew2:

## mdnew3 uses a three step parsing approach

1. parse each line to create an array of blocks
2. then parse each block for children
3. parse each node for inline changes (emphasis, footnotes)

## mdnew4

- changed line parsing algorithm
- modified parse routine to correctly parse blocks

