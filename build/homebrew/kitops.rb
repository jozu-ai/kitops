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
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.13/kitops-linux-arm64.tar.gz"
      sha256 "42b9571631ecbe3e12a83fa38a7c55fcba13c5b302adc058e748d113c99fbed1"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.13/kitops-linux-x86_64.tar.gz"
        sha256 "0c221e6751fa47abd23531969f176acd333c29e67b9d37d4697df67cd6915c30"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.13/kitops-linux-i386.tar.gz"
        sha256 "573df0fb014a109c304ceff2d79cb98d00f7d9533a233fae92008b4319195107"
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
