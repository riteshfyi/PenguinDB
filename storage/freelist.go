package storage

type Freelist struct {
	head uint64

	get func(uint64) BNode  //derefernce a pointer => ezpz
	new func(BNode) uint64  //append a new page
	use func(uint64, BNode) //reuse a page
}

func (fl *Freelist) Total() int {
	return 0
}

func (fl *Freelist) Get(topn int) uint64 {
	assert(0 <= topn && topn < fl.Total())
	node := fl.get(fl.head)
	for flnSize(node) <= topn {
		topn -= flnSize(node)
		next := flnNext(node)
		assert(next != 0)
		node = fl.get(next)
	}
	return flnPtr(node, flnSize(node)-topn-1)
}

func (fl *Freelist) Update(popn int, freed []uint64) {
	assert(popn <= fl.Total())

	if popn == 0 && len(freed) == 0 {
		return
	}

	total := fl.Total()
	reuse := []uint64{}

	for fl.head != 0 && len(reuse)*FREE_LIST_CAP < len(freed) {
		node := fl.get(fl.head)
		freed = append(freed, fl.head)

		if popn >= flnSize(node) {
			popn -= flnSize(node)
		} else {
			remain := flnSize(node) - popn
			popn = 0

			for remain > 0 && len(reuse)&FREE_LIST_CAP < len(freed)+remain {
				remain--
				reuse = append(reuse, flnPtr(node, remain))
			}

			for i := 0; i < remain; i++ {
				freed = append(freed, flnPtr(node, i))
			}
		}

		total -= flnSize(node)
		fl.head = flnNext(node)
	}
	assert(len(reuse)*FREE_LIST_CAP >= len(freed) || fl.head == 0)

	flPush(fl, freed, reuse)

	flnSetTotal(fl.get(fl.head), uint64(total+len(freed)))
}

func flPush(fl *Freelist, freed []uint64, reuse []uint64) {

	for len(freed) > 0 {
		new := BNode{make([]byte, BTREE_PAGE_SIZE)}

		size := len(freed)

		if size > FREE_LIST_CAP {
			size = FREE_LIST_CAP
		}

		flnSetHeader(new, uint16(size), fl.head)

		for i, ptr := range freed[:size] {
			flnSetPtr(new, i, ptr)
		}

		freed = freed[size:]

		if len(reuse) > 0 {
			oldPtr := reuse[0]
			reuse = reuse[1:]
			fl.use(oldPtr, new)
			fl.head = oldPtr
		} else {
			fl.head = fl.new(new)
		}

	}

	assert(len(reuse) == 0)
}
