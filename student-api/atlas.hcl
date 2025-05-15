env "gorm" {
  src = data.external_schema.gorm.url
  url = "env://DB_URL" 
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
