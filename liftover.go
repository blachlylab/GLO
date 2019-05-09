package glo

import (
    "strconv"
)

import "github.com/Workiva/go-datastructures/augmentedtree"



type LiftOverTree struct {
    reference, query string
    contigs map[string]*augmentedtree.Tree
}

type LiftOverKey struct {
    reference, query string
}

type LiftOver struct {
    trees map[LiftOverKey]*LiftOverTree
}

func (self *LiftOver) Init() {
    self.trees = make(map[LiftOverKey]*LiftOverTree)
}

func (self *LiftOver) Lift(reference_build, query_build string, query *ChainInterval) []*ChainInterval {
    var overlaps []*ChainInterval

    // Generate key for accessing the correct tree
    key := new(LiftOverKey)
    key.reference = reference_build
    key.query = query_build

    lotree, lotree_exists := self.trees[*key]
    if lotree_exists {
        atree, atree_exists := lotree.contigs[query.contig]
        if atree_exists {
            for _, res := range (*atree).Query(query) {
                // Use type assertion to specify that the Interval
                // being returned is a *ChainBlock, and a call to
                // GetOverLap() to get the overlapped interval.
                overlap := res.(*ChainLink).GetOverlap(query)
                if (overlap.size() > 0) {
                    overlaps = append(overlaps, overlap)
                }
            }
        }
    }

    return overlaps
}

func (self *LiftOver) LoadChainFile(source, target, fp string) {
    // Initialize new ChainFile object.
    cf := new(ChainFile)
    cf.referenceBuild = source
    cf.queryBuild = target
    cf.filepath = fp
    // Load data
    cf.Load()

    // Generate key for tree mapping
    key := new(LiftOverKey)
    key.reference = source
    key.query = target

    tree := new(LiftOverTree)
    tree.reference = source
    tree.query = target
    tree.contigs = make(map[string]*augmentedtree.Tree)
    for contig, chains := range cf.chains {
        _, exists := tree.contigs[contig]
        if !exists {
            t := augmentedtree.New(1)
            tree.contigs[contig] = &t
        }
        atree := *(tree.contigs[contig])
        for _, chain := range chains {
            for _, link := range chain.links {
                atree.Add(link)
            }
        }
    }

    self.trees[*key] = tree
}


func str2int64(s string) int64 {
    val, _ := strconv.ParseInt(s, 10, 64)
    return val
}
