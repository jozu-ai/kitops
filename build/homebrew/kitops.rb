class Kitops < Formula
  desc "Packaging and versioning system for AI/ML projects"
  homepage "https://KitOps.ml"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-darwin-arm64.tar.gz"
      sha256 "fdc3ca99ce45f4f0e67f04465987c40d954022619e6162cc7347ef18e70ddfa2"
    end
    on_intel do
      url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-darwin-x86_64.tar.gz"
      sha256 "390515b550751fdaed4b7ff137efb4a0fd2debab89af1cc51933deed8359c82d"
    end

  end

  on_linux do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-linux-arm64.tar.gz"
      sha256 "37070ed840603b3606e4866904e9ad981c6cf0932a24cd6b401db0409afde7f2"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-linux-x86_64.tar.gz"
        sha256 "312aa74ec65ca9537d901f13ad43c41fd6923f77fb097bc8a293f67ac0e7fc95"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-linux-i386.tar.gz"
        sha256 "a555dc853f755b84d5021565df2dca404151ca4f93ee553ab42a0852fa5d7f46"
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
