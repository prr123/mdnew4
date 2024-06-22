// parser for markdown files
package mdparser


import (
	"fmt"
)


type MdNode struct {
	ch []*MdNode
	par *MdNode
	typ int
	el string
	blkSt int
	blkEnd int
	txtSt int
	txt []byte
//	att *Attribute
//	prop interface
}

type Block struct {
	linSt int
	linEnd int
	Typ rune
	closed bool
}

type MdParser struct {
	blkMap map[byte]func(ps *MdPState) *MdNode
	charList []byte
//	blkTyp map[string]int
	BlockList []Block
	lines []RLine
	max int
}

type MdPState struct {
	Doc *MdNode //top
	Node *MdNode // current parent node
	Blk *MdNode // prev block
	closed bool
	prevBlk int
	plin int
	nest int
	state int
}

type RLine struct {
	linSt int
	linEnd int
	lintxt []byte
	spCount int
	tbCount int
	txtSt int
	eolChar int
}


func IsAlpha(let byte)(res bool) {
    res = false
    if (let >= 'a' && let <= 'z') || (let >= 'A' && let <= 'Z') { res = true}
    return res
}

func CleanRet (inp *[]byte) {

	ptr := 0
	for i:=0; i< len(*inp); i++ {
		if (*inp)[i] == '\r' {
			if (*inp)[i+1] != '\n' {
				(*inp)[ptr] = '\n'
				ptr++
			}
		} else {
			(*inp)[ptr] = (*inp)[i]
			ptr++
		}
	}
}


func InitParser(inp []byte) (p MdParser) {

	p.blkMap = make(map[byte]func(ps *MdPState) *MdNode)
//	p.blkTyp = make(map[string]int)

	p.blkMap['#'] = p.ParseHeading
	p.blkMap['p'] = p.ParsePar
	p.blkMap[' '] = p.ParseEL
	p.blkMap['`'] = p.ParseCode
	p.blkMap['-'] = p.ParseUL
	p.blkMap['+'] = p.ParseUL
	p.blkMap['*'] = p.ParseUL
	p.blkMap['n'] = p.ParseOL
	p.blkMap['>'] = p.ParseQuote

	p.charList = []byte{' ','\t','#','+','*','>'}

	p.lines = GetLines(inp)
	p.max = len(p.lines)
	return p
}

func (p MdParser) Lines() (ls []RLine) {
	return p.lines
}

func (p MdParser) LinInfo(l RLine) (start int, end int) {
	start = l.linSt
	end = l.linEnd
	return start, end
}

func InitParseState(inp []byte) (pstate *MdPState) {

	var mdDoc MdNode
	mdDoc.blkSt = 0
	mdDoc.blkEnd = len(inp)
//	mdDoc.typ = Doc
	mdDoc.ch = nil
	mdDoc.par = nil

	var ps MdPState
	ps.Doc = &mdDoc
	ps.Blk = nil
	ps.Node = &mdDoc
	ps.closed = true
	ps.nest = 0
//	ps.state = NB
	return &ps
}

func (p *MdParser)ParseCode(ps *MdPState) *MdNode {
	fmt.Printf("parsing code\n")
	return nil
}

func  (p *MdParser)ParseQuote(ps *MdPState) *MdNode {
	fmt.Printf("parsing code\n")
	return nil
}

/*
type MdNode struct {
	ch []*MdNode
	par *MdNode
	typ string
	blkSt int
	blkEnd int
	txtSt int
	txt []byte
//	att *Attribute
//	prop interface
}
*/

func  (p *MdParser)ParseUL(ps *MdPState) *MdNode {
	fmt.Printf("parsing UL: %d nest: %d\n", ps.state, ps.nest)

	l := p.lines[ps.plin]

//fmt.Printf("line nest: %d\n", l.nest)
	blk := &MdNode{}

	nest :=0
	if l.spCount > 0 {
		nest = l.spCount/4
	} else {
		nest = l.tbCount
	}
	// check whether there is a UL element
	if !ps.closed {
		blk.el = "ul"
		blk.par = ps.Node
		blk.blkSt= l.linSt
		blk.blkEnd = -1
		ps.Blk = blk
//		ps.state = UL
	} else {
		if nest > ps.nest {
			blk.el = "ul"
			blk.par = ps.Blk
			blk.blkSt= l.linSt
			blk.blkEnd = -1
			ps.Blk.ch = append(ps.Blk.ch, blk)
			ps.Blk = blk
			ps.nest++
			PrintNode(blk, "nest")
		}
		if nest == ps.nest {
//			ps.Blk = blk.par
			blk = ps.Blk
//			ps.nest--
//			PrintNode(blk, "reversion")
		}
		if nest < ps.nest {
			ps.Blk = ps.Blk.par
			blk = ps.Blk
			ps.nest--
			PrintNode(blk, "reversion")
		}
	}

	liblk := &MdNode{
			el: "li",
			par: blk,
			blkSt: l.linSt,
			blkEnd: -1,
		}

	state:=0
	loop := true
	for i:=1; i<len(l.lintxt); i++ {
		let := l.lintxt[i]
		switch state {
		case 0:
			if let == ' ' {state = 1}
		case 1:
			if let == ' ' {break}
			if let != ' ' {
				state = 2
				liblk.txtSt = i
				liblk.txt = l.lintxt[i:]
				break
			}
		case 2:
			loop = false
		default:
			return nil
		}
		if !loop {break}
	}
	blk.ch = append(blk.ch,liblk)
	ps.closed = false
	return blk
}

func (p *MdParser)ParseOL(ps *MdPState) *MdNode {
	fmt.Printf("parsing OL\n")

	return nil
}

func (p *MdParser)ParseHeading(ps *MdPState) *MdNode{
	l := p.lines[ps.plin]
	fmt.Printf("parsing heading: %s\n", string(l.lintxt))

	hdlev :=0
	state:=0
	txtst:=-1

	fin := false
	for i:=0; i<len(l.lintxt); i++ {

		let := l.lintxt[i]
		switch state {
		case 0:
			if let == '#' {hdlev++}
			if let == ' ' {state = 1}
		case 1:
			if let != ' ' {
				state = 2
				txtst = i
			}
		case 2:
			fin = true
		default:
		}
		if fin {break}
	}

	head := fmt.Sprintf("h%d",hdlev)
	txtSt := l.linSt + txtst
	blk := MdNode{
//		: HD,
		el: head,
		par: ps.Node,
		blkSt: l.linSt,
		blkEnd: l.linEnd,
		txtSt: txtSt,
		txt: l.lintxt[txtst:],
	}
//	ps.state = 
	return &blk
}

func (p *MdParser)ParseEL(ps *MdPState) *MdNode{
	fmt.Println("parsing empty line")
	l := p.lines[ps.plin]
    blk := MdNode{
        el: "br",
        par: ps.Node,
        blkSt: l.linSt,
        blkEnd: l.linEnd,
    }

	for i:=0; i<len(l.lintxt); i++ {
		let := l.lintxt[i]
		if let != ' ' {
//fmt.Printf("not a empty line: %q\n",let)
			return nil
		}
	}

	ps.closed = true
	return &blk
}

func (p *MdParser)ParsePar(ps *MdPState) *MdNode{
	fmt.Println("parsing paragraph")
//fmt.Printf("ps.close: %t\n%v\n",ps.closed, ps)

	l := p.lines[ps.plin]
	eoBlk:= false
	if l.lintxt[len(l.lintxt)-1]== ' ' && l.lintxt[len(l.lintxt)-2] == ' ' {
//fmt.Println("end of par 2ws")
		l.lintxt  =  l.lintxt[:len(l.lintxt)-2]
		eoBlk = true
	}

	blk := &MdNode{}
	if ps.closed {
		blk = &MdNode{
				el: "p",
				par: ps.Node,
				blkSt: l.linSt,
				blkEnd: l.linEnd,
				txtSt: l.linSt,
				txt: l.lintxt,
			}
	} else {
		blk = ps.Blk
		blk.blkEnd = l.linEnd
		blk.txt = append(ps.Blk.txt, ' ')
		blk.txt = append(ps.Blk.txt, l.lintxt...)
	}

//fmt.Printf("par p return: %v\np.Blk:%v\n", p, p.Blk)
	ps.closed = false
	if eoBlk {ps.closed = true}
//	ps.state = Base
	return blk
}

func CloseBlk(ps *MdPState) {
	fmt.Println("closing block")

	if ps.Blk != nil {
		ps.Node.ch = append(ps.Node.ch, ps.Blk)
		ps.Blk = nil
	}
	return
}

func GetLines (inp []byte) (linList []RLine){

	linSt:=0
	linList = make([]RLine,0,128)

	for i:=0; i< len(inp); i++ {
		if inp[i] != '\n' { continue}

		txtst := linSt
		newLine := RLine {
				linSt: linSt,
				linEnd: i,
				lintxt: inp[linSt:i],
				eolChar: 0,
			}
		if linSt == i  {newLine.eolChar = 1}
		if i-linSt >2 {
			if inp[i-2] == ' ' && inp[i-1] == ' ' {newLine.eolChar = 2}
		}

//			ind := linSt
			spCount :=0
			tbCount :=0
			for j:=linSt; j<i-1; j++ {
				state :=0
				switch state {
				case 0:
					if inp[j] == ' ' {
						state =1
						spCount++
						break
					}
					if inp[j] == '\t' {
						state = 2
						tbCount++
						break
					}
					txtst = j
					state = 3
				case 1:
					if inp[j] == ' ' {
						spCount++
					} else {
						txtst = j
						state = 3
					}
				case 2:
					if inp[j] == '\t' {
						tbCount++
					} else {
						txtst = j
						state = 3
					}
				default:
				}
				if state == 3 {break}
			}
		newLine.spCount = spCount
		newLine.tbCount = tbCount
		newLine.txtSt = txtst
		linList = append(linList,newLine)
		linSt = i+1
	}
	return linList
}


func (p *MdParser)Parse (ps *MdPState) (err error){

	linList := p.lines
//	res:=&MdNode{}

	linNum := len(linList)
	prevBlk := -1
	oldBlkTyp := ' '
	newBlkTyp := 'N'
	for i:=0; i< linNum; i++ {
//fmt.Printf("\n**** line: %d\n")
		line := linList[i]
		if line.linSt == line.linEnd {
			if !p.BlockList[prevBlk].closed {
				p.BlockList[prevBlk].closed = true
			}

			b := Block{
				Typ: 'E',
				linSt: i,
				linEnd: i,
				closed: true,
			}
			p.BlockList = append(p.BlockList, b)
			prevBlk = len(p.BlockList) -1
			oldBlkTyp = 'E'
			continue
		}
		linLen := line.linEnd-line.linSt
		flet := line.lintxt[0]
		switch flet {
		case '#':
			if linLen < 3 {newBlkTyp = 'R'; break;}
			newBlkTyp = '#'

		case '-', '*', '+':
			if linLen < 3 {newBlkTyp = 'R'; break;}
			if line.lintxt[1] == ' ' {newBlkTyp = 'U'; break;}
			if line.lintxt[1] == flet && line.lintxt[2] == flet {newBlkTyp = 'H'; break;}
			newBlkTyp = 'P'

		case '1','2','3','4','5','6','7','8','9':
			newBlkTyp = 'O'

		case '>':
			newBlkTyp = 'Q'

		case ' ':
			if linLen < 5 {newBlkTyp = 'R'; break;}

		case '\t':
			if linLen < 2 {newBlkTyp = 'R'; break;}

		case '`':
			if linLen < 3 {newBlkTyp = 'R'; break;}
			for k:=1; k< 3; k++ {
				if line.lintxt[k] != '`' {newBlkTyp = 'R'; break;}
			}
			newBlkTyp = 'C'

		default:
			switch oldBlkTyp {
			case 'Q': newBlkTyp = 'Q'
			case 'U': newBlkTyp = 'U'
			case 'O': newBlkTyp = 'O'
			default: newBlkTyp = 'P'
			}
		}

		if oldBlkTyp == newBlkTyp && !p.BlockList[prevBlk].closed{
			p.BlockList[prevBlk].linEnd=i
			if line.eolChar == 2 {
				p.BlockList[prevBlk].closed = true
			}
		} else {
			b := Block{
				Typ: newBlkTyp,
				linSt: i,
				linEnd: i,
				closed: false,
			}
			p.BlockList = append(p.BlockList, b)
			prevBlk = len(p.BlockList) -1
			oldBlkTyp = newBlkTyp
		}
	}
	return nil
}

func PrintBlock(bl []Block) {
	fmt.Println("************ Block List ***************")
	for i:=0; i<len(bl); i++ {
		b := bl[i]
		fmt.Printf("  [%4d]: %q<%d,%d>\n", i +1, b.Typ, b.linSt + 1, b.linEnd +1)
	}
	fmt.Println("********** End Block List *************")

}

func PrintNode(n *MdNode, title string) {

	fmt.Printf("\n******** Node %s ***********\n", title)
	if n == nil {
		fmt.Println("no node")
		fmt.Printf("****** End Node %s *********\n\n", title)
		return
	}
	fmt.Printf("Typ: %d\n", n.typ)
 	fmt.Printf("st: %d end: %d\n", n.blkSt, n.blkEnd)
	fmt.Printf("children: %d\n", len(n.ch))
	if n.par == nil {
		fmt.Printf("parent: none\n")
	} else {
		fmt.Printf("parent: %d\n", n.par.typ)
	}
	fmt.Printf("txt: %s\n", n.txt)

//	if par == nil {return}
	fmt.Printf("Children [%d]\n", len(n.ch))
	if len(n.ch) == 0 {
		fmt.Printf("****** End Node %s *********\n\n", title)
		return
	}
	for i:= 0; i< len(n.ch); i++ {
		cNode := n.ch[i]
		str := fmt.Sprintf("child: %d", i +1)
//fmt.Printf("** %s **\n", str)
		PrintNode(cNode, str)
	}

	fmt.Printf("****** End Node %s *********\n\n", title)

}

func (p *MdParser)PrintLines() {

	lines := p.lines
	fmt.Println("******* Lines *******")
	for i:=0; i<len(lines); i++ {
		l :=lines[i]
		fmt.Printf("--[%3d]: (%4d %4d %4d %2d %2d %1d) %s\n",i+1, l.linSt, l.linEnd, l.txtSt, l.spCount, l.tbCount, l.eolChar, string(l.lintxt))
	}
	fmt.Println("***** End Lines *****")

}
