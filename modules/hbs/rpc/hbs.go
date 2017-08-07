package rpc

import (
	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/common/rpc"
	"github.com/Cepave/open-falcon-backend/modules/hbs/service"
)

func (t *Hbs) GetExpressions(req model.NullRpcRequest, reply *model.ExpressionResponse) (err error) {
	defer rpc.HandleError(&err)()

	exps, err := service.Expressions()
	for _, ne := range exps {
		oe := &model.Expression{
			Id:         ne.ID,
			Metric:     ne.Metric,
			Tags:       ne.Tags,
			Func:       ne.Func,
			Operator:   ne.Operator,
			RightValue: ne.RightValue,
			MaxStep:    ne.MaxStep,
			Priority:   ne.Priority,
			Note:       ne.Note,
			ActionId:   ne.ActionID,
		}
		reply.Expressions = append(reply.Expressions, oe)
	}
	return err
}

func (t *Hbs) GetStrategies(req model.NullRpcRequest, reply *model.StrategiesResponse) (err error) {
	defer rpc.HandleError(&err)()

	hostStrategies, err := service.Strategies()
	for _, nhs := range hostStrategies {
		ohs := &model.HostStrategy{
			Hostname: nhs.Hostname,
		}
		for _, ns := range nhs.Strategies {
			os := model.Strategy{
				Id:         ns.ID,
				Metric:     ns.Metric,
				Tags:       ns.Tags,
				Func:       ns.Func,
				Operator:   ns.Operator,
				RightValue: ns.RightValue,
				MaxStep:    ns.MaxStep,
				Priority:   ns.Priority,
				Note:       ns.Note,
				Tpl: &model.Template{
					Id:       ns.Tpl.ID,
					Name:     ns.Tpl.Name,
					ParentId: ns.Tpl.ParentID,
					ActionId: ns.Tpl.ActionID,
					Creator:  ns.Tpl.Creator,
				},
			}
			ohs.Strategies = append(ohs.Strategies, os)
		}
		reply.HostStrategies = append(reply.HostStrategies, ohs)
	}
	return err
}
