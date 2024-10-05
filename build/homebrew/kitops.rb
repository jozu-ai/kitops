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
      url "https://github.com/brett-hodges/kitops/releases/download/v0.4.12/kitops-linux-arm64.tar.gz"
      sha256 "ab4f0a443e80a30232d40912787df898bc4f6135f66f7b85ed18ad2f3db59af5"
    end
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.12/kitops-linux-x86_64.tar.gz"
        sha256 "395bb7c119b6bb4a537269f854c8ddf74a21cd9137c100be1da57a7b3d4b857b"
      else
        url "https://github.com/brett-hodges/kitops/releases/download/v0.4.12/kitops-linux-i386.tar.gz"
        sha256 "9d0b81649b7615054674bdd9808ca04aa18f03bd285fa75c1a9b12f31c928548"
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
