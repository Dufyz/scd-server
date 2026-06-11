package errors

import "errors"

var (
	ErrMessageFKChatId = errors.New("foreign key constraint failed: chat_id does not exist")
	ErrMessageNotFound = errors.New("message not found")
)
