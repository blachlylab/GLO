package goLiftOver

import (
    "os"
    "bufio"
    "strings"
)


type ChainFile struct {
    source_build, target_build, fp string
    chains map[string][]*Chain
}

func (cf *ChainFile) add(chain *Chain) {
    contig := chain.source_name
    cf.chains[contig] = append(cf.chains[contig], chain)
}

func (cf *ChainFile) Load() {
    cf.chains = make(map[string][]*Chain)

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