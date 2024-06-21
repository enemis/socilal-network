package post

type PostStatusValidationRule struct {
}

func NewPostStatusValidationRule() *PostStatusValidationRule {
	return &PostStatusValidationRule{}
}

func (r PostStatusValidationRule) IsSplittableField() bool {
	return false
}

func (r PostStatusValidationRule) Rule() string {
	return "post_status"
}

func (r PostStatusValidationRule) Message() string {
	return "Invalid post status"
}

func (r PostStatusValidationRule) Field() string {
	return ""
}

func (r PostStatusValidationRule) Callback() func(interface{}) bool {
	return func(val interface{}) bool {
		_, ok := val.(PostStatus)
		if !ok {
			return false
		}

		return true
	}
}
