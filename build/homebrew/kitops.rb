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
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.18/kitops-linux-arm64.tar.gz"
      sha256 "bcf92b2cdd0006c4c35fbf138bd36d1229252bef56190038ceef8e49d103e0c5"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.18/kitops-linux-x86_64.tar.gz"
        sha256 "6a725d1f8cf70c5a14cc91bbcc9bbd6039c2f4ce7c31d2a4c541252b4829aa90"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.18/kitops-linux-i386.tar.gz"
        sha256 "e355038d809c0e9e5a237a8a7f7dcf64abd8a0577ec05583cd2ac2876e7d4373"
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
