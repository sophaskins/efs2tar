package efs

type CylinderGroup struct {
	blocks []BasicBlock
	fs     *Filesystem
}

func (cg CylinderGroup) Inodes() []Inode {
	numInodes := cg.fs.SuperBlock().CGInodeSize
	inodes := make([]Inode, numInodes)
	for i := 0; i < int(numInodes); i++ {
		inodes[i] = cg.blocks[i].ToInode()
	}

	return inodes
}

func (cg CylinderGroup) NumBlocks() int {
	return len(cg.blocks)
}

func (cg CylinderGroup) DataNodes() {

}
