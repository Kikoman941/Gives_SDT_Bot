package data

import "errors"

var (
	ERROR_MEMBER_ALREADY_EXIST  = errors.New(`ERROR #23505 duplicate key value violates unique constraint "gives_members_un"`)
	ERROR_NO_ADMINS_FOR_REFRESH = errors.New("no admins for refresh")
)
