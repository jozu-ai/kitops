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
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.17/kitops-linux-arm64.tar.gz"
      sha256 "1d1563e2a93fe12ae4b0a02a22aaaf199062bcd50fe5703ae69c5923445f4642"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.17/kitops-linux-x86_64.tar.gz"
        sha256 "d4b1ecf1b4a8c114b86ff37e29c8a0fc0fc7968e3b827d08cc6b7192778be499"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.17/kitops-linux-i386.tar.gz"
        sha256 "3975f0f3eea94874423b23f23d298558a6653dcfb7a30f471e1cddb58d9fa115"
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
