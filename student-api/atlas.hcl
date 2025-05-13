env "gorm" {
  src = data.external_schema.gorm.url
  url = "postgres://postgres:ocYHVBzBVQYoQg50iAIE@localhost:5432/studentdb?sslmode=disable"
  dev = "docker://postgres/15"
  migration {
    dir = "file://internal/db/migrations"
  }
}

data "external_schema" "gorm" {
  program = [
    "go", "run", "-mod=mod", "ariga.io/atlas-provider-gorm", "load",
    "--path", "./internal/model",
    "--dialect", "postgres"
  ]
}
