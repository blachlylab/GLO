package goLiftOver

import (
    "fmt"
    "strconv"
)

import "github.com/Workiva/go-datastructures/augmentedtree"



type LiftOver struct {
    source, target string
    contigs map[string]augmentedtree.Tree
}

func str2int64(s string) int64 {
    val, _ := strconv.ParseInt(s, 10, 64)
    return val
}



func main() {

    // Initialize ChainFile object.
    cf := new(ChainFile)
    cf.source_build = "hg19"
    cf.target_build = "hg38"
    cf.fp = "hg19ToHg38.over.chain"
    cf.Load()


    at := augmentedtree.New(1)
    for _, chain := range cf.chains["chrX"] {
        for _, block := range chain.blocks {
            at.Add(block)
        }
    }
    fmt.Println(at)

    target := ChainInterval{contig: "chrX", start: 115149220, end: 115149335}
    o := at.Query(target)
    fmt.Println(o)

}