provider "aws" {

  access_key = "mock_access_key"
  secret_key = "mock_secret_key"
  region     = "eu-west-2"


  skip_credentials_validation = true
  skip_requesting_account_id  = true
  skip_metadata_api_check     = true
  s3_use_path_style           = true

  endpoints {
    sns = "http://localhost:4566"
  }
}

terraform {
  backend "local" {}
}
