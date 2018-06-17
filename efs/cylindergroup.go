package efs

type CylinderGroup struct {
	blocks []BasicBlock
	fs     *Filesystem
}

func (cg CylinderGroup) NumBlocks() int {
	return len(cg.blocks)
}

func (cg CylinderGroup) DataNodes() {

}
