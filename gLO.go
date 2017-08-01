package main

import (
    "fmt"
    "os"
    //"os/exec"
    "bufio"
    "strings"
    "strconv"
)

import "github.com/Workiva/go-datastructures/augmentedtree"

type ContigChains struct {
    contig string
    chains []*Chain
}

func (cc *ContigChains) add(chain *Chain) {
    //fmt.Printf("*ContigChains.add(%s)\n", chain)
    cc.chains = append(cc.chains, chain)
}

type LiftOver struct {
    chains map[string]*ContigChains
}

func (lo *LiftOver) add_chain(contig string, chain *Chain) {
    //fmt.Printf("*LiftOver.add_chain(%s, %s)\n", contig, chain)
    cc, exists := lo.chains[contig]
    if !exists {
        cc = new(ContigChains)
        cc.contig = contig
        lo.chains[contig] = cc
    }
    cc.add(chain)
}

func str2int64(s string) int64 {
    val, _ := strconv.ParseInt(s, 10, 64)
    return val
}

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

func (c *Chain) load_data(s *bufio.Scanner){
    //fmt.Println("*Chain.load_data()")
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

func load_chain_file(fp string) *LiftOver {
    //fmt.Println("Declare variable")
    liftover := new(LiftOver)
    //fmt.Println("Initialize map")
    liftover.chains = make(map[string]*ContigChains)
    //fmt.Println("Profit")

    f, err := os.Open(fp)
    if err != nil {
        panic(err)
    } 
    defer f.Close()

    var line string
    var chains []*Chain
    var chain *Chain

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line = strings.TrimSpace(scanner.Text())
        // Only handle non-blank lines
        if (len(line) > 0) {
            if len(line) > 6 && line[0:5] == "chain" {
                //fmt.Println(line)
                if chain != nil  {
                    //fmt.Printf("Add chain to ContigChain for %s\n", chain.source_name)

                    liftover.add_chain(chain.source_name, chain)
                    //liftover.chains[chain.source_name].add(chain)
                    //for _, x := range chain.blocks {
                    //    fmt.Println(x)
                    //}
                }
                chain = new(Chain)
                chain.fromString(line)
                //fmt.Println(chain)
                chains = append(chains, chain)
            } else if string(line[0]) != "#" {
                chain.load_data(scanner)
            }

            //cols = strings.Split(line, "\t")
        }
    }

    if chain != nil {
        //fmt.Printf("Add final chain to ContigChain for %s\n", chain.source_name)
        liftover.chains[chain.source_name].add(chain)
    }
    return liftover
}

func main() {
    lo := load_chain_file("hg19ToHg38.over.chain")
    //fmt.Println(lo)
    
    for i, chain := range lo.chains["chrX"].chains {
        fmt.Printf("[%d]\t%s\n", i, chain)
        fmt.Println("-------------------------------->---------->------")
    }

    cb := lo.chains["chrX"].chains[165].blocks[0]
    fmt.Println(cb.GetOverlap("chrX", 115149220, 115149335))

    at := augmentedtree.New(1)
    //at := augmentedtree.newTree(1)
    //fmt.Println(cb)
    //at.Add(cb)
    for _, chain := range lo.chains["chrX"].chains {
        for _, block := range chain.blocks {
            at.Add(block)
            fmt.Printf("Tree.Len() %d\n", at.Len())
        }
    }
    fmt.Println(at)

    target := ChainInterval{contig: "chrX", start: 115149220, end: 115149335}
    o := at.Query(target)
    fmt.Println(o)

}