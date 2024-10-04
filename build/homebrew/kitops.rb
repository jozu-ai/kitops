class Kitops < Formula
  desc "Packaging and versioning system for AI/ML projects"
  homepage "https://KitOps.ml"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.3/kitops-darwin-arm64.tar.gz"
      sha256 "0818eff95a867b09c9478ac8ca7b52d95b60cee2851403f9c1d2285516d08adf"
    end
    on_intel do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.3/kitops-darwin-x86_64.tar.gz"
      sha256 "8e3c0c429c488353a4c5bfde178207f36b677f2b1fd76b4fc22f7336e11b6693"
    end

  end

  on_linux do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.3/kitops-linux-arm64.tar.gz"
      sha256 "60ef16e914935465b40b5e6cb36a18728b91020ceabfa12ef5d5709683bbd500"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.3/kitops-linux-x86_64.tar.gz"
        sha256 "07dcfd6f6019b9a2c2b1ff94e2b508c060d2deaf77f459467b73c5fc67a50a2b"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.3/kitops-linux-i386.tar.gz"
        sha256 "d1ed78e7536956caa99a6b5209b52e6536254ff37cc11def5e6634bb4dcb8cc7"
      end
    end
  end

  def install
    bin.install "kit"
  end

  test do
    expected_version = "Version: 0.4.3"
    actual_version = shell_output("#{bin}/kit version").strip
    assert_match expected_version, actual_version
  end
end
