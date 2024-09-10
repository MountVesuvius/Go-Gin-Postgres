package dto

type AuthenticateUser struct {
    Email string
    Password string
}

type GetUserById struct {
    Id string
}

type Body struct {
    Token string
}

// wip
type DisplayUser struct {
    Name string
    Email string
    Role string
}
