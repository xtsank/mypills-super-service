package command

import "github.com/google/uuid"

type RemoveDosageRuleCmd struct {
	RuleID uuid.UUID
}

func NewRemoveDosageRuleCmd(id uuid.UUID) *RemoveDosageRuleCmd {
	return &RemoveDosageRuleCmd{
		RuleID: id,
	}
}
