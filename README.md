# Go LiftOver

Once the package is imported, the liftover capability can be easily utilized
by first initializing the LiftOver struct for the build-to-build mappings
that will be lifted over, and then passing it requests as needed, e.g.


// Initialize new object
liftover := new(LiftOver)
liftover.Init()

// Load in the hg19 to hg38 liftover
liftover.LoadChainFile("hg19", "hg38", "hg19ToHg38.over.chain")


// Create a target for the liftover, request for it to be lifted
// over from hg19 to hg38.
target := ChainInterval{Contig: "chrX", Start: 115149220, End: 115149335}
o := liftover.Lift("hg19", "hg38", &target)
fmt.Println(o)

