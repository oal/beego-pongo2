package pongo2

import (
	"github.com/astaxie/beego"
	p2 "github.com/flosch/pongo2"
)

type tagURLNode struct {
	objectEvaluators []p2.INodeEvaluator
}

func (node *tagURLNode) Execute(ctx *p2.ExecutionContext) (string, error) {
	args := make([]string, len(node.objectEvaluators))
	for i, ev := range node.objectEvaluators {
		obj, err := ev.Evaluate(ctx)
		if err != nil {
			return "", err
		}
		args[i] = obj.String()
	}

	url := beego.UrlFor(args[0], args[1:]...)
	return url, nil
}

// tagURLForParser implements a {% urlfor %} tag.
//
// urlfor takes one argument for the controller, as well as any number of key/value pairs for additional URL data.
// Example: {% url "UserController.View" ":slug" "oal" %}
func tagURLForParser(doc *p2.Parser, start *p2.Token, arguments *p2.Parser) (p2.INodeTag, error) {
	if (arguments.Count()-1)%2 != 0 {
		return nil, arguments.Error("URL takes one argument for the controller and any number of optional pairs of key/value pairs.", nil)
	}

	evals := []p2.INodeEvaluator{}
	for arguments.Remaining() > 0 {
		eval, err := arguments.ParseExpression()
		evals = append(evals, eval)
		if err != nil {
			return nil, arguments.Error(err.Error(), nil)
		}
	}

	return &tagURLNode{evals}, nil
}

func init() {
	p2.RegisterTag("urlfor", tagURLForParser)
}
