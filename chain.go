package goLiftOver

import (
    "fmt"
    "strings"
    "os"
    "bufio"
)

// Represents a contig:start-end structure.
type ChainInterval struct {
    contig string
    start int64
    end int64
}

// String output function for ChainInterval type
func (ci *ChainInterval) String() string {
    return fmt.Sprintf("%s:%d-%d", ci.contig, ci.start, ci.end)
}

// Represents the source -> target mapping between
// two ChainIntervals, e.g.
// chrA:10000-20000 -> chrB:20123-30123
type ChainBlock struct {
    source *ChainInterval
    target *ChainInterval
}

// String output fuction for ChainBlock type
func (cb *ChainBlock) String() string {
    return fmt.Sprintf("%s -> %s", cb.source, cb.target)
}


// GetOverlap returns a ChainInterval object representing the
// overlapped interval at the target contig
func (cb *ChainBlock) GetOverlap(contig string, start, end int64) *ChainInterval {
    ci := new(ChainInterval)
    if contig != cb.source.contig {
        // No overlap due to contig mismatch
        return nil
    }

    ci.contig = cb.target.contig


    var start_adj int64 = 0

    if start > cb.source.start {
        start_adj = start - cb.source.start
    }
    ci.start = cb.target.start + start_adj


    size := end - start
    if end > cb.source.end {
        size -= (end - cb.source.end)
    }
    ci.end = ci.start + size

    return ci
}


// The Chain type represents a UCSC chain object, including all
// the fields from the header line and each block of mappings
// for that chain.
type Chain struct {
    score int64
    source_name string
    source_size int64
    source_strand string
    source_start int64
    source_end int64
    target_name string
    target_size int64  
    target_strand string
    target_start int64
    target_end int64
    id string
    blocks []*ChainBlock
}

// String output function for Chain type.
func (c *Chain) String() string {
    var output []string
    output = append(output, fmt.Sprintf("%s:%d-%d to %s:%d-%d", c.source_name, c.source_start, c.source_end, c.target_name, c.target_start, c.target_end))
    for _, x := range c.blocks {
        output = append(output, fmt.Sprintf("> %s", x))
    }
    return strings.Join(output, "\n")
}


// Populates the target Chain struct from the data in the input string.
func (c *Chain) fromString(s string) {
    //chain 20851231461 chr1 249250621 + 10000 249240621 chr1 248956422 + 10000 248946422 2
    cols := strings.Split(s, " ")
    c.score = str2int64(cols[1])
    c.source_name = cols[2]
    c.source_size = str2int64(cols[3])
    c.source_strand = cols[4]
    c.source_start = str2int64(cols[5])
    c.source_end = str2int64(cols[6])
    c.target_name = cols[7]
    c.target_size = str2int64(cols[8])
    c.target_strand = cols[9]
    c.target_start = str2int64(cols[10])
    c.target_end = str2int64(cols[11])
    if len(cols) == 13 {
        c.id = cols[12]
    }
}

// Loads the data blocks (stored as ChainBlocks) for the target Chain
// object, using the input Scanner object.

func (c *Chain) load_blocks(s *bufio.Scanner){
    var cols []string
    var block *ChainBlock
    var si, ti *ChainInterval

    cols = strings.Split(strings.TrimSpace(s.Text()), "\t")

    sfrom := c.source_start
    tfrom := c.target_start
    for len(cols) == 3 {
        size := str2int64(cols[0])
        sgap := str2int64(cols[1])
        tgap := str2int64(cols[2])

        block = new(ChainBlock)

        si = new(ChainInterval)
        si.contig = c.source_name
        si.start = sfrom
        si.end = sfrom + size
        block.source = si

        ti = new(ChainInterval)
        ti.contig = c.target_name
        ti.start = tfrom
        ti.end = tfrom + size
        block.target = ti

        sfrom += size + sgap
        tfrom += size + tgap

        c.blocks = append(c.blocks, block)
        //fmt.Printf(">[%d]\t%s\n", len(cols), cols)
        if !s.Scan() {
            //fmt.Println("break")
            break
        }
        cols = strings.Split(strings.TrimSpace(s.Text()), "\t")
        //fmt.Printf("number of cols read in: %d\n", len(cols))
    }
    if len(cols) != 1 {
        fmt.Printf("Error: Expected line with a single value, got \"%s\"\n", cols)
        os.Exit(1)
    }

    size := str2int64(cols[0])
    block = new(ChainBlock)
        
    si = new(ChainInterval)
    si.contig = c.source_name
    si.start = sfrom
    si.end = sfrom + size
    block.source = si

    ti = new(ChainInterval)
    ti.contig = c.target_name
    ti.start = tfrom
    ti.end = tfrom + size
    block.target = ti

    c.blocks = append(c.blocks, block)
}