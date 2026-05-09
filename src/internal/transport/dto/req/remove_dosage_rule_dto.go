package req

import "github.com/google/uuid"

type RemoveDosageRuleDto struct {
	RuleID uuid.UUID `json:"rule_id" binding:"required"`
}
