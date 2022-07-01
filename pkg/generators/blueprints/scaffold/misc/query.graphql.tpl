
  {{ToLowerCamel .Name}}ById(id: ID!): {{.Name}} @isAuthenticated
  {{ToLowerCamel .NamePlural}}(
    page: PageArgs
    filter: {{.Name}}Filter
  ): {{.Name}}Connection @isAuthenticated
}
