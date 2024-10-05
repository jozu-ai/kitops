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
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.14/kitops-linux-arm64.tar.gz"
      sha256 "fb53c16d50750fb101cf4d2f875cb2655f16f34c85b8ad39da249754872f7af6"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.14/kitops-linux-x86_64.tar.gz"
        sha256 "c28f30c9e3b8542b1ce1b66b12402295f41e3942e42689bc3d4298f727a4f64f"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.14/kitops-linux-i386.tar.gz"
        sha256 "0089de4a65e8745610b695f4c3f133188682fa302021904fe96482689ac01335"
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
