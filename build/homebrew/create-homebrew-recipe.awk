# Usage: awk -f create-homebrew-recipe.awk homebrew-metadata.txt

BEGIN {
    template = "kitops.rb.template";
    recipe = "kitops.rb";
}

{
  # Read relevant fields from input file
  # (intended to be homebrew-metadata.txt)
  shas[$3]=$1;
  urls[$3]=$2;
  version[$3]=$4;
}

END {
  # Read a line from template, replace special fields, 
  # and print result to recipe file
  while ((getline ln < template) > 0)
  {
    sub(/url @@darwin-arm64/, "url " urls["darwin-arm64"], ln);
    sub(/sha256 @@darwin-arm64/, "sha256 " shas["darwin-arm64"], ln);

    sub(/url @@darwin-x86_64/, "url " urls["darwin-x86_64"], ln);
    sub(/sha256 @@darwin-x86_64/, "sha256 " shas["darwin-x86_64"], ln);

    sub(/url @@linux-arm64/, "url "  urls["linux-arm64"], ln);
    sub(/sha256 @@linux-arm64/, "sha256 " shas["linux-arm64"], ln);

    sub(/url @@linux-x86_64/, "url "  urls["linux-x86_64"], ln);
    sub(/sha256 @@linux-x86_64/, "sha256 " shas["linux-x86_64"], ln);

    sub(/url @@linux-i386/, "url "  urls["linux-i386"], ln);
    sub(/sha256 @@linux-i386/, "sha256 " shas["linux-i386"], ln);
   
    sub(/@@version/, version["darwin-arm64"], ln);
 
    print(ln) > recipe;
  }

  # Close template and recipe fields
  close(recipe);
  close(template);
}