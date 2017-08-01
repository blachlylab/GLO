package main

import (
    "fmt"
    "os"
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
    cc.chains = append(cc.chains, chain)
}



type ChainFile struct {
    source_build, target_build, fp string
    chains map[string]*ContigChains
}

func (cf *ChainFile) add(chain *Chain) {
    contig := chain.source_name
    cc, exists := cf.chains[contig]
    if !exists {
        // Create missing entry
        cc = new(ContigChains)
        cc.contig = contig
        cf.chains[contig] = cc
    }
    cc.add(chain)
}


func str2int64(s string) int64 {
    val, _ := strconv.ParseInt(s, 10, 64)
    return val
}



func (cf *ChainFile) Load() {
    cf.chains = make(map[string]*ContigChains)

    f, err := os.Open(cf.fp)
    if err != nil {
        panic(err)
    } 
    defer f.Close()

    var line string
    var chain *Chain

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line = strings.TrimSpace(scanner.Text())
        // Only handle non-blank lines
        if (len(line) > 0) {
            if len(line) > 6 && line[0:5] == "chain" {
                if chain != nil  {
                    // Store this chain before moving to the next
                    // chain generation
                    cf.add(chain)
                }
                chain = new(Chain)
                chain.fromString(line)

            } else if string(line[0]) != "#" {
                // Load mapping blocks for this Chain
                chain.load_blocks(scanner)
            }
        }
    }

    if chain != nil {
        // Store the last chain.
        cf.add(chain)
    }
}
func main() {

    // Initialize ChainFile object.
    cf := new(ChainFile)
    cf.source_build = "hg19"
    cf.target_build = "hg38"
    cf.fp = "hg19ToHg38.over.chain"
    cf.Load()


    at := augmentedtree.New(1)
    for _, chain := range cf.chains["chrX"].chains {
        for _, block := range chain.blocks {
            at.Add(block)
        }
    }
    fmt.Println(at)

    target := ChainInterval{contig: "chrX", start: 115149220, end: 115149335}
    o := at.Query(target)
    fmt.Println(o)

}