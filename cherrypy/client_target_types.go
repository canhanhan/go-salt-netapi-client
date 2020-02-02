package cherrypy

/*
TargetType indicates Salt API which targing mode to use.

https://docs.saltstack.com/en/latest/ref/clients/index.html#salt.client.LocalClient.cmd
https://docs.saltstack.com/en/latest/topics/targeting/index.html#advanced-targeting-methods

See the constants available in this file for possible values.
*/
type TargetType string

const (
	// Glob Bash glob completion
	Glob TargetType = "glob"
	// PCRE Perl style regular expression
	PCRE = "pcre"
	// List Python list of hosts
	List = "list"
	// Grain Match based on a grain comparison
	Grain = "grain"
	// GrainPCRE Grain comparison with a regex
	GrainPCRE = "grain_pcre"
	// Pillar data comparison
	Pillar = "pillar"
	// PillarPCRE pillar data comparison with a regex
	PillarPCRE = "pillar_pcre"
	// NodeGroup Match on nodegroup
	NodeGroup = "nodegroup"
	// Range Use a Range server for matching
	Range = "range"
	// Compound a compound match string
	Compound = "compound"
	// IPCIDR match based on Subnet (CIDR notation) or IPv4 address.
	IPCIDR = "ipcidr"
)

var targetTypes = map[string]TargetType{
	"glob":        Glob,
	"pcre":        PCRE,
	"list":        List,
	"grain":       Grain,
	"grain_pcre":  GrainPCRE,
	"pillar":      Pillar,
	"pillar_pcre": PillarPCRE,
	"nodegroup":   NodeGroup,
	"range":       Range,
	"compound":    Compound,
	"ipcidr":      IPCIDR,
}

type Target interface {
	GetTarget() interface{}
	GetType() TargetType
}

type ListTarget struct {
	Targets []string
}

func (t ListTarget) GetTarget() interface{} {
	return t.Targets
}

func (t ListTarget) GetType() TargetType {
	return List
}

type ExpressionTarget struct {
	Expression string
	Type       TargetType
}

func (t ExpressionTarget) GetTarget() interface{} {
	return t.Expression
}

func (t ExpressionTarget) GetType() TargetType {
	return t.Type
}
