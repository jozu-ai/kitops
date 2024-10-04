class Kitops < Formula
  desc "Packaging and versioning system for AI/ML projects"
  homepage "https://KitOps.ml"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.0/kitops-darwin-arm64.tar.gz"
      sha256 "685da60be06fa2c485b9e4aa6a9b211cdf53ac0444d06fd4ed1c96ad3ee1215b"
    end
    on_intel do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.0/kitops-darwin-x86_64.tar.gz"
      sha256 "a7ca1c19a4b0ad74ed55201aaa5661b225cd44a086d19a32793ebc5721bc0c2f"
    end

  end

  on_linux do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.0/kitops-linux-arm64.tar.gz"
      sha256 "4baf05ae4ea7fc7ef91d2c63ff0e75ba1c24933bb8b241fd4ea409ea2ed8b273"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.0/kitops-linux-x86_64.tar.gz"
        sha256 "98605a9d740ba0ca5ce12001fe998064df1c3dbc807e278df4051f89ebca049b"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.0/kitops-linux-i386.tar.gz"
        sha256 "22b76f8bfbe941b9cda90eb3a785738169d0dee4f56260bf7359d6e3b9f78ee1"
      end
    end
  end

  def install
    bin.install "kit"
  end

  test do
    expected_version = "Version: 0.4.0"
    actual_version = shell_output("#{bin}/kit version").strip
    assert_match expected_version, actual_version
  end
end
