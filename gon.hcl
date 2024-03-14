source = ["./dist/darwin_amd64/kit"]
bundle_id = "com.jozu-ai.kitops"

apple_id {
  username = "@env:AC_USERNAME"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "@env:APPLICATION_IDENTITY"
}
