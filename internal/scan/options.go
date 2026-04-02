package scan

// Options controls the inventory engine.
type Options struct {
	Root            string
	Hidden          bool
	MaxDepth        int // 0 = unlimited
	FollowSymlinks  bool
	NoGitignore     bool
}
