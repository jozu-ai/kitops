class Kitops < Formula
  desc "Packaging and versioning system for AI/ML projects"
  homepage "https://KitOps.ml"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-darwin-arm64.tar.gz"
      sha256 "8d62c78ddc55ea6efdf9c33f7392f177715289d6d7c2688c4ff6bdbd773cc490"
    end
    on_intel do
      url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-darwin-x86_64.tar.gz"
      sha256 "29816c5dd8a76e7f0af56d9f2d50e097943c4cf47ade694dd919229beacbaf11"
    end

  end

  on_linux do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-linux-arm64.tar.gz"
      sha256 "6eef018f94513b74e3710c7033941d01172d00b9f82b33abebf34a18a75d6750"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-linux-x86_64.tar.gz"
        sha256 "af7ae9699a5764b529e2883bb824c340f12399ae6580549d3f085cde2d93a905"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-linux-i386.tar.gz"
        sha256 "58fb2d183ff7dc1038e7a5267fcc2556e7315a1e161b1db457e9e31167d9270c"
      end
    end
  end

  def install
    bin.install "kit"
  end

  test do
    expected_version = "Version: main"
    actual_version = shell_output("#{bin}/kit version").strip
    assert_match expected_version, actual_version
  end
end
