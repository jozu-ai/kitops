class Kitops < Formula
  desc "Packaging and versioning system for AI/ML projects"
  homepage "https://KitOps.ml"
  license "Apache-2.0"

  on_macos do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-darwin-arm64.tar.gz"
      sha256 "72218d6541da87cdac27e1dd8c994646470731742217aad52728a37b1f907dc8"
    end
    on_intel do
      url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-darwin-x86_64.tar.gz"
      sha256 "3c7b82321bed488886f5b5f699b66078541ab1a8fc311c083710f80013b670e9"
    end

  end

  on_linux do
    on_arm do
      url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-linux-arm64.tar.gz"
      sha256 "9d92a6f005681e3a494b16be3073ed2ec752ba339c69eb74e5e9783857ef1f0f"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-linux-x86_64.tar.gz"
        sha256 "e9a9c8d7bb78f2478731439469dda503b70ec344e249c7366d1a8c4dfabdac4a"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/main/kitops-linux-i386.tar.gz"
        sha256 "aad06aef59209e5abdc3f36047b16ba7756491c841b8984311e9ccaae37434d3"
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
