type {{.Name}}Edge {
  cursor: String
  node: {{.Name}}!
}

type {{.Name}}Connection {
  edges: [{{.Name}}Edge!]!
  pageInfo: PageInfo!
}

type {{.Name}} {
  {{ range $index, $field := GraphqlFields .Fields false }}{{GraphqlField $field false}}
  {{ end }}
}

input {{.Name}}CreateInput {
  {{ range $index, $field := GraphqlFields .Fields true }}{{GraphqlField $field false}}
  {{ end }}
}

input {{.Name}}UpdateInput {
  {{ range $index, $field := GraphqlFields .Fields true }}{{GraphqlField $field true}}
  {{ end }}
}

input {{.Name}}Filter {
  ids: [ID!]
}
