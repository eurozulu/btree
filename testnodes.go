package btree

func testThreeNode() *node[int, string] {
	return &node[int, string]{
		entries: []*entry[int, string]{{20, "twenty"}},
		children: []*node[int, string]{
			{entries: []*entry[int, string]{{10, "ten"}}},
			{entries: []*entry[int, string]{{30, "thirty"}}},
		},
	}
}

func testNodeWithChildren() *node[int, string] {
	entries := testEntries()
	half := len(entries) / 2
	return &node[int, string]{
		entries: []*entry[int, string]{entries[half]},
		children: []*node[int, string]{
			{entries: entries[:half]},
			{entries: entries[half+1:]},
		},
	}
}

func testEntries() []*entry[int, string] {
	names := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	entries := make([]*entry[int, string], len(names))
	for i, s := range names {
		entries[i] = &entry[int, string]{
			key:   i,
			value: s,
		}
	}
	return entries
}

func buildTestNode(key int, value string) RootNode[int, string] {
	e := &entry[int, string]{key, value}
	es := []*entry[int, string]{e}
	return &node[int, string]{
		entries:  es,
		children: nil,
	}
}

func BuildTestTreeRoot() RootNode[int, string] {
	return &node[int, string]{
		entries: []*entry[int, string]{
			{40, "forty"},
		},
		children: []*node[int, string]{
			{
				entries: []*entry[int, string]{
					{20, "twenty"},
				},
				children: []*node[int, string]{
					{
						entries: []*entry[int, string]{
							{10, "ten"},
						},
						children: []*node[int, string]{
							{
								entries: []*entry[int, string]{
									{1, "one"},
									{3, "three"},
									{9, "nine"},
								},
							},
							{
								entries: []*entry[int, string]{
									{11, "eleven"},
									{12, "twelve"},
									{13, "fiveteen"},
								},
							},
						},
					},
					{
						entries: []*entry[int, string]{
							{30, "thirty"},
						},
						children: []*node[int, string]{
							{
								entries: []*entry[int, string]{
									{24, "twenty four"},
									{29, "twenty nine"},
								},
							},
							{
								entries: []*entry[int, string]{
									{31, "thirty one"},
									{32, "thirty two"},
									{34, "thirty four"},
								},
							},
						},
					},
				},
			},
			{
				entries: []*entry[int, string]{
					{60, "sixty"},
					{80, "eighty"},
				},
				children: []*node[int, string]{
					{
						entries: []*entry[int, string]{
							{50, "fifty"},
							{55, "fifty five"},
						},
						children: []*node[int, string]{
							{
								entries: []*entry[int, string]{
									{49, "fourty nine"},
								},
							},
							{
								entries: []*entry[int, string]{
									{51, "fifty one"},
									{53, "fifty three"},
								},
							},
							{
								entries: []*entry[int, string]{
									{58, "fifty eight"},
									{59, "fifty nine"},
								},
							},
						},
					},
					{
						entries: []*entry[int, string]{
							{70, "seventy"},
						},
						children: []*node[int, string]{
							{
								entries: []*entry[int, string]{
									{69, "sixty nine"},
								},
							},
							{
								entries: []*entry[int, string]{
									{71, "seventy one"},
									{72, "seventy two"},
								},
							},
						},
					},
					{
						entries: []*entry[int, string]{
							{85, "eighty five"},
							{90, "ninty"},
						},
						children: []*node[int, string]{
							{
								entries: []*entry[int, string]{
									{83, "eighty three"},
								},
							},
							{
								entries: []*entry[int, string]{
									{87, "eighty seven"},
									{89, "eighty nine"},
								},
							},
							{
								entries: []*entry[int, string]{
									{91, "ninty one"},
									{92, "ninty two"},
								},
							},
						},
					},
				},
			},
		},
	}
}
