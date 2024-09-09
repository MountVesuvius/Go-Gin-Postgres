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
