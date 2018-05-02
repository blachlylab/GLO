package glo

import (
    "fmt"
    "strings"
    "os"
    "bufio"
)

// Represents a contig:start-end structure.
type ChainInterval struct {
    Contig string
    Start int64
    End int64
}

// String output function for ChainInterval type
func (ci *ChainInterval) String() string {
    return fmt.Sprintf("%s:%d-%d", ci.Contig, ci.Start, ci.End)
}

// Represents the source -> target mapping between
// two ChainIntervals, e.g.
// chrA:10000-20000 -> chrB:20123-30123
type ChainBlock struct {
    Source *ChainInterval
    Target *ChainInterval
}

// String output fuction for ChainBlock type
func (cb *ChainBlock) String() string {
    return fmt.Sprintf("%s -> %s", cb.Source, cb.Target)
}


// GetOverlap returns a ChainInterval object representing the
// overlapped with the input region.
func (cb *ChainBlock) GetOverlap(region *ChainInterval) *ChainInterval { 
    ci := new(ChainInterval)
    if region.Contig != cb.Source.Contig {
        // No overlap due to contig mismatch
        return nil
    }

    ci.Contig = cb.Target.Contig

    var start_adj int64 = 0

    if region.Start > cb.Source.Start {
        start_adj = region.Start - cb.Source.Start
    }
    ci.Start = cb.Target.Start + start_adj


    size := region.End - region.Start
    if region.End > cb.Source.End {
        size -= (region.End - cb.Source.End)
    }
    ci.End = ci.Start + size

    return ci
}


// The Chain type represents a UCSC chain object, including all
// the fields from the header line and each block of mappings
// for that chain.
type Chain struct {
    Score int64
    SourceName string
    SourceSize int64
    SourceStrand string
    SourceStart int64
    SourceEnd int64
    TargetName string
    TargetSize int64  
    TargetStrand string
    TargetStart int64
    TargetEnd int64
    ID string
    Blocks []*ChainBlock
}

// String output function for Chain type.
func (c *Chain) String() string {
    var output []string
    output = append(output, 
        fmt.Sprintf("%s:%d-%d to %s:%d-%d", c.SourceName, c.SourceStart, 
            c.SourceEnd, c.TargetName, c.TargetStart, c.TargetEnd))
    for _, x := range c.Blocks {
        output = append(output, fmt.Sprintf("> %s", x))
    }
    return strings.Join(output, "\n")
}


// Populates the target Chain struct from the data in the input string.
func (c *Chain) FromString(s string) {
    //chain 20851231461 chr1 249250621 + 10000 249240621 chr1 248956422 + 10000 248946422 2
    cols := strings.Split(s, " ")
    c.Score = str2int64(cols[1])
    c.SourceName = strings.ToLower(cols[2])
    c.SourceSize = str2int64(cols[3])
    c.SourceStrand = cols[4]
    c.SourceStart = str2int64(cols[5])
    c.SourceEnd = str2int64(cols[6])
    c.TargetName = strings.ToLower(cols[7])
    c.TargetSize = str2int64(cols[8])
    c.TargetStrand = cols[9]
    c.TargetStart = str2int64(cols[10])
    c.TargetEnd = str2int64(cols[11])
    if len(cols) == 13 {
        c.ID = cols[12]
    }
}

// Chain.load_blocks uses the input bufio.Reader to read mapping
// block from the file, until the end of the chain is found. The
// mapping blocks are added to the Chain as ChainBlock instances.
func (c *Chain) load_blocks(reader *bufio.Reader){
    var cols []string
    var block *ChainBlock
    var si, ti *ChainInterval

    line, err := reader.ReadString('\n')
    if err != nil {
        fmt.Printf("Chain.load_blocks() error: %s\n", err)
        os.Exit(1)
    }
    cols = strings.Split(strings.TrimSpace(string(line)), "\t")

    sfrom := c.SourceStart
    tfrom := c.TargetStart
    for len(cols) == 3 {
        size := str2int64(cols[0])
        sgap := str2int64(cols[1])
        tgap := str2int64(cols[2])

        block = new(ChainBlock)

        si = new(ChainInterval)
        si.Contig = c.SourceName
        si.Start = sfrom
        si.End = sfrom + size
        block.Source = si

        ti = new(ChainInterval)
        ti.Contig = c.TargetName
        ti.Start = tfrom
        ti.End = tfrom + size
        block.Target = ti

        sfrom += size + sgap
        tfrom += size + tgap

        c.Blocks = append(c.Blocks, block)
     
        _, p_err := reader.Peek(1)
        if p_err != nil {
            // EOF
            break
        }
        line, err = reader.ReadString('\n')
        if err != nil {
            fmt.Printf("Chain.load_blocks() error: %s\n", err)
            os.Exit(1)           
        }
        cols = strings.Split(strings.TrimSpace(string(line)), "\t")
    }

    if len(cols) != 1 {
        fmt.Printf("Error: Expected line with a single value, got \"%s\"\n", cols)
        os.Exit(1)
    }

    size := str2int64(cols[0])
    block = new(ChainBlock)
        
    si = new(ChainInterval)
    si.Contig = c.SourceName
    si.Start = sfrom
    si.End = sfrom + size
    block.Source = si

    ti = new(ChainInterval)
    ti.Contig = c.TargetName
    ti.Start = tfrom
    ti.End = tfrom + size
    block.Target = ti

    c.Blocks = append(c.Blocks, block)

}

