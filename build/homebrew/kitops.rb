class Kitops < Formula
  desc "Packaging and versioning system for AI/ML projects"
  homepage "https://KitOps.ml"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url 
      sha256 
    end
    on_intel do
      url 
      sha256 
    end

  end

  on_linux do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/v.0.4.14/kitops-linux-arm64.tar.gz"
      sha256 "a8739d3980b4ced1a4926c1bb3594624b94e51a8d3de7310e51528a3ad89f38f"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v.0.4.14/kitops-linux-x86_64.tar.gz"
        sha256 "9189535d6cd46dd3b305d617b6f06360e9ec902f66909993422bb7efe849f6fa"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v.0.4.14/kitops-linux-i386.tar.gz"
        sha256 "dc820e7ed190b0be05e246900866086213a417d4c25816d8cf133ff42a4da35c"
      end
    end
  end

  def install
    bin.install "kit"
  end

  test do
    expected_version = "Version: "
    actual_version = shell_output("#{bin}/kit version").strip
    assert_match expected_version, actual_version
  end
end
