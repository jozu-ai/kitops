class Kitops < Formula
  desc "Packaging and versioning system for AI/ML projects"
  homepage "https://KitOps.ml"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.1/kitops-darwin-arm64.tar.gz"
      sha256 "86d1c9371cad5f1f4f0a6023918af8ee0365ba64b5d1e4ac9499febe39531efa"
    end
    on_intel do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.1/kitops-darwin-x86_64.tar.gz"
      sha256 "ca4e274e60fb1d1918e16e1423881fd17b2b0d85627a539a7b4e5ad1bdb1d6ab"
    end

  end

  on_linux do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.1/kitops-linux-arm64.tar.gz"
      sha256 "69c4734021a849a5ef3307b87b9cf6a755cee7b87342b63e75afa7167379508a"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.1/kitops-linux-x86_64.tar.gz"
        sha256 "16ffde2a6cde3eaba368d4b10036e0ff26fbd50d062a2f95c4b901a1d7b4d989"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.1/kitops-linux-i386.tar.gz"
        sha256 "8c6080d80548a6db953ef40a785b0843378cca12eb6a142b7626ca1d1c4cbb18"
      end
    end
  end

  def install
    bin.install "kit"
  end

  test do
    expected_version = "Version: 0.4.1"
    actual_version = shell_output("#{bin}/kit version").strip
    assert_match expected_version, actual_version
  end
end
