package efs

// CylinderGroup is an EFS concept but doesn't actually
// end up being used in this implementation. I suspect
// caching entire CylinderGroups together might be a
// useful way to do it if you were actually doing a
// filesystem implementation instead of a hacky "just
// give me a tarball" tool :P
type CylinderGroup struct {
	blocks []BasicBlock
	fs     *Filesystem
}
