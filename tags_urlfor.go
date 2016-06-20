package pongo2

import (
	"bytes"

	"github.com/astaxie/beego"
	p2 "gopkg.in/flosch/pongo2.v3"
)

type tagURLForNode struct {
	objectEvaluators []p2.IEvaluator
}

func (node *tagURLForNode) Execute(ctx *p2.ExecutionContext, buffer *bytes.Buffer) *p2.Error {
	args := make([]string, len(node.objectEvaluators))
	for i, ev := range node.objectEvaluators {
		obj, err := ev.Evaluate(ctx)
		if err != nil {
			return err
		}
		args[i] = obj.String()
	}

	params := make([]interface{}, len(args)-1)
	for i := range params {
		params[i] = args[i+1]
	}

	url := beego.URLFor(args[0], params...)

	buffer.WriteString(url)
	return nil
}

// tagURLForParser implements a {% urlfor %} tag.
//
// urlfor takes one argument for the controller, as well as any number of key/value pairs for additional URL data.
// Example: {% urlfor "UserController.View" ":slug" "oal" %}
func tagURLForParser(doc *p2.Parser, start *p2.Token, arguments *p2.Parser) (p2.INodeTag, *p2.Error) {
	evals := []p2.IEvaluator{}
	for arguments.Remaining() > 0 {
		expr, err := arguments.ParseExpression()
		evals = append(evals, expr)
		if err != nil {
			return nil, err
		}
	}

	if (len(evals)-1)%2 != 0 {
		return nil, arguments.Error("URL takes one argument for the controller and any number of optional pairs of key/value pairs.", nil)
	}

	return &tagURLForNode{evals}, nil
}

func init() {
	p2.RegisterTag("urlfor", tagURLForParser)
}
