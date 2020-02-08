package snowflake

// ParseTime returns an int64 unix timestamp in milliseconds of the snowflake ID time.
func (n *Node) ParseTime(id int64) int64 {
	return (int64(id) >> n.timeShift) + n.epoch.UnixNano()/1000000
}

// ParseServerID returns an int64 of the snowflake ID node number.
func (n *Node) ParseServerID(id int64) int64 {
	return int64(id) & n.nodeMask >> n.nodeShift
}

// ParseStep returns an int64 of the snowflake step (or sequence) number.
func (n *Node) ParseStep(id int64) int64 {
	return int64(id) & n.stepMask
}
