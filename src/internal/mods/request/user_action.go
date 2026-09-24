package request

type ActionCreateReq struct {
	UserID     uint   `json:"userId" label:"用户ID" validate:"required"`
	TargetID   uint   `json:"targetId" label:"目标ID" validate:"required"`
	ActionType string `json:"actionType" label:"操作类型" validate:"required"`
	TargetType string `json:"targetType" label:"目标类型" validate:"required"`
}
