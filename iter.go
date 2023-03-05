package gocrop

type CroppableLoadIterator struct {
	paths   []string
	current int
}

func newCroppableLoadIterator(paths []string) *CroppableLoadIterator {
	return &CroppableLoadIterator{
		paths: paths,
	}
}

func (i *CroppableLoadIterator) Reset() {
	i.current = 0
}

func (i *CroppableLoadIterator) Current() (*Croppable, error) {
	path := i.paths[i.current]
	return LoadCroppable(path)
}

func (i *CroppableLoadIterator) Next() {
	i.current += 1
}

func (i *CroppableLoadIterator) Valid() bool {
	return i.current < len(i.paths)
}
