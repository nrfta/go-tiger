
  create{{.Name}}(input: {{.Name}}CreateInput!): {{.Name}} @isAuthenticated
  update{{.Name}}(id: ID!, input: {{.Name}}UpdateInput!): {{.Name}} @isAuthenticated
  delete{{.Name}}(id: ID!): {{.Name}} @isAuthenticated
}
