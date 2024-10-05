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
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.11/kitops-linux-arm64.tar.gz"
      sha256 "ea2ceb6537d3f8fa6ca8415aef51d8824d1aff78170e560a0cd557e8610dc411"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.11/kitops-linux-x86_64.tar.gz"
        sha256 "7c871e94bb4a0eb48a58302ecb7af9c42fbd25c7384c9c4ac86811ba7560a0c0"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.11/kitops-linux-i386.tar.gz"
        sha256 "27d6e81fe5b8155d20a3a3e43d3c44741b68a7e3e462e5554f349320e20b69bf"
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
