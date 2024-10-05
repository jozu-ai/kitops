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
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.19/kitops-linux-arm64.tar.gz"
      sha256 "56cebc02a0ce8c2ddad0caf668226986398dfe84b4513b600a12b4797d1a0927"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.19/kitops-linux-x86_64.tar.gz"
        sha256 "f93c6d554d810c6587a952adf91d76a8f5b1ae96d5360bd9cf551a4cc09e59c3"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.19/kitops-linux-i386.tar.gz"
        sha256 "c8420c08e8731148f464e22b4c849232de7116ee04c4cc5c0e3bd1613b87cf36"
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
