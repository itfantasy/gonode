package gonode

const (
	// version
	VERSION string = "v0.9.0.1"
	// reg channel
	CHAN_REG string = "gonode_reg"
	// log channel
	CHAN_LOG string = "gonode_log"
	// monitor channel
	CHAN_MONI string = "gonode_moni"
	// the supervisor role
	SUPERVISOR string = "supervisor"
	// when you set the backends to allnodes,
	// the node will try to conn to everynode in the cluster
	ALLNODES string = "allnodes"
)
