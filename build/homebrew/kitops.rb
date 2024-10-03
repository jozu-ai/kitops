class Kitops < Formula
  desc "Packaging and versioning system for AI/ML projects"
  homepage "https://KitOps.ml"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-darwin-arm64.tar.gz"
      sha256 "6486306f3365e64ebe1c94ba1eb0a90840cbaff908c248661ca533e08f87beed"
    end
    on_intel do
      url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-darwin-x86_64.tar.gz"
      sha256 "624d1a06af1820b43338a94b9d20acb33f3c21349e25b07f4f61a256b71cbf37"
    end

  end

  on_linux do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-linux-arm64.tar.gz"
      sha256 "f94ff1877439984599661174436850122a78e0e52b013fca038649ea57adb87b"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-linux-x86_64.tar.gz"
        sha256 "6298a6adf420c32a84a4e781112a591aab0ab0db7dd1734d36d06bd115a502df"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-linux-i386.tar.gz"
        sha256 "fd4eee8c01faeb488423f4b47c771db25aded8af08c4529d2fa1422840b26f53"
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
