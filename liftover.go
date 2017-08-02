package glo

import (
    "fmt"
    "strconv"
)

import "github.com/Workiva/go-datastructures/augmentedtree"



type LiftOverTree struct {
    Source, Target string
    Contigs map[string]*augmentedtree.Tree
}

type LiftOverKey struct {
    Source, Target string
}

type LiftOver struct {
    Trees map[LiftOverKey]*LiftOverTree
}

func (self *LiftOver) Init() {
    self.Trees = make(map[LiftOverKey]*LiftOverTree)
}

func (self *LiftOver) Lift(source_build, target_build string, target *ChainInterval) []*ChainInterval {
    var overlaps []*ChainInterval

    // Generate key for accessing the correct tree
    key := new(LiftOverKey)
    key.Source = source_build
    key.Target = target_build

    lotree, lotree_exists := self.Trees[*key]
    if lotree_exists {
        atree, atree_exists := lotree.Contigs[target.Contig]
        if atree_exists {
            for _, res := range (*atree).Query(target) {
                // Use type assertion to specify that the Interval
                // being returned is a *ChainBlock.
                match := res.(*ChainBlock).Target
                overlaps = append(overlaps, match)
            }
        }
    }

    return overlaps
}

func (self *LiftOver) LoadChainFile(source, target, fp string) {
    // Initialize new ChainFile object.
    fmt.Println("Initializing new ChainFile..")
    cf := new(ChainFile)
    cf.SourceBuild = source
    cf.TargetBuild = target
    cf.Filepath = fp
    // Load data
    fmt.Println("Loading data into ChainFile..")
    cf.Load()
    fmt.Println("Done.")

    // Generate key for tree mapping
    key := new(LiftOverKey)
    key.Source = source
    key.Target = target

    tree := new(LiftOverTree)
    tree.Source = source
    tree.Target = target
    tree.Contigs = make(map[string]*augmentedtree.Tree)
    for contig, chains := range cf.Chains {
        _, exists := tree.Contigs[contig]
        if !exists {
            t := augmentedtree.New(1)
            tree.Contigs[contig] = &t
        }
        atree := *(tree.Contigs[contig])
        for _, chain := range chains {
            for _, block := range chain.Blocks {
                atree.Add(block)
            }
        }
    }

    self.Trees[*key] = tree



}


func str2int64(s string) int64 {
    val, _ := strconv.ParseInt(s, 10, 64)
    return val
}
